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
	"fmt"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v5/models"
)

func userExpectation() map[string]models.Users {
	initStructuredExpected()
	res := StructuredToUserMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%s", k, v.Username)
			res[key] = models.Users{v}
		}
	}
	return res
}

func generateUserConfig(config string) (string, error) {
	f, err := os.CreateTemp("/tmp", "user")
	if err != nil {
		return "", err
	}
	err = prepareTestFile(config, f.Name())
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestUser(t *testing.T) {
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
	user neo password $6$k6y3o.eP$JlKBxxHSwRv6J.C0/D7cV91 groups one
	user thomas insecure-password white-rabbit groups one,two
	user anderson insecure-password hello groups two

userlist empty

userlist add_test
	group G3
	group G4
	user switch insecure-password not-like-this groups G3,G4

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
			m := make(map[string]models.Users)

			// fetch tests users
			_, users, err := c.GetUsers("first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if users == nil {
				t.Errorf("No users configurations found in userlist, expected 2")
			}
			if len(users) != 2 {
				t.Errorf("Expected 2 users in the userlist, found %v", len(users))
			}
			for _, v := range users {
				m[fmt.Sprintf("userlist/first/%s", v.Username)] = models.Users{v}
			}

			_, users, err = c.GetUsers("second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if users == nil {
				t.Errorf("No users configurations found in userlist, expected 3")
			}
			if len(users) != 3 {
				t.Errorf("Expected 3 users in the userlist, found %v", len(users))
			}
			for _, v := range users {
				m[fmt.Sprintf("userlist/second/%s", v.Username)] = models.Users{v}
			}
			checkUsers(t, m)

			_, users, err = c.GetUsers("empty", "")
			if err != nil {
				t.Error(err.Error())
			}
			if users == nil {
				t.Errorf("Expected an empty result instead of nil")
			}
			if len(users) != 0 {
				t.Errorf("Expected 0 users in the userlist, found %v", len(users))
			}

			_, users, err = c.GetUsers("fake", "")
			if err == nil {
				t.Errorf("Fetching users from a non existing userlist didn't throw an error")
			}
			if users != nil {
				t.Errorf("Users found in userlist, expected 0")
			}

			// fetch test, single user - userlist first
			_, user, err := c.GetUser("tiger", "first", "")
			clear(m)
			if err != nil {
				t.Error(err.Error())
			}
			if user == nil {
				t.Errorf("Expected an user instead of nil")
			}
			m["userlist/first/tiger"] = models.Users{user}

			_, user, err = c.GetUser("scott", "first", "")
			if err != nil {
				t.Error(err.Error())
			}
			if user == nil {
				t.Errorf("Expected an user instead of nil")
			}
			m["userlist/first/scott"] = models.Users{user}

			_, user, err = c.GetUser("dummy", "first", "")
			if user != nil {
				t.Errorf("Expected nil - instead got user")
			}

			_, user, err = c.GetUser("fake", "first", "")
			if user != nil {
				t.Errorf("Expected nil - instead got user")
			}

			// fetch test, single user - userlist second
			_, user, err = c.GetUser("neo", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if user == nil {
				t.Errorf("Expected an user instead of nil")
			}
			m["userlist/second/neo"] = models.Users{user}

			_, user, err = c.GetUser("thomas", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if user == nil {
				t.Errorf("Expected an user instead of nil")
			}
			m["userlist/second/thomas"] = models.Users{user}

			_, user, err = c.GetUser("anderson", "second", "")
			if err != nil {
				t.Error(err.Error())
			}
			if user == nil {
				t.Errorf("Expected an user instead of nil")
			}
			m["userlist/second/anderson"] = models.Users{user}

			_, user, err = c.GetUser("third", "second", "")
			if user != nil {
				t.Errorf("Expected nil - instead got user")
			}

			_, user, err = c.GetUser("fake", "second", "")
			if user != nil {
				t.Errorf("Expected nil - instead got user")
			}

			_, users, err = c.GetUsers("empty", "")
			if err != nil {
				t.Error(err.Error())
			}
			if len(users) != 0 {
				t.Errorf("Expected 0 users in the userlist, found %v", len(users))
			}

			checkUsers(t, m)

			// test add
			securePassword := false
			add := models.User{
				Username:       "morpheus",
				Password:       "dreams",
				SecurePassword: &securePassword,
				Groups:         "G3,G4",
			}
			if c.CreateUser("add_test", &add, "", 1) != nil {
				t.Errorf("Adding a new user request failed")
			}

			// test replace
			securePassword = true
			edit := models.User{
				Username:       "trinity",
				Password:       "$6$k6y3o.eP$JlKBxxHSwRv6J.C0/D7cV91",
				SecurePassword: &securePassword,
				Groups:         "zion,io",
			}
			if c.EditUser("trinity", "replace_test", &edit, "", 2) != nil {
				t.Errorf("Replacing an existing user request failed")
			}

			_, u, err := c.GetUser("trinity", "replace_test", "")
			require.NoError(t, err)

			assert.Equal(t, edit.Username, u.Username)
			assert.Equal(t, edit.Password, u.Password)
			assert.Equal(t, *edit.SecurePassword, *u.SecurePassword)
			assert.Equal(t, edit.Groups, u.Groups)

			// test delete
			if c.DeleteUser("trinity", "replace_test", "", 3) != nil {
				t.Errorf("Deleting an existing user request failed")
			}
			if c.DeleteUser("", "delete_test", "", 4) == nil {
				t.Errorf("Attempt to delete an empty user didn't throw an error")
			}
			if c.DeleteUser("fake", "delete_test", "", 4) == nil {
				t.Errorf("Attempt to delete an non existing user didn't throw an error")
			}
		})
	}
}

func checkUsers(t *testing.T, got map[string]models.Users) {
	exp := userExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if g.Username == w.Username {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}
