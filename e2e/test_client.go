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

package e2e

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path"
	goruntime "runtime"
	"strings"
	"testing"
	"time"

	clientnative "github.com/haproxytech/client-native/v2"
	"github.com/haproxytech/client-native/v2/configuration"
	"github.com/haproxytech/client-native/v2/runtime"
)

type ClientResponse struct {
	Client         *clientnative.HAProxyClient
	Cmd            *exec.Cmd
	TmpDir         string
	HAProxyVersion string
}

func GetClient(t *testing.T) (*ClientResponse, error) {
	cmd := exec.Command("haproxy", "-v")

	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	output := strings.Split(out.String(), "/n")[0]
	parts := strings.Split(output, " ")
	if len(parts) < 3 {
		return nil, errors.New("incorrect haproxy -v output")
	}
	version := strings.Split(parts[2], "-")[0]
	parts = strings.Split(version, ".")
	if len(parts) < 2 {
		return nil, errors.New("incorrect haproxy -v output")
	}
	version = fmt.Sprintf("%s.%s", parts[0], parts[1])

	_, file, _, _ := goruntime.Caller(1)     //nolint:dogsled
	_, filename, _, _ := goruntime.Caller(0) //nolint:dogsled
	testName := strings.ReplaceAll(path.Dir(file), path.Dir(filename), "")

	tmpPath := path.Join(os.TempDir(), "client-native/", testName)
	socketPath := path.Join(tmpPath, "runtime.sock")
	err = os.MkdirAll(tmpPath, 0777)
	if err != nil {
		return nil, err
	}

	cmd = exec.Command("haproxy", "-f", "haproxy.cfg")
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, fmt.Sprintf("SOCK_PATH=%s", socketPath))

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	HAProxyCFG := "haproxy.cfg"
	confClient := configuration.Client{}
	err = confClient.Init(configuration.ClientParams{
		ConfigurationFile:      HAProxyCFG,
		PersistentTransactions: false,
		TransactionDir:         os.TempDir(),
		Haproxy:                "haproxy",
	})
	if err != nil {
		return nil, err
	}

	var errSock error
	// sometimes init is faster than haproxy/OS is able to create socket
	// we need to wait a bit to be sure socket is ready
	start := time.Now()
	t.Logf("%s checking if %s exists", start.Format("15:04:05.000"), socketPath)
	for start := time.Now(); time.Since(start) < 10*time.Second; {
		api, errST := net.Dial("unix", socketPath)
		if errST == nil {
			_ = api.Close()
			break
		}
		t.Logf("waiting for socket %s to be created", socketPath)
		time.Sleep(10 * time.Millisecond)
	}
	if _, errSock = os.Stat(socketPath); os.IsNotExist(err) {
		return nil, errSock
	}
	end := time.Now()
	t.Logf("%s done", end.Format("15:04:05.000"))

	runtimeClient := runtime.Client{}
	err = runtimeClient.InitWithSockets(map[int]string{
		0: socketPath,
	})
	if err != nil {
		return nil, err
	}
	nativeAPI := &clientnative.HAProxyClient{
		Configuration: &confClient,
		Runtime:       &runtimeClient,
	}
	return &ClientResponse{
		Client:         nativeAPI,
		Cmd:            cmd,
		TmpDir:         tmpPath,
		HAProxyVersion: version,
	}, nil
}
