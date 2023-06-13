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

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetHTTPErrorRules(t *testing.T) { //nolint:gocognit,gocyclo
	v, checks, err := clientTest.GetHTTPErrorRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) != 1 {
		t.Errorf("%v http-error rules returned, expected 1", len(checks))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range checks {
		switch *r.Index {
		case 0:
			if r.Type != "status" {
				t.Errorf("%v: Action not status: %v", *r.Index, r.Type)
			}
			if r.Status != 400 {
				t.Errorf("%v: Status not 400: %v", *r.Index, r.Type)
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
			if len(r.ReturnHeaders) != 0 {
				t.Errorf(`%v: len(ReturnHeaders) not 0: %v`, *r.Index, len(r.ReturnHeaders))
			}
		}
	}

	_, checks, err = clientTest.GetHTTPErrorRules("defaults", "", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) != 2 {
		t.Errorf("%v http-error rules returned, expected 2", len(checks))
	}

	for _, r := range checks {
		switch *r.Index {
		case 0:
			if r.Type != "status" {
				t.Errorf("%v: Action not status: %v", *r.Index, r.Type)
			}
			if r.Status != 503 {
				t.Errorf("%v: Status not 503: %v", *r.Index, r.Type)
			}
			if *r.ReturnContentType != "\"application/json\"" {
				t.Errorf(`%v: ReturnContentType not \"application/json\": %v`, *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "file" {
				t.Errorf("%v: ReturnContentFormat not file: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != "/test/503" {
				t.Errorf("%v: ReturnContent not /test/503: %v", *r.Index, r.ReturnContent)
			}
			if len(r.ReturnHeaders) != 0 {
				t.Errorf("%v: len(ReturnHeaders) not 0: %v", *r.Index, len(r.ReturnHeaders))
			}
		case 1:
			if r.Type != "status" {
				t.Errorf("%v: Action not status: %v", *r.Index, r.Type)
			}
			if r.Status != 429 {
				t.Errorf("%v: Status not 429: %v", *r.Index, r.Type)
			}
			if *r.ReturnContentType != "application/json" {
				t.Errorf("%v: ReturnContentType not application/json: %v", *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "file" {
				t.Errorf("%v: ReturnContentFormat not file: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != "/test/429" {
				t.Errorf("%v: ReturnContent not /test/429: %v", *r.Index, r.ReturnContent)
			}
			if len(r.ReturnHeaders) != 0 {
				t.Errorf("%v: len(ReturnHeaders) not 0: %v", *r.Index, len(r.ReturnHeaders))
			}
		}
	}

	_, checks, err = clientTest.GetHTTPErrorRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(checks) != 2 {
		t.Errorf("%v http-error rules returned, expected 2", len(checks))
	}

	for _, r := range checks {
		switch *r.Index {
		case 0:
			if r.Type != "status" {
				t.Errorf("%v: Action not status: %v", *r.Index, r.Type)
			}
			if r.Status != 200 {
				t.Errorf("%v: Status not 200: %v", *r.Index, r.Type)
			}
			if *r.ReturnContentType != "\"text/plain\"" {
				t.Errorf(`%v: ReturnContentType not \"text/plain\": %v`, *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "string" {
				t.Errorf("%v: ReturnContentFormat not string: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != "\"My content\"" {
				t.Errorf(`%v: ReturnContent not "\"My content\" %v`, *r.Index, r.ReturnContent)
			}
			if len(r.ReturnHeaders) != 1 {
				t.Errorf("%v: len(ReturnHeaders) not 1: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[0].Name != "Some-Header" {
				t.Errorf("%v: ReturnHeaders[0] name not Some-Header: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[0].Fmt != "value" {
				t.Errorf("%v: ReturnHeaders[0] fmt not value: %v", *r.Index, len(r.ReturnHeaders))
			}
		case 1:
			if r.Type != "status" {
				t.Errorf("%v: Action not status: %v", *r.Index, r.Type)
			}
			if r.Status != 503 {
				t.Errorf("%v: Status not 503: %v", *r.Index, r.Type)
			}
			if *r.ReturnContentType != "application/json" {
				t.Errorf("%v: ReturnContentType not application/json: %v", *r.Index, *r.ReturnContentType)
			}
			if r.ReturnContentFormat != "string" {
				t.Errorf("%v: ReturnContentFormat not string: %v", *r.Index, r.ReturnContentFormat)
			}
			if r.ReturnContent != "\"My content\"" {
				t.Errorf(`%v: ReturnContent not "\"My content\" %v`, *r.Index, r.ReturnContent)
			}
			if len(r.ReturnHeaders) != 2 {
				t.Errorf("%v: len(ReturnHeaders) not 2: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[0].Name != "Additional-Header" {
				t.Errorf("%v: ReturnHeaders[0] name not Additional-Header: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[0].Fmt != "value1" {
				t.Errorf("%v: ReturnHeaders[0] fmt not value1: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[1].Name != "Some-Header" {
				t.Errorf("%v: ReturnHeaders[1] name not Some-Header: %v", *r.Index, len(r.ReturnHeaders))
			} else if *r.ReturnHeaders[1].Fmt != "value" {
				t.Errorf("%v: ReturnHeaders[1] fmt not value: %v", *r.Index, len(r.ReturnHeaders))
			}
		}
	}
}

func TestGetHTTPErrorRule(t *testing.T) {
	v, check, err := clientTest.GetHTTPErrorRule(0, "backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *check.Index != 0 {
		t.Errorf("http-error rule index not 0: %v", *check.Index)
	}
	if check.Type != "status" {
		t.Errorf("%v: Action not status: %v", *check.Index, check.Type)
	}
	if check.Status != 200 {
		t.Errorf("%v: Status not 200: %v", *check.Index, check.Type)
	}
	_, err = check.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetHTTPErrorRule(3, "backend", "test", "")
	if err == nil {
		t.Error("no http-error rules in backend section named test - expected an error")
	}

	_, check, err = clientTest.GetHTTPErrorRule(1, "defaults", "", "")
	if err != nil {
		t.Error(err.Error())
	}
	if *check.Index != 1 {
		t.Errorf("http-error rule index not 1: %v", *check.Index)
	}
	if check.Type != "status" {
		t.Errorf("%v: Action not status: %v", *check.Index, check.Type)
	}
	if check.Status != 429 {
		t.Errorf("%v: Status not 429: %v", *check.Index, check.Type)
	}
}

func TestCreateEditDeleteHTTPErrorRule(t *testing.T) {
	id := int64(1)
	r := &models.HTTPErrorRule{
		Index:               &id,
		Type:                "status",
		Status:              429,
		ReturnContentType:   misc.StringP("application/json"),
		ReturnContentFormat: "file",
		ReturnContent:       "/test/429",
		ReturnHeaders: []*models.ReturnHeader{
			{
				Name: misc.StringP("Some-Header"),
				Fmt:  misc.StringP("value"),
			},
		},
	}

	// TestCreateHTTPErrorRule
	err := clientTest.CreateHTTPErrorRule("backend", "test_2", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetHTTPErrorRule(1, "backend", "test_2", "")
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
		fmt.Printf("Created HTTP error rule: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP error rule: %v\n", string(givenJSONB))
		t.Error("Created HTTP error rule not equal to given HTTP error rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditHTTPErrorRule
	err = clientTest.EditHTTPErrorRule(1, "backend", "test_2", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetHTTPErrorRule(1, "backend", "test_2", "")
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
		fmt.Printf("Created HTTP error rule: %v\n", string(ondiskJSONB))
		fmt.Printf("Given HTTP error rule: %v\n", string(givenJSONB))
		t.Error("Created HTTP error rule not equal to given HTTP error rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteHTTPErrorRule
	err = clientTest.DeleteHTTPErrorRule(0, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPErrorRule(0, "frontend", "test", "")
	if err == nil {
		t.Error("deleting http-error rule failed - http-error rule 0 still exists")
	}

	err = clientTest.DeleteHTTPErrorRule(1, "defaults", "", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetHTTPErrorRule(1, "defaults", "", "")
	if err == nil {
		t.Error("deleting http-error rule failed - http-error rule 1 still exists")
	}

	err = clientTest.DeleteHTTPErrorRule(3, "backend", "test_2", "", version)
	if err == nil {
		t.Error("deleting http-error rule that does not exist - expected an error")
		version++
	}
}
