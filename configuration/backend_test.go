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
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetBackends(t *testing.T) { //nolint:gocognit,gocyclo
	v, backends, err := clientTest.GetBackends("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(backends) != 2 {
		t.Errorf("%v backends returned, expected 2", len(backends))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, b := range backends {
		if b.Name != "test" && b.Name != "test_2" {
			t.Errorf("Expected only test or test_2 backend, %v found", b.Name)
		}
		optionValue := "disabled"
		if b.Name == "test" {
			optionValue = "enabled"
		}
		if b.Name == "test_2" {
			if b.From != "test_defaults_2" {
				t.Errorf("%v: From not test_defaults_2: %v", b.Name, b.From)
			}
		}
		if b.BindProcess != "all" {
			t.Errorf("%v: BindProcess not all: %v", b.Name, b.BindProcess)
		}
		if b.AdvCheck != "httpchk" {
			t.Errorf("%v: AdvCheck.Method not HEAD: %v", b.Name, b.AdvCheck)
		}
		if b.HttpchkParams.Method != "HEAD" {
			t.Errorf("%v: HttpchkParams.Method not HEAD: %v", b.Name, b.HttpchkParams.Method)
		}
		if b.HttpchkParams.URI != "/" {
			t.Errorf("%v: HttpchkParams.URI not HEAD: %v", b.Name, b.HttpchkParams.URI)
		}
		if b.Mode != "http" {
			t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
		}
		if *b.Balance.Algorithm != "roundrobin" {
			t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
		}
		if b.HashType.Method != "consistent" {
			t.Errorf("%v: HashType.Method not consistent: %v", b.Name, b.HashType.Method)
		}
		if b.HashType.Function != "sdbm" {
			t.Errorf("%v: HashType.Function not sdbm: %v", b.Name, b.HashType.Function)
		}
		if b.HashType.Modifier != "avalanche" {
			t.Errorf("%v: HashType.Modifier not avalanche: %v", b.Name, b.HashType.Modifier)
		}
		if b.HTTPConnectionMode != "http-keep-alive" {
			t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
		}
		if *b.Forwardfor.Enabled != "enabled" {
			t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
		}
		if *b.DefaultServer.Fall != 2000 {
			t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer.Fall)
		}
		if *b.DefaultServer.Rise != 4000 {
			t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
		}
		if *b.DefaultServer.Inter != 5000 {
			t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
		}
		if *b.DefaultServer.HealthCheckPort != 8888 {
			t.Errorf("%v: DefaultServer.HealthCheckPort not 8888: %v", b.Name, *b.DefaultServer.HealthCheckPort)
		}
		if *b.Cookie.Name != "BLA" {
			t.Errorf("%v: HTTPCookie Name not BLA: %v", b.Name, b.Cookie)
		}
		if b.Cookie.Type != "rewrite" {
			t.Errorf("%v: HTTPCookie Type not rewrite %v", b.Name, b.Cookie.Type)
		}
		if !b.Cookie.Httponly {
			t.Errorf("%v: HTTPCookie Httponly not true %v", b.Name, b.Cookie.Httponly)
		}
		if !b.Cookie.Nocache {
			t.Errorf("%v: HTTPCookie Nocache not false %v", b.Name, b.Cookie.Nocache)
		}
		if *b.CheckTimeout != 2000 {
			t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
		}
		if *b.ServerTimeout != 3000 {
			t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
		}
		if b.Tcpka != "enabled" {
			t.Errorf("%v: Tcpka not enabled: %v", b.Name, b.Tcpka)
		}
		if b.Srvtcpka != "enabled" {
			t.Errorf("%v: Srvtcpka not enabled: %v", b.Name, b.Srvtcpka)
		}
		if b.Checkcache != optionValue {
			t.Errorf("%v: Checkcache not %s: %v", b.Name, optionValue, b.Checkcache)
		}
		if b.IndependentStreams != optionValue {
			t.Errorf("%v: IndependentStreams not %s: %v", b.Name, optionValue, b.IndependentStreams)
		}
		if b.Nolinger != optionValue {
			t.Errorf("%v: Nolinger not %s: %v", b.Name, optionValue, b.Nolinger)
		}
		if b.Originalto == nil {
			t.Errorf("%v: Originalto is nil, expected not nil", b.Name)
		} else {
			if *b.Originalto.Enabled != "enabled" {
				t.Errorf("%v: Originalto.Enabled is not enabled: %v", b.Name, *b.Originalto.Enabled)
			}
			if b.Name == "test" && b.Originalto.Except != "" {
				t.Errorf("%v: Originalto.Except is not empty: %v", b.Name, b.Originalto.Except)
			}
			if b.Name == "test_2" && b.Originalto.Except != "127.0.0.1" {
				t.Errorf("%v: Originalto.Except is not 127.0.0.1: %v", b.Name, b.Originalto.Except)
			}
			if b.Originalto.Header != "X-Client-Dst" {
				t.Errorf("%v: Originalto.Header is not X-Client-Dst: %v", b.Name, b.Originalto.Header)
			}
		}
		if b.Persist != optionValue {
			t.Errorf("%v: Persist not %s: %v", b.Name, optionValue, b.Persist)
		}
		if b.PreferLastServer != optionValue {
			t.Errorf("%v: PreferLastServer not %s: %v", b.Name, optionValue, b.PreferLastServer)
		}
		if b.SpopCheck != optionValue {
			t.Errorf("%v: SpopCheck not %s: %v", b.Name, optionValue, b.SpopCheck)
		}
		if b.TCPSmartConnect != optionValue {
			t.Errorf("%v: TCPSmartConnect not %s: %v", b.Name, optionValue, b.TCPSmartConnect)
		}
		if b.Transparent != optionValue {
			t.Errorf("%v: Transparent not %s: %v", b.Name, optionValue, b.Transparent)
		}
		if b.SpliceAuto != optionValue {
			t.Errorf("%v: SpliceAuto not %s: %v", b.Name, optionValue, b.SpliceAuto)
		}
		if b.SpliceRequest != optionValue {
			t.Errorf("%v: SpliceRequest not %s: %v", b.Name, optionValue, b.SpliceRequest)
		}
		if b.SpliceResponse != optionValue {
			t.Errorf("%v: SpliceResponse not %s: %v", b.Name, optionValue, b.SpliceResponse)
		}
		if b.SrvtcpkaCnt == nil {
			t.Errorf("%v: SrvtcpkaCnt is nil", b.Name)
		} else if *b.SrvtcpkaCnt != 10 {
			t.Errorf("%v: SrvtcpkaCnt not 10: %v", b.Name, *b.SrvtcpkaCnt)
		}
		if b.SrvtcpkaIdle == nil {
			t.Errorf("%v: SrvtcpkaIdle is nil", b.Name)
		} else if *b.SrvtcpkaIdle != 10000 {
			t.Errorf("%v: SrvtcpkaIdle not 10000: %v", b.Name, *b.SrvtcpkaIdle)
		}
		if b.SrvtcpkaIntvl == nil {
			t.Errorf("%v: SrvtcpkaIntvl is nil", b.Name)
		} else if *b.SrvtcpkaIntvl != 10000 {
			t.Errorf("%v: SrvtcpkaIntvl not 10000: %v", b.Name, *b.SrvtcpkaIntvl)
		}
		if b.StatsOptions == nil {
			t.Errorf("%v: StatsOptions is nil", b.Name)
		}
		if b.StatsOptions.StatsShowModules != true {
			t.Error("StatsShowModules not set")
		}
		if b.StatsOptions.StatsRealm != true {
			t.Error("StatsRealm not set")
		}
		if b.StatsOptions.StatsRealmRealm == nil {
			t.Errorf("%v: StatsRealmRealm is nil", b.Name)
		} else if *b.StatsOptions.StatsRealmRealm != `HAProxy\\ Statistics` {
			t.Errorf("%v: StatsRealmRealm not 'HAProxy Statistics': %v", b.Name, *b.StatsOptions.StatsRealmRealm)
		}
		if len(b.StatsOptions.StatsAuths) != 2 {
			t.Errorf("%v: StatsAuths expected 2 instances got: %v", b.Name, len(b.StatsOptions.StatsAuths))
		}
		if b.StatsOptions.StatsAuths[0].User == nil {
			t.Errorf("%v: StatsAuths 0 User is nil", b.Name)
		} else if *b.StatsOptions.StatsAuths[0].User != "admin" {
			t.Errorf("%v: StatsAuths 0 User not admin: %v", b.Name, *b.StatsOptions.StatsAuths[0])
		}
		if b.StatsOptions.StatsAuths[0].Passwd == nil {
			t.Errorf("%v: StatsAuths 0 Passwd is nil", b.Name)
		} else if *b.StatsOptions.StatsAuths[0].Passwd != "AdMiN123" {
			t.Errorf("%v: StatsAuths 0 Passwd not AdMiN123: %v", b.Name, *b.StatsOptions.StatsAuths[0].Passwd)
		}
		if b.StatsOptions.StatsAuths[1].User == nil {
			t.Errorf("%v: StatsAuths 1 User is nil", b.Name)
		} else if *b.StatsOptions.StatsAuths[1].User != "admin2" {
			t.Errorf("%v: StatsAuths 1 User not admin2: %v", b.Name, *b.StatsOptions.StatsAuths[1].User)
		}
		if b.StatsOptions.StatsAuths[1].Passwd == nil {
			t.Errorf("%v: StatsAuths 1 Passwd is nil", b.Name)
		} else if *b.StatsOptions.StatsAuths[1].Passwd != "AdMiN1234" {
			t.Errorf("%v: StatsAuths 1 Passwd not AdMiN1234: %v", b.Name, *b.StatsOptions.StatsAuths[1].Passwd)
		}
		if len(b.StatsOptions.StatsHTTPRequests) != 2 {
			t.Errorf("%v: StatsHTTPRequests expected 2 instances got: %v", b.Name, len(b.StatsOptions.StatsHTTPRequests))
		}
		if b.StatsOptions.StatsHTTPRequests[0].Type == nil {
			t.Errorf("%v: StatsHTTPRequests 0 Type is nil", b.Name)
		} else if *b.StatsOptions.StatsHTTPRequests[0].Type != "auth" {
			t.Errorf("%v: StatsHTTPRequests 0 Type not auth: %v", b.Name, *b.StatsOptions.StatsHTTPRequests[0].Type)
		}
		if b.StatsOptions.StatsHTTPRequests[0].Realm != `HAProxy\\ Statistics` {
			t.Errorf("%v: StatsHTTPRequests 0 Realm not 'HAProxy Statistics': %v", b.Name, b.StatsOptions.StatsHTTPRequests[0].Realm)
		}
		if b.StatsOptions.StatsHTTPRequests[1].Type == nil {
			t.Errorf("%v: StatsHTTPRequests 1 Type is nil", b.Name)
		} else if *b.StatsOptions.StatsHTTPRequests[1].Type != "allow" {
			t.Errorf("%v: StatsHTTPRequests 1 Type not allow: %v", b.Name, *b.StatsOptions.StatsHTTPRequests[1].Type)
		}
		if b.StatsOptions.StatsHTTPRequests[1].Cond != "if" {
			t.Errorf("%v: StatsHTTPRequests 1 Cond not if: %v", b.Name, b.StatsOptions.StatsHTTPRequests[1].Cond)
		}
		if b.StatsOptions.StatsHTTPRequests[1].CondTest != "something" {
			t.Errorf("%v: StatsHTTPRequests 1 CondTest not something: %v", b.Name, b.StatsOptions.StatsHTTPRequests[1].CondTest)
		}

		if b.EmailAlert == nil {
			t.Error("EmailAlert is nil")
		} else if *b.EmailAlert.From != "prod01@example.com" {
			t.Errorf("EmailAlert.From is not prod01@example.com: %v", *b.EmailAlert.From)
		} else if *b.EmailAlert.To != "sre@example.com" {
			t.Errorf("EmailAlert.To is not sre@example.com: %v", *b.EmailAlert.To)
		} else if b.EmailAlert.Level != "warning" {
			t.Errorf("EmailAlert.Level is not warning: %v", b.EmailAlert.Level)
		} else if b.EmailAlert.Myhostname != "prod01" {
			t.Errorf("EmailAlert.Myhostname is not prod01: %v", b.EmailAlert.Myhostname)
		} else if *b.EmailAlert.Mailers != "localmailer1" {
			t.Errorf("EmailAlert.Mailers is not localmailer1: %v", *b.EmailAlert.Mailers)
		}

	}
}

