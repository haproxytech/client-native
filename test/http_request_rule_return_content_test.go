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

package test

import (
	"testing"

	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/configuration/options"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

// TestSerializeReturnContentValidation guards against the silent-drop bug where
// an unquoted multi-word return content serializes to an invalid HAProxy config
// line (e.g. `... string Missing Auth`) that fails to parse and is silently
// dropped on the next config reload, while the API reports success.
//
// It covers all serialize sites that emit return content: http-request and
// http-response deny/return, and http-error status.
func TestSerializeReturnContentValidation(t *testing.T) {
	opt := &options.ConfigurationOptions{}

	cases := []struct {
		name    string
		format  string
		content string
		wantErr bool
	}{
		{"unquoted_multiword", "string", "Missing Auth", true},
		{"quoted_multiword", "string", `"Missing Auth"`, false},
		{"escaped_space", "string", `Missing\ Auth`, false},
		{"single_word", "string", "MissingAuth", false},
		{"empty_content", "string", "", false},
		// Content is never emitted when the format is empty or default-errorfiles,
		// so it cannot break the round-trip and must not be rejected.
		{"empty_format_multiword", "", "Missing Auth", false},
		{"default_errorfiles_multiword", "default-errorfiles", "Missing Auth", false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// http-request deny
			_, err := configuration.SerializeHTTPRequestRule(models.HTTPRequestRule{
				Type:                "deny",
				ReturnContentType:   misc.Ptr("text/html"),
				ReturnContentFormat: tc.format,
				ReturnContent:       tc.content,
			}, opt)
			assertContentErr(t, "http-request deny", tc.wantErr, err)

			// http-request return
			_, err = configuration.SerializeHTTPRequestRule(models.HTTPRequestRule{
				Type:                "return",
				ReturnContentType:   misc.Ptr("text/html"),
				ReturnContentFormat: tc.format,
				ReturnContent:       tc.content,
			}, opt)
			assertContentErr(t, "http-request return", tc.wantErr, err)

			// http-response deny
			_, err = configuration.SerializeHTTPResponseRule(models.HTTPResponseRule{
				Type:                "deny",
				ReturnContentType:   misc.Ptr("text/html"),
				ReturnContentFormat: tc.format,
				ReturnContent:       tc.content,
			}, opt)
			assertContentErr(t, "http-response deny", tc.wantErr, err)

			// http-response return
			_, err = configuration.SerializeHTTPResponseRule(models.HTTPResponseRule{
				Type:                "return",
				ReturnContentType:   misc.Ptr("text/html"),
				ReturnContentFormat: tc.format,
				ReturnContent:       tc.content,
			}, opt)
			assertContentErr(t, "http-response return", tc.wantErr, err)

			// http-error status (403 is a valid error status code)
			_, err = configuration.SerializeHTTPErrorRule(models.HTTPErrorRule{
				Type:                "status",
				Status:              403,
				ReturnContentType:   misc.Ptr("text/html"),
				ReturnContentFormat: tc.format,
				ReturnContent:       tc.content,
			})
			assertContentErr(t, "http-error status", tc.wantErr, err)
		})
	}
}

func assertContentErr(t *testing.T, site string, wantErr bool, err error) {
	t.Helper()
	if wantErr {
		require.Error(t, err, "%s: value should be rejected, not silently dropped", site)
	} else {
		require.NoError(t, err, "%s: value should be accepted", site)
	}
}
