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
//

package configuration

import (
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

// The "on" step sd/format values must keep the exact form they were written
// in: a re-serialized configuration is compared textually against the
// original by API consumers, so quotes must not be added or dropped.
func TestLogProfileStepQuoteRoundTrip(t *testing.T) {
	config := `log-profile myprof
  log-tag "custom-tag"
  on error format "%ci: error"
  on connect drop
  on any sd "custom-sd"
  on http-req sd test-me
  on http-res sd "test me"
`
	p, err := parser.New(options.String(config))
	if err != nil {
		t.Fatal(err)
	}

	lp, err := ParseLogProfile(p, "myprof")
	if err != nil {
		t.Fatal(err)
	}

	// Quoted values containing spaces are unquoted in the model, anything
	// else is kept verbatim.
	wantSteps := []struct{ format, sd string }{
		{format: "%ci: error"},
		{},
		{sd: `"custom-sd"`},
		{sd: "test-me"},
		{sd: "test me"},
	}
	if len(lp.Steps) != len(wantSteps) {
		t.Fatalf("got %d steps, want %d", len(lp.Steps), len(wantSteps))
	}
	for i, want := range wantSteps {
		if lp.Steps[i].Format != want.format {
			t.Errorf("step %d: format = %q, want %q", i, lp.Steps[i].Format, want.format)
		}
		if lp.Steps[i].Sd != want.sd {
			t.Errorf("step %d: sd = %q, want %q", i, lp.Steps[i].Sd, want.sd)
		}
	}

	if err := SerializeLogProfile(p, lp); err != nil {
		t.Fatal(err)
	}

	got := p.String()
	for _, line := range strings.Split(strings.TrimSpace(config), "\n") {
		line = strings.TrimSpace(line)
		if !strings.Contains(got, line) {
			t.Errorf("line %q not found in serialized config:\n%s", line, got)
		}
	}
}

// Metadata is stored in the section comment of the log-profile section
// itself, not in another section type sharing the name.
func TestLogProfileMetadataRoundTrip(t *testing.T) {
	config := `log-profile myprof
  log-tag "tag1"
`
	p, err := parser.New(options.String(config))
	if err != nil {
		t.Fatal(err)
	}

	lp, err := ParseLogProfile(p, "myprof")
	if err != nil {
		t.Fatal(err)
	}

	lp.Metadata = map[string]any{"comment": "managed"}
	if err := SerializeLogProfile(p, lp); err != nil {
		t.Fatalf("SerializeLogProfile: %v", err)
	}

	got, err := ParseLogProfile(p, "myprof")
	if err != nil {
		t.Fatal(err)
	}
	if got.Metadata == nil || got.Metadata["comment"] != "managed" {
		t.Errorf("metadata not round-tripped, got %#v", got.Metadata)
	}
}
