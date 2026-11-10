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

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

// TestAllMissingKeywordsCoexist parses one configuration containing all 14
// integrated settings to catch section/registration interference.
func TestAllMissingKeywordsCoexist(t *testing.T) {
	const backingFile = "/tmp/coexist-ring.buf"
	config := strings.Join([]string{
		"global",
		"  key-base /etc/haproxy/keys",
		"  log-tag gateway-a",
		"  quick-exit",
		"  tune.ssl.hard-maxrecord 16384",
		"  debug.counters on",
		"",
		"defaults",
		"  max-session-srv-conns 5",
		"  timeout client-hs 5s",
		"  option spop-check",
		"  option tcplog clf",
		"  rate-limit sessions 10",
		"",
		"frontend fe",
		"  max-session-srv-conns 6",
		"  timeout client-hs 6s",
		"  option tcplog clf",
		"  rate-limit sessions 11",
		"  quic-initial do-log",
		"  tcp-request content set-src src",
		"",
		"log-forward lf",
		"  option host append",
		"",
		"ring r",
		"  backing-file " + backingFile,
		"",
	}, "\n")

	p, err := parser.New(parser_options.String(config))
	require.NoError(t, err)

	global, err := ParseGlobalSection(p)
	require.NoError(t, err)
	require.NotNil(t, global.SslOptions)
	require.Equal(t, "/etc/haproxy/keys", global.SslOptions.KeyBase)
	require.Equal(t, "gateway-a", global.LogTag)
	require.NotNil(t, global.DebugOptions)
	require.True(t, global.DebugOptions.QuickExit)
	require.Equal(t, "enabled", global.DebugOptions.Counters)
	require.NotNil(t, global.TuneSslOptions)
	require.NotNil(t, global.TuneSslOptions.HardMaxrecord)
	require.Equal(t, int64(16384), *global.TuneSslOptions.HardMaxrecord)

	defaults := &models.DefaultsBase{}
	require.NoError(t, ParseSection(defaults, parser.Defaults, parser.DefaultsSectionName(p), p))
	require.Equal(t, int64(5), *defaults.MaxSessionSrvConns)
	require.Equal(t, int64(5000), *defaults.ClientHsTimeout)
	require.Equal(t, "enabled", defaults.SpopCheck)
	require.True(t, defaults.Tcplog)
	require.True(t, defaults.TCPClflog)
	require.Equal(t, int64(10), *defaults.RateLimitSessions)

	frontend := &models.FrontendBase{}
	require.NoError(t, ParseSection(frontend, parser.Frontends, "fe", p))
	require.Equal(t, int64(6), *frontend.MaxSessionSrvConns)
	require.Equal(t, int64(6000), *frontend.ClientHsTimeout)
	require.True(t, frontend.Tcplog)
	require.True(t, frontend.TCPClflog)
	require.Equal(t, int64(11), *frontend.RateLimitSessions)

	logForward := &models.LogForward{LogForwardBase: models.LogForwardBase{Name: "lf"}}
	require.NoError(t, ParseLogForward(p, logForward))
	require.Equal(t, "append", logForward.Host)

	ring := &models.Ring{RingBase: models.RingBase{Name: "r"}}
	require.NoError(t, ParseRingSection(p, ring))
	require.Equal(t, backingFile, ring.BackingFile)

	quicRules, err := ParseQUICInitialRules(FrontendParentName, "fe", p)
	require.NoError(t, err)
	require.Len(t, quicRules, 1)
	require.Equal(t, "do-log", quicRules[0].Type)

	tcpRules, err := ParseTCPRequestRules(FrontendParentName, "fe", p)
	require.NoError(t, err)
	require.Len(t, tcpRules, 1)
	require.Equal(t, models.TCPRequestRuleTypeContent, tcpRules[0].Type)
	require.Equal(t, models.TCPRequestRuleActionSetDashSrc, tcpRules[0].Action)
	require.Equal(t, "src", tcpRules[0].Expr)

	rendered := p.String()
	require.Contains(t, rendered, "key-base /etc/haproxy/keys")
	require.Contains(t, rendered, "log-tag gateway-a")
	require.Contains(t, rendered, "quick-exit")
	require.Contains(t, rendered, "tune.ssl.hard-maxrecord 16384")
	require.Contains(t, rendered, "debug.counters on")
	require.Contains(t, rendered, "max-session-srv-conns 5")
	require.Contains(t, rendered, "max-session-srv-conns 6")
	require.Contains(t, rendered, "timeout client-hs 5s")
	require.Contains(t, rendered, "timeout client-hs 6s")
	require.Contains(t, rendered, "option spop-check")
	require.Contains(t, rendered, "option tcplog clf")
	require.Contains(t, rendered, "rate-limit sessions 10")
	require.Contains(t, rendered, "rate-limit sessions 11")
	require.Contains(t, rendered, "quic-initial do-log")
	require.Contains(t, rendered, "tcp-request content set-src src")
	require.Contains(t, rendered, "option host append")
	require.Contains(t, rendered, "backing-file "+backingFile)
}
