// Copyright 2025 HAProxy Technologies
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

	"github.com/google/go-cmp/cmp"
)

func Test_serializeAcmeVars(t *testing.T) {
	tests := []struct {
		vars    map[string]string
		want    string
		wantErr bool
	}{
		{
			vars:    map[string]string{"foo": "bar", "ApiKey": "FEFF,==\""},
			want:    `"foo=bar,ApiKey=FEFF\,==\""`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got, err := serializeAcmeVars(tt.vars)
			if tt.wantErr != (err != nil) {
				t.Errorf("serializeAcmeVars() got error '%v', wantErr=%v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("serializeAcmeVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_ParseAcmeVars(t *testing.T) {
	tests := []struct {
		vars string
		want map[string]string
	}{
		{
			vars: `"foo=bar,ApiKey=FEFF\,==\""`,
			want: map[string]string{"foo": "bar", "ApiKey": "FEFF,==\""},
		},
	}
	for _, tt := range tests {
		t.Run(tt.vars, func(t *testing.T) {
			got := ParseAcmeVars(tt.vars)
			if !cmp.Equal(got, tt.want) {
				t.Errorf("parseAcmeVars() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
