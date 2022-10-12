// Copyright 2019 HAProxy Technologies
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
	"io/ioutil"
	"testing"

	"github.com/haproxytech/client-native/v5/models"
)

func generateTCPCheckConfig(config string) (string, error) {
	f, err := ioutil.TempFile("/tmp", "tcp_check")
	if err != nil {
		return "", err
	}
	err = prepareTestFile(config, f.Name())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func (self Counter) current() int64 {
	return self.count
}

func (self *Counter) increment() int64 {
	self.count++
	return self.count
}

func TestGetTCPCheck(t *testing.T) {
	config := `# _version=1
global
	daemon

defaults
	maxconn 2000
	mode tcp
	option tcp-check

backend test
	mode tcp
	option tcp-check
	balance roundrobin
	bind-process all

	tcp-check connect
	tcp-check comment GET\ phase
	tcp-check send GET\ /\ HTTP/1.0\r\n
	tcp-check send Host:\ haproxy.org\r\n
	tcp-check send \r\n
	tcp-check expect rstring (2..|3..)

	tcp-check connect port 443 ssl
	tcp-check send GET\ /\ HTTP/2.0\r\n
	tcp-check send Host:\ haproxy.org\r\n
	tcp-check send \r\n comment
	tcp-check expect rstring (2..|3..)

	tcp-check send-binary 50494e470d0a comment send-binary-comment
	tcp-check expect binary 2b504F4e47 comment expect-binary-comment

	tcp-check send-binary-lf 50494e470d0a
	tcp-check expect binary-lf 2b504F4e47

	tcp-check set-var(req.my_var) req.fhdr(user-agent),lower
	tcp-check unset-var(req.my_var)

	tcp-check send-lf fmt
	tcp-check send-lf fmt comment this-is-the-comment

	tcp-check set-var-fmt(check.name) "%H"
	tcp-check set-var-fmt(txn.from) "addr=%[src]:%[src_port]"

	tcp-check connect port 443 addr 192.168.0.1 send-proxy via-socks4 ssl sni sni-value alpn http/1.1,http/1.0 proto HTTP linger
	tcp-check expect min-recv 1 comment my-comment ok-status L60K error-status L6RSP tout-status L6TOUT on-success on-success-fmt on-error on-error-fmt status-code status-code-expr rstring (2..|3..)
`
	configFile, err := generateConfig(config)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = deleteTestFile(configFile)
	}()

	tests := []struct {
		name              string
		configurationFile string
		want              int64
		wantErr           bool
	}{
		{
			name:              "tcp-checks",
			configurationFile: configFile,
			want:              1,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := prepareClient(tt.configurationFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareClient error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			_, checks, err := c.GetTCPChecks("backend", "test", "")
			if err != nil {
				t.Error(err.Error())
			}

			if checks == nil {
				t.Errorf("No tcp-checks found, expected 10")
			}

			counter := Counter{0}

			_, check, err := c.GetTCPCheck(0, "backend", "test", "")
			if check.Action != models.TCPCheckActionConnect {
				t.Errorf("Check action %v returned, expected %v", check.Action, models.TCPCheckActionConnect)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionComment {
				t.Errorf("Check action %v returned, expected %v", check.Action, models.TCPCheckActionComment)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `GET\ /\ HTTP/1.0\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `Host:\ haproxy.org\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionExpect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionExpect)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionConnect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionConnect)
			}
			if check.PortString != "443" {
				t.Errorf("tcp-check connect port returned %v, expected %v", check.PortString, "443")
			}
			if !check.Ssl {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `GET\ /\ HTTP/2.0\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `Host:\ haproxy.org\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSend {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSend)
			}
			if check.Data != `\r\n` {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionExpect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionExpect)
			}
			if check.Match != "rstring" {
				t.Errorf("tcp-check expect match %v returned, expected %v", check.Match, "rstring")
			}
			if check.Pattern != "(2..|3..)" {
				t.Errorf("tcp-check expect pattern %v returned, expected %v", check.Pattern, "(2..|3..)")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSendDashBinary {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSendDashBinary)
			}
			if check.HexString != "50494e470d0a" {
				t.Errorf("tcp-check send-binary hex-string is invalid")
			}
			if check.CheckComment != "send-binary-comment" {
				t.Errorf("tcp-check send-binary comment %v returned, expected %v", check.CheckComment, "send-binary-comment")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionExpect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionExpect)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSendDashBinaryDashLf {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSendDashBinaryDashLf)
			}
			if check.HexFmt != "50494e470d0a" {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionExpect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionExpect)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSetDashVar {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSetDashVar)
			}
			if check.VarScope != "req" {
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionUnsetDashVar {
				t.Errorf("tcp-check action data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSendDashLf {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSendDashLf)
			}
			if check.Fmt != "fmt" {
				t.Errorf("tcp-check %v - fmt data is invalid", models.TCPCheckActionSendDashLf)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSendDashLf {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSendDashLf)
			}
			if check.Fmt != "fmt" {
				t.Errorf("tcp-check action fmt data is invalid")
			}
			if check.CheckComment != "this-is-the-comment" {
				t.Errorf("tcp-check action comment data is invalid")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSetDashVarDashFmt {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSetDashVarDashFmt)
			}
			if check.VarScope != "check" {
				t.Errorf("tcp-check set-var-fmt scope returned %v, expected %v", check.VarScope, "check")
			}
			if check.VarName != "name" {
				t.Errorf("tcp-check set-var-fmt name returned %v, expected %v", check.VarName, "name")
			}
			if check.VarFmt != `"%H"` {
				t.Errorf("tcp-check set-var-fmt format returned %v, expected %v", check.VarFmt, `"%H"`)
			}
			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionSetDashVarDashFmt {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionSetDashVarDashFmt)
			}
			if check.VarScope != "txn" {
				t.Errorf("tcp-check set-var-fmt scope returned %v, expected %v", check.VarScope, "txn")
			}
			if check.VarName != "from" {
				t.Errorf("tcp-check set-var-fmt name returned %v, expected %v", check.VarName, "from")
			}
			if check.VarFmt != `"addr=%[src]:%[src_port]"` {
				t.Errorf("tcp-check set-var-fmt fmt returned %v, expected %v", check.VarFmt, `"addr=%[src]:%[src_port]"`)
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionConnect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionConnect)
			}
			if check.PortString != "443" {
				t.Errorf("tcp-check connect port returned %v, expected %v", check.PortString, "443")
			}
			if check.Addr != "192.168.0.1" {
				t.Errorf("tcp-check connect addr returned %v, expected %v", check.Addr, "192.168.0.1")
			}
			if check.SendProxy == false {
				t.Errorf("tcp-check connect send-proxy missing")
			}
			if check.ViaSocks4 == false {
				t.Errorf("tcp-check connect via-socks4 missing")
			}
			if check.Ssl == false {
				t.Errorf("tcp-check connect ssl missing")
			}
			if check.Sni != "sni-value" {
				t.Errorf("tcp-check connect sni returned %v, expected %v", check.Sni, "sni-value")
			}
			if check.Alpn != "http/1.1,http/1.0" {
				t.Errorf("tcp-check connect sni returned %v, expected %v", check.Alpn, "http/1.1,http/1.0")
			}
			if check.Proto != "HTTP" {
				t.Errorf("tcp-check connect proto returned %v, expected %v", check.Proto, "HTTP")
			}
			if check.Linger == false {
				t.Errorf("tcp-check connect linger missing")
			}

			_, check, err = c.GetTCPCheck(counter.increment(), "backend", "test", "")
			if check.Action != models.TCPCheckActionExpect {
				t.Errorf("tcp-check action %v returned, expected %v", check.Action, models.TCPCheckActionExpect)
			}
			if check.MinRecv != 1 {
				t.Errorf("tcp-check expect min-recv returned %v, expected %v", check.MinRecv, 1)
			}
			if check.CheckComment != "my-comment" {
				t.Errorf("tcp-check expect comment returned %v, expected %v", check.CheckComment, "my-comment")
			}
			if check.OkStatus != "L60K" {
				t.Errorf("tcp-check expect ok-status returned %v, expected %v", check.OkStatus, "L60K")
			}
			if check.ErrorStatus != "L6RSP" {
				t.Errorf("tcp-check expect error-status returned %v, expected %v", check.ErrorStatus, "L6RSP")
			}
			if check.ToutStatus != "L6TOUT" {
				t.Errorf("tcp-check expect tout-status returned %v, expected %v", check.ToutStatus, "L6TOUT")
			}
			if check.OnSuccess != "on-success-fmt" {
				t.Errorf("tcp-check expect on-success returned %v, expected %v", check.OnSuccess, "on-success-fmt")
			}
			if check.OnError != "on-error-fmt" {
				t.Errorf("tcp-check expect on-error returned %v, expected %v", check.OnError, "on-error-fmt")
			}
			if check.StatusCode != "status-code-expr" {
				t.Errorf("tcp-check expect status-code returned %v, expected %v", check.OnError, "status-code-expr")
			}
			if check.Match != "rstring" {
				t.Errorf("tcp-check expect match-code returned %v, expected %v", check.Match, "rstring")
			}
			if check.Pattern != "(2..|3..)" {
				t.Errorf("tcp-check expect pattern returned %v, expected %v", check.Pattern, "(2..|3..)")
			}
		})
	}
}
