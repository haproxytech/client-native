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
	"reflect"
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
)

func TestOptionForwardedParseAttributes(t *testing.T) {
	line := "option forwarded proto host by by_port for for_port"
	want := forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
		f.Proto = true
		f.Host = true
		f.By = true
		f.ByPort = true
		f.For = true
		f.ForPort = true
	})

	defaults, frontend, backend := parseOptionForwardedSections(t, line)
	assertForwardedModel(t, "defaults", defaults.Forwarded, want)
	assertForwardedModel(t, "frontend", frontend.Forwarded, want)
	assertForwardedModel(t, "backend", backend.Forwarded, want)
}

func TestOptionForwardedParseExpressions(t *testing.T) {
	line := "option forwarded proto host-expr %[req.hdr(host)] by-expr %[dst] by_port-expr %[dst_port] for-expr %[src] for_port-expr %[src_port]"
	want := forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
		f.Proto = true
		f.HostExpr = "%[req.hdr(host)]"
		f.ByExpr = "%[dst]"
		f.ByPortExpr = "%[dst_port]"
		f.ForExpr = "%[src]"
		f.ForPortExpr = "%[src_port]"
	})

	defaults, frontend, backend := parseOptionForwardedSections(t, line)
	assertForwardedModel(t, "defaults", defaults.Forwarded, want)
	assertForwardedModel(t, "frontend", frontend.Forwarded, want)
	assertForwardedModel(t, "backend", backend.Forwarded, want)
}

func TestOptionForwardedParseNoOption(t *testing.T) {
	want := forwardedModel(models.ForwardedEnabledDisabled, nil)

	defaults, frontend, backend := parseOptionForwardedSections(t, "no option forwarded")
	assertForwardedModel(t, "defaults", defaults.Forwarded, want)
	assertForwardedModel(t, "frontend", frontend.Forwarded, want)
	assertForwardedModel(t, "backend", backend.Forwarded, want)
}

func TestOptionForwardedRenderAttributes(t *testing.T) {
	model := forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
		f.Proto = true
		f.Host = true
		f.By = true
		f.ByPort = true
		f.For = true
		f.ForPort = true
	})
	p := newOptionForwardedParser(t)

	serializeOptionForwardedSections(t, p, model)

	assertConfigContains(t, p.String(), "option forwarded proto host by by_port for for_port")
}

func TestOptionForwardedRenderNoOption(t *testing.T) {
	p := newOptionForwardedParser(t)

	serializeOptionForwardedSections(t, p, forwardedModel(models.ForwardedEnabledDisabled, nil))

	assertConfigContains(t, p.String(), "no option forwarded")
}

func TestOptionForwardedUpdateEnabledAttributes(t *testing.T) {
	p := newOptionForwardedParserWithLine(t, "option forwarded proto host")
	model := forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
		f.ByExpr = "%[dst]"
		f.ForPort = true
	})

	serializeOptionForwardedSections(t, p, model)

	config := p.String()
	assertConfigContains(t, config, "option forwarded by-expr %[dst] for_port")
	assertConfigLineCount(t, config, "option forwarded proto host", 0)
}

func TestOptionForwardedUpdateEnabledToNoOption(t *testing.T) {
	p := newOptionForwardedParserWithLine(t, "option forwarded proto host")

	serializeOptionForwardedSections(t, p, forwardedModel(models.ForwardedEnabledDisabled, nil))

	config := p.String()
	assertConfigContains(t, config, "no option forwarded")
	assertConfigLineCount(t, config, "option forwarded proto host", 0)
}

func TestOptionForwardedUpdateNoOptionToEnabledAttributes(t *testing.T) {
	p := newOptionForwardedParserWithLine(t, "no option forwarded")
	model := forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
		f.HostExpr = "%[req.hdr(host)]"
		f.For = true
	})

	serializeOptionForwardedSections(t, p, model)

	config := p.String()
	assertConfigContains(t, config, "option forwarded host-expr %[req.hdr(host)] for")
	assertConfigLineCount(t, config, "no option forwarded", 0)
}

