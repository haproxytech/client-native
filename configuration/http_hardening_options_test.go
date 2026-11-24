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

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/params"
	config_options "github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func TestBindHTTPHardeningRoundTrip(t *testing.T) {
	config := `frontend f
  bind :80 http-hardening auto
  bind :81
  bind :82 http-hardening trusted
`
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatal(err)
	}
	binds, err := ParseBinds("frontend", "f", p)
	if err != nil {
		t.Fatal(err)
	}
	if len(binds) != 3 {
		t.Fatalf("got %d binds, want 3", len(binds))
	}
	if binds[0].HTTPHardening != "auto" {
		t.Fatalf("explicit auto lost: %#v", binds[0].BindParams)
	}
	if binds[1].HTTPHardening != "" {
		t.Fatalf("omitted reference became %q", binds[1].HTTPHardening)
	}
	if binds[2].HTTPHardening != "trusted" {
		t.Fatalf("named reference lost: %#v", binds[2].BindParams)
	}

	optionsString := params.BindOptionsString(serializeBindParams(binds[0].BindParams, binds[0].Name, ":80", &config_options.ConfigurationOptions{}))
	if !strings.Contains(optionsString, "http-hardening auto") {
		t.Fatalf("auto not serialized: %q", optionsString)
	}
	optionsString = params.BindOptionsString(serializeBindParams(binds[1].BindParams, binds[1].Name, ":81", &config_options.ConfigurationOptions{}))
	if strings.Contains(optionsString, "http-hardening") {
		t.Fatalf("omitted reference serialized: %q", optionsString)
	}
}

func TestServerHTTPHardeningRoundTrip(t *testing.T) {
	config := `backend b
  server s1 127.0.0.1:80 http-hardening auto
  server s2 127.0.0.1:81
  server s3 127.0.0.1:82 http-hardening trusted
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
	if servers[0].HTTPHardening != "auto" {
		t.Fatalf("explicit auto lost: %#v", servers[0].ServerParams)
	}
	if servers[1].HTTPHardening != "" {
		t.Fatalf("omitted reference became %q", servers[1].HTTPHardening)
	}
	if servers[2].HTTPHardening != "trusted" {
		t.Fatalf("named reference lost: %#v", servers[2].ServerParams)
	}

	optionsString := params.ServerOptionsString(SerializeServerParams(servers[0].ServerParams, &config_options.ConfigurationOptions{}))
	if !strings.Contains(optionsString, "http-hardening auto") {
		t.Fatalf("auto not serialized: %q", optionsString)
	}
	optionsString = params.ServerOptionsString(SerializeServerParams(servers[1].ServerParams, &config_options.ConfigurationOptions{}))
	if strings.Contains(optionsString, "http-hardening") {
		t.Fatalf("omitted reference serialized: %q", optionsString)
	}
}

func TestHTTPHardeningSharedConsumers(t *testing.T) {
	serverOptions, err := params.ParseServerOptions([]string{"http-hardening", "auto"})
	if err != nil {
		t.Fatal(err)
	}
	var serverParams models.ServerParams
	parseServerParams(serverOptions, &serverParams)
	if serverParams.HTTPHardening != "auto" {
		t.Fatalf("unexpected server params: %#v", serverParams)
	}

	serverOptions, err = params.ParseServerOptions([]string{"http-hardening", "hardening1"})
	if err != nil {
		t.Fatal(err)
	}
	serverParams = models.ServerParams{}
	parseServerParams(serverOptions, &serverParams)
	if serverParams.HTTPHardening != "hardening1" {
		t.Fatalf("unexpected server params: %#v", serverParams)
	}

	bindOptions, err := params.ParseBindOptions([]string{"http-hardening", "auto"})
	if err != nil {
		t.Fatal(err)
	}
	bindParams, _ := parseBindParams(bindOptions)
	if bindParams.HTTPHardening != "auto" {
		t.Fatalf("unexpected bind params: %#v", bindParams)
	}
}

func TestHTTPHardeningBeforeConsumers(t *testing.T) {
	config := `defaults d
  default-server http-hardening auto

peers p
  peer p1 127.0.0.1:1000

ring r
  server s 127.0.0.1:1001 http-hardening auto

frontend f
  bind :80 http-hardening auto

backend b
  server s 127.0.0.1:80 http-hardening hardening1

http-hardening hardening1
  all on
`
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatal(err)
	}
	out := p.String()
	hardening := strings.Index(out, "http-hardening hardening1")
	for _, consumer := range []string{"defaults d", "peers p", "ring r", "frontend f", "backend b"} {
		index := strings.Index(out, consumer)
		if index < 0 || hardening < 0 || hardening > index {
			t.Fatalf("%s does not follow the hardening section:\n%s", consumer, out)
		}
	}
}
