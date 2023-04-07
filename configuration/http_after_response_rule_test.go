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
	"testing"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

func TestGetHTTPAfterResponseRules(t *testing.T) {
	v, hRules, err := clientTest.GetHTTPAfterResponseRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if expected := 6; len(hRules) != expected {
		t.Errorf("%v http after response rules returned, expected %d", len(hRules), expected)
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range hRules {
		switch *r.Index {
		case 0:
			if actionType := "set-header"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if headerName := "Strict-Transport-Security"; r.HdrName != headerName {
				t.Errorf("%v: HdrName not %s: %v", *r.Index, headerName, r.HdrName)
			}
			if headerFmt := `"max-age=31536000"`; r.HdrFormat != headerFmt {
				t.Errorf("%v: HdrFmt not %s: %v", *r.Index, headerFmt, r.HdrName)
			}
		case 1:
			if actionType := "replace-header"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if headerName := "Set-Cookie"; r.HdrName != headerName {
				t.Errorf("%v: HdrName not %s: %v", *r.Index, headerName, r.HdrName)
			}
			if headerMatch := "(C=[^;]*);(.*)"; r.HdrMatch != headerMatch {
				t.Errorf("%v: HdrMatch not %s: %v", *r.Index, headerMatch, r.HdrName)
			}
			if headerFmt := `\1;ip=%bi;\2`; r.HdrFormat != headerFmt {
				t.Errorf("%v: HdrFmt not %s: %v", *r.Index, headerFmt, r.HdrName)
			}
		case 2:
			if actionType := "replace-value"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if headerName := "Cache-control"; r.HdrName != headerName {
				t.Errorf("%v: HdrName not %s: %v", *r.Index, headerName, r.HdrName)
			}
			if headerMatch := "^public$"; r.HdrMatch != headerMatch {
				t.Errorf("%v: HdrMatch not %s: %v", *r.Index, headerMatch, r.HdrName)
			}
			if headerFmt := "private"; r.HdrFormat != headerFmt {
				t.Errorf("%v: HdrFmt not %s: %v", *r.Index, headerFmt, r.HdrName)
			}
		case 3:
			if actionType := "set-status"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if status := int64(503); r.Status != status {
				t.Errorf("%v: status not %d: %v", *r.Index, status, r.HdrName)
			}
			if reason := fmt.Sprintf("%q", "SlowDown"); r.StatusReason != reason {
				t.Errorf("%v: reason not %s: %v", *r.Index, reason, r.HdrName)
			}
		case 4:
			if actionType := "set-var"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if varName := "last_redir"; r.VarName != varName {
				t.Errorf("%v: VarName not %s: %v", *r.Index, varName, r.VarName)
			}
			if scope := "sess"; r.VarScope != scope {
				t.Errorf("%v: VarName not %s: %v", *r.Index, scope, r.VarScope)
			}
			if expr := "res.hdr(location)"; r.VarExpr != expr {
				t.Errorf("%v: VarExpr not %s: %v", *r.Index, expr, r.VarExpr)
			}
		case 5:
			if actionType := "unset-var"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if varName := "last_redir"; r.VarName != varName {
				t.Errorf("%v: VarName not %s: %v", *r.Index, varName, r.VarName)
			}
			if scope := "sess"; r.VarScope != scope {
				t.Errorf("%v: VarName not %s: %v", *r.Index, scope, r.VarScope)
			}
		}
	}

	_, hRules, err = clientTest.GetHTTPAfterResponseRules("backend", "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	if len(hRules) > 0 {
		t.Errorf("%v HTTP After Response Rules returned, expected 0", len(hRules))
	}
}

func TestCreateHTTPAfterResponseRule(t *testing.T) {
	v, err := clientTest.GetVersion("")
	if err != nil {
		t.Error(err.Error())
	}

	tx, err := clientTest.StartTransaction(v)
	if err != nil {
		t.Error(err.Error())
	}

	har := &models.HTTPAfterResponseRule{
		Index:      misc.Int64P(0),
		StrictMode: "on",
		Type:       "strict-mode",
	}
	if err = clientTest.CreateHTTPAfterResponseRule("backend", "test", har, tx.ID, 0); err != nil {
		t.Error(err.Error())
	}

	_, found, err := clientTest.GetHTTPAfterResponseRule(0, "backend", "test", tx.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected, got := har.Type, found.Type; expected != got {
		t.Error(fmt.Errorf("expected type %s, got %s", expected, got))
	}

	if expected, got := har.StrictMode, found.StrictMode; expected != got {
		t.Error(fmt.Errorf("expected strict-mode %s, got %s", expected, got))
	}
}

func TestDeleteHTTPAfterResponseRule(t *testing.T) {
	v, err := clientTest.GetVersion("")
	if err != nil {
		t.Error(err.Error())
	}

	tx, err := clientTest.StartTransaction(v)
	if err != nil {
		t.Error(err.Error())
	}

	_, hRules, err := clientTest.GetHTTPAfterResponseRules("frontend", "test", tx.ID)
	if err != nil {
		t.Error(err.Error())
	}

	for index := range hRules {
		if err = clientTest.DeleteHTTPResponseRule(int64(index), "frontend", "test", tx.ID, 0); err != nil {
			t.Error(err.Error())
		}
	}
}
