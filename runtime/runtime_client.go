// Copyright 2019 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package runtime

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"

	native_errors "github.com/haproxytech/client-native/v5/errors"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/haproxytech/client-native/v5/runtime/options"
)

// Client handles multiple HAProxy clients
type client struct {
	haproxyVersion *HAProxyVersion
	options        options.RuntimeOptions
	runtimes       []SingleRuntime
}

const (
	// Event though tune.buffsize default value is 16384,
	// it can be changed at build time. Because of that, it is more sensible
	// to have a smaller value, since it is not possible
	// to check this value at runtime.
	maxBufSize = 8192
)

func (c *client) initWithSockets(ctx context.Context, opt options.RuntimeOptions) error {
	socketPath := opt.Sockets
	c.runtimes = make([]SingleRuntime, 0)
	for process, path := range socketPath {
		runtime := SingleRuntime{}
		err := runtime.Init(ctx, path, 0, process)
		if err != nil {
			return err
		}
		c.runtimes = append(c.runtimes, runtime)
	}
	_, _ = c.GetVersion()
	return nil
}

func (c *client) initWithMasterSocket(ctx context.Context, opt options.RuntimeOptions) error {
	masterSocketPath := opt.MasterSocketData.MasterSocketPath
	nbproc := c.options.MasterSocketData.Nbproc
	if nbproc == 0 {
		nbproc = 1
	}
	if masterSocketPath == "" {
		return fmt.Errorf("master socket not configured")
	}
	c.runtimes = make([]SingleRuntime, nbproc)
	for i := 1; i <= nbproc; i++ {
		runtime := SingleRuntime{}
		err := runtime.Init(ctx, masterSocketPath, i, i)
		if err != nil {
			return err
		}
		c.runtimes[i-1] = runtime
	}
	_, _ = c.GetVersion()
	return nil
}

// GetStats returns stats from the socket
func (c *client) GetStats() models.NativeStats {
	result := make(models.NativeStats, len(c.runtimes))
	for index, runtime := range c.runtimes {
		result[index] = runtime.GetStats()
	}
	return result
}

// GetInfo returns info from the socket
func (c *client) GetInfo() (models.ProcessInfos, error) {
	result := models.ProcessInfos{}
	for _, runtime := range c.runtimes {
		i := runtime.GetInfo()
		result = append(result, &i)
	}
	return result, nil
}

var versionSync sync.Once //nolint:gochecknoglobals

// GetVersion returns info from the socket
func (c *client) GetVersion() (HAProxyVersion, error) {
	var err error
	versionSync.Do(func() {
		version := &HAProxyVersion{}
		for _, runtime := range c.runtimes {
			var response string
			response, err = runtime.ExecuteRaw("show info")
			if err != nil {
				return
			}
			for _, line := range strings.Split(response, "\n") {
				if strings.HasPrefix(line, "Version: ") {
					err = version.ParseHAProxyVersion(strings.TrimPrefix(line, "Version: "))
					if err != nil {
						return
					}
					c.haproxyVersion = version
					return
				}
			}
		}
		err = fmt.Errorf("version data not found")
	})
	if err != nil {
		return HAProxyVersion{}, err
	}

	if c.haproxyVersion == nil {
		return HAProxyVersion{}, fmt.Errorf("version data not found")
	}
	return *c.haproxyVersion, err
}

func (c *client) IsVersionBiggerOrEqual(minimumVersion *HAProxyVersion) bool {
	return IsBiggerOrEqual(minimumVersion, c.haproxyVersion)
}

