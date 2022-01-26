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

	"github.com/haproxytech/client-native/v3/runtime"
)

func TestHAProxyVersion(t *testing.T) {
	tests := map[string]runtime.HAProxyVersion{
		"2.3-dev6-2f6f36-13": {
			Major:   2,
			Minor:   3,
			Patch:   0,
			Commit:  "2f6f36",
			Version: "2.3-dev6-2f6f36-13",
		},
		"2.0.18-a42e6c-11": {
			Major:   2,
			Minor:   0,
			Patch:   18,
			Commit:  "a42e6c",
			Version: "2.0.18-a42e6c-11",
		},
		"2.2.5-34b2b10": {
			Major:   2,
			Minor:   2,
			Patch:   5,
			Commit:  "34b2b10",
			Version: "2.2.5-34b2b10",
		},
	}
	for version, result := range tests {
		t.Run(version, func(t *testing.T) {
			v := &runtime.HAProxyVersion{}
			err := v.ParseHAProxyVersion(version)
			if err != nil {
				t.Error(err)
			}

			if v.Major != result.Major {
				t.Fail()
				t.Logf("Major value does not match [%d] != [%d]", v.Major, result.Major)
			}
			if v.Minor != result.Minor {
				t.Fail()
				t.Logf("Minor value does not match [%d] != [%d]", v.Minor, result.Minor)
			}
			if v.Patch != result.Patch {
				t.Fail()
				t.Logf("Patch value does not match [%d != [%d]", v.Patch, result.Patch)
			}
			if v.Commit != result.Commit {
				t.Fail()
				t.Logf("Commit value does not match [%s] != [%s]", v.Commit, result.Commit)
			}
			if v.Version != result.Version {
				t.Fail()
				t.Logf("Version value does not match [%s] != [%s]", v.Version, result.Version)
			}
		})
	}
}

func TestHAProxyVersion_IsBiggerOrEqual(t *testing.T) {
	type fields struct {
		Major int
		Minor int
		Patch int
	}
	type args struct {
		minimumVersion runtime.HAProxyVersion
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Return true when same HAProxy versions",
			fields: fields{
				Major: 2,
				Minor: 4,
				Patch: 0,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 4,
				Patch: 0,
			}},
			want: true,
		},
		{
			name: "Return true when minimum major > current major",
			fields: fields{
				Major: 2,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 3,
			}},
			want: true,
		},
		{
			name: "Return true when major and minor are same",
			fields: fields{
				Major: 2,
				Minor: 4,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 4,
			}},
			want: true,
		},
		{
			name: "Return true when majors are same and minimum minor > current minor",
			fields: fields{
				Major: 2,
				Minor: 2,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 3,
			}},
			want: true,
		},
		{
			name: "Return true when majors, minors are same and patch > current patch",
			fields: fields{
				Major: 2,
				Minor: 4,
				Patch: 0,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 4,
				Patch: 2,
			}},
			want: true,
		},
		{
			name: "Return false when minimum major < current major",
			fields: fields{
				Major: 3,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
			}},
			want: false,
		},
		{
			name: "Return false when majors are same but minor < current minor",
			fields: fields{
				Major: 2,
				Minor: 4,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 1,
			}},
			want: false,
		},
		{
			name: "Return false when majors, minors are same but patch < current patch",
			fields: fields{
				Major: 2,
				Minor: 4,
				Patch: 2,
			},
			args: args{minimumVersion: runtime.HAProxyVersion{
				Major: 2,
				Minor: 4,
				Patch: 1,
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &runtime.HAProxyVersion{
				Major: tt.fields.Major,
				Minor: tt.fields.Minor,
				Patch: tt.fields.Patch,
			}
			if got := v.IsBiggerOrEqual(tt.args.minimumVersion); got != tt.want {
				t.Errorf("HAProxyVersion.IsBiggerOrEqual() = %v, want %v", got, tt.want)
			}
		})
	}
}
