/*
Copyright 2019 HAProxy Technologies

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

package common //nolint:testpackage

import (
	"reflect"
	"testing"
)

func TestStringSplitWithCommentIgnoreEmpty(t *testing.T) {
	tests := []struct {
		name        string
		s           string
		wantData    []string
		wantComment string
	}{
		{"simple", "  http-request deny deny_status 400 # deny", []string{"http-request", "deny", "deny_status", "400"}, "deny"},
		{"empty comment", "  http-request deny deny_status 400 # ", []string{"http-request", "deny", "deny_status", "400"}, ""},
		{"no comment", "  http-request deny deny_status 400  ", []string{"http-request", "deny", "deny_status", "400"}, ""},
		{"single quote", "  acl 'fdsfsdfsd sdf s f' abc\\ def ", []string{"acl", "'fdsfsdfsd sdf s f'", "abc\\ def"}, ""},
		{"escaped space", `  acl abc\ def `, []string{"acl", `abc\ def`}, ""},
		{"escaped space with #", `  acl abc\ def #comment`, []string{"acl", `abc\ def`}, "comment"},
		{"escaped space with # formated", `  acl abc\ def # comment`, []string{"acl", `abc\ def`}, "comment"},
		{"single quote escaped double", `  acl 'fdsfsdfsd "sdf" s f' abc\ def `, []string{"acl", `'fdsfsdfsd "sdf" s f'`, `abc\ def`}, ""},
		{"double quote escaped single", `  acl "fdsfsdfsd 'sdf' s f" abc\ def `, []string{"acl", `"fdsfsdfsd 'sdf' s f"`, `abc\ def`}, ""},
		{"escaped single", `  abc \'def ghi\' jkl`, []string{"abc", `\'def`, `ghi\'`, `jkl`}, ""},
		{"escaped double", `  abc \"def ghi\" jkl`, []string{"abc", `\"def`, `ghi\"`, `jkl`}, ""},
		{"escaped double", `  abc "defghi\"" jkl`, []string{"abc", `"defghi\""`, `jkl`}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotData, gotComment := StringSplitWithCommentIgnoreEmpty(tt.s)
			if !reflect.DeepEqual(gotData, tt.wantData) {
				t.Errorf("StringSplitWithCommentIgnoreEmpty() gotData = %v, want %v", gotData, tt.wantData)
			}
			if gotComment != tt.wantComment {
				t.Errorf("StringSplitWithCommentIgnoreEmpty() gotComment = %v, want %v", gotComment, tt.wantComment)
			}
		})
	}
}

func TestCutRight(t *testing.T) {
	tests := []struct {
		name   string
		s      string
		sep    string
		before string
		after  string
		found  bool
	}{
		{
			name:   "simple",
			s:      "ab:cd:ef",
			sep:    ":",
			before: "ab:cd",
			after:  "ef",
			found:  true,
		},
		{
			name:   "multi-characters-sep",
			s:      "ab:-:cd:-:ef",
			sep:    ":-:",
			before: "ab:-:cd",
			after:  "ef",
			found:  true,
		},
		{
			name:   "sep-at-start",
			s:      ":abcd",
			sep:    ":",
			before: "",
			after:  "abcd",
			found:  true,
		},
		{
			name:   "sep-at-end",
			s:      "abcd:",
			sep:    ":",
			before: "abcd",
			after:  "",
			found:  true,
		},
		{
			name:   "no-sep",
			s:      "abcd",
			sep:    ":",
			before: "abcd",
			after:  "",
			found:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			before, after, found := CutRight(tt.s, tt.sep)
			if before != tt.before {
				t.Errorf("Part before separator doesn't match expected: %q != %q", before, tt.before)
				return
			}
			if after != tt.after {
				t.Errorf("Part after separator doesn't match expected: %q != %q", after, tt.after)
				return
			}
			if found != tt.found {
				t.Errorf("Found result doesn't match expected: %v != %v", found, tt.found)
				return
			}
		})
	}
}

func TestSmartJoin(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "simple1",
			args: []string{"set-param", "abc", "if", "cond_test"},
			want: "set-param abc if cond_test",
		},
		{
			name: "with empty strings",
			args: []string{"set-param", "abc", "", ""},
			want: "set-param abc",
		},
		{
			name: "all empty strings",
			args: []string{"", "", "", ""},
			want: "",
		},
		{
			name: "mixed",
			args: []string{"", "header", "", "zero"},
			want: " header zero",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SmartJoin(tt.args...); got != tt.want {
				t.Errorf("SmartJoin() = '%v', want '%v'", got, tt.want)
			}
		})
	}
}
