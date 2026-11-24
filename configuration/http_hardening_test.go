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
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

func TestParseHTTPHardening(t *testing.T) {
	config := `http-hardening base
  all on

http-hardening child from base
  reject-lowercase-method off # note
  all-delete on
`
	p, err := parser.New(options.String(config))
	if err != nil {
		t.Fatal(err)
	}

	profile, err := ParseHTTPHardening(p, "child")
	if err != nil {
		t.Fatal(err)
	}
	if profile.From != "base" {
		t.Fatalf("from = %q, want base", profile.From)
	}
	if profile.RejectLowercaseMethod != "off" || profile.AllDelete != "on" {
		t.Fatalf("unexpected checks: %#v", profile)
	}
	if profile.All != "" || profile.RejectDuplicateHostHeader != "" {
		t.Fatalf("unset checks were populated: %#v", profile)
	}
}

func TestSerializeHTTPHardeningRoundTrip(t *testing.T) {
	p, err := parser.New(options.String("http-hardening base\n\nhttp-hardening child from base\n  all on\n"))
	if err != nil {
		t.Fatal(err)
	}
	profile, err := ParseHTTPHardening(p, "child")
	if err != nil {
		t.Fatal(err)
	}

	profile.Metadata = map[string]any{"comment": "managed"}
	profile.All = "on"
	profile.RejectLowercaseMethod = "off"
	profile.RejectDuplicateHostHeader = "on"
	if err := SerializeHTTPHardening(p, profile); err != nil {
		t.Fatal(err)
	}

	reparsed, err := ParseHTTPHardening(p, "child")
	if err != nil {
		t.Fatal(err)
	}
	if reparsed.From != "base" {
		t.Fatalf("from = %q, want base", reparsed.From)
	}
	if reparsed.Metadata["comment"] != "managed" {
		t.Fatalf("metadata lost: %#v", reparsed.Metadata)
	}
	if reparsed.All != "on" || reparsed.RejectLowercaseMethod != "off" || reparsed.RejectDuplicateHostHeader != "on" {
		t.Fatalf("checks lost: %#v", reparsed)
	}
}

func TestSerializeHTTPHardeningClearsOmittedFields(t *testing.T) {
	config := `http-hardening base
  all on

http-hardening child from base # {"comment":"managed"}
  all off
`
	p, err := parser.New(options.String(config))
	if err != nil {
		t.Fatal(err)
	}
	profile, err := ParseHTTPHardening(p, "child")
	if err != nil {
		t.Fatal(err)
	}
	if profile.From != "base" || profile.All != "off" || profile.Metadata == nil {
		t.Fatalf("fixture did not parse as expected: %#v", profile)
	}

	profile.From = ""
	profile.All = ""
	profile.Metadata = nil
	if err := SerializeHTTPHardening(p, profile); err != nil {
		t.Fatal(err)
	}

	reparsed, err := ParseHTTPHardening(p, "child")
	if err != nil {
		t.Fatal(err)
	}
	if reparsed.From != "" || reparsed.All != "" || reparsed.Metadata != nil {
		t.Fatalf("omitted fields were not cleared: %#v", reparsed)
	}
}

func TestValidateHTTPHardeningParent(t *testing.T) {
	p, err := parser.New(options.String("http-hardening base\n\nhttp-hardening child\n"))
	if err != nil {
		t.Fatal(err)
	}

	if err := validateHTTPHardeningParent(p, "child", "base"); err != nil {
		t.Fatal(err)
	}
	if err := validateHTTPHardeningParent(p, "child", "trusted"); err != nil {
		t.Fatal(err)
	}
	if err := validateHTTPHardeningParent(p, "child", "auto"); err == nil {
		t.Fatal("auto parent was accepted")
	}
	if err := validateHTTPHardeningParent(p, "child", "child"); err == nil {
		t.Fatal("self parent was accepted")
	}
	if err := validateHTTPHardeningParent(p, "child", "missing"); err == nil {
		t.Fatal("missing parent was accepted")
	}
}

func TestSerializeHTTPHardeningRejectsNil(t *testing.T) {
	p, err := parser.New()
	if err != nil {
		t.Fatal(err)
	}
	if err := SerializeHTTPHardening(p, nil); err == nil {
		t.Fatal("nil profile was accepted")
	}
}
