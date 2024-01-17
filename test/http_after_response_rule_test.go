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

package test

import (
	_ "embed"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func hTTPAfterResponseRuleExpectation() map[string]models.HTTPAfterResponseRules {
	initStructuredExpected()
	res := StructuredToHTTPAfterResponseRuleMap()
	return res
}

func TestGetHTTPAfterResponseRules(t *testing.T) {
	mrules := make(map[string]models.HTTPAfterResponseRules)
	v, hRules, err := clientTest.GetHTTPAfterResponseRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	mrules["frontend/test"] = hRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, hRules, err = clientTest.GetHTTPAfterResponseRules(configuration.BackendParentName, "test_2", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["backend/test_2"] = hRules

	checkHTTPAfterResponseRules(t, mrules)
}

func checkHTTPAfterResponseRules(t *testing.T, got map[string]models.HTTPAfterResponseRules) {
	exp := hTTPAfterResponseRuleExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if *g.Index == *w.Index {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
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
	if err = clientTest.CreateHTTPAfterResponseRule(configuration.BackendParentName, "test", har, tx.ID, 0); err != nil {
		t.Error(err.Error())
	}

	_, found, err := clientTest.GetHTTPAfterResponseRule(0, configuration.BackendParentName, "test", tx.ID)
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

	_, hRules, err := clientTest.GetHTTPAfterResponseRules(configuration.FrontendParentName, "test", tx.ID)
	if err != nil {
		t.Error(err.Error())
	}

	for range hRules {
		if err = clientTest.DeleteHTTPAfterResponseRule(int64(0), configuration.FrontendParentName, "test", tx.ID, 0); err != nil {
			t.Error(err.Error())
		}
	}
}
