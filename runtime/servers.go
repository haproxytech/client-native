package runtime

import (
	"fmt"
	"strings"
)

//SetServerAddr set ip [port] for server
//possible responses:
// - HTTPLikeStatusOK
// - HTTPLikeStatusNotModified
// - HTTPLikeStatusNotFound
// - HTTPLikeStatusInternalServerError
// TODO see if we can have here permission denied !!
func (s *SingleRuntime) SetServerAddr(backend, server string, ip string, port int) (RuntimeResponse, error) {
	var cmd string
	if port > 0 {
		cmd = fmt.Sprintf("set server %s/%s addr %s port %d", backend, server, ip, port)
	} else {
		cmd = fmt.Sprintf("set server %s/%s addr %s", backend, server, ip)
	}
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusInternalServerError}, err
	}
	//No such server.
	//IP changed from '127.0.0.1' to '127.0.0.2' by 'stats socket command'
	//no need to change the addr by 'stats socket command'
	//IP changed from '127.0.0.2' to '127.0.0.3', port changed from '80' to '8080' by 'stats socket command'
	//IP changed from '127.0.0.3' to '127.0.0.4', no need to change the port by 'stats socket command'
	//no need to change the addr, port changed from '5851' to '80' by 'stats socket command'
	//no need to change the addr, no need to change the port by 'stats socket command'
	if strings.Contains(rawdata, "No such server") {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusNotFound}, fmt.Errorf("server not found")
	}
	status := HTTPLikeStatusOK
	if !strings.Contains(rawdata, "changed from") {
		status = HTTPLikeStatusNotModified
	}
	return RuntimeResponse{HTTPLikeStatus: status, Message: rawdata}, nil
}

//SetServerState set ip [port] for server
//possible responses:
// - HTTPLikeStatusOK
// - HTTPLikeStatusNotFound
// - HTTPLikeStatusInternalServerError
// TODO see if we can have here permission denied !!
func (s *SingleRuntime) SetServerState(backend, server string, state string) (RuntimeResponse, error) {
	if !ServerStateValid(state) {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusBadRequest}, fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s state %s", backend, server, state)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusInternalServerError}, err
	}
	//No such server.
	//empty response
	if strings.Contains(rawdata, "No such server") {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusNotFound}, fmt.Errorf("server not found")
	}
	return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusOK, Message: rawdata}, nil
}

//SetServerWeight set ip [port] for server
//possible responses:
// - HTTPLikeStatusOK
// - HTTPLikeStatusNotFound
// - HTTPLikeStatusInternalServerError
// TODO see if we can have here permission denied !!
func (s *SingleRuntime) SetServerWeight(backend, server string, weight string) (RuntimeResponse, error) {
	if !ServerWeightValid(weight) {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusBadRequest}, fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s weight %s", backend, server, weight)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusInternalServerError}, err
	}
	//No such server.
	//empty response
	if strings.Contains(rawdata, "No such server") {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusNotFound}, fmt.Errorf("server not found")
	}
	return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusOK, Message: rawdata}, nil
}

//SetServerHealth set health for server
//possible responses:
// - HTTPLikeStatusOK
// - HTTPLikeStatusNotFound
// - HTTPLikeStatusInternalServerError
// TODO see if we can have here permission denied !!
func (s *SingleRuntime) SetServerHealth(backend, server string, health string) (RuntimeResponse, error) {
	if !ServerHealthValid(health) {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusBadRequest}, fmt.Errorf("bad request")
	}
	cmd := fmt.Sprintf("set server %s/%s health %s", backend, server, health)
	rawdata, err := s.ExecuteRaw(cmd)
	if err != nil {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusInternalServerError}, err
	}
	//No such server.
	//empty response
	if strings.Contains(rawdata, "No such server") {
		return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusNotFound}, fmt.Errorf("server not found")
	}
	return RuntimeResponse{HTTPLikeStatus: HTTPLikeStatusOK, Message: rawdata}, nil
}
