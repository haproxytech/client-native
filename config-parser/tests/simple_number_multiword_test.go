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
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers/simple"
	"github.com/stretchr/testify/require"
)

func TestSimpleNumberMultiword(t *testing.T) {
	t.Run("matches both words", func(t *testing.T) {
		assertNumberResult(t, &simple.Number{Name: "rate-limit sessions"}, "rate-limit sessions 10", "rate-limit sessions 10")
		assertNumberResult(t, &simple.Number{Name: "rate-limit sessions"}, "rate-limit sessions 0", "rate-limit sessions 0")
		assertNumberResult(t, &simple.Number{Name: "rate-limit sessions"}, "rate-limit sessions 10 # comment", "rate-limit sessions 10 # comment")
	})

	t.Run("bounds unsigned input", func(t *testing.T) {
		assertNumberResult(t, &simple.Number{Name: "rate-limit sessions"}, "rate-limit sessions 9223372036854775807", "rate-limit sessions 9223372036854775807")
	})

	t.Run("rejects bad input", func(t *testing.T) {
		for _, line := range []string{
			"rate-limit foo 10",
			"rate limits 10",
			"rate-limit 10",
			"rate-limit sessions",
			"rate-limit sessions abc",
			"rate-limit sessions 9223372036854775808",
		} {
			t.Run(line, func(t *testing.T) {
				parser := &simple.Number{Name: "rate-limit sessions"}
				parser.Init()
				require.Error(t, ProcessLine(line, parser))
			})
		}
	})

	t.Run("single word unchanged", func(t *testing.T) {
		assertNumberResult(t, &simple.Number{Name: "backlog"}, "backlog 10", "backlog 10")
		assertNumberResult(t, &simple.Number{Name: "backlog"}, "backlog 0", "backlog 0")
		assertNumberResult(t, &simple.Number{Name: "backlog"}, "backlog 10 # comment", "backlog 10 # comment")

		for _, line := range []string{"backlog", "backlog abc", "backlogged 10"} {
			t.Run(line, func(t *testing.T) {
				parser := &simple.Number{Name: "backlog"}
				parser.Init()
				require.Error(t, ProcessLine(line, parser))
			})
		}
	})
}

func assertNumberResult(t *testing.T, parser *simple.Number, line, want string) {
	t.Helper()

	parser.Init()
	require.NoError(t, ProcessLine(line, parser))
	result, err := parser.Result()
	require.NoError(t, err)
	got := result[0].Data
	if result[0].Comment != "" {
		got += " # " + result[0].Comment
	}
	require.Equal(t, want, got)
}
