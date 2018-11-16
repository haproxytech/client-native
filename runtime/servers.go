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
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return err
	}
	if len(rawdata) > 1 {
		switch rawdata[1] {
		case '3', '2', '1', '0':
			return fmt.Errorf(rawdata[3:])
		}
	}
	return nil
}

//SetServerState set state for server
func (s *SingleRuntime) SetServerState(backend, server string, state string) error {
	if !ServerStateValid(state) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return err
	}
	if len(rawdata) > 1 {
		switch rawdata[1] {
		case '3', '2', '1', '0':
			return fmt.Errorf(rawdata[3:])
		}
	}
	return nil
}

//SetServerWeight set weight for server
func (s *SingleRuntime) SetServerWeight(backend, server string, weight string) error {
	if !ServerWeightValid(weight) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s weight %s", backend, server, weight)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return err
	}
	if len(rawdata) > 1 {
		switch rawdata[1] {
		case '3', '2', '1', '0':
			return fmt.Errorf(rawdata[3:])
		}
	}
	return nil
}

//SetServerHealth set health for server
func (s *SingleRuntime) SetServerHealth(backend, server string, health string) error {
	if !ServerHealthValid(health) {
		return fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s health %s", backend, server, health)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return err
	}
	if len(rawdata) > 1 {
		switch rawdata[1] {
		case '3', '2', '1', '0':
			return fmt.Errorf(rawdata[3:])
		}
	}
	return nil
}
