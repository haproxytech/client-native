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

	"github.com/haproxytech/client-native/v2/models"
)

func TestGetHTTPRequestRules(t *testing.T) { //nolint:gocognit,gocyclo
	v, hRules, err := client.GetHTTPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(hRules) != 29 {
		t.Errorf("%v http request rules returned, expected 26", len(hRules))
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
				t.Errorf("%v: HdrFormat not [ssl_fc]: %v", *r.Index, r.HdrFormat)
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
				t.Errorf("%v: VarPattern not req.fhdr(user-agent),lower: %v", *r.Index, r.VarExpr)
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
			if r.MapValuefmt != "%[req.hdr(X-Value)]" {
				t.Errorf("%v: MapValuefmt not %%[req.hdr(X-Value)]: %v", *r.Index, r.MapValuefmt)
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
			if r.Type != "cache-use" {
				t.Errorf("%v: Type not cache-use: %v", *r.Index, r.Type)
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
		case 6:
			if r.Type != "disable-l7-retry" {
				t.Errorf("%v: Type not disable-l7-retry: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 7:
			if r.Type != "early-hint" {
				t.Errorf("%v: Type not early-hint: %v", *r.Index, r.Type)
			}
			if r.HintName != "hint-name" {
				t.Errorf("%v: HintName not hint-name: %v", *r.Index, r.MapFile)
			}
			if r.HintFormat != "%[src]" {
				t.Errorf("%v: HintFormat not %%[src]: %v", *r.Index, r.MapFile)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 8:
			if r.Type != "replace-uri" {
				t.Errorf("%v: Type not replace-uri: %v", *r.Index, r.Type)
			}
			if r.URIMatch != "^http://(.*)" {
				t.Errorf("%v: URIMatch not ^http://(.*): %v", *r.Index, r.MapFile)
			}
			if r.URIFmt != "https://1" {
				t.Errorf("%v: URIFmt not https://1: %v", *r.Index, r.MapKeyfmt)
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
			if r.Type != "do-resolve" {
				t.Errorf("%v: Type not do-resolve: %v", *r.Index, r.Type)
			}
			if r.VarName != "txn.myip" {
				t.Errorf("%v: VarName not txn.myip: %v", *r.Index, r.VarName)
			}
			if r.Resolvers != "mydns" {
				t.Errorf("%v: Resolvers not mydns: %v", *r.Index, r.Resolvers)
			}
			if r.Protocol != "ipv4" {
				t.Errorf("%v: Protocol not ipv4: %v", *r.Index, r.Protocol)
			}
			if r.Expr != "hdr(Host),lower" {
				t.Errorf("%v: Expr not hdr(Host),lower: %v", *r.Index, r.Expr)
			}
		case 12:
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
		case 13:
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
		case 14:
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
		case 15:
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
		case 16:
			if r.Type != "set-method" {
				t.Errorf("%v: Type not set-method: %v", *r.Index, r.Type)
			}
			if r.MethodFmt != "POST" {
				t.Errorf("%v: MethodFmt not 0: %v", *r.Index, r.MethodFmt)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 17:
			if r.Type != "set-priority-class" {
				t.Errorf("%v: Type not set-priority-class: %v", *r.Index, r.Type)
			}
			if r.Expr != "req.hdr(class)" {
				t.Errorf("%v: Expr not req.hdr(class): %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 18:
			if r.Type != "set-priority-offset" {
				t.Errorf("%v: Type not set-priority-offset: %v", *r.Index, r.Type)
			}
			if r.Expr != "req.hdr(offset)" {
				t.Errorf("%v: Expr not req.hdr(offset): %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 19:
			if r.Type != "set-src" {
				t.Errorf("%v: Type not set-src: %v", *r.Index, r.Type)
			}
			if r.Expr != "req.hdr(src)" {
				t.Errorf("%v: Expr not 0: %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 20:
			if r.Type != "set-src-port" {
				t.Errorf("%v: Type not set-src-port: %v", *r.Index, r.Type)
			}
			if r.Expr != "req.hdr(port)" {
				t.Errorf("%v: Expr not req.hdr(port): %v", *r.Index, r.Expr)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 21:
			if r.Type != "wait-for-handshake" {
				t.Errorf("%v: Type not wait-for-handshake: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 22:
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
		case 23:
			if r.Type != "silent-drop" {
				t.Errorf("%v: Type not silent-drop: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 24:
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
		case 25:
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
		case 26:
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
		case 27:
			if r.Type != "use-service" {
				t.Errorf("%v: Type not use-service: %v", *r.Index, r.Type)
			}
			if r.ServiceName != "svrs" {
				t.Errorf("%v: ServiceName not svrs: %v", *r.Index, r.ServiceName)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		case 28:
			if r.Type != "return" {
				t.Errorf("%v: Type not return: %v", *r.Index, r.Type)
			}
			if *r.ReturnStatusCode != 200 {
				t.Errorf("%v: ReturnStatusCode not 200: %v", *r.Index, *r.ReturnStatusCode)
			}
			if *r.ReturnContentType != `"text/plain"` {
				t.Errorf("%v: ReturnContentType not text/plain: %v", *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "string" {
				t.Errorf("%v: ReturnContentFormat not string: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != `"My content"` {
				t.Errorf(`%v: ReturnContent not "My content": %v`, *r.Index, r.ReturnContent)
			}
			if len(r.ReturnHeaders) != 1 {
				t.Errorf("%v: ReturnHeaders not length 1: %v", *r.Index, len(r.ReturnHeaders))
			}
			if *r.ReturnHeaders[0].Name != "Some-Header" {
				t.Errorf("%v: ReturnHeaders[0].Name not Some-Header: %v", *r.Index, *r.ReturnHeaders[0].Name)
			}
			if *r.ReturnHeaders[0].Fmt != "value" {
				t.Errorf("%v: ReturnHeaders[0].Fmt not value: %v", *r.Index, *r.ReturnHeaders[0].Fmt)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		default:
			t.Errorf("Expext only http-request 0 to 28, %v found", *r.Index)
		}
	}

	_, hRules, err = client.GetHTTPRequestRules("backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules) > 2 {
		t.Errorf("%v HTTP Request Ruless returned, expected 2", len(hRules))
	}

	for _, r := range hRules {
		switch *r.Index {
		case 0:
			if r.Type != "set-dst" {
				t.Errorf("%v: Type not set-dst: %v", *r.Index, r.Type)
			}
			if r.Expr != "hdr(x-dst)" {
				t.Errorf("%v: Expr not hdr(x-dst): %v", *r.Index, r.VarExpr)
			}
		case 1:
			if r.Type != "set-dst-port" {
				t.Errorf("%v: Type not set-dst-port: %v", *r.Index, r.Type)
			}
			if r.Expr != "int(4000)" {
				t.Errorf("%v: Expr not v: %v", *r.Index, r.Expr)
			}
		default:
			t.Errorf("Expext only http-request 0 to %v, %v found", *r.Index, len(hRules)-1)
		}
	}

	_, hRules, err = client.GetHTTPRequestRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules) > 0 {
		t.Errorf("%v HTTP Request Ruless returned, expected 0", len(hRules))
	}
}

func TestGetHTTPRequestRule(t *testing.T) {
	v, r, err := client.GetHTTPRequestRule(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *r.Index != 0 {
		t.Errorf("HTTP Request Rule Index not 0, %v found", *r.Index)
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

	_, _, err = client.GetHTTPRequestRule(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
	}

	_, r, err = client.GetHTTPRequestRule(0, "frontend", "test_2", "")
	if err != nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
	}
	if r.Type != "capture" {
		t.Errorf("%v: Type not allow: %v", *r.Index, r.Type)
	}
	if r.CaptureLen != 10 {
		t.Errorf("%v: Wrong len parameter for capture: %v", *r.Index, r.CaptureLen)
	}

	_, r, err = client.GetHTTPRequestRule(1, "frontend", "test_2", "")
	if err != nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
	}
	if *r.CaptureID != 0 {
		t.Errorf("%v: Wrong slotIndex: %v", *r.Index, *r.CaptureID)
	}
}

func TestCreateEditDeleteHTTPRequestRule(t *testing.T) {
	id := int64(1)

	// TestCreateHTTPRequestRule
	var redirCode int64 = 301
	r := &models.HTTPRequestRule{
		Index:      &id,
		Type:       "redirect",
		RedirCode:  &redirCode,
		RedirValue: "http://www.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "location",
	}

	err := client.CreateHTTPRequestRule("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := client.GetHTTPRequestRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created HTTP request rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Created HTTP request rule not equal to given HTTP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPRequestRule
	r = &models.HTTPRequestRule{
		Index:      &id,
		Type:       "redirect",
		RedirCode:  &redirCode,
		RedirValue: "http://www1.%[hdr(host)]%[capture.req.uri]",
		RedirType:  "scheme",
	}

	err = client.EditHTTPRequestRule(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = client.GetHTTPRequestRule(1, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited HTTP request rule: %v\n", ondiskR)
		fmt.Printf("Given HTTP request rule: %v\n", r)
		t.Error("Edited HTTP request rule not equal to given HTTP request rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPRequest
	err = client.DeleteHTTPRequestRule(29, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetHTTPRequestRule(29, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPRequestRule failed, HTTP Request Rule 29 still exists")
	}

	err = client.DeleteHTTPRequestRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
		version++
	}
}
