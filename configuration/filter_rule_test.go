// Copyright 2026 HAProxy Technologies
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
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func mustGetFilterRules(t *testing.T, c Configuration, parentType, parentName string) (int64, models.FilterRules) {
	t.Helper()
	v, rules, err := c.GetFilterRules(parentType, parentName, "")
	require.NoError(t, err)
	return v, rules
}

func filterSelectorText(selector *models.FilterSelector) string {
	if selector.ID == "" {
		return selector.Class
	}
	return selector.Class + "/" + selector.ID
}

func filterRuleSelectors(rules models.FilterRules) []string {
	selectors := make([]string, 0, len(rules))
	for _, rule := range rules {
		parts := make([]string, 0, len(rule.Filters))
		for _, selector := range rule.Filters {
			parts = append(parts, filterSelectorText(selector))
		}
		selectors = append(selectors, rule.Type+" "+strings.Join(parts, " "))
	}
	return selectors
}

func TestFilterRuleCRUD(t *testing.T) {
	parents := []struct {
		parentType string
		parentName string
	}{
		{FrontendParentName, "fe"},
		{BackendParentName, "be"},
		{DefaultsParentName, "named"},
	}
	for _, parent := range parents {
		t.Run(parent.parentType+"/"+parent.parentName, func(t *testing.T) {
			c, _ := newFilterTestClient(t)
			parentType, parentName := parent.parentType, parent.parentName

			first := &models.FilterRule{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "A"}}, Metadata: map[string]interface{}{"k": "v"}}
			v, _ := mustGetFilterRules(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterRule(0, parentType, parentName, first, "", v))

			second := &models.FilterRule{Type: "disable", Filters: []*models.FilterSelector{{Class: "trace"}}}
			v, _ = mustGetFilterRules(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterRule(1, parentType, parentName, second, "", v))

			middle := &models.FilterRule{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "B"}, {Class: "trace", ID: "A"}}}
			v, _ = mustGetFilterRules(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterRule(1, parentType, parentName, middle, "", v))

			v, rules := mustGetFilterRules(t, c, parentType, parentName)
			require.Equal(t, []string{"enable trace/A", "enable trace/B trace/A", "disable trace"}, filterRuleSelectors(rules))
			require.Equal(t, []*models.FilterSelector{{Class: "trace", ID: "B"}, {Class: "trace", ID: "A"}}, rules[1].Filters)

			_, got, err := c.GetFilterRule(1, parentType, parentName, "")
			require.NoError(t, err)
			require.Equal(t, "enable", got.Type)

			edited := &models.FilterRule{Type: "disable", Filters: []*models.FilterSelector{{Class: "cache-store", ID: "C"}}}
			require.NoError(t, c.EditFilterRule(1, parentType, parentName, edited, "", v))
			v, rules = mustGetFilterRules(t, c, parentType, parentName)
			require.Equal(t, []string{"enable trace/A", "disable cache-store/C", "disable trace"}, filterRuleSelectors(rules))

			require.NoError(t, c.DeleteFilterRule(0, parentType, parentName, "", v))
			v, rules = mustGetFilterRules(t, c, parentType, parentName)
			require.Equal(t, []string{"disable cache-store/C", "disable trace"}, filterRuleSelectors(rules))

			reversed := models.FilterRules{second, edited}
			require.NoError(t, c.ReplaceFilterRules(parentType, parentName, reversed, "", v))
			v, rules = mustGetFilterRules(t, c, parentType, parentName)
			require.Equal(t, []string{"disable trace", "disable cache-store/C"}, filterRuleSelectors(rules))

			require.NoError(t, c.ReplaceFilterRules(parentType, parentName, models.FilterRules{}, "", v))
			v, rules = mustGetFilterRules(t, c, parentType, parentName)
			require.Empty(t, rules)

			require.NoError(t, c.CreateFilterRule(0, parentType, parentName, first, "", v))
			_, got, err = c.GetFilterRule(0, parentType, parentName, "")
			require.NoError(t, err)
			require.Equal(t, map[string]interface{}{"k": "v"}, got.Metadata)

			v, _ = mustGetFilterRules(t, c, parentType, parentName)
			_, _, err = c.GetFilterRule(99, parentType, parentName, "")
			require.Error(t, err)
			require.Error(t, c.EditFilterRule(99, parentType, parentName, edited, "", v))
			require.Error(t, c.DeleteFilterRule(99, parentType, parentName, "", v))
			require.Error(t, c.CreateFilterRule(0, parentType, parentName, first, "", v+100))
		})
	}

	t.Run("missing parent", func(t *testing.T) {
		c, _ := newFilterTestClient(t)
		_, _, err := c.GetFilterRules(FrontendParentName, "doesnotexist", "")
		require.Error(t, err)
	})
}

func TestFilterRuleMixedOrder(t *testing.T) {
	c, path := newFilterTestClient(t)

	lines := []*models.FilterRule{
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "A"}}},
		{Type: "disable", Filters: []*models.FilterSelector{{Class: "trace"}}},
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "B"}, {Class: "trace", ID: "A"}}},
	}
	v, _ := mustGetFilterRules(t, c, FrontendParentName, "fe")
	for i, rule := range lines {
		require.NoError(t, c.CreateFilterRule(int64(i), FrontendParentName, "fe", rule, "", v))
		v, _ = mustGetFilterRules(t, c, FrontendParentName, "fe")
	}

	// Edit the middle operation without collapsing its neighbours.
	edited := &models.FilterRule{Type: "disable", Filters: []*models.FilterSelector{{Class: "http-comp-req"}}}
	require.NoError(t, c.EditFilterRule(1, FrontendParentName, "fe", edited, "", v))
	v, rules := mustGetFilterRules(t, c, FrontendParentName, "fe")
	require.Equal(t, []string{"enable trace/A", "disable http-comp-req", "enable trace/B trace/A"}, filterRuleSelectors(rules))

	// Reload the persisted configuration in a fresh client.
	reloaded, err := New(
		context.Background(),
		options.ConfigurationFile(path),
		options.HAProxyBin("echo"),
		options.UseModelsValidation,
		options.TransactionsDir(t.TempDir()),
	)
	require.NoError(t, err)
	_, rules = mustGetFilterRules(t, reloaded, FrontendParentName, "fe")
	require.Equal(t, []string{"enable trace/A", "disable http-comp-req", "enable trace/B trace/A"}, filterRuleSelectors(rules))

	// Whole-list replacement keeps the reversed order, duplicates included.
	reversed := models.FilterRules{
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "B"}, {Class: "trace", ID: "A"}}},
		{Type: "disable", Filters: []*models.FilterSelector{{Class: "http-comp-req"}}},
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "A"}}},
	}
	v, _ = mustGetFilterRules(t, reloaded, FrontendParentName, "fe")
	require.NoError(t, reloaded.ReplaceFilterRules(FrontendParentName, "fe", reversed, "", v))
	_, rules = mustGetFilterRules(t, reloaded, FrontendParentName, "fe")
	require.Equal(t, []string{"enable trace/B trace/A", "disable http-comp-req", "enable trace/A"}, filterRuleSelectors(rules))
}
