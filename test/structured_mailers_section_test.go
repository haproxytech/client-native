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
//

package test

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

func TestGetStructuredMailersSections(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.MailersSections)
	v, mailers, err := clientTest.GetStructuredMailersSections("")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.Equal(t, 1, len(mailers), "mailers sections: found %d, expected 1", len(mailers))

	m[mailers[0].Name] = models.MailersSections{mailers[0]}

	checkStructuredMailersSection(t, m)
}

func TestGetStructuredMailersSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	m := make(map[string]models.MailersSections)

	v, section, err := clientTest.GetStructuredMailersSection(testMailersName, "")
	require.NoError(t, err)
	m[testMailersName] = models.MailersSections{section}

	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	checkStructuredMailersSection(t, m)

	_, _, err = clientTest.GetStructuredMailersSection("doesnotexist", "")
	require.Error(t, err, "Should throw error, non existent mailers section")
}

func checkStructuredMailersSection(t *testing.T, got map[string]models.MailersSections) {
	exp := mailersSectionExpectation()
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

func TestCreateEditDeleteStructuredMailersSection(t *testing.T) {
	clientTest, filename, err := getTestClient()
	require.NoError(t, err)
	defer os.Remove(filename)
	version := int64(1)

	ms := &models.MailersSection{
		MailersSectionBase: models.MailersSectionBase{
			Name:    "newMailer",
			Timeout: misc.ParseTimeout("30s"),
		},
		MailerEntries: map[string]models.MailerEntry{
			"entry1": {
				Name:    "entry1",
				Address: "192.168.1.1",
				Port:    9200,
			},
		},
	}

	err = clientTest.CreateStructuredMailersSection(ms, "", version)
	require.NoError(t, err)
	version++

	v, mailers, err := clientTest.GetStructuredMailersSection("newMailer", "")
	require.NoError(t, err)
	require.True(t, mailers.Equal(*ms), "mailers_section=%s - diff %v", mailers.Name, cmp.Diff(*mailers, *ms))
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.CreateStructuredMailersSection(ms, "", version)
	require.Error(t, err, "Should throw error MailersSection already exists")

	// Modify the section.
	ms.Timeout = misc.ParseTimeout("40s")
	ms.MailerEntries = nil
	err = clientTest.EditStructuredMailersSection("newMailer", ms, "", version)
	require.NoError(t, err)
	version++

	// Check if the modification was effective.
	v, mailers, err = clientTest.GetStructuredMailersSection("newMailer", "")
	require.NoError(t, err)
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)
	require.True(t, mailers.Equal(*ms), "mailers_section=%s - diff %v", mailers.Name, cmp.Diff(*mailers, *ms))

	// Delete the section.
	err = clientTest.DeleteMailersSection("newMailer", "", version)
	require.NoError(t, err)
	version++

	v, _ = clientTest.GetVersion("")
	require.Equal(t, version, v, "Version %v returned, expected %v", v, version)

	err = clientTest.DeleteMailersSection("newMailer", "", 999999)
	require.Error(t, err, "Should throw error, non existent frontend")
	require.ErrorIs(t, err, configuration.ErrVersionMismatch, "Should throw configuration.ErrVersionMismatch error")

	_, _, err = clientTest.GetStructuredMailersSection("newMailer", "")
	require.Error(t, err, "DeleteMailersSection failed: newMailer still exists")

	err = clientTest.DeleteMailersSection("doesnotexist", "", version)
	require.Error(t, err, "Should throw error, non existent MailersSection")
}
