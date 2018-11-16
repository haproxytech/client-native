package runtime

import (
	"fmt"
	"net"
	"strings"
	"time"
)

//TaskResponse ...
type TaskResponse struct {
	result string
	err    error
}

//Task has command to execute on runtime api, and response channel for result
type Task struct {
	command  string
	response chan TaskResponse
}

//Client handles multiple HAProxy clients
type Client struct {
	SocketsPath []string
}

//SingleRuntime handles one runtime API
type SingleRuntime struct {
	socketOpen       bool
	jobs             chan Task
	socketPath       string
	autoReconnect    bool
	runtimeAPIsocket net.Conn
}

//Init must be given path to runtime socket
func (s *SingleRuntime) Init(socketPath string, autoReconnect bool) error {
	s.socketPath = socketPath
	s.autoReconnect = autoReconnect
	s.jobs = make(chan Task)
	s.socketConnect()
	go s.handleIncommingJobs()
	return nil
}

func (s *SingleRuntime) socketConnect() error {
	var err error
	s.runtimeAPIsocket, err = net.Dial("unix", s.socketPath)
	if err != nil {
		if s.autoReconnect {
			go func() {
				time.Sleep(time.Second * 1)
				s.socketConnect()
			}()
		}
		return err
	}
	s.socketOpen = true
	_, err = s.runtimeAPIsocket.Write([]byte(fmt.Sprintf("prompt\n")))
	if err != nil {
		return err
	}
	_, err = s.runtimeAPIsocket.Write([]byte(fmt.Sprintf("set severity-output number\n")))
	if err != nil {
		return err
	}
	return nil
}

func (s *SingleRuntime) handleIncommingJobs() {
	for {
		select {
		case job := <-s.jobs:
			result, err := s.readFromSocket(s.runtimeAPIsocket, job.command)
			if err != nil {
				job.response <- TaskResponse{err: err}
			} else {
				job.response <- TaskResponse{result: result}
			}
		case <-time.After(time.Duration(60) * time.Second):
		}
	}
}

func (s *SingleRuntime) readFromSocket(c net.Conn, command string) (string, error) {
	if !s.socketOpen {
		return "", fmt.Errorf("no connection")
	}
	_, err := c.Write([]byte(fmt.Sprintf("%s\n", command)))
	if err != nil {
		s.socketOpen = false
		c.Close()
		return "", err
	}
	time.Sleep(1e9)
	bufferSize := 1024
	buf := make([]byte, bufferSize)
	var data strings.Builder
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			break
		}
		data.Write(buf[0:n])
		if n < bufferSize {
			break
		}
	}
	result := strings.TrimSuffix(data.String(), "\n> ")
	result = strings.TrimSuffix(result, "\n")
	return result, nil
}

func (s *SingleRuntime) readFromSocketClean(command string) (string, error) {
	c, err := net.Dial("unix", s.socketPath)
	if err != nil {
		return "", err
	}
	defer c.Close()

	_, err = c.Write([]byte(fmt.Sprintf("%s\n", command)))
	if err != nil {
		return "", nil
	}
	time.Sleep(1e9)
	buf := make([]byte, 1024)
	var data strings.Builder
	for {
		n, err := c.Read(buf[:])
		if err != nil {
			break
		}
		data.Write(buf[0:n])
	}
	return data.String(), nil
}

//ExecuteRaw executes command on runtime API and returns raw result
func (s *SingleRuntime) ExecuteRaw(command string) (string, error) {
	//allow one retry if connection breaks temporarily
	return s.executeRaw(command, 1)
}

func (s *SingleRuntime) executeRaw(command string, retry int) (string, error) {
	response := make(chan TaskResponse)
	Task := Task{
		command:  command,
		response: response,
	}
	s.jobs <- Task
	select {
	case rsp := <-response:
		if rsp.err != nil && retry > 0 {
			if !s.socketOpen || s.runtimeAPIsocket == nil {
				s.socketConnect()
			}
			retry--
			return s.executeRaw(command, retry)
		}
		return rsp.result, rsp.err
	case <-time.After(time.Duration(30) * time.Second):
		return "", fmt.Errorf("timeout reached")
	}
}
