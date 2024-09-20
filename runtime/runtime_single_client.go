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
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/haproxytech/client-native/v6/runtime/options"
)

const (
	taskTimeout = 30 * time.Second

	statsSocket  socketType = "stats"
	masterSocket socketType = "master"
)

type socketType string

// TaskResponse ...
type TaskResponse struct {
	err    error
	result string
}

// Task has command to execute on runtime api, and response channel for result
type Task struct {
	command  string
	response chan TaskResponse
	socket   socketType
}

// SingleRuntime handles one runtime API
type SingleRuntime struct {
	jobs             chan Task
	socketPath       string
	masterWorkerMode bool
}

func (s SingleRuntime) IsValid() bool {
	return s.socketPath != ""
}

// Init must be given path to runtime socket and a flag to indicate if it's in master-worker mode.
func (s *SingleRuntime) Init(ctx context.Context, socketPath string, masterWorkerMode bool, opt ...options.RuntimeOptions) error {
	var runtimeOptions options.RuntimeOptions
	if len(opt) > 0 {
		runtimeOptions = opt[0]
	}

	s.socketPath = socketPath
	s.jobs = make(chan Task)
	s.masterWorkerMode = masterWorkerMode
	go s.handleIncomingJobs(ctx)
	if !runtimeOptions.DoNotCheckRuntimeOnInit {
		if runtimeOptions.AllowDelayedStartMax != nil {
			now := time.Now()
			var err error
			for {
				if _, err = s.ExecuteRaw("help"); err == nil {
					break
				}
				time.Sleep(*runtimeOptions.AllowDelayedStartTick)
				if time.Since(now) > *runtimeOptions.AllowDelayedStartMax {
					return fmt.Errorf("cannot connect to runtime API %s within %s: %w", socketPath, *runtimeOptions.AllowDelayedStartMax, err)
				}
			}
		} else {
			// check if we have a valid socket
			if _, err := s.ExecuteRaw("help"); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *SingleRuntime) handleIncomingJobs(ctx context.Context) {
	for {
		select {
		case job, ok := <-s.jobs:
			if !ok {
				return
			}
			result, err := s.readFromSocket(job.command, job.socket)
			if err != nil {
				job.response <- TaskResponse{err: err}
			} else {
				job.response <- TaskResponse{result: result}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *SingleRuntime) readFromSocket(command string, socket socketType) (string, error) {
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

	var fullCommand string

	switch socket {
	case statsSocket:
		fullCommand = fmt.Sprintf("set severity-output number;%s\n", command)
		if s.masterWorkerMode {
			fullCommand = fmt.Sprintf("set severity-output number;@%v %s;quit\n", 1, command)
		}
	case masterSocket:
		fullCommand = command + ";quit"
	}

	_, err = api.Write([]byte(fullCommand))
	if err != nil {
		return "", err
	}
	// return "", nil

	if api == nil {
		return "", errors.New("no connection")
	}
	bufferSize := 1024
	buf := make([]byte, bufferSize)
	var data strings.Builder
	for {
		n, readErr := api.Read(buf)
		if readErr != nil {
			break
		}
		data.Write(buf[0:n])
	}

	result := strings.TrimSuffix(data.String(), "\n> ")
	result = strings.TrimSuffix(result, "\n")
	result = strings.TrimSpace(result)
	return result, nil
}

// ExecuteRaw executes command on runtime API and returns raw result
func (s *SingleRuntime) ExecuteRaw(command string) (string, error) {
	// allow one retry if connection breaks temporarily
	return s.executeRaw(command, 1, statsSocket)
}

// Execute executes command on runtime API
func (s *SingleRuntime) Execute(command string) error {
	rawdata, err := s.ExecuteRaw(command)
	if err != nil {
		return fmt.Errorf("%w [%s]", err, command)
	}
	if len(rawdata) > 4 {
		switch rawdata[0:4] {
		case "[3]:", "[2]:", "[1]:", "[0]:":
			return fmt.Errorf("[%c] %s [%s]", rawdata[1], rawdata[4:], command)
		}
	}
	return nil
}

func (s *SingleRuntime) ExecuteWithResponse(command string) (string, error) {
	rawdata, err := s.ExecuteRaw(command)
	if err != nil {
		return "", fmt.Errorf("%w [%s]", err, command)
	}
	if len(rawdata) > 4 {
		switch rawdata[0:4] {
		case "[3]:", "[2]:", "[1]:", "[0]:":
			return "", fmt.Errorf("[%c] %s [%s]", rawdata[1], rawdata[4:], command)
		}
	}
	return rawdata, nil
}

func (s *SingleRuntime) ExecuteMaster(command string) (string, error) {
	// allow one retry if connection breaks temporarily
	return s.executeRaw(command, 1, masterSocket)
}

func (s *SingleRuntime) executeRaw(command string, retry int, socket socketType) (string, error) {
	response := make(chan TaskResponse)
	task := Task{
		command:  command,
		response: response,
		socket:   socket,
	}
	s.jobs <- task
	rsp := <-response
	if rsp.err != nil && retry > 0 {
		retry--
		return s.executeRaw(command, retry, socket)
	}
	return rsp.result, rsp.err
}
