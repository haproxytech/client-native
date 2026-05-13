// Copyright 2020 HAProxy Technologies
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

package runtime_test

import (
	"testing"

	"github.com/Masterminds/semver/v3"
	"github.com/haproxytech/client-native/v6/runtime"
)

func TestHAProxyVersion(t *testing.T) {
	type want struct {
		Major, Minor, Patch uint64
		Commit              string
		Original            string
	}
	tests := map[string]want{
		"2.3-dev6-2f6f36-13": {
			Major:    2,
			Minor:    3,
			Patch:    0,
			Commit:   "2f6f36",
			Original: "2.3-dev6-2f6f36-13",
		},
		"2.0.18-a42e6c-11": {
			Major:    2,
			Minor:    0,
			Patch:    18,
			Commit:   "a42e6c",
			Original: "2.0.18-a42e6c-11",
		},
		"2.2.5-34b2b10": {
			Major:    2,
			Minor:    2,
			Patch:    5,
			Commit:   "34b2b10",
			Original: "2.2.5-34b2b10",
		},
	}
	for version, result := range tests {
		t.Run(version, func(t *testing.T) {
			v := &runtime.HAProxyVersion{}
			err := v.ParseHAProxyVersion(version)
			if err != nil {
				t.Fatal(err)
			}
			if v.Major() != result.Major {
				t.Errorf("Major: got [%d], want [%d]", v.Major(), result.Major)
			}
			if v.Minor() != result.Minor {
				t.Errorf("Minor: got [%d], want [%d]", v.Minor(), result.Minor)
			}
			if v.Patch() != result.Patch {
				t.Errorf("Patch: got [%d], want [%d]", v.Patch(), result.Patch)
			}
			if v.Commit != result.Commit {
				t.Errorf("Commit: got [%s], want [%s]", v.Commit, result.Commit)
			}
			if v.Original() != result.Original {
				t.Errorf("Original: got [%s], want [%s]", v.Original(), result.Original)
			}
		})
	}
}

func TestHAProxyVersion_IsBiggerOrEqual(t *testing.T) {
	tests := []struct {
		name    string
		current string
		minimum string
		want    bool
	}{
		{
			name:    "Return true when same HAProxy versions",
			current: "2.4.0",
			minimum: "2.4.0",
			want:    true,
		},
		{
			name:    "Return false when minimum major > current major",
			current: "2.0.0",
			minimum: "3.0.0",
			want:    false,
		},
		{
			name:    "Return true when major and minor are same",
			current: "2.4.0",
			minimum: "2.4.0",
			want:    true,
		},
		{
			name:    "Return false when majors are same and minimum minor > current minor",
			current: "2.2.0",
			minimum: "2.3.0",
			want:    false,
		},
		{
			name:    "Return false when majors, minors are same and patch > current patch",
			current: "2.4.0",
			minimum: "2.4.2",
			want:    false,
		},
		{
			name:    "Return true when minimum major < current major",
			current: "3.0.0",
			minimum: "2.0.0",
			want:    true,
		},
		{
			name:    "Return true when majors are same but minor < current minor",
			current: "2.4.0",
			minimum: "2.1.0",
			want:    true,
		},
		{
			name:    "Return true when majors, minors are same but patch < current patch",
			current: "2.4.2",
			minimum: "2.4.1",
			want:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			current := &runtime.HAProxyVersion{Version: semver.MustParse(tt.current)}
			minimum := &runtime.HAProxyVersion{Version: semver.MustParse(tt.minimum)}
			if got := runtime.IsBiggerOrEqual(minimum, current); got != tt.want {
				t.Errorf("IsBiggerOrEqual(%s, %s) = %v, want %v", tt.minimum, tt.current, got, tt.want)
			}
		})
	}
}
