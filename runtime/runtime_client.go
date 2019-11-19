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

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/models"
)

//Client handles multiple HAProxy clients
type Client struct {
	runtimes []SingleRuntime
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

//Init must be given path to runtime socket and nbproc that is not 0 when in master worker mode
//
//Deprecated: use InitWithSockets or InitWithMasterSocket instead
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
	return nil
}

func (c *Client) InitWithMasterSocket(masterSocketPath string, nbproc int) error {
	if nbproc == 0 {
		nbproc = 1
	}
	if masterSocketPath == "" {
		return fmt.Errorf("Master socket not configured")
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
	return nil
}

//GetStats returns stats from the socket
func (c *Client) GetStats() models.NativeStats {
	result := make(models.NativeStats, len(c.runtimes))
	for index, runtime := range c.runtimes {
		result[index] = runtime.GetStats()
	}
	return result
}

//GetInfo returns info from the socket
func (c *Client) GetInfo() (models.ProcessInfos, error) {
	result := models.ProcessInfos{}
	for _, runtime := range c.runtimes {
		i := runtime.GetInfo()
		result = append(result, &i)
	}
	return result, nil
}

//SetFrontendMaxConn set maxconn for frontend
func (c *Client) SetFrontendMaxConn(frontend string, maxconn int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetFrontendMaxConn(frontend, maxconn)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerAddr set ip [port] for server
func (c *Client) SetServerAddr(backend, server string, ip string, port int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAddr(backend, server, ip, port)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerState set state for server
func (c *Client) SetServerState(backend, server string, state string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerState(backend, server, state)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerWeight set weight for server
func (c *Client) SetServerWeight(backend, server string, weight string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerWeight(backend, server, weight)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerHealth set health for server
func (c *Client) SetServerHealth(backend, server string, health string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerHealth(backend, server, health)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//EnableAgentCheck enable agent check for server
func (c *Client) EnableAgentCheck(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//DisableAgentCheck disable agent check for server
func (c *Client) DisableAgentCheck(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.DisableAgentCheck(backend, server)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//EnableServer marks server as UP
func (c *Client) EnableServer(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.EnableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//DisableServer marks server as DOWN for maintenance
func (c *Client) DisableServer(backend, server string) error {
	for _, runtime := range c.runtimes {
		err := runtime.DisableServer(backend, server)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerAgentAddr set agent-addr for server
func (c *Client) SetServerAgentAddr(backend, server string, addr string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentAddr(backend, server, addr)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//SetServerAgentSend set agent-send for server
func (c *Client) SetServerAgentSend(backend, server string, send string) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerAgentSend(backend, server, send)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//GetServerState returns server runtime state
func (c *Client) GetServersState(backend string) (models.RuntimeServers, error) {
	var prevRs models.RuntimeServers
	var rs models.RuntimeServers
	for _, runtime := range c.runtimes {
		rs, _ = runtime.GetServersState(backend)
		if prevRs == nil {
			continue
		}
		if !cmp.Equal(rs, prevRs) {
			return nil, fmt.Errorf("Servers states differ in multiple runtime APIs")
		}
	}
	return rs, nil
}

//GetServerState returns server runtime state
func (c *Client) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	var prevRs *models.RuntimeServer
	var rs *models.RuntimeServer
	for _, runtime := range c.runtimes {
		rs, _ = runtime.GetServerState(backend, server)
		if prevRs == nil {
			continue
		}
		if !cmp.Equal(*rs, *prevRs) {
			return nil, fmt.Errorf("Server states differ in multiple runtime APIs")
		}
	}
	return rs, nil
}

//SetServerCheckPort set health heck port for server
func (c *Client) SetServerCheckPort(backend, server string, port int) error {
	for _, runtime := range c.runtimes {
		err := runtime.SetServerCheckPort(backend, server, port)
		if err != nil {
			return fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return nil
}

//Show tables show tables from runtime API and return it structured, if process is 0, return for all processes
func (c *Client) ShowTables(process int) (models.StickTables, error) {
	tables := models.StickTables{}
	for _, runtime := range c.runtimes {
		if process == 0 || runtime.process == process {
			t, err := runtime.ShowTables()
			if err != nil {
				return nil, fmt.Errorf("%s %s", runtime.socketPath, err)
			}
			tables = append(tables, t...)
		}
	}
	return tables, nil
}

//GetTableEntries returns all entries for specified table in the given process with filters and a key
func (c *Client) GetTableEntries(name string, process int, filter []string, key string) (models.StickTableEntries, error) {
	var entries models.StickTableEntries
	var err error
	for _, runtime := range c.runtimes {
		if runtime.process != process {
			continue
		}
		entries, err = runtime.GetTableEntries(name, filter, key)
		if err != nil {
			return nil, fmt.Errorf("%s %s", runtime.socketPath, err)
		}
		break
	}
	return entries, nil
}

//Show table show tables {name} from runtime API associated with process id and return it structured
func (c *Client) ShowTable(name string, process int) (*models.StickTable, error) {
	var table *models.StickTable
	var err error
	for _, runtime := range c.runtimes {
		if runtime.process != process {
			continue
		}
		table, err = runtime.ShowTable(name)
		if err != nil {
			return nil, fmt.Errorf("%s %s", runtime.socketPath, err)
		}
	}
	return table, nil
}

//ExecuteRaw does not procces response, just returns its values for all processes
func (c *Client) ExecuteRaw(command string) ([]string, error) {
	result := make([]string, len(c.runtimes))
	for index, runtime := range c.runtimes {
		r, err := runtime.ExecuteRaw(command)
		if err != nil {
			return nil, fmt.Errorf("%s %s", runtime.socketPath, err)
		}
		result[index] = r
	}
	return result, nil
}
