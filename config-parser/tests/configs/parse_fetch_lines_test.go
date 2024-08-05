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
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

func TestParseFecthResultLines(t *testing.T) { //nolint:gocognit
	tests := []struct {
		Name, Config string
	}{
		{"configBasic1", configBasic1},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			p, err := parser.New(options.String(config.Config))
			if err != nil {
				t.Fatalf(err.Error())
			}
			lines, err := p.GetResult(parser.Frontends, "http", "bind")
			if err != nil {
				t.Fatalf(err.Error())
			}
			if lines[0].Data != "bind 0.0.0.0:80 name bind_1" {
				t.Fatalf("Unexpected line: %s", lines[0].Data)
			}
			if lines[1].Data != "bind :::80 v4v6 name bind_2" {
				t.Fatalf("Unexpected line: %s", lines[1].Data)
			}
		})
	}
}
