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

package test

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/models"
)

func sequenceSelector(class, id string) *models.FilterSelector {
	return &models.FilterSelector{Class: class, ID: id}
}

func filterInstanceLists() (models.FilterConfigs, models.FilterRules, models.FilterSequences) {
	configs := models.FilterConfigs{
		{Class: "trace", ID: "A", Enabled: true, TraceQuiet: true, Metadata: map[string]interface{}{"k": "v"}},
		{Class: "http-comp-req"},
	}
	rules := models.FilterRules{
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "A"}}},
		{Type: "disable", Filters: []*models.FilterSelector{{Class: "trace"}}},
		{Type: "enable", Filters: []*models.FilterSelector{{Class: "trace", ID: "B"}, {Class: "trace", ID: "A"}}},
	}
	sequences := models.FilterSequences{
		{
			Side: "request",
			Rules: []*models.FilterSequenceRule{
				{Instance: sequenceSelector("trace", "A"), Direction: "before", Target: sequenceSelector("trace", "B")},
			},
		},
	}
	return configs, rules, sequences
}

func TestStructuredFilterInstances(t *testing.T) {
	c, _, err := getTestClient()
	require.NoError(t, err)

	configs, rules, sequences := filterInstanceLists()
	frontend := &models.Frontend{
		FrontendBase:       models.FrontendBase{Name: "filtered", Mode: "http"},
		FilterList:         models.Filters{{Type: "comp-req"}},
		FilterConfigList:   configs,
		FilterRuleList:     rules,
		FilterSequenceList: sequences,
	}
	v, _, err := c.GetStructuredFrontends("")
	require.NoError(t, err)
	require.NoError(t, c.CreateStructuredFrontend(frontend, "", v))

	v, got, err := c.GetStructuredFrontend("filtered", "")
	require.NoError(t, err)
	require.Equal(t, configs, got.FilterConfigList)
	require.Equal(t, rules, got.FilterRuleList)
	require.Equal(t, sequences, got.FilterSequenceList)
	require.Equal(t, frontend.FilterList, got.FilterList)

	// Clearing the new lists must not touch the legacy filter list.
	got.FilterConfigList = nil
	got.FilterRuleList = nil
	got.FilterSequenceList = nil
	require.NoError(t, c.EditStructuredFrontend("filtered", got, "", v))
	v, got, err = c.GetStructuredFrontend("filtered", "")
	require.NoError(t, err)
	require.Empty(t, got.FilterConfigList)
	require.Empty(t, got.FilterRuleList)
	require.Empty(t, got.FilterSequenceList)
	require.Equal(t, frontend.FilterList, got.FilterList)

	configs, rules, sequences = filterInstanceLists()
	backend := &models.Backend{
		BackendBase:        models.BackendBase{Name: "filtered_be", Mode: "http"},
		FilterConfigList:   configs,
		FilterRuleList:     rules,
		FilterSequenceList: sequences,
	}
	v, _, err = c.GetStructuredBackends("")
	require.NoError(t, err)
	require.NoError(t, c.CreateStructuredBackend(backend, "", v))
	_, gotBackend, err := c.GetStructuredBackend("filtered_be", "")
	require.NoError(t, err)
	require.Equal(t, configs, gotBackend.FilterConfigList)
	require.Equal(t, rules, gotBackend.FilterRuleList)
	require.Equal(t, sequences, gotBackend.FilterSequenceList)

	defaults := &models.Defaults{
		DefaultsBase:       models.DefaultsBase{Name: "filtered_defaults", Mode: "http"},
		FilterConfigList:   configs,
		FilterRuleList:     rules,
		FilterSequenceList: sequences,
	}
	v, _, err = c.GetStructuredDefaultsSections("")
	require.NoError(t, err)
	require.NoError(t, c.CreateStructuredDefaultsSection(defaults, "", v))
	_, gotDefaults, err := c.GetStructuredDefaultsSection("filtered_defaults", "")
	require.NoError(t, err)
	require.Equal(t, configs, gotDefaults.FilterConfigList)
	require.Equal(t, rules, gotDefaults.FilterRuleList)
	require.Equal(t, sequences, gotDefaults.FilterSequenceList)

	// The named defaults section stays independent from the frontend.
	_, got, err = c.GetStructuredFrontend("filtered", "")
	require.NoError(t, err)
	require.Empty(t, got.FilterConfigList)
}

