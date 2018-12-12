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

//SingleRuntime handles one runtime API
type SingleRuntime struct {
	jobs       chan Task
	socketPath string
}

//Init must be given path to runtime socket
func (s *SingleRuntime) Init(socketPath string, autoReconnect bool) error {
	s.socketPath = socketPath
	s.jobs = make(chan Task)
	go s.handleIncommingJobs()
	return nil
}

//Reload closes connection and open new one
func (s *SingleRuntime) Reload() error {
	_, err := s.readFromSocket("quit")
	return err
}

func (s *SingleRuntime) handleIncommingJobs() {
	for {
		select {
		case job := <-s.jobs:
			result, err := s.readFromSocket(job.command)
			if err != nil {
				job.response <- TaskResponse{err: err}
			} else {
				job.response <- TaskResponse{result: result}
			}
		}
	}
}

func (s *SingleRuntime) readFromSocket(command string) (string, error) {
	api, err := net.Dial("unix", s.socketPath)
	if err != nil {
		return "", err
	}
	_, err = api.Write([]byte(fmt.Sprintf("set severity-output number;%s\n", command)))
	if err != nil {
		return "", err
	}
	return "", nil

	if api == nil {
		return "", fmt.Errorf("no connection")
	}
	time.Sleep(1e9)
	bufferSize := 1024
	buf := make([]byte, bufferSize)
	var data strings.Builder
	for {
		n, err := api.Read(buf[:])
		if err != nil {
			break
		}
		data.Write(buf[0:n])
		if n < bufferSize {
			break
		}
	}
	api.Close()
	if err != nil {
		return "", err
	}
	result := strings.TrimSuffix(data.String(), "\n> ")
	result = strings.TrimSuffix(result, "\n")
	return result, nil
}

func (s *SingleRuntime) readFromSocketClean(command string) (string, error) {
	api, err := net.Dial("unix", s.socketPath)
	if err != nil {
		return "", err
	}
	defer api.Close()

	_, err = api.Write([]byte(fmt.Sprintf("%s\n", command)))
	if err != nil {
		return "", nil
	}
	time.Sleep(1e9)
	buf := make([]byte, 1024)
	var data strings.Builder
	for {
		n, err := api.Read(buf[:])
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

//Execute executes command on runtime API
func (s *SingleRuntime) Execute(command string) error {
	rawdata, err := s.ExecuteRaw(command)
	if err != nil {
		return fmt.Errorf("%s [%s]", err.Error(), command)
	}
	if len(rawdata) > 1 {
		switch rawdata[1] {
		case '3', '2', '1', '0':
			return fmt.Errorf("[%c] %s [%s]", rawdata[1], rawdata[3:], command)
		}
	}
	return nil
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
			retry--
			return s.executeRaw(command, retry)
		}
		return rsp.result, rsp.err
	case <-time.After(time.Duration(30) * time.Second):
		return "", fmt.Errorf("timeout reached")
	}
}
