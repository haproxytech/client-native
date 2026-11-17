/*
Copyright 2026 HAProxy Technologies

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

package parser_test

import (
	"slices"
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	"github.com/haproxytech/client-native/v6/config-parser/options"
)

const diffBase = `global
  maxconn 100

defaults d1
  mode http

defaults d2
  mode tcp

frontend fe from d1
  maxconn 50
  bind :80
  default_backend be

backend be from d1
  balance roundrobin
  # full-line comment
  server s1 127.0.0.1:80 # line comment
  server s2 127.0.0.1:81
  this-is-unknown foo
`

func newDiffParser(t *testing.T, cfg string) parser.Parser {
	t.Helper()
	p, err := parser.New(options.String(cfg))
	if err != nil {
		t.Fatalf("parser.New: %v", err)
	}
	return p
}

func TestDiffSection(t *testing.T) {
	tests := []struct {
		name     string
		other    string
		section  parser.Section
		sname    string
		stop     bool // fn returns false on first call
		wantOK   bool
		want     []string
		nonEmpty bool // only assert that something was reported
	}{
		{name: "identical", other: diffBase, section: parser.Backends, sname: "be", wantOK: true},
		{
			name:    "comments only",
			other:   strings.NewReplacer("# full-line comment", "# changed", "# line comment", "# other").Replace(diffBase),
			section: parser.Backends, sname: "be", wantOK: true,
		},
		{
			name:    "server changed",
			other:   strings.Replace(diffBase, "127.0.0.1:81", "127.0.0.1:82", 1),
			section: parser.Backends, sname: "be", wantOK: true, want: []string{"server"},
		},
		{
			name:    "server added",
			other:   strings.Replace(diffBase, "  server s2", "  server s3 127.0.0.1:83\n  server s2", 1),
			section: parser.Backends, sname: "be", wantOK: true, want: []string{"server"},
		},
		{
			name:    "server removed",
			other:   strings.Replace(diffBase, "  server s2 127.0.0.1:81\n", "", 1),
			section: parser.Backends, sname: "be", wantOK: true, want: []string{"server"},
		},
		{
			name:    "balance and server changed",
			other:   strings.NewReplacer("balance roundrobin", "balance leastconn", "127.0.0.1:81", "127.0.0.1:82").Replace(diffBase),
			section: parser.Backends, sname: "be", wantOK: true, want: []string{"balance", "server"},
		},
		{
			name:    "stop after first",
			other:   strings.NewReplacer("balance roundrobin", "balance leastconn", "127.0.0.1:81", "127.0.0.1:82").Replace(diffBase),
			section: parser.Backends, sname: "be", stop: true, wantOK: true, want: []string{"balance"},
		},
		{
			name:    "from changed",
			other:   strings.NewReplacer("backend be from d1", "backend be from d2", "127.0.0.1:81", "127.0.0.1:82").Replace(diffBase),
			section: parser.Backends, sname: "be", wantOK: true, want: []string{"from", "server"},
		},
		{
			name:    "section missing",
			other:   strings.Replace(diffBase, "backend be from", "backend other from", 1),
			section: parser.Backends, sname: "be", wantOK: false,
		},
		{
			name:    "frontend maxconn",
			other:   strings.Replace(diffBase, "maxconn 50", "maxconn 60", 1),
			section: parser.Frontends, sname: "fe", wantOK: true, want: []string{"maxconn"},
		},
		{
			name:    "global maxconn",
			other:   strings.Replace(diffBase, "maxconn 100", "maxconn 200", 1),
			section: parser.Global, sname: parser.GlobalSectionName, wantOK: true, want: []string{"maxconn"},
		},
		{
			name:    "unknown keyword",
			other:   strings.Replace(diffBase, "this-is-unknown foo", "this-is-unknown bar", 1),
			section: parser.Backends, sname: "be", wantOK: true, nonEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := newDiffParser(t, diffBase)
			b := newDiffParser(t, tt.other)
			var got []string
			ok := parser.DiffSection(a, b, tt.section, tt.sname, func(attr string) bool {
				got = append(got, attr)
				return !tt.stop
			})
			if ok != tt.wantOK {
				t.Fatalf("ok = %v, want %v", ok, tt.wantOK)
			}
			if tt.nonEmpty {
				if len(got) == 0 {
					t.Fatal("expected a reported attribute, got none")
				}
				return
			}
			if !slices.Equal(got, tt.want) {
				t.Fatalf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDiffSectionSameParser(t *testing.T) {
	p := newDiffParser(t, diffBase)
	ok := parser.DiffSection(p, p, parser.Backends, "be", func(attr string) bool {
		t.Errorf("unexpected diff %q", attr)
		return true
	})
	if !ok {
		t.Fatal("ok = false, want true")
	}
}

func TestSectionTypes(t *testing.T) {
	got := parser.SectionTypes(newDiffParser(t, "global\n  maxconn 1\n\ndefaults d\n  mode http\n\nfrontend fe\n  bind :80\n\nbackend be\n  balance roundrobin\n"))
	for _, want := range []parser.Section{parser.Global, parser.Defaults, parser.Frontends, parser.Backends, parser.Comments} {
		if !slices.Contains(got, want) {
			t.Errorf("missing %q in %v", want, got)
		}
	}
	if slices.Contains(got, parser.Resolvers) {
		t.Errorf("unexpected %q in %v", parser.Resolvers, got)
	}
}