// Reloads HAProxy's configuration file. Similar to SIGUSR2. Returns the startup logs.
func (c *client) Reload() (string, error) {
	var status, logs string

	if c.options.MasterSocketData == nil {
		return "", fmt.Errorf("cannot reload: not connected to a master socket")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 7}) {
		return "", fmt.Errorf("cannot reload: requires HAProxy 2.7 or later")
	}

	for _, runtime := range c.runtimes {
		output, err := runtime.ExecuteMaster("reload")
		if err != nil {
			return "", fmt.Errorf("cannot reload: %w", err)
		}
		parts := strings.SplitN(output, "\n--\n", 2)
		if len(parts) == 1 {
			// No startup logs. This happens when HAProxy is compiled without USE_SHM_OPEN.
			status = output[:len(output)-1]
		} else {
			status, logs = parts[0], parts[1]
		}
		switch status {
		case "Success=1":
			// Do nothing.
		case "Success=0":
			return logs, fmt.Errorf("failed to reload configuration")
		default:
			return logs, fmt.Errorf("reload: unknown status: %s", status)
		}
	}

	return logs, nil
}

// GetMapsPath returns runtime map file path or map id
func (c *client) GetMapsPath(name string) (string, error) {
	// we can refer to runtime map with either id or path
	if strings.HasPrefix(name, "#") { // id
		return name, nil
	}
	// if not id then sanitize filename
	name = misc.SanitizeFilename(name)

	// CLI
	if c.options.MapsDir != nil && *c.options.MapsDir != "" {
		ext := filepath.Ext(name)
		if ext != ".map" {
			name = fmt.Sprintf("%s%s", name, ".map")
		}
		p := filepath.Join(*c.options.MapsDir, name) // path
		return p, nil
	}
	// config
	maps, _ := c.ShowMaps()
	for _, m := range maps {
		basename := filepath.Base(m.File)
		if strings.TrimSuffix(basename, filepath.Ext(basename)) == name {
			return m.File, nil // path from config
		}
	}
	return "", fmt.Errorf("maps dir doesn't exist or not specified. Either use `maps-dir` CLI option or reload HAProxy if map section exists in config file")
}

