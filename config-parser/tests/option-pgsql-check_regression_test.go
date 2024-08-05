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

package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/haproxytech/client-native/v6/config-parser/parsers"
)

func TestOptionPgsqlCheckRegression(t *testing.T) {
	data := [][2]string{
		{`option pgsql-check user postgresuser`, `option pgsql-check user postgresuser`},
	}

	parser := &parsers.OptionPgsqlCheck{}

	for _, d := range data {

		line := strings.TrimSpace(d[0])
		expected := d[1]

		err := ProcessLine(line, parser)
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		result, err := parser.Result()
		if err != nil {
			t.Errorf(err.Error())
			return
		}

		var actual string
		if result[0].Comment == "" {
			actual = result[0].Data
		} else {
			actual = fmt.Sprintf("%s # %s", result[0].Data, result[0].Comment)
		}

		if line != actual {
			t.Errorf(fmt.Sprintf("error: has [%s] expects [%s]", actual, expected))
		}
	}
}
