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
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/haproxytech/client-native/v6/configuration"
	"github.com/haproxytech/client-native/v6/misc"
	"github.com/haproxytech/client-native/v6/models"
	"github.com/stretchr/testify/require"
)

const testMailersName = "localmailer1"

func mailersSectionExpectation() map[string]models.MailersSections {
	initStructuredExpected()
	res := StructuredToMailersSectionMap()
	// Add individual entries
	for _, vs := range res {
		for _, v := range vs {
			key := v.Name
			res[key] = models.MailersSections{v}
		}
	}
	return res
}

func TestGetMailersSections(t *testing.T) {
	m := make(map[string]models.MailersSections)
	v, mailers, err := clientTest.GetMailersSections("")
	if err != nil {
		t.Error(err)
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	if len(mailers) != 1 {
		t.Errorf("mailers sections: found %d, expected 1", len(mailers))
	}
	m[mailers[0].Name] = models.MailersSections{mailers[0]}

	checkMailersSection(t, m)
}

func TestGetMailersSection(t *testing.T) {
	m := make(map[string]models.MailersSections)

	v, section, err := clientTest.GetMailersSection(testMailersName, "")
	if err != nil {
		t.Error(err)
	}
	m[testMailersName] = models.MailersSections{section}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}
	checkMailersSection(t, m)

	_, _, err = clientTest.GetMailersSection("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent mailers section")
	}
}

func checkMailersSection(t *testing.T, got map[string]models.MailersSections) {
	exp := mailersSectionExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Name == w.Name {
					require.True(t, g.MailersSectionBase.Equal(w.MailersSectionBase), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestCreateEditDeleteMailersSection(t *testing.T) {
	ms := &models.MailersSection{
		MailersSectionBase: models.MailersSectionBase{
			Name:    "newMailer",
			Timeout: misc.ParseTimeout("30s"),
		},
	}

	err := clientTest.CreateMailersSection(ms, "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	v, created, err := clientTest.GetMailersSection("newMailer", "")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(created, ms) {
		fmt.Printf("Created MailersSection: %v\n", created)
		fmt.Printf("Given MailersSection: %+v\n", ms)
		t.Error("Created MailersSection not equal to given MailersSection")
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	err = clientTest.CreateMailersSection(ms, "", version)
	if err == nil {
		t.Error("Should throw error MailersSection already exists")
		version++
	}

	// Modify the section.
	ms.Timeout = misc.ParseTimeout("40s")
	err = clientTest.EditMailersSection("newMailer", ms, "", version)
	if err != nil {
		t.Errorf("EditMailerSection: %v", err)
	} else {
		version++
	}

	// Check if the modification was effective.
	v, created, err = clientTest.GetMailersSection("newMailer", "")
	if err != nil {
		t.Error(err)
	}
	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}
	if *created.Timeout != *ms.Timeout {
		t.Errorf("MailersSection timeout was not modified: got %d, expected %d", *created.Timeout, *ms.Timeout)
	}

	// Delete the section.
	err = clientTest.DeleteMailersSection("newMailer", "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = clientTest.DeleteMailersSection("newMailer", "", 999999)
	if err != nil {
		if confErr, ok := err.(*configuration.ConfError); ok {
			if !confErr.Is(configuration.ErrVersionMismatch) {
				t.Error("Should throw configuration.ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw configuration.ErrVersionMismatch error")
		}
	}

	_, _, err = clientTest.GetMailersSection("newMailer", "")
	if err == nil {
		t.Error("DeleteMailersSection failed: newMailer still exists")
	}

	err = clientTest.DeleteMailersSection("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existent MailersSection")
		version++
	}
}
