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

package test

import (
	"testing"

	"github.com/haproxytech/client-native/v6/models"
)

func TestUserList(t *testing.T) {
	config := `# _version=1
global
	daemon

defaults
	maxconn 2000

userlist first
	group G1 users tiger,scott
	group G2 users xdb,scott
	user tiger password $6$k6y3o.eP$JlKBx9za9667qe4xHSwRv6J.C0/D7cV91
	user scott insecure-password elgato
	user xdb insecure-password hello

userlist second
	group one
	group two
	user neo password $6$k6y3o.eP$JlKBxxHSwRv6J.C0/D7cV91 groups one
	user thomas insecure-password elgato groups one,two
	user anderson insecure-password hello groups two
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
			name:              "userlists",
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

			// fetch all userlists tests
			_, lists, err := c.GetUserLists("")
			if err != nil {
				t.Error(err.Error())
			}
			if lists == nil {
				t.Errorf("No userlist configurations found, expected 2")
			}
			if len(lists) != 2 {
				t.Errorf("Expected 2 userlist configurations, found %v", len(lists))
			}

			// fetch test, first userlist
			_, list, err := c.GetUserList("first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if list == nil {
				t.Errorf("Userlist not found")
			}
			if list.Name != "first" {
				t.Errorf("Userlist name %v returned, expected %v", list.Name, "first")
			}

			// fetch test, second userlist
			_, list, err = c.GetUserList("second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if list == nil {
				t.Errorf("Userlist not found")
			}
			if list.Name != "second" {
				t.Errorf("Userlist name %v returned, expected %v", list.Name, "second")
			}

			// fetch test, nonexisting userlist
			_, list, err = c.GetUserList("fake", "")
			if err == nil {
				t.Errorf("No error thrown for a non existing userlist")
			}
			if list != nil {
				t.Errorf("Non existing userlist found")
			}

			// add tests
			add := models.Userlist{UserlistBase: models.UserlistBase{Name: "third"}}
			if c.CreateUserList(&add, "", 1) != nil {
				t.Errorf("Adding a new userlist failed")
			}
			add = models.Userlist{UserlistBase: models.UserlistBase{Name: "fourth"}}
			if c.CreateUserList(&add, "", 2) != nil {
				t.Errorf("Adding a new userlist failed")
			}
			if c.CreateUserList(&add, "", 2) == nil {
				t.Errorf("Duplicated list creation allowed")
			}

			// delete tests
			if c.DeleteUserList("first", "", 3) != nil {
				t.Errorf("Deleting an existing userlist failed")
			}
			if c.DeleteUserList("second", "", 4) != nil {
				t.Errorf("Deleting an existing userlist failed")
			}
			if c.DeleteUserList("fake", "", 5) == nil {
				t.Errorf("Deleting an non existing userlist succeeded")
			}
		})
	}
}
