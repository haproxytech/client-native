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

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetHTTPResponseRules(t *testing.T) { //nolint:gocognit,gocyclo
	v, hRules, err := clientTest.GetHTTPResponseRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(hRules) != 32 {
		t.Errorf("%v http response rules returned, expected 32", len(hRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range hRules {
		switch *r.Index {
		case 0:
			if r.Type != "allow" {
				t.Errorf("%v: Type not allow: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "src 192.168.0.0/16" {
				t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", *r.Index, r.CondTest)
			}
		case 1:
			if r.Type != "set-header" {
				t.Errorf("%v: Type not set-header: %v", *r.Index, r.Type)
			}
			if r.HdrName != "X-SSL" {
				t.Errorf("%v: HdrName not X-SSL: %v", *r.Index, r.HdrName)
			}
			if r.HdrFormat != "%[ssl_fc]" {
				t.Errorf("%v: HdrValue not [ssl_fc]: %v", *r.Index, r.HdrFormat)
			}
		case 2:
			if r.Type != "set-var" {
				t.Errorf("%v: Type not set-var: %v", *r.Index, r.Type)
			}
			if r.VarName != "my_var" {
				t.Errorf("%v: VarName not my_var: %v", *r.Index, r.VarName)
			}
			if r.VarScope != "req" {
				t.Errorf("%v: VarName not req: %v", *r.Index, r.VarScope)
			}
			if r.VarExpr != "req.fhdr(user-agent),lower" {
				t.Errorf("%v: VarExpr not req.fhdr(user-agent),lower: %v", *r.Index, r.VarExpr)
			}
		case 3:
			if r.Type != "set-map" {
				t.Errorf("%v: Type not set-map: %v", *r.Index, r.Type)
			}
			if r.MapFile != "map.lst" {
				t.Errorf("%v: MapFile not map.lst: %v", *r.Index, r.MapFile)
			}
			if r.MapKeyfmt != "%[src]" {
				t.Errorf("%v: MapKeyfmt not %%[src]: %v", *r.Index, r.MapKeyfmt)
			}
			if r.MapValuefmt != "%[res.hdr(X-Value)]" {
				t.Errorf("%v: MapValuefmt not %%[res.hdr(X-Value)]: %v", *r.Index, r.MapValuefmt)
			}
		case 4:
			if r.Type != "del-map" {
				t.Errorf("%v: Type not del-map: %v", *r.Index, r.Type)
			}
			if r.MapFile != "map.lst" {
				t.Errorf("%v: MapFile not map.lst: %v", *r.Index, r.MapFile)
			}
			if r.MapKeyfmt != "%[src]" {
				t.Errorf("%v: MapKeyfmt not %%[src]: %v", *r.Index, r.MapKeyfmt)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 5:
			if r.Type != "del-acl" {
				t.Errorf("%v: Type not del-acl: %v", *r.Index, r.Type)
			}
			if r.ACLFile != "map.lst" {
				t.Errorf("%v: ACLFile not map.lst: %v", *r.Index, r.ACLFile)
			}
			if r.ACLKeyfmt != "%[src]" {
				t.Errorf("%v: ACLKeyfmt not %%[src]: %v", *r.Index, r.ACLKeyfmt)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 6:
			if r.Type != "cache-store" {
				t.Errorf("%v: Type not cache-store: %v", *r.Index, r.Type)
			}
			if r.CacheName != "cache-name" {
				t.Errorf("%v: CacheName not cache-name: %v", *r.Index, r.MapFile)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 7:
			if r.Type != "sc-add-gpc" {
				t.Errorf("%v: Type not sc-add-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != 0 {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScID)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 8:
			if r.Type != "sc-inc-gpc" {
				t.Errorf("%v: Type not sc-inc-gpc: %v", *r.Index, r.Type)
			}
			if r.ScIdx != 0 {
				t.Errorf("%v: ScIdx not 0: %v", *r.Index, r.ScID)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 9:
			if r.Type != "sc-inc-gpc0" {
				t.Errorf("%v: Type not sc-inc-gpc0: %v", *r.Index, r.Type)
			}
			if r.ScID != 0 {
				t.Errorf("%v: ScID not 0: %v", *r.Index, r.ScID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 10:
			if r.Type != "sc-inc-gpc1" {
				t.Errorf("%v: Type not sc-inc-gpc1: %v", *r.Index, r.Type)
			}
			if r.ScID != 0 {
				t.Errorf("%v: ScID not 0: %v", *r.Index, r.ScID)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 11:
			if r.Type != "sc-set-gpt0" {
				t.Errorf("%v: Type not sc-set-gpt0: %v", *r.Index, r.Type)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScID)
			}
			if r.ScInt != nil {
				t.Errorf("%v: ScInt not nil: %v", *r.Index, *r.ScInt)
			}
			if r.ScExpr != "hdr(Host),lower" {
				t.Errorf("%v: ScExpr not hdr(Host),lower: %v", *r.Index, r.ScExpr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 12:
			if r.Type != "sc-set-gpt0" {
				t.Errorf("%v: Type not sc-set-gpt0: %v", *r.Index, r.Type)
			}
			if r.ScID != 1 {
				t.Errorf("%v: ScID not 1: %v", *r.Index, r.ScID)
			}
			if r.ScInt == nil || *r.ScInt != 20 {
				if r.ScInt == nil {
					t.Errorf("%v: ScInt is nil", *r.Index)
				} else {
					t.Errorf("%v: ScInt not 20: %v", *r.Index, *r.ScInt)
				}
			}
			if r.ScExpr != "" {
				t.Errorf("%v: ScExpr not empty string: %v", *r.Index, r.ScExpr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 13:
			if r.Type != "set-mark" {
				t.Errorf("%v: Type not set-mark: %v", *r.Index, r.Type)
			}
			if r.MarkValue != "20" {
				t.Errorf("%v: MarkValue not 20: %v", *r.Index, r.MarkValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 14:
			if r.Type != "set-nice" {
				t.Errorf("%v: Type not set-nice: %v", *r.Index, r.Type)
			}
			if r.NiceValue != 20 {
				t.Errorf("%v: NiceValue not 20: %v", *r.Index, r.NiceValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 15:
			if r.Type != "set-tos" {
				t.Errorf("%v: Type not set-tos: %v", *r.Index, r.Type)
			}
			if r.TosValue != "0" {
				t.Errorf("%v: TosValue not 0: %v", *r.Index, r.TosValue)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 16:
			if r.Type != "silent-drop" {
				t.Errorf("%v: Type not silent-drop: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 17:
			if r.Type != "unset-var" {
				t.Errorf("%v: Type not unset-var: %v", *r.Index, r.Type)
			}
			if r.VarName != "my_var" {
				t.Errorf("%v: VarName not my_var: %v", *r.Index, r.VarName)
			}
			if r.VarScope != "req" {
				t.Errorf("%v: VarScope not req: %v", *r.Index, r.VarScope)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 18:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 0 {
				t.Errorf("%v: TrackScStickCounter not 0: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackSc0Key)
			}
			if r.TrackScTable != "tr0" {
				t.Errorf("%v: TrackScTable not tr0: %v", *r.Index, r.TrackSc0Table)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 19:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 1 {
				t.Errorf("%v: TrackScStickCounter not 1: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackSc1Key)
			}
			if r.TrackScTable != "tr1" {
				t.Errorf("%v: TrackScTable not tr1: %v", *r.Index, r.TrackSc1Table)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 20:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 2 {
				t.Errorf("%v: TrackScStickCounter not 2: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackSc2Key)
			}
			if r.TrackScTable != "tr2" {
				t.Errorf("%v: TrackScTable not tr2: %v", *r.Index, r.TrackSc2Table)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 21:
			if r.Type != "strict-mode" {
				t.Errorf("%v: Type not strict-mode: %v", *r.Index, r.Type)
			}
			if r.StrictMode != "on" {
				t.Errorf("%v: StrictMode not on: %v", *r.Index, r.StrictMode)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 22:
			if r.Type != "lua" {
				t.Errorf("%v: Type not lua: %v", *r.Index, r.Type)
			}
			if r.LuaAction != "foo" {
				t.Errorf("%v: LuaAction not foo: %v", *r.Index, r.LuaAction)
			}
			if r.LuaParams != "param1 param2" {
				t.Errorf("%v: LuaParams not 'param1 param2': %v", *r.Index, r.LuaParams)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 23:
			if r.Type != "deny" {
				t.Errorf("%v: Type not deny: %v", *r.Index, r.Type)
			}
			if *r.DenyStatus != 400 {
				t.Errorf("%v: DenyStatus not 400: %v", *r.Index, *r.DenyStatus)
			}
			if *r.ReturnContentType != "application/json" {
				t.Errorf("%v: ReturnContentType not application/json: %v", *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "lf-file" {
				t.Errorf("%v: ReturnContentFormat not lf-file: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != "/var/errors.file" {
				t.Errorf(`%v: ReturnContent not "/var/errors.file": %v`, *r.Index, r.ReturnContent)
			}
		case 24:
			if r.Type != "wait-for-body" {
				t.Errorf("%v: Type not wait-for-body: %v", *r.Index, r.Type)
			}
			if *r.WaitTime != 20000 {
				t.Errorf("%v: WaitTime not 20000: %v", *r.Index, *r.WaitTime)
			}
			if *r.WaitAtLeast != 102400 {
				t.Errorf("%v: AtLeast not 102400: %v", *r.Index, *r.WaitAtLeast)
			}
		case 25:
			if r.Type != "set-bandwidth-limit" {
				t.Errorf("%v: Type not set-bandwidth-limit: %v", *r.Index, r.Type)
			}
			if r.BandwidthLimitName != "my-limit" {
				t.Errorf("%v: BandwidthLimitName not my-limit: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "1m" {
				t.Errorf("%v: BandwidthLimitLimit not 1m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "10s" {
				t.Errorf("%v: BandwidthLimitPeriod not 10s: %v", *r.Index, r.BandwidthLimitPeriod)
			}
		case 26:
			if r.Type != "set-bandwidth-limit" {
				t.Errorf("%v: Type not set-bandwidth-limit: %v", *r.Index, r.Type)
			}
			if r.BandwidthLimitName != "my-limit-reverse" {
				t.Errorf("%v: BandwidthLimitName not my-limit-reverse: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "2m" {
				t.Errorf("%v: BandwidthLimitLimit not 2m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "20s" {
				t.Errorf("%v: BandwidthLimitPeriod no 20s: %v", *r.Index, r.BandwidthLimitPeriod)
			}
		case 27:
			if r.Type != "set-bandwidth-limit" {
				t.Errorf("%v: Type not set-bandwidth-limit: %v", *r.Index, r.Type)
			}
			if r.BandwidthLimitName != "my-limit-cond" {
				t.Errorf("%v: BandwidthLimitName not my-limit-cond: %v", *r.Index, r.BandwidthLimitName)
			}
			if r.BandwidthLimitLimit != "3m" {
				t.Errorf("%v: BandwidthLimitLimit not 3m: %v", *r.Index, r.BandwidthLimitLimit)
			}
			if r.BandwidthLimitPeriod != "" {
				t.Errorf("%v: BandwidthLimitPeriod not empty", *r.Index)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 28:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 0 {
				t.Errorf("%v: TrackScStickCounter not 0: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackScKey)
			}
			if r.TrackScTable != "tr0" {
				t.Errorf("%v: TrackScTable not tr0: %v", *r.Index, r.TrackScTable)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *r.Index, r.CondTest)
			}
		case 29:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 1 {
				t.Errorf("%v: TrackScStickCounter not 1: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackScKey)
			}
			if r.TrackScTable != "tr1" {
				t.Errorf("%v: TrackScTable not tr1: %v", *r.Index, r.TrackScTable)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *r.Index, r.CondTest)
			}
		case 30:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 2 {
				t.Errorf("%v: TrackScStickCounter not 2: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackScKey)
			}
			if r.TrackScTable != "tr2" {
				t.Errorf("%v: TrackScTable not tr2: %v", *r.Index, r.TrackScTable)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *r.Index, r.CondTest)
			}
		case 31:
			if r.Type != "track-sc" {
				t.Errorf("%v: Type not track-sc: %v", *r.Index, r.Type)
			}
			if *r.TrackScStickCounter != 5 {
				t.Errorf("%v: TrackScStickCounter not 5: %v", *r.Index, r.TrackScStickCounter)
			}
			if r.TrackScKey != "src" {
				t.Errorf("%v: TrackScKey not src: %v", *r.Index, r.TrackScKey)
			}
			if r.TrackScTable != "test" {
				t.Errorf("%v: TrackScTable not test: %v", *r.Index, r.TrackScTable)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "TRUE" {
				t.Errorf("%v: CondTest not TRUE: %v", *r.Index, r.CondTest)
			}
		default:
			t.Errorf("Expect only http-response 0 to 28, %v found", *r.Index)
		}
	}

	_, hRules, err = clientTest.GetHTTPResponseRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules) > 0 {
		t.Errorf("%v HTTP Response Rules returned, expected 0", len(hRules))
	}
}

func TestGetHTTPResponseRule(t *testing.T) {
	v, r, err := clientTest.GetHTTPResponseRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *r.Index != 0 {
		t.Errorf("HTTPResponse Rule ID not 0, %v found", *r.Index)
	}
	if r.Type != "allow" {
		t.Errorf("%v: Type not allow: %v", *r.Index, r.Type)
	}
	if r.Cond != "if" {
		t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
	}
	if r.CondTest != "src 192.168.0.0/16" {
		t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", *r.Index, r.CondTest)
	}

	_, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPResponseRule(3, "backend", "test2", "")
	if err == nil {
		t.Error("Should throw error, non existent HTTPResponse Rule")
	}

	_, r, err = clientTest.GetHTTPResponseRule(0, "frontend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if r.Type != "capture" {
		t.Errorf("%v: Type not 'capture': %v", *r.Index, r.Type)
	}
	if *r.CaptureID != 0 {
		t.Errorf("%v: Wrong slotID: %v", *r.Index, r.CaptureID)
	}
}

func TestCreateEditDeleteHTTPResponseRule(t *testing.T) {
	id := int64(1)
	// TestCreateHTTPResponseRule
	r := &models.HTTPResponseRule{
		Index:    &id,
		Type:     "set-log-level",
		LogLevel: "alert",
	}

	err := clientTest.CreateHTTPResponseRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPResponseRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created HTTP response rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Created HTTP response rule not equal to given HTTP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPResponseRule
	r = &models.HTTPResponseRule{
		Index:    &id,
		Type:     "set-log-level",
		LogLevel: "warning",
	}

	err = clientTest.EditHTTPResponseRule(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPResponseRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited HTTP response rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP response rule: %v\n", r)
		t.Error("Edited HTTP response rule not equal to given HTTP response rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPResponse
	err = clientTest.DeleteHTTPResponseRule(32, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPResponseRule(32, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPResponseRule failed, HTTPResponse Rule 32 still exists")
	}

	err = clientTest.DeleteHTTPResponseRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent HTTPResponse Rule")
		version++
	}
}

func TestSerializeHTTPResponseRule(t *testing.T) {
	testCases := []struct {
		input          models.HTTPResponseRule
		expectedResult string
	}{
		{
			input: models.HTTPResponseRule{
				Type:                models.HTTPResponseRuleTypeTrackDashSc,
				Cond:                "if",
				CondTest:            "TRUE",
				TrackScKey:          "src",
				TrackScTable:        "tr0",
				TrackScStickCounter: misc.Int64P(3),
			},
			expectedResult: "track-sc3 src table tr0 if TRUE",
		},
		{
			input: models.HTTPResponseRule{
				Type:          models.HTTPResponseRuleTypeTrackDashSc0,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc0Key:   "src",
				TrackSc0Table: "tr0",
			},
			expectedResult: "track-sc0 src table tr0 if TRUE",
		},
		{
			input: models.HTTPResponseRule{
				Type:          models.HTTPResponseRuleTypeTrackDashSc1,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc1Key:   "src",
				TrackSc1Table: "tr1",
			},
			expectedResult: "track-sc1 src table tr1 if TRUE",
		},
		{
			input: models.HTTPResponseRule{
				Type:          models.HTTPResponseRuleTypeTrackDashSc2,
				Cond:          "if",
				CondTest:      "TRUE",
				TrackSc2Key:   "src",
				TrackSc2Table: "tr2",
			},
			expectedResult: "track-sc2 src table tr2 if TRUE",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.expectedResult, func(t *testing.T) {
			tcpType, err := SerializeHTTPResponseRule(testCase.input)
			if err != nil {
				t.Error(err.Error())
			}

			actual := tcpType.String()
			if actual != testCase.expectedResult {
				t.Errorf("Expected %q, got: %q", testCase.expectedResult, actual)
			}
		})
	}
}
