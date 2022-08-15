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

	"github.com/haproxytech/client-native/v4/models"
)

func TestGetMailerEntries(t *testing.T) {
	v, mailerEntries, err := clientTest.GetMailerEntries(testMailersName, "")
	if err != nil {
		t.Error(err)
		return
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if len(mailerEntries) != 2 {
		t.Errorf("found %d mailer entries, expected %d", len(mailerEntries), 2)
	}
}

func TestGetMailerEntry(t *testing.T) {
	v, mailer, err := clientTest.GetMailerEntry("smtp1", testMailersName, "")
	if err != nil {
		t.Error(err)
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if mailer.Name != "smtp1" {
		t.Errorf("found mailer name %s, expected %s", mailer.Name, "smtp1")
	}
	if mailer.Address != "10.0.10.1" {
		t.Errorf("found mailer address %s, expected %s", mailer.Address, "10.0.10.1")
	}
	if mailer.Port != 514 {
		t.Errorf("found mailer port %d, expected %d", mailer.Port, 514)
	}

	_, err = mailer.MarshalBinary()
	if err != nil {
		t.Error(err)
	}

	_, _, err = clientTest.GetMailerEntry("nonexistent", testMailersName, "")
	if err == nil {
		t.Error("Should throw error, non existent mailer entry")
	}
}

func TestCreateEditDeleteMailerEntry(t *testing.T) {
	me := &models.MailerEntry{
		Name:    "smtp3",
		Address: "10.0.10.3",
		Port:    514,
	}

	err := clientTest.CreateMailerEntry(testMailersName, me, "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	v, created, err := clientTest.GetMailerEntry("smtp3", testMailersName, "")
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(created, me) {
		fmt.Printf("Created mailerEntry: %+v\n", created)
		fmt.Printf("Given mailerEntry: %+v\n", me)
		t.Error("Created mailerEntry not equal to given mailerEntry")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = clientTest.CreateMailerEntry(testMailersName, me, "", version)
	if err == nil {
		t.Error("Should throw error mailerEntry already exists")
		version++
	}

	// Edit the mailer entry.
	me.Port = 1024
	err = clientTest.EditMailerEntry("smtp3", testMailersName, me, "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	// Check if the edition was effective.
	v, created, err = clientTest.GetMailerEntry("smtp3", testMailersName, "")
	if err != nil {
		t.Error(err)
	}
	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if !reflect.DeepEqual(created, me) {
		fmt.Printf("Edited mailerEntry: %+v\n", created)
		fmt.Printf("Given mailerEntry: %+v\n", me)
		t.Error("Edited mailerEntry not equal to given mailerEntry")
	}

	// Delete the entry.
	err = clientTest.DeleteMailerEntry("smtp3", testMailersName, "", version)
	if err != nil {
		t.Error(err)
	} else {
		version++
	}

	if v, _ = clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetMailerEntry("smtp3", testMailersName, "")
	if err == nil {
		t.Error("DeleteMailerEntry failed, mailer entry still exists")
	}

	err = clientTest.DeleteMailerEntry("smtp3", testMailersName, "", version)
	if err == nil {
		t.Error("Should throw error, non existent mailer entry")
		version++
	}
}
