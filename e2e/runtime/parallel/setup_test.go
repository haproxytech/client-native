// Copyright 2021 HAProxy Technologies
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
//go:build integration
// +build integration

package parallel_test

import (
	"os"
	"os/exec"
	"testing"

	client_native "github.com/haproxytech/client-native/v6"

	"github.com/haproxytech/client-native/v6/e2e"
	"github.com/stretchr/testify/suite"
)

type ParallelRuntime struct {
	suite.Suite
	cmd            *exec.Cmd
	client         client_native.HAProxyClient
	tmpDir         string
	haproxyVersion string
	socketPath     string
}

func (s *ParallelRuntime) SetupTest() {
	result, err := e2e.GetClient(s.T())
	if err != nil {
		s.FailNow(err.Error())
	}
	s.haproxyVersion = result.HAProxyVersion
	s.cmd = result.Cmd
	s.client = result.Client
	s.tmpDir = result.TmpDir
	s.socketPath = result.SocketPath
}

func (s *ParallelRuntime) TearDownSuite() {
	if err := s.cmd.Process.Kill(); err != nil {
		s.FailNow(err.Error())
	}
	if s.tmpDir != "" {
		err := os.RemoveAll(s.tmpDir)
		if err != nil {
			s.FailNow(err.Error())
		}
	}
}

func TestParallelRuntimes(t *testing.T) {
	suite.Run(t, new(ParallelRuntime))
}