func TestStructuredFilterSequenceCyclePreflight(t *testing.T) {
	c, _, err := getTestClient()
	require.NoError(t, err)

	_, _, sequences := filterInstanceLists()
	frontend := &models.Frontend{
		FrontendBase:       models.FrontendBase{Name: "cycle_preflight", Mode: "http"},
		FilterSequenceList: sequences,
	}
	v, _, err := c.GetStructuredFrontends("")
	require.NoError(t, err)
	require.NoError(t, c.CreateStructuredFrontend(frontend, "", v))

	v, got, err := c.GetStructuredFrontend("cycle_preflight", "")
	require.NoError(t, err)
	got.FilterSequenceList = models.FilterSequences{
		{
			Side: "request",
			Rules: []*models.FilterSequenceRule{
				{Instance: sequenceSelector("trace", "A"), Direction: "before", Target: sequenceSelector("trace", "B")},
				{Instance: sequenceSelector("trace", "B"), Direction: "before", Target: sequenceSelector("trace", "A")},
			},
		},
	}
	require.Error(t, c.EditStructuredFrontend("cycle_preflight", got, "", v))

	_, got, err = c.GetStructuredFrontend("cycle_preflight", "")
	require.NoError(t, err)
	require.Equal(t, sequences, got.FilterSequenceList)

	_, _, err = c.GetStructuredDefaultsSections("")
	require.NoError(t, err)
	defaults := &models.Defaults{
		DefaultsBase: models.DefaultsBase{Name: "cycle_defaults", Mode: "http"},
		FilterSequenceList: models.FilterSequences{
			{
				Side: "request",
				Rules: []*models.FilterSequenceRule{
					{Instance: sequenceSelector("trace", "X"), Direction: "before", Target: sequenceSelector("trace", "Y")},
					{Instance: sequenceSelector("trace", "Y"), Direction: "before", Target: sequenceSelector("trace", "X")},
				},
			},
		},
	}
	v, _, err = c.GetStructuredDefaultsSections("")
	require.NoError(t, err)
	require.Error(t, c.CreateStructuredDefaultsSection(defaults, "", v))
	_, _, err = c.GetStructuredDefaultsSection("cycle_defaults", "")
	require.Error(t, err)
}

func TestLegacyTraceQuiet(t *testing.T) {
	c, path, err := getTestClient()
	require.NoError(t, err)

	maxFwd := int64(0)
	frontend := &models.Frontend{
		FrontendBase: models.FrontendBase{Name: "quiet_trace", Mode: "http"},
		FilterList: models.Filters{
			{
				Type:            "trace",
				TraceName:       "quiet",
				TraceQuiet:      true,
				TraceRndParsing: true,
				TraceMaxFwd:     &maxFwd,
				Metadata:        map[string]interface{}{"k": "v"},
			},
			{Type: "trace", TraceName: "loud"},
		},
	}
	v, _, err := c.GetStructuredFrontends("")
	require.NoError(t, err)
	require.NoError(t, c.CreateStructuredFrontend(frontend, "", v))

	_, got, err := c.GetStructuredFrontend("quiet_trace", "")
	require.NoError(t, err)
	require.Len(t, got.FilterList, 2)
	require.True(t, got.FilterList[0].TraceQuiet)
	require.True(t, got.FilterList[0].TraceRndParsing)
	require.NotNil(t, got.FilterList[0].TraceMaxFwd)
	require.Equal(t, int64(0), *got.FilterList[0].TraceMaxFwd)
	require.Equal(t, map[string]interface{}{"k": "v"}, got.FilterList[0].Metadata)
	require.False(t, got.FilterList[1].TraceQuiet)

	content, err := os.ReadFile(path)
	require.NoError(t, err)
	lines := strings.Split(string(content), "\n")
	quietLines := 0
	for _, line := range lines {
		if !strings.Contains(line, "filter trace") {
			continue
		}
		if strings.Contains(line, " name quiet") {
			require.Contains(t, line, " quiet")
			quietLines++
		}
		if strings.Contains(line, " name loud") {
			require.NotContains(t, line, " quiet")
		}
	}
	require.Equal(t, 1, quietLines)

	serialized := configuration.SerializeFilter(models.Filter{
		Type:        "trace",
		TraceQuiet:  true,
		TraceMaxFwd: &maxFwd,
	}, &options.ConfigurationOptions{})
	require.Contains(t, serialized.Result().Data, "quiet")
	require.Contains(t, serialized.Result().Data, "max-fwd 0")
	parsed := configuration.ParseFilter(serialized)
	require.NotNil(t, parsed)
	require.True(t, parsed.TraceQuiet)
}
