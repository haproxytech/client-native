// Copyright 2024 HAProxy Technologies
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

package misc

import "testing"

func Test_parseTimeout_valid(t *testing.T) {
	tests := []struct {
		tOut              string
		defaultMultiplier int64
		want              int64
	}{
		{"200us", milli, 1},
		{"4000us", milli, 4},
		{"30ms", milli, 30},
		{"30s", milli, 30 * second},
		{"2m", milli, 2 * minute},
		{"1h", milli, 1 * hour},
		{"23d", milli, 23 * day},
		{"1234", milli, 1234},
		{"0", milli, 0},
		{"50", second, 50 * second},
		{"3", minute, 3 * minute},
		{"000", day, 0},
	}
	for _, tt := range tests {
		t.Run(tt.tOut, func(t *testing.T) {
			if got := parseTimeout(tt.tOut, tt.defaultMultiplier); *got != tt.want {
				t.Errorf("parseTimeout() = %d, want %d", *got, tt.want)
			}
		})
	}
}

func Test_parseTimeout_invalid(t *testing.T) {
	tests := []struct {
		tOut              string
		defaultMultiplier int64
	}{
		{"-10", milli},
		{"-10s", milli},
		{"", milli},
		{"lol", milli},
		{"20g", milli},
		{"20M", milli},
		{"-0a", second},
		{"25 m", minute},
	}
	for _, tt := range tests {
		t.Run(tt.tOut, func(t *testing.T) {
			if got := parseTimeout(tt.tOut, tt.defaultMultiplier); got != nil {
				t.Errorf("parseTimeout() = %+v, want %+v", got, nil)
			}
		})
	}
}
