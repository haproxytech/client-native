// Copyright 2022 HAProxy Technologies
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
	"testing"

	"github.com/haproxytech/client-native/v3/misc"
	"github.com/haproxytech/client-native/v3/models"
)

func TestGetHTTPChecks(t *testing.T) { //nolint:gocognit,gocyclo
	v, checks, err := client.GetHTTPChecks("backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(checks) != 14 {
		t.Errorf("%v http checks returned, expected 14", len(checks))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range checks {
		switch *r.Index {
		case 0:
			if r.Type != "connect" {
				t.Errorf("%v: Action not allow: %v", *r.Index, r.Type)
			}
		case 1:
			if r.Type != "send" {
				t.Errorf("%v: Action not send: %v", *r.Index, r.Type)
			}
			if r.Method != "GET" {
				t.Errorf("%v: Method not GET: %v", *r.Index, r.Method)
			}
			if r.URI != "/" {
				t.Errorf("%v: URI not /: %v", *r.Index, r.URI)
			}
			if r.Version != "HTTP/1.1" {
				t.Errorf("%v: Version not HTTP/1.1: %v", *r.Index, r.Version)
			}
			if len(r.CheckHeaders) < 1 {
				t.Errorf("%v: Not enough headers", *r.Index)
			} else {
				if *r.CheckHeaders[0].Name != "host" {
					t.Errorf("%v: Header Name not host: %v", *r.Index, r.CheckHeaders[0].Name)
				}
				if *r.CheckHeaders[0].Fmt != "haproxy.1wt.eu" {
					t.Errorf("%v: Header Fmt not haproxy.1wt.eu: %v", *r.Index, r.CheckHeaders[0].Fmt)
				}
			}
		case 2:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if r.Match != "status" {
				t.Errorf("%v: Match not status: %v", *r.Index, r.Match)
			}
			if r.Pattern != "200-399" {
				t.Errorf("%v: Pattern not 200-399: %v", *r.Index, r.Pattern)
			}
		case 3:
			if r.Type != "connect" {
				t.Errorf("%v: Action not connect: %v", *r.Index, r.Type)
			}
			if *r.Port != 443 {
				t.Errorf("%v: Port not 443: %v", *r.Index, *r.Port)
			}
			if !r.Ssl {
				t.Errorf("%v: Ssl not enabled", *r.Index)
			}
			if r.Sni != "haproxy.1wt.eu" {
				t.Errorf("%v: Sni not haproxy.1wt.eu: %v", *r.Index, r.Sni)
			}
		case 4:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if r.Match != "status" {
				t.Errorf("%v: Match not status: %v", *r.Index, r.Match)
			}
			if r.Pattern != "200,201,300-310" {
				t.Errorf("%v: Pattern not 200,201,300-310: %v", *r.Index, r.Pattern)
			}
		case 5:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if r.Match != "header" {
				t.Errorf("%v: Match not header: %v", *r.Index, r.Match)
			}
			if r.Pattern != "name \"set-cookie\" value -m beg \"sessid=\"" {
				t.Errorf("%v: Pattern not name \"set-cookie\" value -m beg \"sessid=\": %v", *r.Index, r.Pattern)
			}
		case 6:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if !r.ExclamationMark {
				t.Errorf("%v: ExclamationMark not set", *r.Index)
			}
			if r.Match != "string" {
				t.Errorf("%v: Match not string: %v", *r.Index, r.Match)
			}
			if r.Pattern != "SQL\\ Error" {
				t.Errorf("%v: Pattern not SQL\\ Error: %v", *r.Index, r.Pattern)
			}
		case 7:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if !r.ExclamationMark {
				t.Errorf("%v: ExclamationMark not set", *r.Index)
			}
			if r.Match != "rstatus" {
				t.Errorf("%v: Match not string: %v", *r.Index, r.Match)
			}
			if r.Pattern != "^5" {
				t.Errorf("%v: Pattern not ^5: %v", *r.Index, r.Pattern)
			}
		case 8:
			if r.Type != "expect" {
				t.Errorf("%v: Action not expect: %v", *r.Index, r.Type)
			}
			if r.Match != "rstring" {
				t.Errorf("%v: Match not rstring: %v", *r.Index, r.Match)
			}
			if r.Pattern != "<!--tag:[0-9a-f]*--></html>" {
				t.Errorf("%v: Pattern not <!--tag:[0-9a-f]*--></html>: %v", *r.Index, r.Pattern)
			}
		case 9:
			if r.Type != "unset-var" {
				t.Errorf("%v: Type not unset-var: %v", *r.Index, r.Type)
			}
			if r.VarScope != "check" {
				t.Errorf("%v: VarScope not check: %v", *r.Index, r.VarScope)
			}
			if r.VarName != "port" {
				t.Errorf("%v: VarName not port: %v", *r.Index, r.VarName)
			}
		case 10:
			if r.Type != "set-var" {
				t.Errorf("%v: Type not set-var: %v", *r.Index, r.Type)
			}
			if r.VarScope != "check" {
				t.Errorf("%v: VarScope not check: %v", *r.Index, r.VarScope)
			}
			if r.VarName != "port" {
				t.Errorf("%v: VarName not port: %v", *r.Index, r.VarName)
			}
			if r.VarExpr != "int(1234)" {
				t.Errorf("%v: VarExpr not int(1234): %v", *r.Index, r.VarExpr)
			}
		case 11:
			if r.Type != "set-var-fmt" {
				t.Errorf("%v: Type not set-var-fmt: %v", *r.Index, r.Type)
			}
			if r.VarScope != "check" {
				t.Errorf("%v: VarScope not check: %v", *r.Index, r.VarScope)
			}
			if r.VarName != "port" {
				t.Errorf("%v: VarName not port: %v", *r.Index, r.VarName)
			}
			if r.VarExpr != "int(1234)" {
				t.Errorf("%v: VarExpr not int(1234): %v", *r.Index, r.VarExpr)
			}
		case 12:
			if r.Type != "send-state" {
				t.Errorf("%v: Action not send-state: %v", *r.Index, r.Type)
			}
		case 13:
			if r.Type != "disable-on-404" {
				t.Errorf("%v: Action not disable-on-404: %v", *r.Index, r.Type)
			}
		default:
			t.Errorf("Expext only http checks 0 to 31, %v found", *r.Index)
		}
	}

	_, checks, err = client.GetHTTPChecks("defaults", "", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) > 2 {
		t.Errorf("%v HTTP Checks returned, expected 2", len(checks))
	}

	for _, r := range checks {
		switch *r.Index {
		case 0:
			if r.Type != "send-state" {
				t.Errorf("%v: Action not send-state: %v", *r.Index, r.Type)
			}
		case 1:
			if r.Type != "disable-on-404" {
				t.Errorf("%v: Action not disable-on-404: %v", *r.Index, r.Type)
			}
		default:
			t.Errorf("Expext only http-check 0 to %v, %v found", *r.Index, len(checks)-1)
		}
	}

	_, checks, err = client.GetHTTPChecks("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) != 1 {
		t.Errorf("%v HTTP Checks returned, expected 1", len(checks))
	}
}

func TestGetHTTPCheck(t *testing.T) {
	v, check, err := client.GetHTTPCheck(0, "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *check.Index != 0 {
		t.Errorf("HTTP Request Rule Index not 0, %v found", *check.Index)
	}
	if check.Type != "connect" {
		t.Errorf("%v: Action not allow: %v", *check.Index, check.Type)
	}

	_, err = check.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetHTTPCheck(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existant HTTP Request Rule")
	}

	_, check, err = client.GetHTTPCheck(0, "defaults", "", "")
	if err != nil {
		t.Error(err.Error())
	}
	if check.Type != "send-state" {
		t.Errorf("%v: Action not send-state: %v", *check.Index, check.Type)
	}
}

func TestCreateEditDeleteHTTPCheck(t *testing.T) {
	id := int64(1)

	// TestCreateHTTPCheck
	r := &models.HTTPCheck{
		Index:        &id,
		Type:         "send",
		Method:       "GET",
		Version:      "HTTP/1.1",
		URI:          "/",
		CheckHeaders: []*models.ReturnHeader{},
	}

	err := client.CreateHTTPCheck("backend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := client.GetHTTPCheck(1, "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	var givenJSONB []byte
	givenJSONB, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	var ondiskJSONB []byte
	ondiskJSONB, err = ondiskR.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created HTTP check: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP check: %v\n", string(givenJSONB))
		t.Error("Created HTTP check not equal to given HTTP check")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPRequestRule
	r = &models.HTTPCheck{
		Index:   &id,
		Type:    "send",
		Method:  "GET",
		Version: "HTTP/1.1",
		URI:     "/",
		CheckHeaders: []*models.ReturnHeader{
			{
				Name: misc.StringP("Host"),
				Fmt:  misc.StringP("google.com"),
			},
		},
	}

	err = client.EditHTTPCheck(1, "backend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = client.GetHTTPCheck(1, "backend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	givenJSONB, err = r.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	ondiskJSONB, err = ondiskR.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	if string(givenJSONB) != string(ondiskJSONB) {
		fmt.Printf("Created HTTP check: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP check: %v\n", string(givenJSONB))
		t.Error("Created HTTP check not equal to given HTTP check")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPRequest
	err = client.DeleteHTTPCheck(14, "backend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = client.GetHTTPCheck(14, "backend", "test", "")
	if err == nil {
		t.Error("DeleteHTTPCheck failed, HTTP check 13 still exists")
	}

	err = client.DeleteHTTPCheck(5, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existant HTTP Check")
		version++
	}
}
