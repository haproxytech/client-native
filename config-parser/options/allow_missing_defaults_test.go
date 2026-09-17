/*
Copyright 2026 HAProxy Technologies

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package options

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllowMissingDefault(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "empty list",
			input:    []string{},
			expected: []string{},
		},
		{
			name:     "empty list (due to trimming)",
			input:    []string{" ", "    "},
			expected: []string{},
		},
		{
			name:     "names are trimmed",
			input:    []string{"default1", "default2 ", " default3", "   default4  "},
			expected: []string{"default1", "default2", "default3", "default4"},
		},
		{
			name:     "duplicates are removed",
			input:    []string{"default1", "default2", "default1", "default3", "default3"},
			expected: []string{"default1", "default2", "default3"},
		},
		{
			name:     "names are trimmed and duplicates are removed",
			input:    []string{"default1", "default2 ", " default3", "   default4  ", "default2"},
			expected: []string{"default1", "default2", "default3", "default4"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var p Parser
			err := AllowMissingDefaults(tc.input).Set(&p)
			require.NoError(t, err)
			require.ElementsMatch(t, tc.expected, p.AllowMissingDefaults)
		})
	}
}
