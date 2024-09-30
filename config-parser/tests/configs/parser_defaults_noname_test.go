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

	parser "github.com/haproxytech/client-native/v5/config-parser"
	"github.com/haproxytech/client-native/v5/config-parser/options"
)

func TestDefaultsConfigsNoFromFlag(t *testing.T) {
	tests := []struct {
		Name, Config, Result string
	}{
		{"configDefaultsNoName1", configDefaultsNoFlag1, configDefaultsNoFlag1Result},
		{"configDefaultsNoName1", configDefaultsNoFlag2, configDefaultsNoFlag2Result},
	}
	for _, config := range tests {
		t.Run(config.Name, func(t *testing.T) {
			var buffer bytes.Buffer
			buffer.WriteString(config.Config)
			p, err := parser.New(options.Reader(&buffer), options.NoNamedDefaultsFrom)
			if err != nil {
				t.Fatalf(err.Error())
			}
			result := p.String()
			if result != config.Result {
				compare(t, config.Result, result)
				t.Fatalf("configurations does not match")
			}
		})
	}
}
