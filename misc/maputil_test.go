// Copyright 2026 HAProxy Technologies
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
package misc //nolint: revive

import (
	"slices"
	"testing"
)

func TestSortedByKey(t *testing.T) {
	m := map[string]int{"v6": 2, "quicv4": 3, "v4": 1, "quicv6": 4}
	want := []string{"quicv4", "quicv6", "v4", "v6"}
	for range 20 {
		var got []string
		for k, v := range SortedByKey(m) {
			if m[k] != v {
				t.Fatalf("value mismatch for %s: %d", k, v)
			}
			got = append(got, k)
		}
		if !slices.Equal(got, want) {
			t.Fatalf("got %v, want %v", got, want)
		}
	}
	for k := range SortedByKey(m) {
		if k != "quicv4" {
			t.Fatalf("early break yielded %s", k)
		}
		break
	}
}
