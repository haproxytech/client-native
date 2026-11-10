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
	tcp_actions "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/actions"
	tcp_types "github.com/haproxytech/client-native/v6/config-parser/parsers/tcp/types"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestTCPContentSetSrc(t *testing.T) {
	tests := []struct {
		name       string
		parentType string
		section    string
		config     string
		want       models.TCPRequestRule
	}{
		{
			name:       "frontend",
			parentType: FrontendParentName,
			section:    "fe",
			config:     "frontend fe\n  tcp-request content set-src src\n",
			want: models.TCPRequestRule{
				Type:   models.TCPRequestRuleTypeContent,
				Action: models.TCPRequestRuleActionSetDashSrc,
				Expr:   "src",
			},
		},
		{
			name:       "backend conditional",
			parentType: BackendParentName,
			section:    "be",
			config:     "backend be\n  tcp-request content set-src hdr(ip) if TRUE\n",
			want: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeContent,
				Action:   models.TCPRequestRuleActionSetDashSrc,
				Expr:     "hdr(ip)",
				Cond:     "if",
				CondTest: "TRUE",
			},
		},
		{
			name:       "named defaults with comment",
			parentType: DefaultsParentName,
			section:    "named_defaults",
			config:     "defaults named_defaults\n" + `  tcp-request content set-src src unless { src 10.0.0.1 } # {"comment":"conditional"}` + "\n",
			want: models.TCPRequestRule{
				Type:     models.TCPRequestRuleTypeContent,
				Action:   models.TCPRequestRuleActionSetDashSrc,
				Expr:     "src",
				Cond:     "unless",
				CondTest: "{ src 10.0.0.1 }",
				Metadata: map[string]interface{}{"comment": "conditional"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p, err := parser.New(parser_options.String(tt.config))
			require.NoError(t, err)

			rules, err := ParseTCPRequestRules(tt.parentType, tt.section, p)
			require.NoError(t, err)
			require.Len(t, rules, 1)
			require.True(t, rules[0].Equal(tt.want), "got %#v, want %#v", rules[0], tt.want)
			require.Equal(t, tt.want.Metadata, rules[0].Metadata, "metadata")

			serialized, err := SerializeTCPRequestRule(*rules[0], nil)
			require.NoError(t, err)
			content, ok := serialized.(*tcp_types.Content)
			require.True(t, ok, "serialized type = %T, want *tcp_types.Content", serialized)
			_, ok = content.Action.(*tcp_actions.SetSrc)
			require.True(t, ok, "action type = %T, want *tcp_actions.SetSrc", content.Action)

			// model -> serialize -> fresh parse must round-trip.
			line := "  tcp-request " + content.String()
			if comment := content.GetComment(); comment != "" {
				line += " # " + comment
			}
			header := tt.config[:strings.Index(tt.config, "\n")+1]
			fresh, err := parser.New(parser_options.String(header + line + "\n"))
			require.NoError(t, err)
			freshRules, err := ParseTCPRequestRules(tt.parentType, tt.section, fresh)
			require.NoError(t, err)
			require.Len(t, freshRules, 1)
			require.True(t, freshRules[0].Equal(*rules[0]), "round-trip: got %#v, want %#v", freshRules[0], rules[0])
			require.Equal(t, rules[0].Metadata, freshRules[0].Metadata, "round-trip metadata")
		})
	}
}

func TestTCPContentSetSrcInvalid(t *testing.T) {
	for _, parts := range [][]string{
		{"tcp-request", "content", "set-src"},
		{"tcp-request", "content", "set-src", "if", "TRUE"},
		{"tcp-request", "content", "set-src", "unless"},
	} {
		t.Run(strings.Join(parts, " "), func(t *testing.T) {
			c := &tcp_types.Content{}
			require.Error(t, c.Parse(parts, ""))
		})
	}
}

func TestTCPConnectionSessionSetSrcSerialization(t *testing.T) {
	for _, ruleType := range []string{
		models.TCPRequestRuleTypeConnection,
		models.TCPRequestRuleTypeSession,
	} {
		t.Run(ruleType, func(t *testing.T) {
			rule := models.TCPRequestRule{
				Type:     ruleType,
				Action:   models.TCPRequestRuleActionSetDashSrc,
				Expr:     "src",
				Cond:     "if",
				CondTest: "TRUE",
			}
			serialized, err := SerializeTCPRequestRule(rule, nil)
			require.NoError(t, err)

			var action interface {
				String() string
			}
			switch v := serialized.(type) {
			case *tcp_types.Connection:
				action = v.Action
			case *tcp_types.Session:
				action = v.Action
			default:
				t.Fatalf("serialized type = %T", serialized)
			}
			_, ok := action.(*tcp_actions.SetSrc)
			require.True(t, ok, "action type = %T, want *tcp_actions.SetSrc", action)
			require.Equal(t, "set-src src if TRUE", action.String())
		})
	}
}
