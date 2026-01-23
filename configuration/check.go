// Copyright 2025 HAProxy Technologies
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

package configuration

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/haproxytech/client-native/v6/configuration/options"
	shellquote "github.com/kballard/go-shellquote"
)

//nolint:noctx
func checkHaproxyConfiguration(opt options.ConfigurationOptions, path string, transactionID ...string) error {
	var name string
	var args []string

	// Inherit the environment but filter out a few unwanted variables.
	envs := removeFromEnv(os.Environ(), "HAPROXY_STARTUPLOGS_FD",
		"HAPROXY_MWORKER_WAIT_ONLY", "HAPROXY_PROCESSES")

	switch {
	case len(transactionID) > 0 && len(opt.ValidateCmd) > 0:
		w, _ := shellquote.Split(opt.ValidateCmd)
		name = w[0]
		args = w[1:]
		envs = append(envs, "DATAPLANEAPI_TRANSACTION_FILE="+path)
	case opt.MasterWorker:
		name = opt.Haproxy
		args = []string{"-W", "-f", path, "-c"}
		args = addConfigFilesToArgs(args, opt)
	default:
		name = opt.Haproxy
		args = []string{"-f", path, "-c"}
		args = addConfigFilesToArgs(args, opt)
	}

	// #nosec G204
	cmd := exec.Command(name, args...)
	cmd.Env = envs
	var stderr strings.Builder
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		errStr := fmt.Sprintf("%s: %s", err.Error(), parseHAProxyCheckError(stderr.String(), transactionID...))
		return NewConfError(ErrValidationError, errStr)
	}
	return nil
}

func parseHAProxyCheckError(output string, transactionID ...string) string { //nolint:gocognit
	var b strings.Builder

	if len(transactionID) > 0 && transactionID[0] != "" {
		b.WriteString(fmt.Sprintf("err transactionId=%s \n", transactionID[0]))
	}

	for lineWhole := range strings.SplitSeq(output, "\n") {
		line := strings.TrimSpace(lineWhole)
		if strings.HasPrefix(line, "[ALERT]") {
			if strings.HasSuffix(line, "fatal errors found in configuration.") {
				continue
			}
			if strings.Contains(line, "error(s) found in configuration file : ") {
				continue
			}

			parts := strings.Split(line, " : ")
			if len(parts) > 2 && strings.HasPrefix(strings.TrimSpace(parts[1]), "parsing [") {
				fParts := strings.Split(strings.TrimSpace(parts[1]), ":")
				var msgB strings.Builder
				for i := 2; i < len(parts); i++ {
					msgB.WriteString(parts[i])
					msgB.WriteString(" ")
				}
				if len(fParts) > 1 {
					lNo, err := strconv.ParseInt(strings.TrimSuffix(fParts[1], "]"), 10, 64)
					if err == nil {
						b.WriteString(fmt.Sprintf("line=%d msg=\"%s\"\n", lNo, strings.TrimSpace(msgB.String())))
					} else {
						b.WriteString(fmt.Sprintf("msg=\"%s\"\n", strings.TrimSpace(msgB.String())))
					}
				}
			} else if len(parts) > 1 {
				var msgB strings.Builder
				for i := 1; i < len(parts); i++ {
					msgB.WriteString(parts[i])
					msgB.WriteString(" ")
				}
				b.WriteString(fmt.Sprintf("msg=\"%s\"\n", strings.TrimSpace(msgB.String())))
			}
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func addConfigFilesToArgs(args []string, clientParams options.ConfigurationOptions) []string {
	result := make([]string, 0) //nolint: prealloc
	for _, file := range clientParams.ValidateConfigFilesBefore {
		result = append(result, "-f", file)
	}
	result = append(result, args...)

	for _, file := range clientParams.ValidateConfigFilesAfter {
		result = append(result, "-f", file)
	}
	return result
}

// Returns a copy of envs without the unwanted environment variables.
func removeFromEnv(envs []string, unwanted ...string) []string {
	newEnv := make([]string, 0, len(envs))

	for _, v := range envs {
		skip := false
		for _, bad := range unwanted {
			if strings.HasPrefix(v, bad+"=") {
				skip = true
				break
			}
		}
		if !skip {
			newEnv = append(newEnv, v)
		}
	}

	return newEnv
}
