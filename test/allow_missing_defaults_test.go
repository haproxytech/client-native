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

package test

import (
	"strings"
	"testing"

	parser "github.com/haproxytech/client-native/v6/config-parser"
	parser_errors "github.com/haproxytech/client-native/v6/config-parser/errors"
	parser_options "github.com/haproxytech/client-native/v6/config-parser/options"
	"github.com/stretchr/testify/require"
)

const baseConfiguration = `
defaults existing_defaults
  mode http

frontend fe1
  mode http
`

func TestAllowMissingDefaults(t *testing.T) {
	tests := []struct {
		Name          string
		ParserOptions []parser_options.ParserOption
		From          string
		ExpectedErr   error
		ExpectedFrom  string
		ExpectedConf  string
	}{
		{
			Name:          "no option: frontend from missing defaults section fails",
			ParserOptions: nil,
			From:          "missing_defaults",
			ExpectedErr:   parser_errors.ErrFromDefaultsSectionMissing,
		},
		{
			Name:          "empty option list: frontend from missing defaults section fails",
			ParserOptions: []parser_options.ParserOption{parser_options.AllowMissingDefaults([]string{})},
			From:          "missing_defaults",
			ExpectedErr:   parser_errors.ErrFromDefaultsSectionMissing,
		},
		{
			Name: "option with duplicates not containing the missing defaults section fails",
			ParserOptions: []parser_options.ParserOption{
				parser_options.AllowMissingDefaults([]string{"other_defaults", "other_defaults", "another_defaults"}),
			},
			From:        "missing_defaults",
			ExpectedErr: parser_errors.ErrFromDefaultsSectionMissing,
		},
		{
			Name: "option with the missing defaults section succeeds",
			ParserOptions: []parser_options.ParserOption{
				parser_options.AllowMissingDefaults([]string{"other_defaults", "missing_defaults"}),
			},
			From:         "missing_defaults",
			ExpectedFrom: "missing_defaults",
			ExpectedConf: `
defaults existing_defaults
  mode http

frontend fe1 from missing_defaults
  mode http
`,
		},
		{
			Name:          "no option: frontend from existing defaults section succeeds",
			ParserOptions: nil,
			From:          "existing_defaults",
			ExpectedFrom:  "existing_defaults",
			ExpectedConf: `
defaults existing_defaults
  mode http

frontend fe1 from existing_defaults
  mode http
`,
		},
		{
			Name:          "no option: frontend without from succeeds",
			ParserOptions: nil,
			From:          "",
			ExpectedFrom:  "",
			ExpectedConf: `
defaults existing_defaults
  mode http

frontend fe1
  mode http
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			parserOpts := append(tt.ParserOptions, parser_options.Reader(strings.NewReader(baseConfiguration)))
			p, err := parser.New(parserOpts...)
			require.NoError(t, err)

			// set default section on frontend
			err = p.SectionsDefaultsFromSet(parser.Frontends, "fe1", tt.From)

			if tt.ExpectedErr != nil {
				require.ErrorIs(t, err, tt.ExpectedErr)
				return
			}
			require.NoError(t, err)

			// check that frontend reference the default section
			name, err := p.SectionsDefaultsFromGet(parser.Frontends, "fe1")
			require.NoError(t, err)
			require.Equal(t, tt.ExpectedFrom, name)

			// check parser to raw
			result := p.String()
			require.Equal(t, tt.ExpectedConf, result)
		})
	}
}
