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
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

// SetServerAddr set ip [port] for server
func (s *SingleRuntime) SetServerAddr(backend, server string, ip string, port int) error {
	var cmd string
	if port > 0 {
		cmd = fmt.Sprintf("set server %s/%s addr %s port %d", backend, server, ip, port)
	} else {
		cmd = fmt.Sprintf("set server %s/%s addr %s", backend, server, ip)
	}
	return s.Execute(cmd)
}

// SetServerState set state for server
func (s *SingleRuntime) SetServerState(backend, server string, state string) error {
	if !ServerStateValid(state) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	return s.Execute(cmd)
}

// SetServerWeight set weight for server
func (s *SingleRuntime) SetServerWeight(backend, server string, weight string) error {
	if !ServerWeightValid(weight) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s weight %s", backend, server, weight)
	return s.Execute(cmd)
}

// SetServerHealth set health for server
func (s *SingleRuntime) SetServerHealth(backend, server string, health string) error {
	if !ServerHealthValid(health) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s health %s", backend, server, health)
	return s.Execute(cmd)
}

// SetServerCheckPort set health heck port for server
func (s *SingleRuntime) SetServerCheckPort(backend, server string, port int) error {
	if !(port > 0 && port <= 65535) {
		return fmt.Errorf("bad request")
	}
	return s.Execute(fmt.Sprintf("set server %s/%s check-port %d", backend, server, port))
}

// EnableAgentCheck enable agent check for server
func (s *SingleRuntime) EnableAgentCheck(backend, server string) error {
	cmd := fmt.Sprintf("enable agent %s/%s", backend, server)
	return s.Execute(cmd)
}

// DisableAgentCheck disable agent check for server
func (s *SingleRuntime) DisableAgentCheck(backend, server string) error {
	cmd := fmt.Sprintf("disable agent %s/%s", backend, server)
	return s.Execute(cmd)
}

// EnableServer marks server as UP
func (s *SingleRuntime) EnableServer(backend, server string) error {
	cmd := fmt.Sprintf("enable server %s/%s", backend, server)
	return s.Execute(cmd)
}

// DisableServer marks server as DOWN for maintenanc
func (s *SingleRuntime) DisableServer(backend, server string) error {
	cmd := fmt.Sprintf("disable server %s/%s", backend, server)
	return s.Execute(cmd)
}

// SetServerAgentAddr set agent-addr for server
func (s *SingleRuntime) SetServerAgentAddr(backend, server string, addr string) error {
	cmd := fmt.Sprintf("set server %s/%s agent-addr %s", backend, server, addr)
	return s.Execute(cmd)
}

// SetServerAgentSend set agent-send for server
func (s *SingleRuntime) SetServerAgentSend(backend, server string, send string) error {
	cmd := fmt.Sprintf("set server %s/%s agent-send %s", backend, server, send)
	return s.Execute(cmd)
}

// GetServersState returns servers runtime state
func (s *SingleRuntime) GetServersState(backend string) (models.RuntimeServers, error) {
	cmd := fmt.Sprintf("show servers state %s", backend)
	result, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}
	return parseRuntimeServers(result)
}

// GetServersState returns server runtime state
func (s *SingleRuntime) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	cmd := fmt.Sprintf("show servers state %s", backend)
	result, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result, "\n")
	if strings.TrimSpace(lines[0]) != "1" {
		return nil, fmt.Errorf("unsupported output format version, supporting format version 1")
	}

	for _, line := range lines {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "1" {
			continue
		}
		fields := strings.Split(line, " ")
		if fields[3] != server {
			continue
		}
		return parseRuntimeServer(line), nil
	}
	return nil, nil
}

func parseRuntimeServers(output string) (models.RuntimeServers, error) {
	lines := strings.Split(output, "\n")
	result := models.RuntimeServers{}

	if strings.TrimSpace(lines[0]) != "1" {
		return nil, fmt.Errorf("unsupported output format version, supporting format version 1")
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "" || strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "1" {
			continue
		}
		result = append(result, parseRuntimeServer(line))
	}
	return result, nil
}

func parseRuntimeServer(line string) *models.RuntimeServer {
	fields := strings.Split(line, " ")

	if len(fields) < 19 {
		return nil
	}

	p, err := strconv.ParseInt(fields[18], 10, 64)
	var port *int64
	if err == nil {
		port = &p
	}

	admState, _ := misc.GetServerAdminState(fields[6])

	var opState string
	switch fields[5] {
	case "0":
		opState = "down"
	case "3":
		opState = "stopping"
	case "1", "2":
		opState = "up"
	}

	return &models.RuntimeServer{
		Name:             fields[3],
		Address:          fields[4],
		Port:             port,
		ID:               fields[2],
		AdminState:       admState,
		OperationalState: opState,
	}
}
