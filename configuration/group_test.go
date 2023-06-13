// Copyright 2019 HAProxy Technologies
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
	"io/ioutil"
	"testing"

	"github.com/haproxytech/client-native/v5/models"
)

func generateGroupConfig(config string) (string, error) {
	f, err := ioutil.TempFile("/tmp", "group")
	if err != nil {
		return "", err
	}
	err = prepareTestFile(config, f.Name())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestGroup(t *testing.T) {
	config := `# _version=1
global
	daemon

defaults
	maxconn 2000

userlist first
	group G1 users tiger,scott
	group G2 users scott
	user tiger password $6$k6y3o.eP$JlKBx9za9667qe4xHSwRv6J.C0/D7cV91
	user scott insecure-password elgato

userlist second
	group one
	group two
	group three
	user neo password $6$k6y3o.eP$JlKBxxHSwRv6J.C0/D7cV91 groups one
	user thomas insecure-password white-rabbit groups one,two
	user anderson insecure-password hello groups two

userlist empty

userlist add_test
	group G3
	group G4

userlist replace_test
	group zion
	group io
	user trinity insecure-password the-one groups zion

userlist delete_test
	group virus
	user smith insecure-password cloning groups virus

	`
	configFile, err := generateConfig(config)
	if err != nil {
		t.Error(err.Error())
	}
	defer func() {
		_ = deleteTestFile(configFile)
	}()

	tests := []struct {
		name              string
		configurationFile string
		want              int64
		wantErr           bool
	}{
		{
			name:              "user",
			configurationFile: configFile,
			want:              1,
			wantErr:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := prepareClient(tt.configurationFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("prepareClient error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// fetch tests groups
			_, groups, err := c.GetGroups("first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if groups == nil {
				t.Errorf("No groups configurations found in userlist, expected 2")
			}
			if len(groups) != 2 {
				t.Errorf("Expected 2 groups in the userlist, found %v", len(groups))
			}

			_, groups, err = c.GetGroups("second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if groups == nil {
				t.Errorf("No groups configurations found in userlist, expected 3")
			}
			if len(groups) != 3 {
				t.Errorf("Expected 3 groups in the userlist, found %v", len(groups))
			}

			_, groups, err = c.GetGroups("empty", "")
			if err != nil {
				t.Error(err.Error())
			}
			if groups == nil {
				t.Errorf("Expected an empty result instead of nil")
			}
			if len(groups) != 0 {
				t.Errorf("Expected 0 groups in the userlist, found %v", len(groups))
			}

			_, groups, err = c.GetGroups("fake", "")
			if err == nil {
				t.Errorf("Fetching groups from a non existing userlist didn't throw an error")
			}
			if groups != nil {
				t.Errorf("Group found in userlist, expected 0")
			}

			// fetch test, single group - userlist first
			_, group, err := c.GetGroup("G1", "first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if group == nil {
				t.Errorf("Expected a group instead of nil")
			}
			if group.Name != "G1" {
				t.Errorf("Name %v returned, expected %v", group.Name, "G1")
			}
			if group.Users != "tiger,scott" {
				t.Errorf("Users %v returned, expected %v", group.Users, "tiger,scott")
			}

			_, group, err = c.GetGroup("G2", "first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if group == nil {
				t.Errorf("Expected a group instead of nil")
			}
			if group.Name != "G2" {
				t.Errorf("Name %v returned, expected %v", group.Name, "G2")
			}
			if group.Users != "scott" {
				t.Errorf("Users %v returned, expected %v", group.Users, "scott")
			}

			_, group, err = c.GetGroup("G3", "first", "")
			if group != nil {
				t.Errorf("Expected nil - instead got group")
			}

			_, group, err = c.GetGroup("G1000", "first", "")
			if group != nil {
				t.Errorf("Expected nil - instead got group")
			}

			// fetch test, single group - userlist second
			_, group, err = c.GetGroup("one", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if group == nil {
				t.Errorf("Expected a group instead of nil")
			}
			if group.Name != "one" {
				t.Errorf("Name %v returned, expected %v", group.Name, "one")
			}
			if group.Users != "" {
				t.Errorf("Users %v returned, expected %v", group.Users, "")
			}

			_, group, err = c.GetGroup("two", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if group == nil {
				t.Errorf("Expected a group instead of nil")
			}
			if group.Name != "two" {
				t.Errorf("Name %v returned, expected %v", group.Name, "two")
			}
			if group.Users != "" {
				t.Errorf("Users %v returned, expected %v", group.Users, "")
			}

			_, group, err = c.GetGroup("three", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if group == nil {
				t.Errorf("Expected a group instead of nil")
			}
			if group.Name != "three" {
				t.Errorf("Name %v returned, expected %v", group.Name, "three")
			}
			if group.Users != "" {
				t.Errorf("Users %v returned, expected %v", group.Users, "")
			}

			_, group, err = c.GetGroup("four", "second", "")
			if group != nil {
				t.Errorf("Expected nil - instead got group")
			}

			_, group, err = c.GetGroup("G1000", "second", "")
			if group != nil {
				t.Errorf("Expected nil - instead got group")
			}

			_, groups, err = c.GetGroups("empty", "")
			if err != nil {
				t.Error(err.Error())
			}
			if len(groups) != 0 {
				t.Errorf("Expected 0 groups in the userlist, found %v", len(groups))
			}

			// test add
			add := models.Group{
				Name:  "avengers",
				Users: "",
			}
			if c.CreateGroup("add_test", &add, "", 1) != nil {
				t.Errorf("Adding a new group request failed")
			}

			// test replace
			edit := models.Group{
				Name:  "zion",
				Users: "trinity",
			}
			if c.EditGroup("zion", "replace_test", &edit, "", 2) != nil {
				t.Errorf("Replacing an existing group request failed")
			}

			// test delete
			if c.DeleteGroup("G0", "delete_test", "", 3) == nil {
				t.Errorf("Attempt to delete a non existing group didn't throw and error")
			}
			if c.DeleteGroup("G1000", "delete_test", "", 4) == nil {
				t.Errorf("Attempt to delete a non existing group didn't throw an error")
			}
		})
	}
}
