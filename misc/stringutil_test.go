// Copyright 2020 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package misc

import (
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should convert leading dots",
			input: ".....file.map!!@#",
			want:  "_file.map_",
		},
		{
			name:  "Should convert hidden files",
			input: ".hidden",
			want:  "_hidden",
		},
		{
			name:  "Should accept unusual filenames",
			input: ".unusual.",
			want:  "_unusual.",
		},
		{
			name:  "Should sanitize input correctly",
			input: "#1_?a;b/c!&?",
			want:  "_1__a_b_c_",
		},
		{
			name:  "Should return same input when name doesn't contain regex characters",
			input: "abcDEF",
			want:  "abcDEF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeFilename(tt.input); got != tt.want {
				t.Errorf("SanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDashCase(t *testing.T) {
	tests := []struct {
		fieldname string
		dashcase  string
	}{
		{fieldname: "Abortonclose", dashcase: "abortonclose"},
		{fieldname: "AcceptInvalidHTTPRequest", dashcase: "accept-invalid-http-request"},
		{fieldname: "AcceptInvalidHTTPResponse", dashcase: "accept-invalid-http-response"},
		{fieldname: "Backlog", dashcase: "backlog"},
		{fieldname: "Check", dashcase: "check"},
		{fieldname: "Checkcache", dashcase: "checkcache"},
		{fieldname: "Client", dashcase: "client"},
		{fieldname: "ClientFin", dashcase: "client-fin"},
		{fieldname: "Clitcpka", dashcase: "clitcpka"},
		{fieldname: "ClitcpkaCnt", dashcase: "clitcpka-cnt"},
		{fieldname: "Connect", dashcase: "connect"},
		{fieldname: "Contstats", dashcase: "contstats"},
		{fieldname: "DisableH2Upgrade", dashcase: "disable-h2-upgrade"},
		{fieldname: "Disabled", dashcase: "disabled"},
		{fieldname: "DontlogNormal", dashcase: "dontlog-normal"},
		{fieldname: "Dontlognull", dashcase: "dontlognull"},
		{fieldname: "DynamicCookieKey", dashcase: "dynamic-cookie-key"},
		{fieldname: "Enabled", dashcase: "enabled"},
		{fieldname: "ErrorLogFormat", dashcase: "error-log-format"},
		{fieldname: "Fullconn", dashcase: "fullconn"},
		{fieldname: "H1CaseAdjustBogusClient", dashcase: "h1-case-adjust-bogus-client"},
		{fieldname: "H1CaseAdjustBogusServer", dashcase: "h1-case-adjust-bogus-server"},
		{fieldname: "HashBalanceFactor", dashcase: "hash-balance-factor"},
		{fieldname: "HTTPBufferRequest", dashcase: "http-buffer-request"},
		{fieldname: "HTTPIgnoreProbes", dashcase: "http-ignore-probes"},
		{fieldname: "HTTPKeepAlive", dashcase: "http-keep-alive"},
		{fieldname: "HTTPNoDelay", dashcase: "http-no-delay"},
		{fieldname: "HTTPPretendKeepalive", dashcase: "http-pretend-keepalive"},
		{fieldname: "HTTPRequest", dashcase: "http-request"},
		{fieldname: "HTTPUseHtx", dashcase: "http-use-htx"},
		{fieldname: "HTTPUseProxyHeader", dashcase: "http-use-proxy-header"},
		{fieldname: "HttpchkParams", dashcase: "httpchk-params"},
		{fieldname: "Httpslog", dashcase: "httpslog"},
		{fieldname: "IdleCloseOnResponse", dashcase: "idle-close-on-response"},
		{fieldname: "IndependentStreams", dashcase: "independent-streams"},
		{fieldname: "LogFormat", dashcase: "log-format"},
		{fieldname: "LogFormatSd", dashcase: "log-format-sd"},
		{fieldname: "LogHealthChecks", dashcase: "log-health-checks"},
		{fieldname: "LogSeparateErrors", dashcase: "log-separate-errors"},
		{fieldname: "LogTag", dashcase: "log-tag"},
		{fieldname: "MaxKeepAliveQueue", dashcase: "max-keep-alive-queue"},
		{fieldname: "Maxconn", dashcase: "maxconn"},
		{fieldname: "Mode", dashcase: "mode"},
		{fieldname: "MysqlCheckParams", dashcase: "mysql-check-params"},
		{fieldname: "Name", dashcase: "name"},
		{fieldname: "Nolinger", dashcase: "nolinger"},
		{fieldname: "Persist", dashcase: "persist"},
		{fieldname: "PgsqlCheckParams", dashcase: "pgsql-check-params"},
		{fieldname: "PreferLastServer", dashcase: "prefer-last-server"},
		{fieldname: "Queue", dashcase: "queue"},
		{fieldname: "Retries", dashcase: "retries"},
		{fieldname: "RetryOn", dashcase: "retry-on"},
		{fieldname: "Server", dashcase: "server"},
		{fieldname: "ServerFin", dashcase: "server-fin"},
		{fieldname: "SmtpchkParams", dashcase: "smtpchk-params"},
		{fieldname: "SocketStats", dashcase: "socket-stats"},
		{fieldname: "SpliceAuto", dashcase: "splice-auto"},
		{fieldname: "SpliceRequest", dashcase: "splice-request"},
		{fieldname: "SpliceResponse", dashcase: "splice-response"},
		{fieldname: "Srvtcpka", dashcase: "srvtcpka"},
		{fieldname: "SrvtcpkaCnt", dashcase: "srvtcpka-cnt"},
		{fieldname: "Tarpit", dashcase: "tarpit"},
		{fieldname: "TCPSmartAccept", dashcase: "tcp-smart-accept"},
		{fieldname: "TCPSmartConnect", dashcase: "tcp-smart-connect"},
		{fieldname: "Tcpka", dashcase: "tcpka"},
		{fieldname: "Tcplog", dashcase: "tcplog"},
		{fieldname: "Transparent", dashcase: "transparent"},
		{fieldname: "Tunnel", dashcase: "tunnel"},
	}

	for _, tt := range tests {
		t.Run(tt.fieldname, func(t *testing.T) {
			if got := DashCase(tt.fieldname); got != tt.dashcase {
				t.Errorf("DashCase(%s) = %s, want %s", tt.fieldname, got, tt.dashcase)
			}
		})
	}
}