func TestOptionForwardedRenderErrors(t *testing.T) {
	tests := []struct {
		name  string
		model *models.Forwarded
	}{
		{
			name:  "empty enabled",
			model: &models.Forwarded{},
		},
		{
			name:  "invalid enabled",
			model: forwardedModel("unexpected", nil),
		},
		{
			name: "disabled with attributes",
			model: forwardedModel(models.ForwardedEnabledDisabled, func(f *models.Forwarded) {
				f.Proto = true
			}),
		},
		{
			name: "host conflict",
			model: forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
				f.Host = true
				f.HostExpr = "%[req.hdr(host)]"
			}),
		},
		{
			name: "by port conflict",
			model: forwardedModel(models.ForwardedEnabledEnabled, func(f *models.Forwarded) {
				f.ByPort = true
				f.ByPortExpr = "%[dst_port]"
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := newOptionForwardedParser(t)
			err := CreateEditSection(&models.DefaultsBase{Forwarded: tt.model}, parser.Defaults, parser.DefaultSectionName, p, nil)
			if err == nil {
				t.Fatal("CreateEditSection returned nil error")
			}
		})
	}
}

func parseOptionForwardedSections(t *testing.T, line string) (*models.DefaultsBase, *models.FrontendBase, *models.BackendBase) {
	t.Helper()

	config := strings.Join([]string{
		"defaults",
		"  " + line,
		"frontend fe",
		"  " + line,
		"backend be",
		"  " + line,
		"",
	}, "\n")
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}

	defaults := &models.DefaultsBase{}
	if err := ParseSection(defaults, parser.Defaults, parser.DefaultSectionName, p); err != nil {
		t.Fatalf("ParseSection defaults: %v", err)
	}
	frontend := &models.FrontendBase{}
	if err := ParseSection(frontend, parser.Frontends, "fe", p); err != nil {
		t.Fatalf("ParseSection frontend: %v", err)
	}
	backend := &models.BackendBase{}
	if err := ParseSection(backend, parser.Backends, "be", p); err != nil {
		t.Fatalf("ParseSection backend: %v", err)
	}
	return defaults, frontend, backend
}

func serializeOptionForwardedSections(t *testing.T, p parser.Parser, forwarded *models.Forwarded) {
	t.Helper()

	if err := CreateEditSection(&models.DefaultsBase{Forwarded: forwarded}, parser.Defaults, parser.DefaultSectionName, p, nil); err != nil {
		t.Fatalf("CreateEditSection defaults: %v", err)
	}
	if err := CreateEditSection(&models.FrontendBase{Forwarded: forwarded}, parser.Frontends, "fe", p, nil); err != nil {
		t.Fatalf("CreateEditSection frontend: %v", err)
	}
	if err := CreateEditSection(&models.BackendBase{Forwarded: forwarded}, parser.Backends, "be", p, nil); err != nil {
		t.Fatalf("CreateEditSection backend: %v", err)
	}
}

func newOptionForwardedParser(t *testing.T) parser.Parser {
	t.Helper()

	p, err := parser.New(parser_options.String("defaults\nfrontend fe\nbackend be\n"))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

func newOptionForwardedParserWithLine(t *testing.T, line string) parser.Parser {
	t.Helper()

	config := strings.Join([]string{
		"defaults",
		"  " + line,
		"frontend fe",
		"  " + line,
		"backend be",
		"  " + line,
		"",
	}, "\n")
	p, err := parser.New(parser_options.String(config))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

func forwardedModel(enabled string, mutate func(*models.Forwarded)) *models.Forwarded {
	f := &models.Forwarded{Enabled: &enabled}
	if mutate != nil {
		mutate(f)
	}
	return f
}

func assertForwardedModel(t *testing.T, section string, got, want *models.Forwarded) {
	t.Helper()

	if got == nil {
		t.Fatalf("%s Forwarded is nil", section)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s Forwarded = %#v, want %#v", section, got, want)
	}
}

func assertConfigContains(t *testing.T, config, line string) {
	t.Helper()

	count := countConfigLine(config, line)
	if count != 3 {
		t.Fatalf("config contains %q %d times, want 3:\n%s", line, count, config)
	}
}

func assertConfigLineCount(t *testing.T, config, line string, want int) {
	t.Helper()

	if count := countConfigLine(config, line); count != want {
		t.Fatalf("config contains %q %d times, want %d:\n%s", line, count, want, config)
	}
}

func countConfigLine(config, line string) int {
	count := 0
	for _, configLine := range strings.Split(config, "\n") {
		if strings.TrimSpace(configLine) == line {
			count++
		}
	}
	return count
}
