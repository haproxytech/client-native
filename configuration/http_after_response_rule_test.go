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

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

func TestGetHTTPAfterResponseRules(t *testing.T) {
	v, hRules, err := clientTest.GetHTTPAfterResponseRules("frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if expected := 17; len(hRules) != expected {
		t.Errorf("%v http after response rules returned, expected %d", len(hRules), expected)
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, r := range hRules {
		switch *r.Index {
		case 0:
			if actionType := "set-map"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if mapFile := "map.lst"; r.MapFile != mapFile {
				t.Errorf("%v: MapFile not %s: %v", *r.Index, mapFile, r.MapFile)
			}
			if mapKeyFmt := "%[src]"; r.MapKeyfmt != mapKeyFmt {
				t.Errorf("%v: MapKeyfmt not %s: %v", *r.Index, mapKeyFmt, r.MapKeyfmt)
			}
			if mapValueFmt := "%[res.hdr(X-Value)]"; r.MapValuefmt != mapValueFmt {
				t.Errorf("%v: MapValuefmt not %s: %v", *r.Index, mapValueFmt, r.MapValuefmt)
			}
		case 1:
			if actionType := "del-map"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if mapFile := "map.lst"; r.MapFile != mapFile {
				t.Errorf("%v: MapFile not %s: %v", *r.Index, mapFile, r.MapFile)
			}
			if mapKeyFmt := "%[src]"; r.MapKeyfmt != mapKeyFmt {
				t.Errorf("%v: MapKeyfmt not %s: %v", *r.Index, mapKeyFmt, r.MapKeyfmt)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 2:
			if actionType := "del-acl"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if aclFile := "map.lst"; r.ACLFile != aclFile {
				t.Errorf("%v: ACLFile not %s: %v", *r.Index, aclFile, r.MapFile)
			}
			if aclKeyFmt := "%[src]"; r.ACLKeyfmt != aclKeyFmt {
				t.Errorf("%v: ACLKeyfmt not %s: %v", *r.Index, aclKeyFmt, r.MapKeyfmt)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 3:
			if actionType := "sc-add-gpc"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scIdx := int64(0); r.ScIdx != scIdx {
				t.Errorf("%v: ScIdx not %d: %v", *r.Index, scIdx, r.ScID)
			}
			if scId := int64(1); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 4:
			if actionType := "sc-inc-gpc"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scIdx := int64(0); r.ScIdx != scIdx {
				t.Errorf("%v: ScIdx not %d: %v", *r.Index, scIdx, r.ScID)
			}
			if scId := int64(1); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 5:
			if actionType := "sc-inc-gpc0"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scId := int64(0); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 6:
			if actionType := "sc-inc-gpc1"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scId := int64(0); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 7:
			if actionType := "sc-set-gpt0"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scId := int64(1); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if r.ScInt != nil {
				t.Errorf("%v: ScInt not nil: %v", *r.Index, *r.ScInt)
			}
			if expr := "hdr(Host),lower"; r.ScExpr != expr {
				t.Errorf("%v: ScExpr not %s: %v", *r.Index, expr, r.ScExpr)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 8:
			if r.Type != "sc-set-gpt0" {
				t.Errorf("%v: Type not sc-set-gpt0: %v", *r.Index, r.Type)
			}
			if scId := int64(1); r.ScID != scId {
				t.Errorf("%v: ScID not %d: %v", *r.Index, scId, r.ScID)
			}
			if scInt := int64(20); r.ScInt == nil || *r.ScInt != scInt {
				if r.ScInt == nil {
					t.Errorf("%v: ScInt is nil", *r.Index)
				} else {
					t.Errorf("%v: ScInt not %d: %v", *r.Index, scInt, *r.ScInt)
				}
			}
			if r.ScExpr != "" {
				t.Errorf("%v: ScExpr not empty string: %v", *r.Index, r.ScExpr)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 9:
			if actionType := "set-header"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if headerName := "Strict-Transport-Security"; r.HdrName != headerName {
				t.Errorf("%v: HdrName not %s: %v", *r.Index, headerName, r.HdrName)
			}
			if headerFmt := `"max-age=31536000"`; r.HdrFormat != headerFmt {
				t.Errorf("%v: HdrFmt not %s: %v", *r.Index, headerFmt, r.HdrName)
			}
		case 10:
			if actionType := "set-log-level"; r.Type != actionType {
				t.Errorf("%v: Action not %s: %v", *r.Index, actionType, r.Type)
			}
			if logLevel := "silent"; r.LogLevel != logLevel {
				t.Errorf("%v: LogLevel not %s %v", *r.Index, logLevel, r.LogLevel)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: Cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: CondTest not %s: %v", *r.Index, condTest, r.CondTest)
			}
		case 11:
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
		case 12:
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
		case 13:
			if actionType := "set-status"; r.Type != actionType {
				t.Errorf("%v: Type not allow: %v", *r.Index, actionType)
			}
			if status := int64(503); r.Status != status {
				t.Errorf("%v: status not %d: %v", *r.Index, status, r.HdrName)
			}
			if reason := fmt.Sprintf("%q", "SlowDown"); r.StatusReason != reason {
				t.Errorf("%v: reason not %s: %v", *r.Index, reason, r.HdrName)
			}
		case 14:
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
		case 15:
			if actionType := "unset-var"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if varName := "last_redir"; r.VarName != varName {
				t.Errorf("%v: VarName not %s: %v", *r.Index, varName, r.VarName)
			}
			if scope := "sess"; r.VarScope != scope {
				t.Errorf("%v: VarName not %s: %v", *r.Index, scope, r.VarScope)
			}
		case 16:
			if actionType := "sc-set-gpt"; r.Type != actionType {
				t.Errorf("%v: Type not %s: %v", *r.Index, actionType, r.Type)
			}
			if scID := 1; r.ScID != int64(scID) {
				t.Errorf("%v: sc-id not %d: %v", *r.Index, scID, r.ScID)
			}
			if scIdx := 2; r.ScIdx != int64(scIdx) {
				t.Errorf("%v: sc-idx not %d: %v", *r.Index, scIdx, r.ScIdx)
			}
			if cond := "if"; r.Cond != cond {
				t.Errorf("%v: cond not %s: %v", *r.Index, cond, r.Cond)
			}
			if condTest := "FALSE"; r.CondTest != condTest {
				t.Errorf("%v: condtest not %s: %v", *r.Index, condTest, r.CondTest)
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

	for range hRules {
		if err = clientTest.DeleteHTTPAfterResponseRule(int64(0), "frontend", "test", tx.ID, 0); err != nil {
			t.Error(err.Error())
		}
	}
}
