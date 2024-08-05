/*
Copyright 2019 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package configs //nolint:testpackage

import (
	"bytes"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

const configSetEnv = `# _version=1
# HAProxy Technologies
# https://www.haproxy.com/

global
  presetenv processor1 "AMD Ryzen 7 1700 Eight-Core Processor"
  presetenv processor2 "AMD Ryzen 7 1700 Eight-Core Processor"
  presetenv custom "something 1"
  presetenv custom "something 2" # with comment
  setenv processor1 "AMD Ryzen 7 1700 Eight-Core Processor"
  setenv processor2 "AMD Ryzen 7 1700 Eight-Core Processor"
  setenv custom "something 1"
  setenv custom "something 2" # with comment
  master-worker

defaults A
  log global

frontend http from A
  mode http
  bind 0.0.0.0:80 name bind_1
  bind :::80 v4v6 name bind_2
  default_backend default_backend

backend default_backend from A
  mode http
  http-request deny deny_status 400 # deny
`

func TestSetEnv(t *testing.T) {
	tests := []struct {
		Name, Config string
	}{
		{"configBasic1", configSetEnv},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer))
			if err != nil {
				t.Fatalf(err.Error())
			}
			result := p.String()
			// fmt.Println(result)
			if result != config.Config {
				compare(t, config.Config, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}
