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
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"

	native_errors "github.com/haproxytech/client-native/v2/errors"
	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// Client handles multiple HAProxy clients
type Client struct {
	ClientParams
	haproxyVersion *HAProxyVersion
	runtimes       []SingleRuntime
}

type ClientParams struct {
	MapsDir string
}

const (
	// DefaultSocketPath sane default for runtime API socket path
	DefaultSocketPath string = "/var/run/haproxy.sock"
)

// DefaultClient return runtime Client with sane defaults
func DefaultClient() (*Client, error) {
	c := &Client{}
	err := c.Init([]string{DefaultSocketPath}, "", 0)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// Init must be given path to runtime socket and nbproc that is not 0 when in master worker mode
//
// Deprecated: use InitWithSockets or InitWithMasterSocket instead
func (c *Client) Init(socketPath []string, masterSocketPath string, nbproc int) error {
	c.runtimes = make([]SingleRuntime, len(socketPath))
	for index, path := range socketPath {
		runtime := SingleRuntime{}
		err := runtime.Init(path, 0, index)
		if err != nil {
			return err
		}
		c.runtimes[index] = runtime
	}
	if masterSocketPath != "" && nbproc != 0 {
		for i := 1; i <= nbproc; i++ {
			runtime := SingleRuntime{}
			err := runtime.Init(masterSocketPath, i, i)
			if err != nil {
				return err
			}
			c.runtimes = append(c.runtimes, runtime)
		}
	}
	_, _ = c.GetVersion()
	return nil
}

func (c *Client) InitWithSockets(socketPath map[int]string) error {
	c.runtimes = make([]SingleRuntime, 0)
	for process, path := range socketPath {
		runtime := SingleRuntime{}
		err := runtime.Init(path, 0, process)
		if err != nil {
			return err
		}
		c.runtimes = append(c.runtimes, runtime)
	}
	_, _ = c.GetVersion()
	return nil
}

func (c *Client) InitWithMasterSocket(masterSocketPath string, nbproc int) error {
	if nbproc == 0 {
		nbproc = 1
	}
	if masterSocketPath == "" {
		return fmt.Errorf("master socket not configured")
	}
	c.runtimes = make([]SingleRuntime, nbproc)
	for i := 1; i <= nbproc; i++ {
		runtime := SingleRuntime{}
		err := runtime.Init(masterSocketPath, i, i)
		if err != nil {
			return err
		}
		c.runtimes[i-1] = runtime
	}
	_, _ = c.GetVersion()
	return nil
}

// GetStats returns stats from the socket
func (c *Client) GetStats() models.NativeStats {
	result := make(models.NativeStats, len(c.runtimes))
	for index, runtime := range c.runtimes {
		result[index] = runtime.GetStats()
	}
	return result
}

// GetInfo returns info from the socket
func (c *Client) GetInfo() (models.ProcessInfos, error) {
	result := models.ProcessInfos{}
	for _, runtime := range c.runtimes {
		i := runtime.GetInfo()
		result = append(result, &i)
	}
	return result, nil
}

// GetVersion returns info from the socket
func (c *Client) GetVersion() (*HAProxyVersion, error) {
	if c.haproxyVersion != nil {
		return c.haproxyVersion, nil
	}
	version := &HAProxyVersion{}
	for _, runtime := range c.runtimes {
		response, err := runtime.ExecuteRaw("show info")
		if err != nil {
			return nil, err
		}
		for _, line := range strings.Split(response, "\n") {
			if strings.HasPrefix(line, "Version: ") {
				err := version.ParseHAProxyVersion(strings.TrimPrefix(line, "Version: "))
				if err != nil {
					return nil, err
				}
				c.haproxyVersion = version
				return version, nil
			}
		}
	}
	return nil, fmt.Errorf("version data not found")
}

// GetMapsPath returns runtime map file path or map id
func (c *Client) GetMapsPath(name string) (string, error) {
	name = misc.SanitizeFilename(name)

	// we can refer to runtime map with either id or path
	if strings.HasPrefix(name, "#") { // id
		return name, nil
	}
	// CLI
	if c.MapsDir != "" {
		ext := filepath.Ext(name)
		if ext != ".map" {
			name = fmt.Sprintf("%s%s", name, ".map")
		}
		p := filepath.Join(c.MapsDir, name) // path
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
	return "", fmt.Errorf("maps dir doesn't exists or not specified. Either use `maps-dir` CLI option or reload HAProxy if map section exists in config file")
}

// SetFrontendMaxConn set maxconn for frontend
func (c *Client) SetFrontendMaxConn(frontend string, maxconn int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetFrontendMaxConn(frontend, maxconn)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAddr set ip [port] for server
func (c *Client) SetServerAddr(backend, server string, ip string, port int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAddr(backend, server, ip, port)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerState set state for server
func (c *Client) SetServerState(backend, server string, state string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerState(backend, server, state)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerWeight set weight for server
func (c *Client) SetServerWeight(backend, server string, weight string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerWeight(backend, server, weight)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerHealth set health for server
func (c *Client) SetServerHealth(backend, server string, health string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerHealth(backend, server, health)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// EnableAgentCheck enable agent check for server
func (c *Client) EnableAgentCheck(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// DisableAgentCheck disable agent check for server
func (c *Client) DisableAgentCheck(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.DisableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// EnableServer marks server as UP
func (c *Client) EnableServer(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// DisableServer marks server as DOWN for maintenance
func (c *Client) DisableServer(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.DisableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAgentAddr set agent-addr for server
func (c *Client) SetServerAgentAddr(backend, server string, addr string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentAddr(backend, server, addr)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// SetServerAgentSend set agent-send for server
func (c *Client) SetServerAgentSend(backend, server string, send string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentSend(backend, server, send)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// GetServerState returns server runtime state
func (c *Client) GetServersState(backend string) (models.RuntimeServers, error) {
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
func (c *Client) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	var prevRs *models.RuntimeServer
	var rs *models.RuntimeServer
	for _, runtime := range c.runtimes {
		rs, _ = runtime.GetServerState(backend, server)
		if prevRs == nil {
			prevRs = rs
			continue
		}
		if !cmp.Equal(*rs, *prevRs) {
			return nil, fmt.Errorf("server states differ in multiple runtime APIs")
		}
	}
	return rs, nil
}

// SetServerCheckPort set health heck port for server
func (c *Client) SetServerCheckPort(backend, server string, port int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerCheckPort(backend, server, port)
		if err != nil {
			return fmt.Errorf("%s %w", runtime.socketPath, err)
		}
	}
	return nil
}

// Show tables show tables from runtime API and return it structured, if process is 0, return for all processes
func (c *Client) ShowTables(process int) (models.StickTables, error) {
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
func (c *Client) GetTableEntries(name string, process int, filter []string, key string) (models.StickTableEntries, error) {
	var entries models.StickTableEntries
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
func (c *Client) ShowTable(name string, process int) (*models.StickTable, error) {
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
	return table, nil
}

// ExecuteRaw does not procces response, just returns its values for all processes
func (c *Client) ExecuteRaw(command string) ([]string, error) {
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
func (c *Client) ShowMaps() (models.Maps, error) {
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
func (c *Client) CreateMap(file io.Reader, header multipart.FileHeader) (*models.Map, error) {
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

// GetMap returns one structured runtime map file
func (c *Client) GetMap(name string) (*models.Map, error) {
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
func (c *Client) ClearMap(name string, forceDelete bool) error {
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

// ShowMapEntries list all map entries by map file name
func (c *Client) ShowMapEntries(name string) (models.MapEntries, error) {
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
	if len(entries) > 0 {
		return entries, nil
	}
	return nil, lastErr
}

// AddMapPayload adds multiple entries to the map file
func (c *Client) AddMapPayload(name, payload string) error {
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

// AddMapEntry adds an entry into the map file
func (c *Client) AddMapEntry(name, key, value string) error {
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

// GetMapEntry returns one map runtime setting
func (c *Client) GetMapEntry(name, id string) (*models.MapEntry, error) {
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
func (c *Client) SetMapEntry(name, id, value string) error {
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
func (c *Client) DeleteMapEntry(name, id string) error {
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

func (c *Client) ParseMapEntries(output string) models.MapEntries {
	e := ParseMapEntries(output, false)
	return e
}

// ParseMapEntriesFromFile reads entries from file
func (c *Client) ParseMapEntriesFromFile(inputFile io.Reader, hasID bool) models.MapEntries {
	return parseMapEntriesFromFile(inputFile, hasID)
}
