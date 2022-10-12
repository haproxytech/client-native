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

package configuration

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/haproxytech/client-native/v5/misc"
	"github.com/haproxytech/client-native/v5/models"
)

const testMailersName = "localmailer1"

func TestGetMailersSections(t *testing.T) {
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

	if mailers[0].Name != testMailersName {
		t.Errorf("mailers name: found '%s', expected '%s'", mailers[0].Name, testMailersName)
	}
}

func TestGetMailersSection(t *testing.T) {
	v, section, err := clientTest.GetMailersSection(testMailersName, "")
	if err != nil {
		t.Error(err)
	}

	if v != version {
		t.Errorf("found version %d, expected %d", v, version)
	}

	if section.Name != testMailersName {
		t.Errorf("mailers name: found '%s', expected '%s'", section.Name, testMailersName)
	}

	timeout := *misc.ParseTimeout("15s")
	if *section.Timeout != timeout {
		t.Errorf("mailers timeout: found %d, expected %d", *section.Timeout, timeout)
	}

	_, err = section.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	_, _, err = clientTest.GetMailersSection("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existent mailers section")
	}
}

func TestCreateEditDeleteMailersSection(t *testing.T) {
	ms := &models.MailersSection{
		Name:    "newMailer",
		Timeout: misc.ParseTimeout("30s"),
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
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
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
