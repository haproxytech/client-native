// Copyright 2022 HAProxy Technologies
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

package test

import (
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v5/models"
)

func httpErrorSectionExpectation() map[string]models.HTTPErrorsSections {
	initStructuredExpected()
	res := StructuredToHTTPErrorSectionMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.HTTPErrorsSections{v}
		}
	}
	return res
}

var testDataHTTPErrors map[string]models.HTTPErrorsSection

func init() {
	exp := httpErrorSectionExpectation()
	testDataHTTPErrors = map[string]models.HTTPErrorsSection{
		"website-1": *exp["website-1"][0],
		"website-2": *exp["website-2"][0],
		"new": {
			Name: "website-3",
			ErrorFiles: []*models.Errorfile{
				{Code: 400, File: "/etc/haproxy/errorfiles/site3/400.http"},
				{Code: 404, File: "/etc/haproxy/errorfiles/site3/404.http"},
				{Code: 501, File: "/etc/haproxy/errorfiles/site3/501.http"},
			},
		},
		"not-there": {
			Name: "i_am_not_there",
			ErrorFiles: []*models.Errorfile{
				{Code: 400, File: "/etc/haproxy/errorfiles/site3/400.http"},
				{Code: 404, File: "/etc/haproxy/errorfiles/site3/404.http"},
				{Code: 501, File: "/etc/haproxy/errorfiles/site3/501.http"},
			},
		},
		"edit": {
			Name: "website-2",
			ErrorFiles: []*models.Errorfile{
				{Code: 400, File: "/etc/haproxy/errorfiles/site3/400.http"},
				{Code: 404, File: "/etc/haproxy/errorfiles/site3/404.http"},
				{Code: 501, File: "/etc/haproxy/errorfiles/site3/501.http"},
			},
		},
		"no-entries": {
			Name: "website-3",
		},
		"no-name": {
			ErrorFiles: []*models.Errorfile{
				{Code: 400, File: "/etc/haproxy/errorfiles/site3/400.http"},
				{Code: 404, File: "/etc/haproxy/errorfiles/site3/404.http"},
				{Code: 501, File: "/etc/haproxy/errorfiles/site3/501.http"},
			},
		},
		"entry-with-unsupported-code": {
			Name: "website-4",
			ErrorFiles: []*models.Errorfile{
				{Code: 406, File: "/etc/haproxy/errorfiles/site3/400.http"},
				{Code: 404, File: "/etc/haproxy/errorfiles/site3/404.http"},
				{Code: 501, File: "/etc/haproxy/errorfiles/site3/501.http"},
			},
		},
	}
}

func TestGetHTTPErrorsSections(t *testing.T) {
	m := make(map[string]models.HTTPErrorsSections)
	v, sections, err := clientTest.GetHTTPErrorsSections("")

	require.NoError(t, err)
	require.Len(t, sections, 2)
	require.Equal(t, version, v)
	m[""] = sections

	checkHTTPErrorSection(t, m)
}

func checkHTTPErrorSection(t *testing.T, got map[string]models.HTTPErrorsSections) {
	exp := httpErrorSectionExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestGetHTTPErrorsSection(t *testing.T) {
	m := make(map[string]models.HTTPErrorsSections)

	// Test a successful operation.
	v, section, err := clientTest.GetHTTPErrorsSection("website-2", "")
	require.NoError(t, err)
	require.NotNil(t, section)
	require.Equal(t, version, v)
	m["website-2"] = models.HTTPErrorsSections{section}

	checkHTTPErrorSection(t, m)
	// Test expected failures.
	_, _, err = clientTest.GetHTTPErrorsSection(testDataHTTPErrors["not-there"].Name, "")
	assert.Error(t, err, "attempt to retrieve a section that does not exist should fail")
}

func TestCreateHTTPErrorsSection(t *testing.T) {
	// Test a successful operation.
	valid := testDataHTTPErrors["new"]
	err := clientTest.CreateHTTPErrorsSection(&valid, "", version)
	require.NoError(t, err)
	version++

	configVersion, found, err := clientTest.GetHTTPErrorsSection(valid.Name, "")
	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, configVersion, version)
	assert.Equal(t, &valid, found, "added and subsequently retrieved sections should match")

	// Test expected failures.
	existing := testDataHTTPErrors["website-2"]
	err = clientTest.CreateHTTPErrorsSection(&existing, "", version)
	assert.Error(t, err, "attempt to create a section with existing name should fail")

	invalid := testDataHTTPErrors["no-entries"]
	err = clientTest.CreateHTTPErrorsSection(&invalid, "", version)
	assert.Error(t, err, "attempt to create a section without entries should fail")

	invalid = testDataHTTPErrors["no-name"]
	err = clientTest.CreateHTTPErrorsSection(&invalid, "", version)
	assert.Error(t, err, "attempt to create a section without a name should fail")

	invalid = testDataHTTPErrors["entry-with-unsupported-code"]
	err = clientTest.CreateHTTPErrorsSection(&invalid, "", version)
	assert.Error(t, err, "attempt to create a section with an unsupported HTTP status code in an entry should fail")
}

func TestEditHTTPErrorsSection(t *testing.T) {
	// Test a successful operation.
	valid := testDataHTTPErrors["edit"]
	err := clientTest.EditHTTPErrorsSection(valid.Name, &valid, "", version)
	require.NoError(t, err)
	version++

	configVersion, found, err := clientTest.GetHTTPErrorsSection(valid.Name, "")
	assert.NoError(t, err)
	assert.Equal(t, configVersion, version)
	assert.Equal(t, valid, *found, "provided section data should match data in subsequently retrieved section")

	// Test expected failures.
	notThere := testDataHTTPErrors["not-there"]
	err = clientTest.EditHTTPErrorsSection(notThere.Name, &notThere, "", version)
	assert.Error(t, err, "attempt to replace a section that does not exist should fail")

	invalid := testDataHTTPErrors["no-entries"]
	err = clientTest.EditHTTPErrorsSection(invalid.Name, &invalid, "", version)
	assert.Error(t, err, "attempt to replace a section with a one that has no entries should fail")

	invalid = testDataHTTPErrors["no-name"]
	err = clientTest.EditHTTPErrorsSection(invalid.Name, &invalid, "", version)
	assert.Error(t, err, "attempt to replace a section with a one that has no name should fail")

	invalid = testDataHTTPErrors["entry-with-unsupported-code"]
	err = clientTest.EditHTTPErrorsSection(invalid.Name, &invalid, "", version)
	assert.Error(t, err, "attempt to replace a section with a one that has unsupported HTTP status code in an entry should fail")
}

func TestDeleteHTTPErrorsSection(t *testing.T) {
	// Test a successful operation.
	existing := testDataHTTPErrors["website-1"]
	err := clientTest.DeleteHTTPErrorsSection(existing.Name, "", version)
	require.NoError(t, err)
	version++

	_, _, err = clientTest.GetHTTPErrorsSection(existing.Name, "")
	assert.Error(t, err, "retrieving a section that was deleted should fail")

	// Test expected failures.
	notThere := testDataHTTPErrors["not-there"]
	err = clientTest.EditHTTPErrorsSection(notThere.Name, &notThere, "", version)
	assert.Error(t, err, "attempt to delete a section that does not exist should fail")
}
