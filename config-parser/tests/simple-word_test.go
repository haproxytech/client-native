/*
Copyright 2026 HAProxy Technologies

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

package tests

import (
	"errors"
	"testing"

	parsererrors "github.com/haproxytech/client-native/v6/config-parser/errors"
	"github.com/haproxytech/client-native/v6/config-parser/parsers/simple"
	"github.com/haproxytech/client-native/v6/config-parser/types"
)

// A Word never parses a bare keyword, so it must not write one either: an
// empty value renders nothing instead of a line HAProxy rejects.
func TestWordEmptyValueIsNotWritten(t *testing.T) {
	parser := &simple.Word{Name: "crt-base"}
	parser.Init()

	if err := parser.Set(types.StringC{Value: ""}, -1); err != nil {
		t.Fatal(err)
	}
	if _, err := parser.Result(); !errors.Is(err, parsererrors.ErrFetch) {
		t.Errorf("empty value: got %v, want ErrFetch", err)
	}

	if err := parser.Set(types.StringC{Value: "/certs"}, -1); err != nil {
		t.Fatal(err)
	}
	result, err := parser.Result()
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 1 || result[0].Data != "crt-base /certs" {
		t.Errorf("got %+v, want [crt-base /certs]", result)
	}
}
