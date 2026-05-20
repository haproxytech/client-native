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
	"strings"
	"testing"

	"github.com/go-openapi/strfmt"
	"github.com/haproxytech/client-native/v5/models"
	"github.com/stretchr/testify/require"
)

func namePatternStrPtr(s string) *string { return &s }
func namePatternI64Ptr(v int64) *int64   { return &v }
func namePatternBoolPtr(b bool) *bool    { return &b }

type namePatternCase struct {
	label   string
	name    string
	wantErr bool
}

// Cases shared by all six fields. Positive cases use only characters from
// the smaller relaxed class `^[A-Za-z0-9-_${}"]+$` so they also pass the
// MailerEntry pattern (which omits `.` and `:`).
var namePatternCases = []namePatternCase{
	{label: "plain_alphanumeric", name: "mypeer1", wantErr: false},
	{label: "unquoted_variable", name: "${HAPROXY_LOCALPEER}", wantErr: false},
	{label: "quoted_variable", name: `"${HAPROXY_LOCALPEER}"`, wantErr: false},
	{label: "mixed_literal_and_variable", name: "prefix_${VAR}_suffix", wantErr: false},
	{label: "whitespace_rejected", name: "bad name", wantErr: true},
	{label: "slash_rejected", name: "bad/name", wantErr: true},
	{label: "empty_string_rejected", name: "", wantErr: true},
}

type namePatternValidator interface {
	Validate(strfmt.Registry) error
}

func runNamePatternCases(t *testing.T, fieldName string, cases []namePatternCase, build func(string) namePatternValidator) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.label, func(t *testing.T) {
			err := build(tc.name).Validate(strfmt.Default)
			if tc.wantErr {
				require.Error(t, err, "expected error for %s=%q", fieldName, tc.name)
				require.Contains(t, err.Error(), fieldName, "error must mention field %q", fieldName)
				return
			}
			require.NoError(t, err, "unexpected error for %s=%q", fieldName, tc.name)
		})
	}
}

func TestNamePatternPeerEntry(t *testing.T) {
	runNamePatternCases(t, "name", namePatternCases, func(name string) namePatternValidator {
		return &models.PeerEntry{Name: name, Address: namePatternStrPtr("127.0.0.1"), Port: namePatternI64Ptr(4444)}
	})
}

func TestNamePatternNameserver(t *testing.T) {
	runNamePatternCases(t, "name", namePatternCases, func(name string) namePatternValidator {
		return &models.Nameserver{Name: name, Address: namePatternStrPtr("10.0.0.1")}
	})
}

func TestNamePatternMailerEntry(t *testing.T) {
	runNamePatternCases(t, "name", namePatternCases, func(name string) namePatternValidator {
		return &models.MailerEntry{Name: name, Address: "10.0.0.1", Port: 587}
	})
}

func TestNamePatternGroup(t *testing.T) {
	runNamePatternCases(t, "name", namePatternCases, func(name string) namePatternValidator {
		return &models.Group{Name: name}
	})
}

func TestNamePatternUser(t *testing.T) {
	runNamePatternCases(t, "username", namePatternCases, func(name string) namePatternValidator {
		return &models.User{Username: name, Password: "secret", SecurePassword: namePatternBoolPtr(false)}
	})
}

func TestNamePatternFrontendLogTag(t *testing.T) {
	// log_tag is optional; an empty value short-circuits the pattern check,
	// so drop the empty-string negative case for this field.
	cases := make([]namePatternCase, 0, len(namePatternCases))
	for _, tc := range namePatternCases {
		if tc.name == "" {
			continue
		}
		cases = append(cases, tc)
	}
	runNamePatternCases(t, "log_tag", cases, func(name string) namePatternValidator {
		return &models.Frontend{Name: "frontend1", LogTag: name}
	})
}

// Section names were deliberately not relaxed; this guards against an
// accidental future widening of the FrontendBase.Name pattern.
func TestNamePatternSectionNameRejectsVariable(t *testing.T) {
	err := (&models.Frontend{Name: `"${SECTION_NAME}"`}).Validate(strfmt.Default)
	require.Error(t, err, "section name must still reject variable substitution")
	require.Contains(t, strings.ToLower(err.Error()), "name")
}