// SetFrontendMaxConn set maxconn for frontend
func (c *client) SetFrontendMaxConn(frontend string, maxconn int) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetFrontendMaxConn(frontend, maxconn)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// AddServer adds a new server to a backend
func (c *client) AddServer(backend, name, attributes string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 6}) {
		return fmt.Errorf("this operation requires HAProxy 2.6 or later")
	}
	for _, runtime := range c.runtimes {
		err := runtime.AddServer(backend, name, attributes)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// DeleteServer removes a server from a backend
func (c *client) DeleteServer(backend, name string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 6}) {
		return fmt.Errorf("this operation requires HAProxy 2.6 or later")
	}
	for _, runtime := range c.runtimes {
		err := runtime.DeleteServer(backend, name)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAddr set ip [port] for server
func (c *client) SetServerAddr(backend, server string, ip string, port int) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAddr(backend, server, ip, port)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerState set state for server
func (c *client) SetServerState(backend, server string, state string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerState(backend, server, state)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerWeight set weight for server
func (c *client) SetServerWeight(backend, server string, weight string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerWeight(backend, server, weight)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerHealth set health for server
func (c *client) SetServerHealth(backend, server string, health string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerHealth(backend, server, health)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// EnableAgentCheck enable agent check for server
func (c *client) EnableAgentCheck(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// DisableAgentCheck disable agent check for server
func (c *client) DisableAgentCheck(backend, server string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.DisableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// EnableServer marks server as UP
func (c *client) EnableServer(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// DisableServer marks server as DOWN for maintenance
func (c *client) DisableServer(backend, server string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.DisableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAgentAddr set agent-addr for server
func (c *client) SetServerAgentAddr(backend, server string, addr string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentAddr(backend, server, addr)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAgentSend set agent-send for server
func (c *client) SetServerAgentSend(backend, server string, send string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentSend(backend, server, send)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// GetServerState returns server runtime state
func (c *client) GetServersState(backend string) (models.RuntimeServers, error) {
	var prevRs models.RuntimeServers
	var rs models.RuntimeServers
	for _, runtime := range c.runtimes {
		rs, _ = runtime.GetServersState(backend)
		if prevRs == nil {
			prevRs = rs
			continue
		}
		if !cmp.Equal(rs, prevRs) {
			return nil, fmt.Errorf("servers states differ in multiple runtime APIs")
		}
	}
	return rs, nil
}

// GetServerState returns server runtime state
func (c *client) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("no valid runtimes found")
	}
	var prevRs *models.RuntimeServer
	var result *models.RuntimeServer
	for _, runtime := range c.runtimes {
		result, _ = runtime.GetServerState(backend, server)
		if prevRs == nil {
			prevRs = result
			continue
		}
		if !cmp.Equal(*result, *prevRs) {
			return nil, fmt.Errorf("server states differ in multiple runtime APIs")
		}
	}
	if result == nil {
		return nil, fmt.Errorf("no data for %s/%s: %w", backend, server, native_errors.ErrNotFound)
	}
	return result, nil
}

// SetServerCheckPort set health check port for server
func (c *client) SetServerCheckPort(backend, server string, port int) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	for _, runtime := range c.runtimes {
		err := runtime.SetServerCheckPort(backend, server, port)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetTableEntry create or update a stick-table entry in the table.
func (c *client) SetTableEntry(table, key string, dataType models.StickTableEntry, process int) error {
	for _, runtime := range c.runtimes {
		if process == 0 || runtime.process == process {
			err := runtime.SetTableEntry(table, key, dataType)
			if err != nil {
				return fmt.Errorf("%s %w", runtime.socketPath, err)
			}
			return nil
		}
	}
	return fmt.Errorf("no valid runtimes found")
}

// Show tables show tables from runtime API and return it structured, if process is 0, return for all processes
func (c *client) ShowTables(process int) (models.StickTables, error) {
	tables := models.StickTables{}
	for _, runtime := range c.runtimes {
		if process == 0 || runtime.process == process {
			t, err := runtime.ShowTables()
			if err != nil {
				return nil, fmt.Errorf("%s %w", runtime.socketPath, err)
			}
			tables = append(tables, t...)
		}
	}
	return tables, nil
}

// GetTableEntries returns all entries for specified table in the given process with filters and a key
func (c *client) GetTableEntries(name string, process int, filter []string, key string) (models.StickTableEntries, error) {
	entries := models.StickTableEntries{}
	var err error
	for _, runtime := range c.runtimes {
		if runtime.process != process {
			continue
		}
		entries, err = runtime.GetTableEntries(name, filter, key)
		if err != nil {
			return nil, fmt.Errorf("%s %w", runtime.socketPath, err)
		}
		break
	}
	return entries, nil
}

// Show table show tables {name} from runtime API associated with process id and return it structured
func (c *client) ShowTable(name string, process int) (*models.StickTable, error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("no valid runtimes found")
	}
	var table *models.StickTable
	var err error
	for _, runtime := range c.runtimes {
		if runtime.process != process {
			continue
		}
		table, err = runtime.ShowTable(name)
		if err != nil {
			return nil, fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	if table == nil {
		return nil, fmt.Errorf("no data for table %s in process %d: %w", name, process, native_errors.ErrNotFound)
	}
	return table, nil
}

// ExecuteRaw does not process response, just returns its values for all processes
func (c *client) ExecuteRaw(command string) ([]string, error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("no valid runtimes found")
	}
	result := make([]string, len(c.runtimes))
	for index, runtime := range c.runtimes {
		r, err := runtime.ExecuteRaw(command)
		if err != nil {
			return nil, fmt.Errorf("%s %w", runtime.socketPath, err)
		}
		result[index] = r
	}
	return result, nil
}

// ShowMaps returns structured unique map files
func (c *client) ShowMaps() (models.Maps, error) {
	maps := models.Maps{}
	var lastErr error
	for _, runtime := range c.runtimes {
		m, err := runtime.ShowMaps()
		if err != nil {
			lastErr = err
		}

		if len(maps) == 0 {
			maps = append(maps, m...)
		} else {
			// merge unique files from all processes
			for i := 0; i < len(m); i++ {
				exists := false
				for j := 0; j < len(maps); j++ {
					if m[i].File == maps[j].File {
						exists = true
						break
					}
				}
				if !exists {
					maps = append(maps, m[i])
				}
			}
		}
	}
	if len(maps) > 0 {
		return maps, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, nil
}

// CreateMap creates a new map file with its entries
func (c *client) CreateMap(file io.Reader, header multipart.FileHeader) (*models.Map, error) {
	name, err := c.GetMapsPath(header.Filename)
	if err != nil {
		return nil, err
	}
	m, err := CreateMap(name, file)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (c *client) GetMapsDir() (string, error) {
	if c.options.MapsDir == nil {
		return "", fmt.Errorf("maps dir not set")
	}
	return *c.options.MapsDir, nil
}

// GetMap returns one structured runtime map file
func (c *client) GetMap(name string) (*models.Map, error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		m, err := runtime.GetMap(name)
		if m != nil {
			return m, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	return nil, lastErr
}

// ClearMap removes all map entries from the map file. If forceDelete is true, deletes file from disk
func (c *client) ClearMap(name string, forceDelete bool) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	if forceDelete {
		if err := os.Remove(name); err != nil {
			if os.IsNotExist(err) {
				return native_errors.ErrNotFound
			}
			return fmt.Errorf(strings.Join([]string{err.Error(), native_errors.ErrNotFound.Error()}, " "))
		}
	}

	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.ClearMap(name)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// ClearMapVersioned removes all map entries from the map file. If forceDelete is true, deletes file from disk
func (c *client) ClearMapVersioned(name, version string, forceDelete bool) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	if forceDelete {
		if err := os.Remove(name); err != nil {
			if os.IsNotExist(err) {
				return native_errors.ErrNotFound
			}
			return fmt.Errorf(strings.Join([]string{err.Error(), native_errors.ErrNotFound.Error()}, " "))
		}
	}

	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.ClearMapVersioned(name, version)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// ShowMapEntries list all map entries by map file name
func (c *client) ShowMapEntries(name string) (models.MapEntries, error) {
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, err
	}
	entries := models.MapEntries{}
	var lastErr error
	for _, runtime := range c.runtimes {
		m, err := runtime.ShowMapEntries(name)
		if err != nil {
			lastErr = err
		}

		if len(entries) == 0 {
			entries = append(entries, m...)
		} else {
			// merge unique map entries from all processes
			for i := 0; i < len(m); i++ {
				exists := false
				for j := 0; j < len(entries); j++ {
					if m[i].Key == entries[j].Key {
						exists = true
						break
					}
				}
				if !exists {
					entries = append(entries, m[i])
				}
			}
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return entries, nil
}

// ShowMapEntriesVersioned list all map entries by map file name
func (c *client) ShowMapEntriesVersioned(name, version string) (models.MapEntries, error) {
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, err
	}
	entries := models.MapEntries{}
	var lastErr error
	for _, runtime := range c.runtimes {
		m, err := runtime.ShowMapEntriesVersioned(name, version)
		if err != nil {
			lastErr = err
		}

		if len(entries) == 0 {
			entries = append(entries, m...)
		} else {
			// merge unique map entries from all processes
			for i := 0; i < len(m); i++ {
				exists := false
				for j := 0; j < len(entries); j++ {
					if m[i].Key == entries[j].Key {
						exists = true
						break
					}
				}
				if !exists {
					entries = append(entries, m[i])
				}
			}
		}
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return entries, nil
}

// AddMapPayload adds multiple entries to the map file
func (c *client) AddMapPayload(name, payload string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	if len(payload) > maxBufSize {
		return fmt.Errorf("payload exceeds max buffer size of %dB", maxBufSize)
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.AddMapPayload(name, payload)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

func parseMapPayload(entries models.MapEntries, maxBufSize int) (exceededSize bool, payload []string) {
	prevKV := ""
	currKV := ""
	data := ""
	for _, d := range entries {
		if prevKV != "" {
			data += prevKV
			prevKV = ""
		}
		kv := d.Key + " " + d.Value + "\n"
		data += kv
		switch {
		case len(data) < maxBufSize:
			currKV = data
		case len(data) == maxBufSize:
			payload = append(payload, data)
			data = ""
		case len(data) > maxBufSize:
			exceededSize = true
			if currKV == "" {
				currKV = kv
			}
			payload = append(payload, currKV)
			prevKV = d.Key + " " + d.Value + "\n"
			data = ""
			currKV = ""
		}
	}
	if len(currKV) > 0 {
		payload = append(payload, currKV)
	}
	return exceededSize, payload
}

// AddMapPayloadVersioned adds multiple entries to the map file atomically (using `prepare`, `add` and `commit` commands)
// if HAProxy version is 2.4 or higher. Otherwise performs `add map payload` command
func (c *client) AddMapPayloadVersioned(name string, entries models.MapEntries) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	canAtomicUpdate := false
	v := HAProxyVersion{Major: 2, Minor: 4}
	if c.IsVersionBiggerOrEqual(&v) {
		canAtomicUpdate = true
	}
	exceededSize, payload := parseMapPayload(entries, maxBufSize)
	var lastErr error
	for _, runtime := range c.runtimes {
		if canAtomicUpdate && exceededSize {
			var version string
			version, err = runtime.PrepareMap(name)
			if err != nil {
				lastErr = err
				continue
			}
			for i := 0; i < len(payload); i++ {
				err = runtime.AddMapPayloadVersioned(version, name, payload[i])
				if err != nil {
					lastErr = err
					continue
				}
			}
			err = runtime.CommitMap(version, name)
			if err != nil {
				lastErr = err
				continue
			}
		} else {
			err = runtime.AddMapPayload(name, payload[0])
			if err != nil {
				lastErr = err
				continue
			}
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// AddMapEntry adds an entry into the map file
func (c *client) AddMapEntry(name, key, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.AddMapEntry(name, key, value)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// AddMapEntry adds an entry into the map file
func (c *client) AddMapEntryVersioned(version, name, key, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.AddMapEntryVersioned(version, name, key, value)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

func (c *client) PrepareMap(name string) (version string, err error) {
	if len(c.runtimes) == 0 {
		return "", fmt.Errorf("no valid runtimes found")
	}
	name, err = c.GetMapsPath(name)
	if err != nil {
		return "", err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		version, err = runtime.PrepareMap(name)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return "", lastErr
	}
	return version, nil
}

func (c *client) CommitMap(version, name string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err = runtime.CommitMap(version, name)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// GetMapEntry returns one map runtime setting
func (c *client) GetMapEntry(name, id string) (*models.MapEntry, error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		m, err := runtime.GetMapEntry(name, id)
		if m != nil {
			return m, nil
		}
		if err != nil {
			lastErr = err
		}
	}
	return nil, lastErr
}

// SetMapEntry replace the value corresponding to each id in a map
func (c *client) SetMapEntry(name, id, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.SetMapEntry(name, id, value)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

// DeleteMapEntry deletes all the map entries from the map by its id
func (c *client) DeleteMapEntry(name, id string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.DeleteMapEntry(name, id)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

func (c *client) ParseMapEntries(output string) models.MapEntries {
	e := ParseMapEntries(output, false)
	return e
}

// ParseMapEntriesFromFile reads entries from file
func (c *client) ParseMapEntriesFromFile(inputFile io.Reader, hasID bool) models.MapEntries {
	return parseMapEntriesFromFile(inputFile, hasID)
}

// GetACLFile returns a the ACL file by its ID
func (c *client) GetACLFile(id string) (files *models.ACLFile, err error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("missing runtimes, cannot retrieve ACL files")
	}

	files, err = c.runtimes[0].GetACL("#" + id)
	if err != nil {
		err = errors.Wrap(err, "cannot retrieve ACL file for "+id)
	}

	return
}

// GetACLFiles returns all the ACL files
func (c *client) GetACLFiles() (files models.ACLFiles, err error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("missing runtimes, cannot retrieve ACL files")
	}

	files, err = c.runtimes[0].ShowACLS()
	if err != nil {
		err = errors.Wrap(err, "cannot retrieve ACL files")
	}

	return
}

// GetACLFilesEntries returns all the files entries for the ACL file ID
func (c *client) GetACLFilesEntries(id string) (files models.ACLFilesEntries, err error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("missing runtimes, cannot retrieve ACL files")
	}

	files, err = c.runtimes[0].ShowACLFileEntries("#" + id)
	if err != nil {
		err = errors.Wrap(err, "cannot retrieve ACL files entries for "+id)
	}

	return
}

// AddACLFileEntry adds the value for the specified ACL file entry based on its ID
func (c *client) AddACLFileEntry(id, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("missing runtimes, cannot add ACL file entry")
	}
	for _, runtime := range c.runtimes {
		if err := runtime.AddACLFileEntry(id, value); err != nil {
			return errors.Wrap(err, "cannot add ACL files entry for "+id)
		}
	}

	return nil
}

// GetACLFileEntry returns the specified file entry based on value and ACL file ID
func (c *client) GetACLFileEntry(id, value string) (fileEntry *models.ACLFileEntry, err error) {
	if len(c.runtimes) == 0 {
		return nil, fmt.Errorf("missing runtimes, cannot get ACL file entry")
	}

	var fe models.ACLFilesEntries
	if fe, err = c.runtimes[0].ShowACLFileEntries("#" + id); err != nil {
		return nil, errors.Wrap(err, "cannot retrieve ACL file entries, cannot list available ACL files")
	}

	for _, e := range fe {
		if e.ID == value {
			value = e.Value
			break
		}
	}

	if fileEntry, err = c.runtimes[0].GetACLFileEntry(id, value); err != nil {
		err = errors.Wrap(err, "cannot retrieve ACL file entry for "+id)
	}

	return
}

// DeleteACLFileEntry deletes the value for the specified ACL file entry based on its ID
func (c *client) DeleteACLFileEntry(id, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("missing runtimes, cannot delete ACL file entry")
	}
	for _, runtime := range c.runtimes {
		if err := runtime.DeleteACLFileEntry(id, value); err != nil {
			return errors.Wrap(err, "cannot delete ACL files entry for "+id)
		}
	}

	return nil
}

// AddACLAtomic adds multiple entries to the ACL file atomically (using `prepare`, `add` and `commit` commands)
// if HAProxy version is 2.4 or higher.
func (c *client) AddACLAtomic(aclID string, entries models.ACLFilesEntries) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	v := HAProxyVersion{Major: 2, Minor: 4}
	if !c.IsVersionBiggerOrEqual(&v) {
		return fmt.Errorf("not supported for HAProxy versions lower than 2.4 %w", native_errors.ErrGeneral)
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		version, err := runtime.PrepareACL(aclID)
		if err != nil {
			lastErr = err
			continue
		}
		for _, e := range entries {
			err = runtime.AddACLVersioned(version, aclID, e.Value)
			if err != nil {
				lastErr = err
				continue
			}
		}
		err = runtime.CommitACL(version, aclID)
		if err != nil {
			lastErr = err
			continue
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

func (c *client) PrepareACL(name string) (version string, err error) {
	if len(c.runtimes) == 0 {
		return "", fmt.Errorf("no valid runtimes found")
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		version, err = runtime.PrepareACL(name)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return "", lastErr
	}
	return version, nil
}

func (c *client) AddACLVersioned(version, aclID, value string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		err := runtime.AddACLVersioned(version, aclID, value)
		if err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}

func (c *client) CommitACL(version, name string) error {
	if len(c.runtimes) == 0 {
		return fmt.Errorf("no valid runtimes found")
	}
	var lastErr error
	for _, runtime := range c.runtimes {
		if err := runtime.CommitACL(version, name); err != nil {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	return nil
}
