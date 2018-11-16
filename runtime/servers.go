package runtime

import (
	"fmt"
)

//SetServerAddr set ip [port] for server
func (s *SingleRuntime) SetServerAddr(backend, server string, ip string, port int) error {
	var cmd string
	if port > 0 {
		cmd = fmt.Sprintf("set server %s/%s addr %s port %d", backend, server, ip, port)
	} else {
		cmd = fmt.Sprintf("set server %s/%s addr %s", backend, server, ip)
	}
	return s.Execute(cmd)
}

//SetServerState set state for server
func (s *SingleRuntime) SetServerState(backend, server string, state string) error {
	if !ServerStateValid(state) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	return s.Execute(cmd)
}

//SetServerWeight set weight for server
func (s *SingleRuntime) SetServerWeight(backend, server string, weight string) error {
	if !ServerWeightValid(weight) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s weight %s", backend, server, weight)
	return s.Execute(cmd)
}

//SetServerHealth set health for server
func (s *SingleRuntime) SetServerHealth(backend, server string, health string) error {
	if !ServerHealthValid(health) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s health %s", backend, server, health)
	return s.Execute(cmd)
}