func TestGetBackend(t *testing.T) {
	v, b, err := clientTest.GetBackend("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if b.Name != "test" {
		t.Errorf("Expected only test, %v found", b.Name)
	}
	if b.BindProcess != "all" {
		t.Errorf("%v: BindProcess not all: %v", b.Name, b.BindProcess)
	}
	if b.AdvCheck != "httpchk" {
		t.Errorf("%v: AdvCheck.Method not HEAD: %v", b.Name, b.AdvCheck)
	}
	if b.HttpchkParams.Method != "HEAD" {
		t.Errorf("%v: HttpchkParams.Method not HEAD: %v", b.Name, b.HttpchkParams.Method)
	}
	if b.HttpchkParams.URI != "/" {
		t.Errorf("%v: HttpchkParams.URI not HEAD: %v", b.Name, b.HttpchkParams.URI)
	}
	if b.Mode != "http" {
		t.Errorf("%v: Mode not http: %v", b.Name, b.Mode)
	}
	if *b.Balance.Algorithm != "roundrobin" {
		t.Errorf("%v: Balance.Algorithm not roundrobin: %v", b.Name, b.Balance.Algorithm)
	}
	if b.HashType.Method != "consistent" {
		t.Errorf("%v: HashType.Method not consistent: %v", b.Name, b.HashType.Method)
	}
	if b.HashType.Function != "sdbm" {
		t.Errorf("%v: HashType.Function not sdbm: %v", b.Name, b.HashType.Function)
	}
	if b.HashType.Modifier != "avalanche" {
		t.Errorf("%v: HashType.Modifier not avalanche: %v", b.Name, b.HashType.Modifier)
	}
	if b.HTTPConnectionMode != "http-keep-alive" {
		t.Errorf("%v: HTTPConnectionMode not http-keep-alive: %v", b.Name, b.HTTPConnectionMode)
	}
	if b.HTTPRestrictReqHdrNames != "preserve" {
		t.Errorf("%v: HTTPRestrictReqHdrNames not delete: %v", b.Name, b.HTTPRestrictReqHdrNames)
	}
	if *b.Forwardfor.Enabled != "enabled" {
		t.Errorf("%v: Forwardfor not enabled: %v", b.Name, b.Forwardfor)
	}
	if *b.DefaultServer.Fall != 2000 {
		t.Errorf("%v: DefaultServer.Fall not 2000: %v", b.Name, *b.DefaultServer)
	}
	if *b.DefaultServer.Rise != 4000 {
		t.Errorf("%v: DefaultServer.Rise not 4000: %v", b.Name, *b.DefaultServer.Rise)
	}
	if *b.DefaultServer.Inter != 5000 {
		t.Errorf("%v: DefaultServer.Inter not 5000: %v", b.Name, *b.DefaultServer.Inter)
	}
	if *b.DefaultServer.HealthCheckPort != 8888 {
		t.Errorf("%v: DefaultServer.HealthCheckPort not 8888: %v", b.Name, *b.DefaultServer.HealthCheckPort)
	}
	if *b.Cookie.Name != "BLA" {
		t.Errorf("%v: HTTPCookie Name not BLA: %v", b.Name, b.Cookie)
	}
	if b.Cookie.Type != "rewrite" {
		t.Errorf("%v: HTTPCookie Type not rewrite %v", b.Name, b.Cookie.Type)
	}
	if !b.Cookie.Httponly {
		t.Errorf("%v: HTTPCookie Httponly not true %v", b.Name, b.Cookie.Httponly)
	}
	if !b.Cookie.Nocache {
		t.Errorf("%v: HTTPCookie Nocache not false %v", b.Name, b.Cookie.Nocache)
	}
	if *b.CheckTimeout != 2000 {
		t.Errorf("%v: CheckTimeout not 2000: %v", b.Name, *b.CheckTimeout)
	}
	if *b.ServerTimeout != 3000 {
		t.Errorf("%v: ServerTimeout not 3000: %v", b.Name, *b.ServerTimeout)
	}
	if *b.ServerFinTimeout != 1000 {
		t.Errorf("%v: ServerFinTimeout not 1000: %v", b.Name, *b.ServerFinTimeout)
	}
	if *b.TarpitTimeout != 2000 {
		t.Errorf("%v: TarpitTimeout not 2000: %v", b.Name, *b.TarpitTimeout)
	}
	if b.AcceptInvalidHTTPResponse != "disabled" {
		t.Errorf("%v: AcceptInvalidHTTPResponse not disabled: %v", b.Name, b.AcceptInvalidHTTPResponse)
	}
	if b.H1CaseAdjustBogusServer != "disabled" {
		t.Errorf("%v: H1CaseAdjustBogusServer not disabled: %v", b.Name, b.H1CaseAdjustBogusServer)
	}
	if b.Tcpka != "enabled" {
		t.Errorf("%v: Tcpka not enabled: %v", b.Name, b.Tcpka)
	}
	if b.Srvtcpka != "enabled" {
		t.Errorf("%v: Srvtcpka not enabled: %v", b.Name, b.Srvtcpka)
	}
	if b.Compression == nil {
		t.Errorf("%v: Compression is nil", b.Name)
	} else {
		if len(b.Compression.Types) != 2 {
			t.Errorf("%v: len Compression.Types not 2: %v", b.Name, len(b.Compression.Types))
		} else {
			if !(b.Compression.Types[0] == "application/json" || b.Compression.Types[0] != "text/plain") {
				t.Errorf("%v: Compression.Types[0] wrong: %v", b.Name, b.Compression.Types[0])
			}
			if !(b.Compression.Types[1] != "application/json" || b.Compression.Types[0] != "text/plain") {
				t.Errorf("%v: Compression.Types[1] wrong: %v", b.Name, b.Compression.Types[1])
			}
		}
		if b.Compression.Direction != "both" {
			t.Errorf("%v: Compression.Direction not both: %v", b.Name, b.Compression.Direction)
		}
		if b.Compression.AlgoReq != "raw-deflate" {
			t.Errorf("%v: Compression.AlgoReq wrong: %v", b.Name, b.Compression.AlgoReq)
		}
		if len(b.Compression.AlgosRes) != 2 {
			t.Errorf("%v: len Compression.AlgosRes not 2: %v", b.Name, len(b.Compression.AlgosRes))
		} else {
			if !(b.Compression.AlgosRes[0] == "identity" || b.Compression.AlgosRes[0] == "raw-deflate") {
				t.Errorf("%v: Compression.AlgosRes[0] wrong: %v", b.Name, b.Compression.AlgosRes[0])
			}
			if !(b.Compression.AlgosRes[1] != "identity" || b.Compression.AlgosRes[0] == "raw-deflate") {
				t.Errorf("%v: Compression.AlgosRes[1] wrong: %v", b.Name, b.Compression.AlgosRes[1])
			}
		}
		if len(b.Compression.TypesReq) != 2 {
			t.Errorf("%v: len Compression.TypesReq not 2: %v", b.Name, len(b.Compression.TypesReq))
		} else {
			if !(b.Compression.TypesReq[0] == "text/plain" || b.Compression.TypesReq[0] == "application/json") {
				t.Errorf("%v: Compression.TypesReq[0] wrong: %v", b.Name, b.Compression.TypesReq[0])
			}
			if !(b.Compression.TypesReq[1] != "text/plain" || b.Compression.TypesReq[0] == "application/json") {
				t.Errorf("%v: Compression.TypesReq[1] wrong: %v", b.Name, b.Compression.TypesReq[1])
			}
		}
		if len(b.Compression.TypesRes) != 1 {
			t.Errorf("%v: len Compression.TypesRes not 1: %v", b.Name, len(b.Compression.TypesRes))
		} else {
			if !(b.Compression.TypesReq[0] == "text/plain") {
				t.Errorf("%v: Compression.TypesReq[0] wrong: %v", b.Name, b.Compression.TypesReq[0])
			}
		}
	}
	if b.Checkcache != "enabled" {
		t.Errorf("%v: Checkcache not enabled: %v", b.Name, b.Checkcache)
	}
	if b.IndependentStreams != "enabled" {
		t.Errorf("%v: IndependentStreams not enabled: %v", b.Name, b.IndependentStreams)
	}
	if b.Nolinger != "enabled" {
		t.Errorf("%v: Nolinger not enabled: %v", b.Name, b.Nolinger)
	}
	if b.Originalto == nil {
		t.Errorf("%v: Originalto is nil, expected not nil", b.Name)
	} else {
		if *b.Originalto.Enabled != "enabled" {
			t.Errorf("%v: Originalto.Enabled is not enabled: %v", b.Name, *b.Originalto.Enabled)
		}
		if b.Originalto.Except != "" {
			t.Errorf("%v: Originalto.Except is not empty: %v", b.Name, b.Originalto.Except)
		}
		if b.Originalto.Header != "X-Client-Dst" {
			t.Errorf("%v: Originalto.Header is not X-Client-Dst: %v", b.Name, b.Originalto.Header)
		}
	}
	if b.Persist != "enabled" {
		t.Errorf("%v: Persist not enabled: %v", b.Name, b.Persist)
	}
	if b.PreferLastServer != "enabled" {
		t.Errorf("%v: PreferLastServer not enabled: %v", b.Name, b.PreferLastServer)
	}
	if b.SpopCheck != "enabled" {
		t.Errorf("%v: SpopCheck not enabled: %v", b.Name, b.SpopCheck)
	}
	if b.TCPSmartConnect != "enabled" {
		t.Errorf("%v: TCPSmartConnect not enabled: %v", b.Name, b.TCPSmartConnect)
	}
	if b.Transparent != "enabled" {
		t.Errorf("%v: Transparent not enabled: %v", b.Name, b.Transparent)
	}
	if b.SpliceAuto != "enabled" {
		t.Errorf("%v: SpliceAuto not enabled: %v", b.Name, b.SpliceAuto)
	}
	if b.SpliceRequest != "enabled" {
		t.Errorf("%v: SpliceRequest not enabled: %v", b.Name, b.SpliceRequest)
	}
	if b.SpliceResponse != "enabled" {
		t.Errorf("%v: SpliceResponse not enabled: %v", b.Name, b.SpliceResponse)
	}
	if b.SrvtcpkaCnt == nil {
		t.Errorf("%v: SrvtcpkaCnt is nil", b.Name)
	} else if *b.SrvtcpkaCnt != 10 {
		t.Errorf("%v: SrvtcpkaCnt not 10: %v", b.Name, *b.SrvtcpkaCnt)
	}
	if b.SrvtcpkaIdle == nil {
		t.Errorf("%v: SrvtcpkaIdle is nil", b.Name)
	} else if *b.SrvtcpkaIdle != 10000 {
		t.Errorf("%v: SrvtcpkaIdle not 10000: %v", b.Name, *b.SrvtcpkaIdle)
	}
	if b.SrvtcpkaIntvl == nil {
		t.Errorf("%v: SrvtcpkaIntvl is nil", b.Name)
	} else if *b.SrvtcpkaIntvl != 10000 {
		t.Errorf("%v: SrvtcpkaIntvl not 10000: %v", b.Name, *b.SrvtcpkaIntvl)
	}
	if b.StatsOptions == nil {
		t.Errorf("%v: StatsOptions is nil", b.Name)
	}
	if b.StatsOptions.StatsShowModules != true {
		t.Error("StatsShowModules not set")
	}
	if b.StatsOptions.StatsRealm != true {
		t.Error("StatsRealm not set")
	}
	if b.StatsOptions.StatsRealmRealm == nil {
		t.Errorf("%v: StatsRealmRealm is nil", b.Name)
	} else if *b.StatsOptions.StatsRealmRealm != `HAProxy\\ Statistics` {
		t.Errorf("%v: StatsRealmRealm not 'HAProxy Statistics': %v", b.Name, *b.StatsOptions.StatsRealmRealm)
	}
	if len(b.StatsOptions.StatsAuths) != 2 {
		t.Errorf("%v: StatsAuths expected 2 instances got: %v", b.Name, len(b.StatsOptions.StatsAuths))
	}
	if b.StatsOptions.StatsAuths[0].User == nil {
		t.Errorf("%v: StatsAuths 0 User is nil", b.Name)
	} else if *b.StatsOptions.StatsAuths[0].User != "admin" {
		t.Errorf("%v: StatsAuths 0 User not admin: %v", b.Name, *b.StatsOptions.StatsAuths[0])
	}
	if b.StatsOptions.StatsAuths[0].Passwd == nil {
		t.Errorf("%v: StatsAuths 0 Passwd is nil", b.Name)
	} else if *b.StatsOptions.StatsAuths[0].Passwd != "AdMiN123" {
		t.Errorf("%v: StatsAuths 0 Passwd not AdMiN123: %v", b.Name, *b.StatsOptions.StatsAuths[0].Passwd)
	}
	if b.StatsOptions.StatsAuths[1].User == nil {
		t.Errorf("%v: StatsAuths 1 User is nil", b.Name)
	} else if *b.StatsOptions.StatsAuths[1].User != "admin2" {
		t.Errorf("%v: StatsAuths 1 User not admin2: %v", b.Name, *b.StatsOptions.StatsAuths[1].User)
	}
	if b.StatsOptions.StatsAuths[1].Passwd == nil {
		t.Errorf("%v: StatsAuths 1 Passwd is nil", b.Name)
	} else if *b.StatsOptions.StatsAuths[1].Passwd != "AdMiN1234" {
		t.Errorf("%v: StatsAuths 1 Passwd not AdMiN1234: %v", b.Name, *b.StatsOptions.StatsAuths[1].Passwd)
	}
	if len(b.StatsOptions.StatsHTTPRequests) != 2 {
		t.Errorf("%v: StatsHTTPRequests expected 2 instances got: %v", b.Name, len(b.StatsOptions.StatsHTTPRequests))
	}
	if b.StatsOptions.StatsHTTPRequests[0].Type == nil {
		t.Errorf("%v: StatsHTTPRequests 0 Type is nil", b.Name)
	} else if *b.StatsOptions.StatsHTTPRequests[0].Type != "auth" {
		t.Errorf("%v: StatsHTTPRequests 0 Type not auth: %v", b.Name, *b.StatsOptions.StatsHTTPRequests[0].Type)
	}
	if b.StatsOptions.StatsHTTPRequests[0].Realm != `HAProxy\\ Statistics` {
		t.Errorf("%v: StatsHTTPRequests 0 Realm not 'HAProxy Statistics': %v", b.Name, b.StatsOptions.StatsHTTPRequests[0].Realm)
	}
	if b.StatsOptions.StatsHTTPRequests[1].Type == nil {
		t.Errorf("%v: StatsHTTPRequests 1 Type is nil", b.Name)
	} else if *b.StatsOptions.StatsHTTPRequests[1].Type != "allow" {
		t.Errorf("%v: StatsHTTPRequests 1 Type not allow: %v", b.Name, *b.StatsOptions.StatsHTTPRequests[1].Type)
	}
	if b.StatsOptions.StatsHTTPRequests[1].Cond != "if" {
		t.Errorf("%v: StatsHTTPRequests 1 Cond not if: %v", b.Name, b.StatsOptions.StatsHTTPRequests[1].Cond)
	}
	if b.StatsOptions.StatsHTTPRequests[1].CondTest != "something" {
		t.Errorf("%v: StatsHTTPRequests 1 CondTest not something: %v", b.Name, b.StatsOptions.StatsHTTPRequests[1].CondTest)
	}
	if b.LoadServerStateFromFile != "local" {
		t.Errorf("%v: LoadServerStateFromFile not global: %v", b.Name, b.LoadServerStateFromFile)
	}
	if b.ServerStateFileName != "use-backend-name" {
		t.Errorf("%v: ServerStateFileName not: %v", b.Name, b.ServerStateFileName)
	}
	if b.Description != "this is a backend description" {
		t.Errorf("%v: Description not `this is a backend description`: %v", b.Name, b.Description)
	}
	if b.UseFCGIApp != "app-name" {
		t.Errorf("%v: UseFcgiApp not app-name: %v", b.Name, b.UseFCGIApp)
	}
	if !b.Enabled {
		t.Errorf("%v: Enabled not enabled", b.Name)
	}
	if *b.ID != 456 {
		t.Errorf("ID not 456: %v", *b.ID)
	}
	if *b.Errorloc302.Code != 404 {
		t.Errorf("%v: Errorloc302 Code not 404: %v", b.Name, *b.Errorloc302.Code)
	}
	if *b.Errorloc302.URL != "http://www.myawesomesite.com/not_found" {
		t.Errorf("%v: Errorloc302 Code not http://www.myawesomesite.com/not_found: %v", b.Name, *b.Errorloc302.URL)
	}
	if *b.Errorloc303.Code != 404 {
		t.Errorf("%v: Errorloc302 Code not 404: %v", b.Name, *b.Errorloc303.Code)
	}
	if *b.Errorloc303.URL != "http://www.myawesomesite.com/not_found" {
		t.Errorf("%v: Errorloc302 Code not http://www.myawesomesite.com/not_found: %v", b.Name, *b.Errorloc303.URL)
	}
	if len(b.ErrorFiles) != 3 {
		t.Errorf("ErrorFiles not 3: %v", len(b.ErrorFiles))
	} else {
		for _, ef := range b.ErrorFiles {
			if ef.Code == 403 {
				if ef.File != "/test/403.html" {
					t.Errorf("File for %v not 403: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 500 {
				if ef.File != "/test/500.html" {
					t.Errorf("File for %v not 500: %v", ef.Code, ef.File)
				}
			}
			if ef.Code == 429 {
				if ef.File != "/test/429.html" {
					t.Errorf("File for %v not 429: %v", ef.Code, ef.File)
				}
			}
		}
	}
	if len(b.ErrorFilesFromHTTPErrors) != 3 {
		t.Errorf("ErrorFilesFromHTTPErrors not 3: %v", len(b.ErrorFilesFromHTTPErrors))
	} else {
		for _, ef := range b.ErrorFilesFromHTTPErrors {
			if ef.Name == "my_errors" {
				if len(ef.Codes) != 3 {
					t.Errorf("ErrorFilesFromHTTPErrors Codes not 3 for errofiles name %s: %v", ef.Name, len(ef.Codes))
				}
				if ef.Codes[0] != 404 {
					t.Errorf("Codes[0] not 404 for errofiles name %s: %v", ef.Name, ef.Codes[0])
				}
				if ef.Codes[1] != 401 {
					t.Errorf("Codes[1] not 401 for errofiles name %s: %v", ef.Name, ef.Codes[1])
				}
				if ef.Codes[2] != 500 {
					t.Errorf("Codes[2] not 500 for errofiles name %s: %v", ef.Name, ef.Codes[2])
				}
			}
			if ef.Name == "other_errors" {
				if len(ef.Codes) != 0 {
					t.Errorf("ErrorFilesFromHTTPErrors Codes not 0 for errofiles name %s: %v", ef.Name, len(ef.Codes))
				}
			}
			if ef.Name == "another_errors" {
				if len(ef.Codes) != 1 {
					t.Errorf("ErrorFilesFromHTTPErrors Codes not 1 for errofiles name %s: %v", ef.Name, len(ef.Codes))
				}
				if ef.Codes[0] != 501 {
					t.Errorf("Codes[0] not 501 for errofiles name %s: %v", ef.Name, ef.Codes[0])
				}
			}
		}
	}
	if b.EmailAlert == nil {
		t.Error("EmailAlert is nil")
	} else if *b.EmailAlert.From != "prod01@example.com" {
		t.Errorf("EmailAlert.From is not prod01@example.com: %v", *b.EmailAlert.From)
	} else if *b.EmailAlert.To != "sre@example.com" {
		t.Errorf("EmailAlert.To is not sre@example.com: %v", *b.EmailAlert.To)
	} else if b.EmailAlert.Level != "warning" {
		t.Errorf("EmailAlert.Level is not warning: %v", b.EmailAlert.Level)
	} else if b.EmailAlert.Myhostname != "prod01" {
		t.Errorf("EmailAlert.Myhostname is not prod01: %v", b.EmailAlert.Myhostname)
	} else if *b.EmailAlert.Mailers != "localmailer1" {
		t.Errorf("EmailAlert.Mailers is not localmailer1: %v", *b.EmailAlert.Mailers)
	}

	if *b.Fullconn != 11 {
		t.Errorf("Fullconn not 11: %v", *b.Fullconn)
	}
	if *b.MaxKeepAliveQueue != 101 {
		t.Errorf("MaxKeepAliveQueue not 101: %v", *b.MaxKeepAliveQueue)
	}

	if *b.IgnorePersist.Cond != "if" {
		t.Errorf("IgnorePersist Cond not if: %v", *b.IgnorePersist.Cond)
	}
	if *b.IgnorePersist.CondTest != "acl-name" {
		t.Errorf("IgnorePersist CondTest not acl-name: %v", *b.IgnorePersist.CondTest)
	}

	if *b.ForcePersist.Cond != "unless" {
		t.Errorf("ForcePersist Cond not if: %v", *b.ForcePersist.Cond)
	}
	if *b.ForcePersist.CondTest != "acl-name-2" {
		t.Errorf("ForcePersist CondTest not acl-name-2: %v", *b.ForcePersist.CondTest)
	}
	if b.RetryOn != "504 505" {
		t.Errorf("RetryOn CondTest not 504 505: %v", b.RetryOn)
	}
	if *b.HTTPSendNameHeader != "X-My-Awesome-Header" {
		t.Errorf("%v: HTTPSendNameHeader not X-My-Awesome-Header: %v", b.Name, *b.HTTPSendNameHeader)
	}
	if b.PersistRule.RdpCookieName != "name" {
		t.Errorf("%v: PersistRule.RdpCookieName is not name: %v", b.Name, b.PersistRule.RdpCookieName)
	}
	if *b.Source.Address != "192.168.1.222" {
		t.Errorf("Source Address not 192.168.1.222 %v", *b.Source.Address)
	}
	if b.Source.Port != 0 {
		t.Errorf("Source Port not 0: %v", b.Source.Port)
	}
	if b.Source.Usesrc != "hdr_ip" {
		t.Errorf("Source Usesrc not hdr_ip: %v", b.Source.Usesrc)
	}
	if b.Source.AddressSecond != "" {
		t.Errorf("Source Address not empty: %v", b.Source.AddressSecond)
	}
	if b.Source.PortSecond != 0 {
		t.Errorf("Source PortSecond not 0: %v", b.Source.PortSecond)
	}
	if b.Source.Hdr != "hdr" {
		t.Errorf("Source Hdr not hdr: %v", b.Source.Hdr)
	}
	if b.Source.Occ != "occ" {
		t.Errorf("Source Occ not occ: %v", b.Source.Occ)
	}
	if b.Source.Interface != "" {
		t.Errorf("Source Interface not empty: %v", b.Source.Interface)
	}

	_, err = b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetBackend("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent bck")
	}
}

func TestCreateEditDeleteBackend(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)
	cookieName := "BLA"
	balanceAlgorithm := "uri"
	srvtcpkaCnt := int64(10)
	srvtcpkaTimeout := int64(10000)
	statsRealm := "Haproxy Stats"
	b := &models.Backend{
		Name: "created",
		Mode: "http",
		Balance: &models.Balance{
			Algorithm: &balanceAlgorithm,
			URILen:    100,
			URIDepth:  250,
		},
		BindProcess: "4",
		Cookie: &models.Cookie{
			Domains: []*models.Domain{
				{Value: "dom1"},
				{Value: "dom2"},
			},
			Attrs:    []*models.Attr{},
			Dynamic:  true,
			Httponly: true,
			Indirect: true,
			Maxidle:  5,
			Maxlife:  20,
			Name:     &cookieName,
			Nocache:  true,
			Postonly: true,
			Preserve: false,
			Secure:   false,
			Type:     "prefix",
		},
		HashType: &models.HashType{
			Method:   "map-based",
			Function: "crc32",
		},
		DefaultServer: &models.DefaultServer{
			ServerParams: models.ServerParams{
				Fall:  &tOut,
				Inter: &tOut,
			},
		},
		HTTPConnectionMode:   "http-keep-alive",
		HTTPKeepAlive:        "enabled",
		ConnectTimeout:       &tOut,
		ExternalCheck:        "enabled",
		ExternalCheckCommand: "/bin/false",
		ExternalCheckPath:    "/bin",
		Allbackups:           "enabled",
		AdvCheck:             "smtpchk",
		SmtpchkParams: &models.SmtpchkParams{
			Hello:  "HELO",
			Domain: "example.com",
		},
		AcceptInvalidHTTPResponse: "enabled",
		Compression: &models.Compression{
			Offload:   true,
			Direction: "both",
		},
		LogHealthChecks:    "enabled",
		Checkcache:         "enabled",
		IndependentStreams: "enabled",
		Nolinger:           "enabled",
		Originalto: &models.Originalto{
			Enabled: misc.StringP("enabled"),
			Except:  "127.0.0.1",
			Header:  "X-Client-Dst",
		},
		Persist:          "enabled",
		PreferLastServer: "enabled",
		SpopCheck:        "enabled",
		TCPSmartConnect:  "enabled",
		Transparent:      "enabled",
		SpliceAuto:       "enabled",
		SpliceRequest:    "enabled",
		SpliceResponse:   "enabled",
		SrvtcpkaCnt:      &srvtcpkaCnt,
		SrvtcpkaIdle:     &srvtcpkaTimeout,
		SrvtcpkaIntvl:    &srvtcpkaTimeout,
		StatsOptions: &models.StatsOptions{
			StatsShowModules: true,
			StatsRealm:       true,
			StatsRealmRealm:  &statsRealm,
			StatsAuths: []*models.StatsAuth{
				{User: misc.StringP("user1"), Passwd: misc.StringP("pwd1")},
				{User: misc.StringP("user2"), Passwd: misc.StringP("pwd2")},
			},
			StatsHTTPRequests: []*models.StatsHTTPRequest{
				{Type: misc.StringP("allow"), Cond: "if", CondTest: "something"},
				{Type: misc.StringP("auth"), Realm: "haproxy\\ stats"},
			},
		},
		EmailAlert: &models.EmailAlert{
			From:    misc.StringP("prod01@example.com"),
			To:      misc.StringP("sre@example.com"),
			Level:   "warning",
			Mailers: misc.StringP("localmailer1"),
		},
		ErrorFilesFromHTTPErrors: []*models.Errorfiles{
			{Name: "test_errors", Codes: []int64{400}},
			{Name: "test_errors_all"},
		},
		Disabled: true,
	}

	err := clientTest.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("created", "")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSONB []byte
	givenJSONB, err = b.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSONB []byte
	ondiskJSONB, err = backend.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created backend: %v\n", string(ondiskJSONB))
		fmt.Printf("Given backend: %v\n", string(givenJSONB))
		t.Error("Created backend not equal to given backend")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version++
	}
	// TestEditBackend
	tOut = int64(3)
	e := int64(1200000)
	kl := int64(128)
	s := int64(25600)
	backends := []*models.Backend{
		{
			From: "test_defaults",
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
				URILen:    10,
				URIDepth:  25,
			},
			BindProcess: "3",
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom1"},
					{Value: "dom3"},
				},
				Dynamic:  true,
				Httponly: true,
				Indirect: false,
				Maxidle:  150,
				Maxlife:  100,
				Name:     &cookieName,
				Nocache:  false,
				Postonly: false,
				Preserve: true,
				Secure:   true,
				Type:     "rewrite",
			},
			// Use non deprecated option only
			HTTPConnectionMode: "httpclose",
			ConnectTimeout:     &tOut,
			StickTable: &models.ConfigStickTable{
				Expire: &e,
				Keylen: &kl,
				Size:   &s,
				Store:  "gpc0,http_req_rate(40s)",
				Type:   "string",
			},
			AdvCheck: "mysql-check",
			MysqlCheckParams: &models.MysqlCheckParams{
				Username:      "user",
				ClientVersion: "pre-41",
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
			},
		},
		{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
			},
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom1"},
					{Value: "dom2"},
				},
				Name: &cookieName,
			},
			ConnectTimeout: &tOut,
			StickTable:     &models.ConfigStickTable{},
			AdvCheck:       "pgsql-check",
			PgsqlCheckParams: &models.PgsqlCheckParams{
				Username: "user",
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			// Use deprecated option only
			Httpclose: "enabled",
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Header:  "X-Client-Dst",
			},
		},
		{
			Name: "created",
			Mode: "http",
			Balance: &models.Balance{
				Algorithm: &balanceAlgorithm,
			},
			Cookie: &models.Cookie{
				Domains: []*models.Domain{
					{Value: "dom4"},
					{Value: "dom5"},
				},
				Name: &cookieName,
			},
			ConnectTimeout: &tOut,
			StickTable:     &models.ConfigStickTable{},
			AdvCheck:       "httpchk",
			HttpchkParams: &models.HttpchkParams{
				Method: "HEAD",
				URI:    "/",
			},
			HTTPCheck: &models.HTTPCheck{
				Type:    "send",
				Method:  "OPTIONS",
				URI:     "/",
				Version: "HTTP/1.1",
				Index:   misc.Int64P(0),
			},
			Checkcache:         "disabled",
			IndependentStreams: "disabled",
			Nolinger:           "disabled",
			Persist:            "disabled",
			PreferLastServer:   "disabled",
			SpopCheck:          "disabled",
			TCPSmartConnect:    "disabled",
			Transparent:        "disabled",
			SpliceAuto:         "disabled",
			SpliceRequest:      "disabled",
			SpliceResponse:     "disabled",
			SrvtcpkaCnt:        &srvtcpkaCnt,
			SrvtcpkaIdle:       &srvtcpkaTimeout,
			SrvtcpkaIntvl:      &srvtcpkaTimeout,
			StatsOptions: &models.StatsOptions{
				StatsShowModules: true,
				StatsRealm:       true,
				StatsRealmRealm:  &statsRealm,
				StatsAuths: []*models.StatsAuth{
					{User: misc.StringP("new_user1"), Passwd: misc.StringP("new_pwd1")},
					{User: misc.StringP("new_user2"), Passwd: misc.StringP("new_pwd2")},
				},
				StatsHTTPRequests: []*models.StatsHTTPRequest{
					{Type: misc.StringP("allow"), Cond: "if", CondTest: "something_else"},
					{Type: misc.StringP("auth"), Realm: "haproxy\\ stats2"},
				},
			},
			EmailAlert: &models.EmailAlert{
				From:    misc.StringP("prod01@example.com"),
				To:      misc.StringP("sre@example.com"),
				Level:   "warning",
				Mailers: misc.StringP("localmailer1"),
			},
			Originalto: &models.Originalto{
				Enabled: misc.StringP("enabled"),
				Except:  "127.0.0.1",
				Header:  "X-Client-Dst",
			},
		},
	}

	for i, backend := range backends {
		if errB := testBackendUpdate(backend, t); errB != nil {
			t.Errorf("failed update for backend %d: %v", i, errB)
		}
	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("created", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteBackend("created", "", 999999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}

	_, _, err = clientTest.GetBackend("created", "")
	if err == nil {
		t.Error("DeleteBackend failed, bck test still exists")
	}

	err = clientTest.DeleteBackend("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent bck")
		version++
	}
}

