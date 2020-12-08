// Copyright 2019 HAProxy Technologies
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

func TestGetServerAdminState(t *testing.T) {
	tests := []struct {
		args    string
		want    string
		wantErr bool
	}{
		{"", "", true},
		{"-", "", true},
		{"0", "ready", false},
		{"1", "maint", false},
		{"2", "maint", false},
		{"4", "ready", false},
		{"5", "maint", false},
		{"8", "drain", false},
		{"10", "drain", false},
		{"20", "maint", false},
		{"40", "ready", false},
	}
	for _, tt := range tests {
		t.Run(tt.args, func(t *testing.T) {
			got, err := GetServerAdminState(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetServerAdminState() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetServerAdminState() = %v, want %v", got, tt.want)
			}
		})
	}
}
