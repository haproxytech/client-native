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
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/quic"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestQUICInitialDoLog(t *testing.T) {
	lines := []string{
		"  quic-initial do-log",
		`  quic-initial do-log if TRUE # {"comment":"conditional"}`,
		`  quic-initial do-log unless { src 10.0.0.1 }`,
	}
	want := models.QUICInitialRules{
		{Type: "do-log"},
		{
			Type:     "do-log",
			Cond:     "if",
			CondTest: "TRUE",
			Metadata: map[string]interface{}{"comment": "conditional"},
		},
		{
			Type:     "do-log",
			Cond:     "unless",
			CondTest: "{ src 10.0.0.1 }",
		},
	}

	tests := []struct {
		name       string
		parentType string
		section    string
		header     string
	}{
		{"frontend", FrontendParentName, "fe", "frontend fe"},
		{"named defaults", DefaultsParentName, "named_defaults", "defaults named_defaults"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := tt.header + "\n" + strings.Join(lines, "\n") + "\n"
			p, err := parser.New(parser_options.String(config))
			require.NoError(t, err)

			rules, err := ParseQUICInitialRules(tt.parentType, tt.section, p)
			require.NoError(t, err)
			require.Len(t, rules, len(want))
			for i := range want {
				require.True(t, rules[i].Equal(*want[i]), "rule %d: got %#v, want %#v", i, rules[i], want[i])
				require.Equal(t, want[i].Metadata, rules[i].Metadata, "rule %d metadata", i)
			}

			// model -> serialize -> fresh parse must preserve order, conditions and comments.
			var serialized []string
			for _, rule := range rules {
				action, err := SerializeQUICInitialRule(*rule)
				require.NoError(t, err)
				line := "  quic-initial " + action.String()
				if comment := action.GetComment(); comment != "" {
					line += " # " + comment
				}
				serialized = append(serialized, line)
			}
			fresh, err := parser.New(parser_options.String(tt.header + "\n" + strings.Join(serialized, "\n") + "\n"))
			require.NoError(t, err)
			freshRules, err := ParseQUICInitialRules(tt.parentType, tt.section, fresh)
			require.NoError(t, err)
			require.Len(t, freshRules, len(rules))
			for i := range rules {
				require.True(t, freshRules[i].Equal(*rules[i]), "round-trip rule %d: got %#v, want %#v", i, freshRules[i], rules[i])
				require.Equal(t, rules[i].Metadata, freshRules[i].Metadata, "round-trip rule %d metadata", i)
			}
		})
	}
}

func TestQUICInitialDoLogInvalid(t *testing.T) {
	for _, line := range []string{
		"quic-initial do-log bogus token",
		"quic-initial do-log bogus if TRUE",
		"quic-initial do-log if",
	} {
		t.Run(line, func(t *testing.T) {
			p := &quic.Initial{}
			p.Init()
			_, err := p.Parse(line, strings.Fields(line), "")
			require.Error(t, err)
		})
	}
}
