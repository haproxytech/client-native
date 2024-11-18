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
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	native_errors "github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/haproxytech/client-native/v6/runtime/options"
	"golang.org/x/sync/singleflight"
)

// Client handles multiple HAProxy clients
type client struct {
	options options.RuntimeOptions
	runtime *SingleRuntime
}

const (
	// Event though tune.buffsize default value is 16384,
	// it can be changed at build time. Because of that, it is more sensible
	// to have a smaller value, since it is not possible
	// to check this value at runtime.
	maxBufSize = 8192
)

func (c *client) initWithSockets(opt options.RuntimeOptions) error {
	socketPath := opt.Socket

	runtime := &SingleRuntime{}
	masterWorkerMode := false
	err := runtime.Init(socketPath, masterWorkerMode, opt)
	if err != nil {
		return err
	}
	c.runtime = runtime
	_, _ = c.GetVersion()
	return nil
}

func (c *client) initWithMasterSocket(opt options.RuntimeOptions) error {
	masterSocketPath := opt.MasterSocketData.MasterSocketPath

	if masterSocketPath == "" {
		return errors.New("master socket not configured")
	}
	runtime := &SingleRuntime{}
	masterWorkerMode := true
	err := runtime.Init(masterSocketPath, masterWorkerMode, opt)
	if err != nil {
		return err
	}
	c.runtime = runtime

	_, _ = c.GetVersion()
	return nil
}

// GetStats returns stats from the socket
func (c *client) GetStats() models.NativeStats {
	result := c.runtime.GetStats()
	return result
}

// GetInfo returns info from the socket
func (c *client) GetInfo() (models.ProcessInfo, error) {
	result := c.runtime.GetInfo()
	return result, nil
}

var (
	haproxyVersion *HAProxyVersion        //nolint:gochecknoglobals
	versionKey     = "version"            //nolint:gochecknoglobals
	versionSfg     = singleflight.Group{} //nolint:gochecknoglobals
)

// GetVersion returns info from the socket
func (c *client) GetVersion() (HAProxyVersion, error) {
	var err error
	if haproxyVersion != nil {
		return *haproxyVersion, nil
	}
	_, err, _ = versionSfg.Do(versionKey, func() (interface{}, error) {
		version := &HAProxyVersion{}
		var response string
		response, err = c.runtime.ExecuteRaw("show info")
		if err != nil {
			return HAProxyVersion{}, err
		}
		for _, line := range strings.Split(response, "\n") {
			if strings.HasPrefix(line, "Version: ") {
				err = version.ParseHAProxyVersion(strings.TrimPrefix(line, "Version: "))
				if err != nil {
					return HAProxyVersion{}, err
				}
				haproxyVersion = version
				return HAProxyVersion{}, err
			}
			// Starting with HAProxy 3.0, there is no more "Version:" prefix.
			if len(line) > 0 && line[0] >= '3' && line[0] <= '9' {
				err = version.ParseHAProxyVersion(line)
				if err == nil {
					haproxyVersion = version
				}
				return HAProxyVersion{}, err
			}
		}
		err = errors.New("version data not found")
		return HAProxyVersion{}, err // it's dereferenced in IsVersionBiggerOrEqual
	})
	if err != nil {
		return HAProxyVersion{}, err
	}

	if haproxyVersion == nil {
		return HAProxyVersion{}, errors.New("version data not found")
	}

	return *haproxyVersion, err
}

func (c *client) IsVersionBiggerOrEqual(minimumVersion *HAProxyVersion) bool {
	return IsBiggerOrEqual(minimumVersion, haproxyVersion)
}

// Reloads HAProxy's configuration file. Similar to SIGUSR2. Returns the startup logs.
func (c *client) Reload() (string, error) {
	var status, logs string

	if c.options.MasterSocketData == nil {
		return "", errors.New("cannot reload: not connected to a master socket")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 7}) {
		return "", fmt.Errorf("cannot reload: requires HAProxy 2.7 or later but current version is %v", haproxyVersion)
	}

	output, err := c.runtime.ExecuteMaster("reload")
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
		return logs, errors.New("failed to reload configuration")
	default:
		return logs, fmt.Errorf("reload: unknown status: %s", status)
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
	return "", errors.New("maps dir doesn't exist or not specified. Either use `maps-dir` CLI option or reload HAProxy if map section exists in config file")
}

// SetFrontendMaxConn set maxconn for frontend
func (c *client) SetFrontendMaxConn(frontend string, maxconn int) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetFrontendMaxConn(frontend, maxconn)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

