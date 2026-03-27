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

package runtime_test

import (
	"errors"
	"testing"

	"github.com/haproxytech/client-native/v5/runtime"
)

func TestSingleRuntime_Execute(t *testing.T) {
	tests := []struct {
		name    string // description of this test case
		command string
		wantErr error
	}{
		{
			name:    "two commands with semicolon",
			command: "show ssl cert foo;dump ssl foo",
			wantErr: runtime.ErrRuntimeInvalidChar,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s runtime.SingleRuntime
			gotErr := s.Execute(tt.command)
			if !errors.Is(gotErr, tt.wantErr) {
				t.Errorf("got %v, expected %v", gotErr, tt.wantErr)
			}
		})
	}
}
