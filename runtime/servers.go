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

	stderrors "errors"

	"github.com/haproxytech/client-native/v6/errors"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
)

// AddServer adds a new server to a backend
func (s *SingleRuntime) AddServer(backend, name, attributes string) error {
	cmd := fmt.Sprintf("add server %s/%s %s", backend, name, attributes)
	return s.Execute(cmd)
}

// DeleteServer removes a server from a backend
func (s *SingleRuntime) DeleteServer(backend, name string) error {
	cmd := fmt.Sprintf("del server %s/%s", backend, name)
	return s.Execute(cmd)
}

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
		return stderrors.New("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	return s.Execute(cmd)
}

// SetServerWeight set weight for server
func (s *SingleRuntime) SetServerWeight(backend, server string, weight string) error {
	if !ServerWeightValid(weight) {
		return stderrors.New("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s weight %s", backend, server, weight)
	return s.Execute(cmd)
}

// SetServerHealth set health for server
func (s *SingleRuntime) SetServerHealth(backend, server string, health string) error {
	if !ServerHealthValid(health) {
		return stderrors.New("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s health %s", backend, server, health)
	return s.Execute(cmd)
}

// EnableServerHealth enable health check for server
func (s *SingleRuntime) EnableServerHealth(backend, server string) error {
	cmd := fmt.Sprintf("enable health %s/%s", backend, server)
	return s.Execute(cmd)
}

// SetServerCheckPort set health heck port for server
func (s *SingleRuntime) SetServerCheckPort(backend, server string, port int) error {
	if port < 1 || port > 65535 {
		return stderrors.New("bad request")
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

// DisableServer marks server as DOWN for maintenance
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

// SetServerSSL set SSL for server
func (s *SingleRuntime) SetServerSSL(backend, server string, ssl string) error {
	cmd := fmt.Sprintf("set server %s/%s ssl %s", backend, server, ssl)
	return s.Execute(cmd)
}

// GetServersState returns servers runtime state
func (s *SingleRuntime) GetServersState(backend string) (models.RuntimeServers, error) {
	cmd := "show servers state " + backend
	result, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}
	return parseRuntimeServers(result)
}

// GetServersState returns server runtime state
func (s *SingleRuntime) GetServerState(backend, server string) (*models.RuntimeServer, error) {
	cmd := "show servers state " + backend
	result, err := s.ExecuteWithResponse(cmd)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(result, "\n")
	if strings.TrimSpace(lines[0]) != "1" {
		return nil, stderrors.New("unsupported output format version, supporting format version 1")
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
	return nil, fmt.Errorf("server state for %s/%s: %w", backend, server, errors.ErrNotFound)
}

func parseRuntimeServers(output string) (models.RuntimeServers, error) {
	lines := strings.Split(output, "\n")
	result := models.RuntimeServers{}

	if strings.TrimSpace(lines[0]) != "1" {
		return nil, stderrors.New("unsupported output format version, supporting format version 1")
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

	if len(fields) < 25 {
		return nil
	}

	bID, err := strconv.ParseInt(fields[0], 10, 64)
	var backendID *int64
	if err == nil {
		backendID = &bID
	}

	var opState string
	switch fields[5] {
	case "0":
		opState = "down"
	case "3":
		opState = "stopping"
	case "1", "2":
		opState = "up"
	}

	admState, _ := misc.GetServerAdminState(fields[6])

	uW, err := strconv.ParseInt(fields[7], 10, 64)
	var uWeight *int64
	if err == nil {
		uWeight = &uW
	}

	iW, err := strconv.ParseInt(fields[8], 10, 64)
	var iWeight *int64
	if err == nil {
		iWeight = &iW
	}

	lTC, err := strconv.ParseInt(fields[9], 10, 64)
	var lastTimeChange *int64
	if err == nil {
		lastTimeChange = &lTC
	}

	cStatus, err := strconv.ParseInt(fields[10], 10, 64)
	var checkStatus *int64
	if err == nil {
		checkStatus = &cStatus
	}

	cResult, err := strconv.ParseInt(fields[11], 10, 64)
	var checkResult *int64
	if err == nil {
		checkResult = &cResult
	}

	cHealth, err := strconv.ParseInt(fields[12], 10, 64)
	var checkHealth *int64
	if err == nil {
		checkHealth = &cHealth
	}

	cState, err := strconv.ParseInt(fields[13], 10, 64)
	var checkState *int64
	if err == nil {
		checkState = &cState
	}

	aState, err := strconv.ParseInt(fields[14], 10, 64)
	var agentState *int64
	if err == nil {
		agentState = &aState
	}

	bFID, err := strconv.ParseInt(fields[15], 10, 64)
	var backendForcedID *int64
	if err == nil {
		backendForcedID = &bFID
	}

	fID, err := strconv.ParseInt(fields[16], 10, 64)
	var forcedID *int64
	if err == nil {
		forcedID = &fID
	}

	p, err := strconv.ParseInt(fields[18], 10, 64)
	var port *int64
	if err == nil {
		port = &p
	}

	uSSL, err := strconv.ParseBool(fields[20])
	var useSSL *bool
	if err == nil {
		useSSL = &uSSL
	}

	cPort, err := strconv.ParseInt(fields[21], 10, 64)
	var checkPort *int64
	if err == nil {
		checkPort = &cPort
	}

	aPort, err := strconv.ParseInt(fields[24], 10, 64)
	var agentPort *int64
	if err == nil {
		agentPort = &aPort
	}

	return &models.RuntimeServer{
		BackendID:        backendID,
		BackendName:      fields[1],
		ID:               fields[2],
		Name:             fields[3],
		Address:          fields[4],
		AdminState:       admState,
		OperationalState: opState,
		Uweight:          uWeight,
		Iweight:          iWeight,
		LastTimeChange:   lastTimeChange,
		CheckStatus:      checkStatus,
		CheckResult:      checkResult,
		CheckHealth:      checkHealth,
		CheckState:       checkState,
		AgentState:       agentState,
		BackendForcedID:  backendForcedID,
		ForecedID:        forcedID,
		Fqdn:             fields[17],
		Port:             port,
		Srvrecord:        fields[19],
		UseSsl:           useSSL,
		CheckPort:        checkPort,
		CheckAddr:        fields[22],
		AgentAddr:        fields[23],
		AgentPort:        agentPort,
	}
}