// AddServer adds a new server to a backend
func (c *client) AddServer(backend, name, attributes string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 6}) {
		return fmt.Errorf("this operation requires HAProxy 2.6 or later but current version is %v", haproxyVersion)
	}
	err := c.runtime.AddServer(backend, name, attributes)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// DeleteServer removes a server from a backend
func (c *client) DeleteServer(backend, name string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if !c.IsVersionBiggerOrEqual(&HAProxyVersion{Major: 2, Minor: 6}) {
		return errors.New("this operation requires HAProxy 2.6 or later")
	}
	err := c.runtime.DeleteServer(backend, name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetServerAddr set ip [port] for server
func (c *client) SetServerAddr(backend, server string, ip string, port int) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerAddr(backend, server, ip, port)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

// SetServerState set state for server
func (c *client) SetServerState(backend, server string, state string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerState(backend, server, state)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetServerWeight set weight for server
func (c *client) SetServerWeight(backend, server string, weight string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerWeight(backend, server, weight)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetServerHealth set health for server
func (c *client) SetServerHealth(backend, server string, health string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerHealth(backend, server, health)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// EnableAgentCheck enable agent check for server
func (c *client) EnableAgentCheck(backend, server string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.EnableAgentCheck(backend, server)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// DisableAgentCheck disable agent check for server
func (c *client) DisableAgentCheck(backend, server string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.DisableAgentCheck(backend, server)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// EnableServer marks server as UP
func (c *client) EnableServer(backend, server string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.EnableServer(backend, server)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// DisableServer marks server as DOWN for maintenance
func (c *client) DisableServer(backend, server string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.DisableServer(backend, server)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetServerAgentAddr set agent-addr for server
func (c *client) SetServerAgentAddr(backend, server string, addr string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerAgentAddr(backend, server, addr)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetServerAgentSend set agent-send for server
func (c *client) SetServerAgentSend(backend, server string, send string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerAgentSend(backend, server, send)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

// GetServerState returns server runtime state
func (c *client) GetServersState(backend string) (models.RuntimeServers, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	rs, _ := c.runtime.GetServersState(backend)
	if rs == nil {
		return nil, fmt.Errorf("no data for %s: %w", backend, native_errors.ErrNotFound)
	}
	return rs, nil
}

// GetServerState returns server runtime state
func (c *client) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	result, _ := c.runtime.GetServerState(backend, server)
	if result == nil {
		return nil, fmt.Errorf("no data for %s/%s: %w", backend, server, native_errors.ErrNotFound)
	}
	return result, nil
}

// SetServerCheckPort set health check port for server
func (c *client) SetServerCheckPort(backend, server string, port int) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetServerCheckPort(backend, server, port)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// SetTableEntry create or update a stick-table entry in the table.
func (c *client) SetTableEntry(table, key string, dataType models.StickTableEntry) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.SetTableEntry(table, key, dataType)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// Show tables show tables from runtime API and return it structured
func (c *client) ShowTables() (models.StickTables, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	tables, err := c.runtime.ShowTables()
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if tables == nil {
		tables = models.StickTables{}
	}

	return tables, nil
}

// GetTableEntries returns all entries for specified table with filters and a key
func (c *client) GetTableEntries(name string, filter []string, key string) (models.StickTableEntries, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	entries, err := c.runtime.GetTableEntries(name, filter, key)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if entries == nil { // return empty list, not nil, backward compatibility
		entries = models.StickTableEntries{}
	}
	return entries, nil
}

// Show table show tables {name} from runtime API and return it structured
func (c *client) ShowTable(name string) (*models.StickTable, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	table, err := c.runtime.ShowTable(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if table == nil {
		return nil, fmt.Errorf("no data for table %s: %w", name, native_errors.ErrNotFound)
	}
	return table, nil
}

// ExecuteRaw does not process response, just returns its values
func (c *client) ExecuteRaw(command string) (string, error) {
	if !c.runtime.IsValid() {
		return "", errors.New("no valid runtime found")
	}
	result, err := c.runtime.ExecuteRaw(command)
	if err != nil {
		return "", fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return result, nil
}

// ShowMaps returns structured unique map files
func (c *client) ShowMaps() (models.Maps, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	maps, err := c.runtime.ShowMaps()
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if len(maps) > 0 {
		return maps, nil
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
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return m, nil
}

func (c *client) GetMapsDir() (string, error) {
	if c.options.MapsDir == nil {
		return "", errors.New("maps dir not set")
	}
	return *c.options.MapsDir, nil
}

// GetMap returns one structured runtime map file
func (c *client) GetMap(name string) (*models.Map, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	m, err := c.runtime.GetMap(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return m, nil
}

// ClearMap removes all map entries from the map file. If forceDelete is true, deletes file from disk
func (c *client) ClearMap(name string, forceDelete bool) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if forceDelete {
		if err = os.Remove(name); err != nil {
			if os.IsNotExist(err) {
				return native_errors.ErrNotFound
			}
			return fmt.Errorf("%s %s", err.Error(), native_errors.ErrNotFound.Error())
		}
	}

	if err = c.runtime.ClearMap(name); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// ClearMapVersioned removes all map entries from the map file. If forceDelete is true, deletes file from disk
func (c *client) ClearMapVersioned(name, version string, forceDelete bool) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return err
	}
	if forceDelete {
		if err = os.Remove(name); err != nil {
			if os.IsNotExist(err) {
				return native_errors.ErrNotFound
			}
			return fmt.Errorf("%s %s", err.Error(), native_errors.ErrNotFound.Error())
		}
	}

	if err = c.runtime.ClearMapVersioned(name, version); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// ShowMapEntries list all map entries by map file name
func (c *client) ShowMapEntries(name string) (models.MapEntries, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	entries, err := c.runtime.ShowMapEntries(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if entries == nil { // https://github.com/haproxytech/dataplaneapi/issues/234, return an empty list, not nil
		entries = models.MapEntries{}
	}
	return entries, nil
}

// ShowMapEntriesVersioned list all map entries by map file name
func (c *client) ShowMapEntriesVersioned(name, version string) (models.MapEntries, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	entries, err := c.runtime.ShowMapEntriesVersioned(name, version)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if entries == nil {
		entries = models.MapEntries{}
	}
	return entries, nil
}

// AddMapPayload adds multiple entries to the map file
func (c *client) AddMapPayload(name, payload string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if len(payload) > maxBufSize {
		return fmt.Errorf("payload exceeds max buffer size of %dB", maxBufSize)
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.AddMapPayload(name, payload); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

func parseMapPayload(entries models.MapEntries, maxBufSize int) (bool, []string) {
	prevKV := ""
	currKV := ""
	data := ""
	var payload []string
	var exceededSize bool
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
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	canAtomicUpdate := false
	v := HAProxyVersion{Major: 2, Minor: 4}
	if c.IsVersionBiggerOrEqual(&v) {
		canAtomicUpdate = true
	}
	exceededSize, payload := parseMapPayload(entries, maxBufSize)
	if canAtomicUpdate && exceededSize {
		var version string
		version, err = c.runtime.PrepareMap(name)
		if err != nil {
			return fmt.Errorf("%s %w", c.runtime.socketPath, err)
		}
		for i := range payload {
			err = c.runtime.AddMapPayloadVersioned(version, name, payload[i])
			if err != nil {
				return fmt.Errorf("%s %w", c.runtime.socketPath, err)
			}
		}
		err = c.runtime.CommitMap(version, name)
		if err != nil {
			return fmt.Errorf("%s %w", c.runtime.socketPath, err)
		}
	} else {
		err = c.runtime.AddMapPayload(name, payload[0])
		if err != nil {
			return fmt.Errorf("%s %w", c.runtime.socketPath, err)
		}
	}
	return nil
}

// AddMapEntry adds an entry into the map file
func (c *client) AddMapEntry(name, key, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.AddMapEntry(name, key, value); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// AddMapEntry adds an entry into the map file
func (c *client) AddMapEntryVersioned(version, name, key, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.AddMapEntryVersioned(version, name, key, value); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

func (c *client) PrepareMap(name string) (string, error) {
	if !c.runtime.IsValid() {
		return "", errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return "", fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	version, err := c.runtime.PrepareMap(name)
	if err != nil {
		return "", fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return version, nil
}

func (c *client) CommitMap(version, name string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.CommitMap(version, name); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// GetMapEntry returns one map runtime setting
func (c *client) GetMapEntry(name, id string) (*models.MapEntry, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	m, err := c.runtime.GetMapEntry(name, id)
	if err != nil {
		return nil, fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return m, nil
}

// SetMapEntry replace the value corresponding to each id in a map
func (c *client) SetMapEntry(name, id, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.SetMapEntry(name, id, value); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	return nil
}

// DeleteMapEntry deletes all the map entries from the map by its id
func (c *client) DeleteMapEntry(name, id string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	name, err := c.GetMapsPath(name)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	if err = c.runtime.DeleteMapEntry(name, id); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
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
func (c *client) GetACLFile(id string) (*models.ACLFile, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}

	files, err := c.runtime.GetACL("#" + id)
	if err != nil {
		err = fmt.Errorf("cannot retrieve ACL file for %s: %w", id, err)
	}

	return files, err
}

// GetACLFiles returns all the ACL files
func (c *client) GetACLFiles() (models.ACLFiles, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}

	files, err := c.runtime.ShowACLS()
	if err != nil {
		err = fmt.Errorf("cannot retrieve ACL files: %w", err)
	}

	return files, err
}

// GetACLFilesEntries returns all the files entries for the ACL file ID
func (c *client) GetACLFilesEntries(id string) (models.ACLFilesEntries, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}

	files, err := c.runtime.ShowACLFileEntries("#" + id)
	if err != nil {
		err = fmt.Errorf("cannot retrieve ACL files entries for %s: %w", id, err)
	}

	return files, err
}

// AddACLFileEntry adds the value for the specified ACL file entry based on its ID
func (c *client) AddACLFileEntry(id, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.AddACLFileEntry(id, value); err != nil {
		return fmt.Errorf("cannot add ACL files entry for %s: %w", id, err)
	}

	return nil
}

// GetACLFileEntry returns the specified file entry based on value and ACL file ID
func (c *client) GetACLFileEntry(id, value string) (*models.ACLFileEntry, error) {
	if !c.runtime.IsValid() {
		return nil, errors.New("no valid runtime found")
	}
	fe, err := c.runtime.ShowACLFileEntries("#" + id)
	if err != nil {
		return nil, fmt.Errorf("cannot retrieve ACL file entries, cannot list available ACL files: %w", err)
	}

	for _, e := range fe {
		if e.ID == value {
			value = e.Value
			break
		}
	}

	fileEntry, err := c.runtime.GetACLFileEntry(id, value)
	if err != nil {
		err = fmt.Errorf("cannot retrieve ACL file entry for %s: %w", id, err)
	}

	return fileEntry, err
}

// DeleteACLFileEntry deletes the value for the specified ACL file entry based on its ID
func (c *client) DeleteACLFileEntry(id, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.DeleteACLFileEntry(id, value); err != nil {
		return fmt.Errorf("cannot delete ACL files entry for %s: %w", id, err)
	}

	return nil
}

// AddACLAtomic adds multiple entries to the ACL file atomically (using `prepare`, `add` and `commit` commands)
// if HAProxy version is 2.4 or higher.
func (c *client) AddACLAtomic(aclID string, entries models.ACLFilesEntries) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	v := HAProxyVersion{Major: 2, Minor: 4}
	if !c.IsVersionBiggerOrEqual(&v) {
		return fmt.Errorf("not supported for HAProxy versions lower than 2.4 %w", native_errors.ErrGeneral)
	}
	version, err := c.runtime.PrepareACL(aclID)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}
	for _, e := range entries {
		err = c.runtime.AddACLVersioned(version, aclID, e.Value)
		if err != nil {
			return fmt.Errorf("%s %w", c.runtime.socketPath, err)
		}
	}
	err = c.runtime.CommitACL(version, aclID)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) PrepareACL(name string) (string, error) {
	if !c.runtime.IsValid() {
		return "", errors.New("no valid runtime found")
	}
	version, err := c.runtime.PrepareACL(name)
	if err != nil {
		return "", fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return version, nil
}

func (c *client) AddACLVersioned(version, aclID, value string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	err := c.runtime.AddACLVersioned(version, aclID, value)
	if err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) CommitACL(version, name string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.CommitACL(version, name); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) SocketPath() string {
	return c.runtime.socketPath
}

func (c *client) IsStatsSocket() bool {
	return !c.runtime.masterWorkerMode
}

func (c *client) NewCertEntry(filename string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.NewCertEntry(filename); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) SetCertEntry(filename string, payload string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.SetCertEntry(filename, payload); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) CommitCertEntry(filename string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.CommitCertEntry(filename); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) AbortCertEntry(filename string) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}
	if err := c.runtime.AbortCertEntry(filename); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}

func (c *client) AddCrtListEntry(crtList string, entry CrtListEntry) error {
	if !c.runtime.IsValid() {
		return errors.New("no valid runtime found")
	}

	if err := c.runtime.AddCrtListEntry(crtList, entry); err != nil {
		return fmt.Errorf("%s %w", c.runtime.socketPath, err)
	}

	return nil
}
