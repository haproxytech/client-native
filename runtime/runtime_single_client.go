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
	"net"
	"strings"
	"time"
)

const (
	taskTimeout = 30 * time.Second
)

// TaskResponse ...
type TaskResponse struct {
	result string
	err    error
}

// Task has command to execute on runtime api, and response channel for result
type Task struct {
	command  string
	response chan TaskResponse
}

// SingleRuntime handles one runtime API
type SingleRuntime struct {
	jobs       chan Task
	socketPath string
	worker     int
	process    int
}

// Init must be given path to runtime socket and worker number. If in master-worker mode,
// give the path to the master socket path, and non 0 number for workers. Process is for
// nbproc > 1. In master-worker mode it's the same as the worker number, but when having
// multiple stats socket lines bound to processes then use the correct process number
func (s *SingleRuntime) Init(socketPath string, worker int, process int) error {
	s.socketPath = socketPath
	s.jobs = make(chan Task)
	s.worker = worker
	s.process = process
	go s.handleIncommingJobs()
	return nil
}

func (s *SingleRuntime) handleIncommingJobs() {
	for job := range s.jobs {
		result, err := s.readFromSocket(job.command)
		if err != nil {
			job.response <- TaskResponse{err: err}
		} else {
			job.response <- TaskResponse{result: result}
		}
	}
}

func (s *SingleRuntime) readFromSocket(command string) (string, error) {
	var api net.Conn
	var err error

	if api, err = net.DialTimeout("unix", s.socketPath, taskTimeout); err != nil {
		return "", err
	}
	defer func() {
		_ = api.Close()
	}()
	if err = api.SetDeadline(time.Now().Add(taskTimeout)); err != nil {
		return "", err
	}

	fullCommand := fmt.Sprintf("set severity-output number;%s\n", command)
	if s.worker > 0 {
		fullCommand = fmt.Sprintf("@%v set severity-output number;@%v %s;quit\n", s.worker, s.worker, command)
	}
	_, err = api.Write([]byte(fullCommand))
	if err != nil {
		return "", err
	}
	// return "", nil

	if api == nil {
		return "", fmt.Errorf("no connection")
	}
	bufferSize := 1024
	buf := make([]byte, bufferSize)
	var data strings.Builder
	for {
		if n, readErr := api.Read(buf); readErr != nil {
			break
		} else {
			data.Write(buf[0:n])
		}
	}

	result := strings.TrimSuffix(data.String(), "\n> ")
	result = strings.TrimSuffix(result, "\n")
	return result, nil
}

// ExecuteRaw executes command on runtime API and returns raw result
func (s *SingleRuntime) ExecuteRaw(command string) (string, error) {
	// allow one retry if connection breaks temporarily
	return s.executeRaw(command, 1)
}

// Execute executes command on runtime API
func (s *SingleRuntime) Execute(command string) error {
	rawdata, err := s.ExecuteRaw(command)
	if err != nil {
		return fmt.Errorf("%w [%s]", err, command)
	}
	if len(rawdata) > 5 {
		switch rawdata[1:5] {
		case "[3]:", "[2]:", "[1]:", "[0]:":
			return fmt.Errorf("[%c] %s [%s]", rawdata[2], rawdata[5:], command)
		}
	}
	return nil
}

func (s *SingleRuntime) ExecuteWithResponse(command string) (string, error) {
	rawdata, err := s.ExecuteRaw(command)
	if err != nil {
		return "", fmt.Errorf("%w [%s]", err, command)
	}
	if len(rawdata) > 5 {
		switch rawdata[1:5] {
		case "[3]:", "[2]:", "[1]:", "[0]:":
			return "", fmt.Errorf("[%c] %s [%s]", rawdata[2], rawdata[5:], command)
		}
	}
	if len(rawdata) > 1 {
		return rawdata[1:], nil
	}
	return "", nil
}

func (s *SingleRuntime) executeRaw(command string, retry int) (string, error) {
	response := make(chan TaskResponse)
	task := Task{
		command:  command,
		response: response,
	}
	s.jobs <- task
	select {
	case rsp := <-response:
		if rsp.err != nil && retry > 0 {
			retry--
			return s.executeRaw(command, retry)
		}
		return rsp.result, rsp.err
	case <-time.After(taskTimeout):
		return "", fmt.Errorf("timeout reached")
	}
}
