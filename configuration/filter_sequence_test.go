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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v6/models"
)

func mustGetFilterSequences(t *testing.T, c Configuration, parentType, parentName string) (int64, models.FilterSequences) {
	t.Helper()
	v, sequences, err := c.GetFilterSequences(parentType, parentName, "")
	require.NoError(t, err)
	return v, sequences
}

func filterSelector(class, id string) *models.FilterSelector {
	return &models.FilterSelector{Class: class, ID: id}
}

func filterSequenceLines(sequences models.FilterSequences) []string {
	lines := make([]string, 0, len(sequences))
	for _, sequence := range sequences {
		line := sequence.Side
		for _, rule := range sequence.Rules {
			line += " " + filterSelectorText(rule.Instance) + " " + rule.Direction + " " + filterSelectorText(rule.Target)
		}
		lines = append(lines, line)
	}
	return lines
}

func TestFilterSequenceCRUD(t *testing.T) {
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

			first := &models.FilterSequence{
				Side: "request",
				Rules: []*models.FilterSequenceRule{
					{Instance: filterSelector("trace", "A"), Direction: "before", Target: filterSelector("trace", "B")},
				},
				Metadata: map[string]interface{}{"k": "v"},
			}
			v, _ := mustGetFilterSequences(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterSequence(0, parentType, parentName, first, "", v))

			second := &models.FilterSequence{
				Side: "response",
				Rules: []*models.FilterSequenceRule{
					{Instance: filterSelector("trace", "C"), Direction: "after", Target: filterSelector("trace", "D")},
				},
			}
			v, _ = mustGetFilterSequences(t, c, parentType, parentName)
			require.NoError(t, c.CreateFilterSequence(1, parentType, parentName, second, "", v))

			v, sequences := mustGetFilterSequences(t, c, parentType, parentName)
			require.Equal(t, []string{"request trace/A before trace/B", "response trace/C after trace/D"}, filterSequenceLines(sequences))

			_, got, err := c.GetFilterSequence(0, parentType, parentName, "")
			require.NoError(t, err)
			require.Equal(t, map[string]interface{}{"k": "v"}, got.Metadata)

			edited := &models.FilterSequence{
				Side: "global",
				Rules: []*models.FilterSequenceRule{
					{Instance: filterSelector("trace", "E"), Direction: "before", Target: filterSelector("trace", "F")},
				},
			}
			require.NoError(t, c.EditFilterSequence(0, parentType, parentName, edited, "", v))
			v, sequences = mustGetFilterSequences(t, c, parentType, parentName)
			require.Equal(t, []string{"global trace/E before trace/F", "response trace/C after trace/D"}, filterSequenceLines(sequences))

			require.NoError(t, c.DeleteFilterSequence(0, parentType, parentName, "", v))
			v, sequences = mustGetFilterSequences(t, c, parentType, parentName)
			require.Equal(t, []string{"response trace/C after trace/D"}, filterSequenceLines(sequences))

			reversed := models.FilterSequences{second, edited}
			require.NoError(t, c.ReplaceFilterSequences(parentType, parentName, reversed, "", v))
			v, sequences = mustGetFilterSequences(t, c, parentType, parentName)
			require.Equal(t, []string{"response trace/C after trace/D", "global trace/E before trace/F"}, filterSequenceLines(sequences))

			require.NoError(t, c.ReplaceFilterSequences(parentType, parentName, models.FilterSequences{}, "", v))
			v, sequences = mustGetFilterSequences(t, c, parentType, parentName)
			require.Empty(t, sequences)

			require.NoError(t, c.CreateFilterSequence(0, parentType, parentName, first, "", v))
			v, _ = mustGetFilterSequences(t, c, parentType, parentName)
			_, _, err = c.GetFilterSequence(99, parentType, parentName, "")
			require.Error(t, err)
			require.Error(t, c.EditFilterSequence(99, parentType, parentName, edited, "", v))
			require.Error(t, c.DeleteFilterSequence(99, parentType, parentName, "", v))
			require.Error(t, c.CreateFilterSequence(0, parentType, parentName, first, "", v+100))
		})
	}

	t.Run("missing parent", func(t *testing.T) {
		c, _ := newFilterTestClient(t)
		_, _, err := c.GetFilterSequences(FrontendParentName, "doesnotexist", "")
		require.Error(t, err)
	})
}

func TestFilterSequenceMutationCycles(t *testing.T) {
	c, _ := newFilterTestClient(t)

	valid := &models.FilterSequence{
		Side: "request",
		Rules: []*models.FilterSequenceRule{
			{Instance: filterSelector("trace", "A"), Direction: "before", Target: filterSelector("trace", "B")},
		},
	}
	v, _ := mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.NoError(t, c.CreateFilterSequence(0, FrontendParentName, "fe", valid, "", v))

	cyclicEdit := &models.FilterSequence{
		Side: "request",
		Rules: []*models.FilterSequenceRule{
			{Instance: filterSelector("trace", "A"), Direction: "before", Target: filterSelector("trace", "B")},
			{Instance: filterSelector("trace", "B"), Direction: "before", Target: filterSelector("trace", "A")},
		},
	}
	v, _ = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Error(t, c.EditFilterSequence(0, FrontendParentName, "fe", cyclicEdit, "", v))
	_, sequences := mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Equal(t, []string{"request trace/A before trace/B"}, filterSequenceLines(sequences))

	cyclicInsert := &models.FilterSequence{
		Side: "request",
		Rules: []*models.FilterSequenceRule{
			{Instance: filterSelector("trace", "A"), Direction: "before", Target: filterSelector("trace", "B")},
			{Instance: filterSelector("trace", "B"), Direction: "before", Target: filterSelector("trace", "A")},
		},
	}
	v, _ = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Error(t, c.CreateFilterSequence(1, FrontendParentName, "fe", cyclicInsert, "", v))
	_, sequences = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Equal(t, []string{"request trace/A before trace/B"}, filterSequenceLines(sequences))

	v, _ = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Error(t, c.ReplaceFilterSequences(FrontendParentName, "fe", models.FilterSequences{cyclicEdit}, "", v))
	_, sequences = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Equal(t, []string{"request trace/A before trace/B"}, filterSequenceLines(sequences))

	// Dangling references are accepted when the graph stays acyclic.
	dangling := models.FilterSequences{
		{
			Side: "request",
			Rules: []*models.FilterSequenceRule{
				{Instance: filterSelector("trace", "X"), Direction: "before", Target: filterSelector("trace", "Y")},
			},
		},
	}
	v, _ = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.NoError(t, c.ReplaceFilterSequences(FrontendParentName, "fe", dangling, "", v))
	_, sequences = mustGetFilterSequences(t, c, FrontendParentName, "fe")
	require.Equal(t, []string{"request trace/X before trace/Y"}, filterSequenceLines(sequences))
}
