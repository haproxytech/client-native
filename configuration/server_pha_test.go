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
	"strings"
	"testing"

	strfmt "github.com/go-openapi/strfmt"
	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func TestServerPhaRoundTrip(t *testing.T) {
	config := `backend b
  server s1 127.0.0.1:80 pha
  server s2 127.0.0.1:81 no-pha
  server s3 127.0.0.1:82
`
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatal(err)
	}
	servers, err := ParseServers("backend", "b", p)
	if err != nil {
		t.Fatal(err)
	}
	if len(servers) != 3 {
		t.Fatalf("got %d servers, want 3", len(servers))
	}
	if servers[0].Pha != "enabled" {
		t.Fatalf("pha not parsed: %#v", servers[0].ServerParams)
	}
	if servers[1].Pha != "disabled" {
		t.Fatalf("no-pha not parsed: %#v", servers[1].ServerParams)
	}
	if servers[2].Pha != "" {
		t.Fatalf("omitted pha became %q", servers[2].Pha)
	}

	optionsString := params.ServerOptionsString(SerializeServerParams(servers[0].ServerParams, &config_options.ConfigurationOptions{}))
	if !strings.Contains(optionsString, "pha") {
		t.Fatalf("pha not serialized: %q", optionsString)
	}
	optionsString = params.ServerOptionsString(SerializeServerParams(servers[1].ServerParams, &config_options.ConfigurationOptions{}))
	if !strings.Contains(optionsString, "no-pha") {
		t.Fatalf("no-pha not serialized: %q", optionsString)
	}
	optionsString = params.ServerOptionsString(SerializeServerParams(servers[2].ServerParams, &config_options.ConfigurationOptions{}))
	if strings.Contains(optionsString, "pha") {
		t.Fatalf("omitted pha serialized: %q", optionsString)
	}
}

func TestPhaRuntimePreparation(t *testing.T) {
	enabled := models.ServerParams{Pha: "enabled"}.PrepareFieldsForRuntimeAddServer()
	if enabled.Pha != "enabled" {
		t.Fatalf("enabled PHA was dropped: %#v", enabled)
	}
	disabled := models.ServerParams{Pha: "disabled"}.PrepareFieldsForRuntimeAddServer()
	if disabled.Pha != "" {
		t.Fatalf("disabled PHA was sent to runtime: %#v", disabled)
	}
}

func TestPhaModelValidation(t *testing.T) {
	invalid := models.ServerParams{Pha: "maybe"}
	if err := invalid.Validate(strfmt.Default); err == nil {
		t.Fatal("invalid pha value passed validation")
	}
	valid := models.ServerParams{Pha: "enabled"}
	if err := valid.Validate(strfmt.Default); err != nil {
		t.Fatal(err)
	}
}
