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

	"github.com/stretchr/testify/require"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/types"
	"github.com/haproxytech/client-native/v6/runtime"
)

const runtimeApplyBase = `global
  maxconn 100

defaults d
  mode http

frontend fe from d
  maxconn 50
  bind :80
  default_backend be

backend be from d
  balance roundrobin
  # comment
  server s1 10.0.0.1:80 weight 10 check
  server s2 10.0.0.2:80 weight 10
`

func newRuntimeApplyParser(t testing.TB, cfg string) parser.Parser {
	t.Helper()
	p, err := parser.New(parser_options.String(cfg))
	require.NoError(t, err)
	return p
}

// testServerItems turns the server lines of a section into items: the
// address first, then one option per server parameter.
func testServerItems(p parser.Parser, t parser.Section, section string) ([]item, error) {
	data, err := p.Get(t, section, "server")
	if err != nil {
		return nil, err
	}
	servers, _ := data.([]types.Server)
	items := make([]item, 0, len(servers))
	for _, srv := range servers {
		it := item{name: srv.Name, options: []option{{"address", srv.Address}}}
		for _, param := range srv.Params {
			name, args, _ := strings.Cut(param.String(), " ")
			it.options = append(it.options, option{name, args})
		}
		items = append(items, it)
	}
	return items, nil
}

// stepLog records the labels of the steps that were run.
type stepLog []string

func (l *stepLog) step(label string) step {
	return func(runtime.Runtime) error {
		*l = append(*l, label)
		return nil
	}
}

func (l *stepLog) itemField(name string) itemApplier {
	return func(_ string, _, cur item) ([]step, bool) {
		return []step{l.step(name + " " + cur.name)}, true
	}
}

func (l *stepLog) itemOp(op string) func(string, item) ([]step, bool) {
	return func(_ string, it item) ([]step, bool) {
		return []step{l.step(op + " " + it.name)}, true
	}
}

func testRuntimeRegistry(t *testing.T, l *stepLog) runtimeRegistry {
	t.Helper()
	return runtimeRegistry{
		parser.Frontends: {
			"maxconn": {apply: func(c change) ([]step, bool) {
				if c.old == nil || c.new == nil {
					t.Errorf("maxconn change has nil data: %+v", c)
				}
				return []step{l.step("maxconn " + c.section)}, true
			}},
		},
		parser.Backends: {
			"server": {object: &objectSpec{
				items:  testServerItems,
				fields: map[string]itemApplier{"weight": l.itemField("weight"), "check": l.itemField("check")},
				add:    l.itemOp("add"),
				del:    l.itemOp("del"),
			}},
		},
	}
}

func runSteps(t *testing.T, steps []step) {
	t.Helper()
	for _, s := range steps {
		require.NoError(t, s(nil))
	}
}

func TestRuntimeStepsPackageRegistry(t *testing.T) {
	for name, cfg := range map[string]string{
		"identical":    runtimeApplyBase,
		"comment only": strings.Replace(runtimeApplyBase, "# comment", "# other comment", 1),
	} {
		t.Run(name, func(t *testing.T) {
			steps, ok := runtimeSteps(registry, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, cfg))
			require.True(t, ok)
			require.Empty(t, steps)
		})
	}
}

func TestRuntimeSteps(t *testing.T) {
	r := strings.Replace
	tests := []struct {
		name   string
		cfg    string
		mutate func(*objectSpec)
		wantOK bool
		want   []string
	}{
		{name: "identical", cfg: runtimeApplyBase, wantOK: true},
		{name: "no registry entry", cfg: r(runtimeApplyBase, "balance roundrobin", "balance leastconn", 1)},
		{name: "backend added", cfg: runtimeApplyBase + "\nbackend be2\n  balance roundrobin\n"},
		{name: "backend removed", cfg: r(runtimeApplyBase, "backend be from d", "backend other from d", 1)},
		{name: "frontend maxconn", cfg: r(runtimeApplyBase, "maxconn 50", "maxconn 60", 1), wantOK: true, want: []string{"maxconn fe"}},
		{name: "server added", cfg: runtimeApplyBase + "  server s3 10.0.0.3:80\n", wantOK: true, want: []string{"add s3"}},
		{name: "server removed", cfg: r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "", 1), wantOK: true, want: []string{"del s2"}},
		{name: "weight changed", cfg: r(runtimeApplyBase, "weight 10 check", "weight 20 check", 1), wantOK: true, want: []string{"weight s1"}},
		{name: "address changed", cfg: r(runtimeApplyBase, "10.0.0.1:80", "10.0.0.9:80", 1)},
		{name: "two options changed", cfg: r(runtimeApplyBase, "weight 10 check", "weight 20", 1), wantOK: true, want: []string{"weight s1", "check s1"}},
		{name: "servers reordered", cfg: r(runtimeApplyBase, "  server s1 10.0.0.1:80 weight 10 check\n  server s2 10.0.0.2:80 weight 10\n", "  server s2 10.0.0.2:80 weight 10\n  server s1 10.0.0.1:80 weight 10 check\n", 1)},
		{name: "server inserted before kept", cfg: r(runtimeApplyBase, "  server s1", "  server s0 10.0.0.5:80\n  server s1", 1)},
		{name: "duplicate server names", cfg: runtimeApplyBase + "  server s2 10.0.0.2:80 weight 10\n"},
		{name: "nil add", cfg: runtimeApplyBase + "  server s3 10.0.0.3:80\n", mutate: func(s *objectSpec) { s.add = nil }},
		{name: "nil del", cfg: r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "", 1), mutate: func(s *objectSpec) { s.del = nil }},
		{name: "options reordered", cfg: r(runtimeApplyBase, "weight 10 check", "check weight 10", 1)},
		{
			name:   "mixed change order",
			cfg:    r(r(runtimeApplyBase, "  server s2 10.0.0.2:80 weight 10\n", "  server s3 10.0.0.3:80\n", 1), "weight 10 check", "weight 20 check", 1),
			wantOK: true, want: []string{"del s2", "weight s1", "add s3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var l stepLog
			reg := testRuntimeRegistry(t, &l)
			if tt.mutate != nil {
				tt.mutate(reg[parser.Backends]["server"].object)
			}
			steps, ok := runtimeSteps(reg, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, tt.cfg))
			require.Equal(t, tt.wantOK, ok)
			if !ok {
				require.Nil(t, steps)
				return
			}
			runSteps(t, steps)
			require.Equal(t, tt.want, []string(l))
		})
	}
}

func TestRuntimeStepsGlobalBeforeBackends(t *testing.T) {
	reg := runtimeRegistry{
		parser.Backends: {"server": {object: &objectSpec{
			items: func(parser.Parser, parser.Section, string) ([]item, error) {
				t.Fatal("backend items must not be read")
				return nil, nil
			},
		}}},
	}
	cfg := strings.NewReplacer("maxconn 100", "maxconn 200", "weight 10 check", "weight 20 check").Replace(runtimeApplyBase)
	_, ok := runtimeSteps(reg, newRuntimeApplyParser(t, runtimeApplyBase), newRuntimeApplyParser(t, cfg))
	require.False(t, ok)
}

func TestDiffObjectItemGet(t *testing.T) {
	it := item{name: "s1", options: []option{{"weight", "1"}, {"check", ""}, {"weight", "2"}}}
	args, ok := it.get("weight")
	require.True(t, ok)
	require.Equal(t, "2", args)
	_, ok = it.get("maxconn")
	require.False(t, ok)
}
