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

package sorter_test

import (
	"testing"

	parserErrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/sorter"
)

func TestSortWithFrom(t *testing.T) {
	tests := []struct {
		name          string
		sections      []sorter.Section
		result        []sorter.Section
		ExpectedError error
	}{
		{"empty", []sorter.Section{}, []sorter.Section{}, nil},
		{"1", []sorter.Section{{"1", ""}}, []sorter.Section{{"1", ""}}, nil},
		{"2", []sorter.Section{{"1", ""}, {"2", ""}}, []sorter.Section{{"1", ""}, {"2", ""}}, nil},
		{"2a", []sorter.Section{{"1", ""}, {"2", "1"}}, []sorter.Section{{"1", ""}, {"2", "1"}}, nil},
		{"2b", []sorter.Section{{"2", "1"}, {"1", ""}}, []sorter.Section{{"1", ""}, {"2", "1"}}, nil},
		{"2 self depend", []sorter.Section{{"2", "2"}, {"1", "1"}}, []sorter.Section{{"2", "2"}, {"1", "1"}}, parserErrors.ErrCircularDependency},
		{"with non existent sorter.Section", []sorter.Section{{"2", "3"}, {"1", ""}}, []sorter.Section{{"1", ""}, {"2", "3"}}, parserErrors.ErrFromDefaultsSectionMissing},
		{"with non existent sorter.Section 2", []sorter.Section{{"0", "2"}, {"2", "3"}, {"1", ""}}, []sorter.Section{{"1", ""}, {"2", "3"}, {"0", "2"}}, parserErrors.ErrFromDefaultsSectionMissing},
		{"3", []sorter.Section{{"0", ""}, {"2", "1"}, {"1", ""}}, []sorter.Section{{"0", ""}, {"1", ""}, {"2", "1"}}, nil},
		{"4", []sorter.Section{{"0", ""}, {"3", "1"}, {"2", "1"}, {"1", ""}}, []sorter.Section{{"0", ""}, {"1", ""}, {"2", "1"}, {"3", "1"}}, nil},
		{"4a", []sorter.Section{{"0", ""}, {"3", "0"}, {"2", "1"}, {"1", ""}}, []sorter.Section{{"0", ""}, {"1", ""}, {"2", "1"}, {"3", "0"}}, nil},
		{"err 1", []sorter.Section{{"0", "1"}, {"1", "0"}}, []sorter.Section{{"0", "1"}, {"1", "0"}}, parserErrors.ErrCircularDependency},
		{"err 2", []sorter.Section{{"0", "1"}, {"1", "2"}, {"2", "3"}, {"3", "0"}}, []sorter.Section{{"0", "1"}, {"1", "2"}, {"2", "3"}, {"3", "0"}}, parserErrors.ErrCircularDependency},
		{"malicious", []sorter.Section{{"0", "3"}, {"1", "2"}, {"2", "3"}, {"3", "1"}}, []sorter.Section{{"0", "3"}, {"1", "2"}, {"2", "3"}, {"3", "1"}}, parserErrors.ErrCircularDependency},
		{"ok but reorder", []sorter.Section{{"0", "3"}, {"1", ""}, {"2", ""}, {"3", ""}}, []sorter.Section{{"1", ""}, {"2", ""}, {"3", ""}, {"0", "3"}}, nil},
		{"ok but reorder 2", []sorter.Section{{"0", "3"}, {"1", "2"}, {"2", ""}, {"3", ""}}, []sorter.Section{{"2", ""}, {"1", "2"}, {"3", ""}, {"0", "3"}}, nil},
		{"ok but reorder 2", []sorter.Section{{"0", "3"}, {"1", "0"}, {"2", "1"}, {"3", ""}}, []sorter.Section{{"3", ""}, {"0", "3"}, {"1", "0"}, {"2", "1"}}, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sorter.Sort(tt.sections)
			if tt.ExpectedError == nil {
				if err != nil {
					t.Errorf("SortWithFrom() error = %v, expectedError %v", err, tt.ExpectedError)
				}
			} else {
				if err != tt.ExpectedError { //nolint:errorlint
					t.Errorf("SortWithFrom() error = %v, expectedError %v", err, tt.ExpectedError)
				}
			}
			if err == nil {
				if !Equal(tt.sections, tt.result) {
					t.Errorf("SortWithFrom() src = %v, result %v", tt.sections, tt.result)
				}
			}
		})
	}
}
