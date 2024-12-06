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
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func quicInitialRuleExpectation() map[string]models.QUICInitialRules {
	initStructuredExpected()
	res := StructuredToQUICInitialRuleMap()
	return res
}

func TestGetQUICInitialRules(t *testing.T) {
	mrules := make(map[string]models.QUICInitialRules)
	v, hRules, err := clientTest.GetQUICInitialRules(configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	mrules["frontend/test"] = hRules

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	_, hRules, err = clientTest.GetQUICInitialRules(configuration.DefaultsParentName, "test_defaults", "")
	if err != nil {
		t.Error(err.Error())
	}
	mrules["defaults/test_defaults"] = hRules

	checkQUICInitialRules(t, mrules)
}

func checkQUICInitialRules(t *testing.T, got map[string]models.QUICInitialRules) {
	exp := quicInitialRuleExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		i := 0
		for _, w := range want {
			require.True(t, v[i].Equal(*w), "k=%s - diff %v", k, cmp.Diff(*v[i], *w))
			i++
		}
	}
}

func TestCreateQUICInitialRule(t *testing.T) {
	v, err := clientTest.GetVersion("")
	if err != nil {
		t.Error(err.Error())
	}

	tx, err := clientTest.StartTransaction(v)
	if err != nil {
		t.Error(err.Error())
	}

	har := &models.QUICInitialRule{
		Type: "reject",
	}
	if err = clientTest.CreateQUICInitialRule(0, configuration.FrontendParentName, "test", har, tx.ID, 0); err != nil {
		t.Error(err.Error())
	}

	_, found, err := clientTest.GetQUICInitialRule(0, configuration.FrontendParentName, "test", tx.ID)
	if err != nil {
		t.Error(err.Error())
	}

	if expected, got := har.Type, found.Type; expected != got {
		t.Error(fmt.Errorf("expected type %s, got %s", expected, got))
	}
}

func TestDeleteQUICInitialRule(t *testing.T) {
	v, err := clientTest.GetVersion("")
	if err != nil {
		t.Error(err.Error())
	}

	tx, err := clientTest.StartTransaction(v)
	if err != nil {
		t.Error(err.Error())
	}

	_, hRules, err := clientTest.GetQUICInitialRules(configuration.FrontendParentName, "test", tx.ID)
	if err != nil {
		t.Error(err.Error())
	}

	for range hRules {
		if err = clientTest.DeleteQUICInitialRule(int64(0), configuration.FrontendParentName, "test", tx.ID, 0); err != nil {
			t.Error(err.Error())
		}
	}
}