func testBackendUpdate(b *models.Backend, t *testing.T) error {
	err := clientTest.EditBackend("created", b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("created", "")
	if err != nil {
		return err
	}

	if !compareBackends(backend, b, t) {
		fmt.Printf("Edited bck: %+v\n", backend)
		fmt.Printf("Given bck: %+v\n", b)
		return fmt.Errorf("edited backend not equal to given backend")
	}

	if v != version {
		return fmt.Errorf("version %v returned, expected %v", v, version)
	}
	return nil
}

func compareBackends(x, y *models.Backend, t *testing.T) bool { //nolint:gocognit,gocyclo
	if *x.Balance.Algorithm != *y.Balance.Algorithm {
		return false
	}
	if x.Balance.HdrName != y.Balance.HdrName {
		return false
	}
	if x.Balance.HdrUseDomainOnly != y.Balance.HdrUseDomainOnly {
		return false
	}
	if x.Balance.RandomDraws != y.Balance.RandomDraws {
		return false
	}
	if x.Balance.RdpCookieName != y.Balance.RdpCookieName {
		return false
	}
	if x.Balance.URIDepth != y.Balance.URIDepth {
		return false
	}
	if x.Balance.URILen != y.Balance.URILen {
		return false
	}
	if x.Balance.URIWhole != y.Balance.URIWhole {
		return false
	}
	if x.Balance.URLParam != y.Balance.URLParam {
		return false
	}
	if x.Balance.URLParamCheckPost != y.Balance.URLParamCheckPost {
		return false
	}
	if x.Balance.URLParamMaxWait != y.Balance.URLParamMaxWait {
		return false
	}

	x.Balance = nil
	y.Balance = nil

	if *x.Cookie.Name != *y.Cookie.Name {
		return false
	}
	if len(x.Cookie.Domains) != len(y.Cookie.Domains) {
		return false
	}
	if x.Cookie.Domains[0].Value != y.Cookie.Domains[0].Value {
		return false
	}
	if x.Cookie.Dynamic != y.Cookie.Dynamic {
		return false
	}
	if x.Cookie.Httponly != y.Cookie.Httponly {
		return false
	}
	if x.Cookie.Indirect != y.Cookie.Indirect {
		return false
	}
	if x.Cookie.Maxidle != y.Cookie.Maxidle {
		return false
	}
	if x.Cookie.Maxlife != y.Cookie.Maxlife {
		return false
	}
	if x.Cookie.Nocache != y.Cookie.Nocache {
		return false
	}
	if x.Cookie.Postonly != y.Cookie.Postonly {
		return false
	}
	if x.Cookie.Preserve != y.Cookie.Preserve {
		return false
	}
	if x.Cookie.Secure != y.Cookie.Secure {
		return false
	}
	if x.Cookie.Type != y.Cookie.Type {
		return false
	}

	x.Cookie = nil
	y.Cookie = nil

	if x.BindProcess != y.BindProcess {
		return false
	}

	if !reflect.DeepEqual(x.DefaultServer, y.DefaultServer) {
		return false
	}

	x.DefaultServer = nil
	y.DefaultServer = nil

	if !reflect.DeepEqual(x.HttpchkParams, y.HttpchkParams) {
		return false
	}

	x.HttpchkParams = nil
	y.HttpchkParams = nil

	if !cmp.Equal(x.HTTPCheck, y.HTTPCheck, cmpopts.EquateEmpty()) {
		t.Errorf("Diff in HTTPChecK %s", cmp.Diff(x.HTTPCheck, y.HTTPCheck, cmpopts.EquateEmpty()))
		return false
	}

	x.HTTPCheck = nil
	y.HTTPCheck = nil

	if !reflect.DeepEqual(x.StickTable, y.StickTable) {
		return false
	}

	x.StickTable = nil
	y.StickTable = nil

	if !reflect.DeepEqual(x.Redispatch, y.Redispatch) {
		return false
	}

	x.Redispatch = nil
	y.Redispatch = nil

	if !reflect.DeepEqual(x.Forwardfor, y.Forwardfor) {
		return false
	}

	x.Forwardfor = nil
	y.Forwardfor = nil

	if !reflect.DeepEqual(x.SmtpchkParams, y.SmtpchkParams) {
		return false
	}

	x.SmtpchkParams = nil
	y.SmtpchkParams = nil

	if !reflect.DeepEqual(x.MysqlCheckParams, y.MysqlCheckParams) {
		return false
	}

	x.MysqlCheckParams = nil
	y.MysqlCheckParams = nil

	if !reflect.DeepEqual(x.PgsqlCheckParams, y.PgsqlCheckParams) {
		return false
	}

	x.PgsqlCheckParams = nil
	y.PgsqlCheckParams = nil

	if !reflect.DeepEqual(x.Originalto, y.Originalto) {
		return false
	}

	x.Originalto = nil
	y.Originalto = nil

	// Due to deprecated fields Httpclose, HTTPKeepAlive, HTTPServerClose
	// in favor of HTTPConnectionMode
	// If HTTPConnectionMode is set in original backend
	// - Httpclose, HTTPKeepAlive, HTTPServerClose will be set in updated backend, even if not present in original backend
	// If HTTPConnectionMode is unset in original backend:
	// - it will be set in updated backend
	switch y.HTTPConnectionMode {
	case "http-keep-alive":
		if x.HTTPKeepAlive != "enabled" {
			return false
		}
		x.HTTPKeepAlive = ""
	case "http-server-close":
		if x.HTTPServerClose != "enabled" {
			return false
		}
		x.HTTPServerClose = ""
	case "httpclose":
		if x.Httpclose != "enabled" {
			return false
		}
		x.Httpclose = ""
	case "":
		x.HTTPConnectionMode = ""
	}

	return reflect.DeepEqual(x, y)
}

func TestCreateEditDeleteBackendHTTPConnectionMode(t *testing.T) {
	// TestCreateBackend
	tOut := int64(5)

	// Backend with HTTPConnectionMode only
	b := &models.Backend{
		Name: "special-httpconnectionmode",
		Mode: "http",
		DefaultServer: &models.DefaultServer{
			ServerParams: models.ServerParams{
				Fall:  &tOut,
				Inter: &tOut,
			},
		},
		HTTPConnectionMode: "http-keep-alive",
	}

	err := clientTest.CreateBackend(b, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, backend, err := clientTest.GetBackend("special-httpconnectionmode", "")
	if err != nil {
		t.Error(err.Error())
	}

	if backend.HTTPConnectionMode != "http-keep-alive" {
		t.Errorf("Created backend is not correct for HTTPConnectionMode: %s", backend.HTTPConnectionMode)
	}
	if backend.HTTPKeepAlive != "enabled" {
		t.Errorf("Created backend is not correct for HTTPKeepAlive: %s", backend.HTTPConnectionMode)
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateBackend(b, "", version)
	if err == nil {
		t.Error("Should throw error bck already exists")
		version++
	}

	type testinput struct {
		backend                    *models.Backend
		expectedHTTPConnectionMode string
		expectedHTTPKeepAlive      string
		expectedHttpclose          string
		exptectedHTTPServerClose   string
	}
	// TestEditBackend
	inputs := []testinput{
		{
			// Update HTTPConnectionMode
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "httpclose",
			},
			expectedHTTPConnectionMode: "httpclose",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "enabled",
		},
		{
			// Use only deprecated option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPServerClose: "enabled",
			},
			expectedHTTPConnectionMode: "http-server-close",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "enabled",
			expectedHttpclose:          "",
		},
		{
			// Use both - Priority on HTTPConnection
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "http-keep-alive",
				HTTPServerClose:    "enabled",
			},
			expectedHTTPConnectionMode: "http-keep-alive",
			expectedHTTPKeepAlive:      "enabled",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "",
		},
		{
			// no option with deprecated option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPServerClose: "disabled",
			},
			expectedHTTPConnectionMode: "", // not possible to set no option with this field
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "disabled",
			expectedHttpclose:          "",
		},
		{
			// set back with HTTPConnectionMode
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "httpclose",
			},
			expectedHTTPConnectionMode: "httpclose",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "enabled",
		},
		{
			// remove option
			backend: &models.Backend{
				Name: "special-httpconnectionmode",
				Mode: "http",
				DefaultServer: &models.DefaultServer{
					ServerParams: models.ServerParams{
						Fall:  &tOut,
						Inter: &tOut,
					},
				},
				HTTPConnectionMode: "",
			},
			expectedHTTPConnectionMode: "",
			expectedHTTPKeepAlive:      "",
			exptectedHTTPServerClose:   "",
			expectedHttpclose:          "",
		},
	}

	for i, input := range inputs {
		err := clientTest.EditBackend("special-httpconnectionmode", input.backend, "", version)
		if err != nil {
			t.Error(err.Error())
		} else {
			version++
		}

		_, backend, err := clientTest.GetBackend("special-httpconnectionmode", "")
		if err != nil {
			t.Error(err.Error())
		}

		if backend.HTTPConnectionMode != input.expectedHTTPConnectionMode {
			t.Errorf("Updated backend %d is not correct for HTTPConnectionMode: %s", i, backend.HTTPConnectionMode)
		}
		if backend.HTTPKeepAlive != input.expectedHTTPKeepAlive {
			t.Errorf("Updated backend %d is not correct for HTTPKeepAlive: %s", i, backend.HTTPConnectionMode)
		}
		if backend.HTTPServerClose != input.exptectedHTTPServerClose {
			t.Errorf("Updated backend  %d is not correct for HTTPServerClose: %s", i, backend.HTTPServerClose)
		}
		if backend.Httpclose != input.expectedHttpclose {
			t.Errorf("Updated backend %d is not correct for Httpclose: %s", i, backend.Httpclose)
		}

	}

	// TestDeleteBackend
	err = clientTest.DeleteBackend("special-httpconnectionmode", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}
}
