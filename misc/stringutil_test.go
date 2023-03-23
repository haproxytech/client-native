// Copyright 2020 HAProxy Technologies
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package misc

import (
	"testing"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "Should preserve # and . in input",
			input: "#4_abc.#",
			want:  "#4_abc.#",
		},
		{
			name:  "Should convert leading dots",
			input: ".....file.map!!@#",
			want:  "_file.map_#",
		},
		{
			name:  "Should convert hidden files",
			input: ".hidden",
			want:  "_hidden",
		},
		{
			name:  "Should accept unusual filenames",
			input: ".unusual.",
			want:  "_unusual.",
		},
		{
			name:  "Should sanitize input correctly",
			input: "#1_?a;b/c!&?",
			want:  "#1__a_b_c_",
		},
		{
			name:  "Should return same input when name doesn't contain regex characters",
			input: "abcDEF",
			want:  "abcDEF",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SanitizeFilename(tt.input); got != tt.want {
				t.Errorf("SanitizeFilename() = %v, want %v", got, tt.want)
			}
		})
	}
}
