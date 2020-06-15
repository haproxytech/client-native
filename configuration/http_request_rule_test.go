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

	"github.com/haproxytech/models/v2"
)

func TestGetHTTPRequestRules(t *testing.T) {
	v, hRules, err := client.GetHTTPRequestRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(hRules) != 12 {
		t.Errorf("%v http request rules returned, expected 12", len(hRules))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range hRules {
		if *r.Index == 0 {
			if r.Type != "allow" {
				t.Errorf("%v: Type not allow: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "src 192.168.0.0/16" {
				t.Errorf("%v: CondTest not src 192.168.0.0/16: %v", *r.Index, r.CondTest)
			}
		} else if *r.Index == 1 {
			if r.Type != "set-header" {
				t.Errorf("%v: Type not set-header: %v", *r.Index, r.Type)
			}
			if r.HdrName != "X-SSL" {
				t.Errorf("%v: HdrName not X-SSL: %v", *r.Index, r.HdrName)
			}
			if r.HdrFormat != "%[ssl_fc]" {
				t.Errorf("%v: HdrFormat not [ssl_fc]: %v", *r.Index, r.HdrFormat)
			}
		} else if *r.Index == 2 {
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
		} else if *r.Index == 3 {
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
		} else if *r.Index == 4 {
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
		} else if *r.Index == 5 {
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
		} else if *r.Index == 6 {
			if r.Type != "disable-l7-retry" {
				t.Errorf("%v: Type not disable-l7-retry: %v", *r.Index, r.Type)
			}
			if r.Cond != "if" {
				t.Errorf("%v: Cond not if: %v", *r.Index, r.Cond)
			}
			if r.CondTest != "FALSE" {
				t.Errorf("%v: CondTest not FALSE: %v", *r.Index, r.CondTest)
			}
		} else if *r.Index == 7 {
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
		} else if *r.Index == 8 {
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
		} else if *r.Index == 9 {
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
		} else if *r.Index == 10 {
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
		} else if *r.Index == 11 {
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
		} else {
			t.Errorf("Expext only http-request 0 to 11, %v found", *r.Index)
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
		if *r.Index == 0 {
			if r.Type != "set-dst" {
				t.Errorf("%v: Type not set-dst: %v", *r.Index, r.Type)
			}
			if r.Expr != "hdr(x-dst)" {
				t.Errorf("%v: Expr not hdr(x-dst): %v", *r.Index, r.VarExpr)
			}
		} else if *r.Index == 1 {
			if r.Type != "set-dst-port" {
				t.Errorf("%v: Type not set-dst-port: %v", *r.Index, r.Type)
			}
			if r.Expr != "int(4000)" {
				t.Errorf("%v: Expr not v: %v", *r.Index, r.Expr)
			}
		} else {
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
	r := &models.HTTPRequestRule{
		Index:      &id,
		Type:       "redirect",
		RedirCode:  301,
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
		RedirCode:  302,
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
	err = client.DeleteHTTPRequestRule(12, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetHTTPRequestRule(12, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPRequestRule failed, HTTP Request Rule 12 still exists")
	}

	err = client.DeleteHTTPRequestRule(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
		version++
	}
}
