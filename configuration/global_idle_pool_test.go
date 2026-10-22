// Copyright 2026 HAProxy Technologies
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

package configuration

import (
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
)

func TestGlobalIdlePoolSharedRoundTrip(t *testing.T) {
	tests := map[string]string{
		"enabled":  "on",
		"disabled": "off",
		"full":     "full",
	}

	for api, haproxy := range tests {
		t.Run(api, func(t *testing.T) {
			p, err := parser.New(parser_options.String("global\n  tune.idle-pool.shared " + haproxy + "\n"))
			if err != nil {
				t.Fatalf("parser.New: %v", err)
			}
			global, err := ParseGlobalSection(p)
			if err != nil {
				t.Fatalf("ParseGlobalSection: %v", err)
			}
			if global.TuneOptions == nil || global.TuneOptions.IdlePoolShared != api {
				t.Fatalf("parsed IdlePoolShared = %#v, want %q", global.TuneOptions, api)
			}

			out, err := parser.New(parser_options.String("global\n"))
			if err != nil {
				t.Fatalf("parser.New: %v", err)
			}
			if err := serializeTuneOptions(out, &models.TuneOptions{IdlePoolShared: api}, nil); err != nil {
				t.Fatalf("serializeTuneOptions: %v", err)
			}
			line := "tune.idle-pool.shared " + haproxy
			if !strings.Contains(out.String(), line) {
				t.Fatalf("serialized config missing %q:\n%s", line, out.String())
			}

			fresh, err := parser.New(parser_options.String(out.String()))
			if err != nil {
				t.Fatalf("parser.New: %v", err)
			}
			reparsed, err := ParseGlobalSection(fresh)
			if err != nil {
				t.Fatalf("reparse: %v", err)
			}
			if reparsed.TuneOptions == nil || reparsed.TuneOptions.IdlePoolShared != api {
				t.Fatalf("reparsed IdlePoolShared = %#v, want %q", reparsed.TuneOptions, api)
			}
		})
	}
}

func TestOtherOnOffSerializationRejectsFull(t *testing.T) {
	p, err := parser.New(parser_options.String("global\n"))
	if err != nil {
		t.Fatal(err)
	}
	if err := serializeOnOffOption(p, "tune.listener.multi-queue", "full"); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(p.String(), "tune.listener.multi-queue") {
		t.Fatalf("serialized unsupported full value:\n%s", p.String())
	}
}

func TestTuneOptionsIdlePoolSharedEnum(t *testing.T) {
	for _, valid := range []string{"enabled", "disabled", "full"} {
		if err := (&models.TuneOptions{IdlePoolShared: valid}).Validate(strfmt.Default); err != nil {
			t.Errorf("Validate(%q) = %v, want nil", valid, err)
		}
	}
	if err := (&models.TuneOptions{IdlePoolShared: "invalid"}).Validate(strfmt.Default); err == nil {
		t.Error("Validate(invalid) = nil, want error")
	}
}
