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
	"github.com/stretchr/testify/require"
)

func TestStructuredUserList(t *testing.T) {
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
	require.NoError(t, err)

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
			require.Equal(t, tt.wantErr, (err != nil), "prepareClient error = %v, wantErr %v", err, tt.wantErr)

			// fetch all userlists tests
			_, lists, err := c.GetStructuredUserLists("")
			require.NoError(t, err)
			require.NotNil(t, lists, "No userlist configurations found, expected 2")
			require.Equal(t, 2, len(lists), "Expected 2 userlist configurations, found %v", len(lists))

			// fetch test, first userlist
			_, list, err := c.GetStructuredUserList("first", "")
			require.NoError(t, err)
			require.NotNil(t, list, "Userlist not found")
			require.Equal(t, "first", list.Name, "Userlist name %v returned, expected %v", list.Name, "first")

			require.Equal(t, 2, len(list.Groups), "Userlist %v, %d groups returned, expected 2", list.Name, len(list.Groups))
			_, ok := list.Groups["G1"]
			require.True(t, ok, "Userlist %v, missing user G1", list.Name)

			require.Equal(t, 3, len(list.Users), "Userlist %v, %d users returned, expected 2", list.Name, len(list.Users))
			_, ok = list.Users["scott"]
			require.True(t, ok, "Userlist %v, missing user scott", list.Name)

			// fetch test, second userlist
			_, list, err = c.GetStructuredUserList("second", "")
			require.NoError(t, err)
			require.NotNil(t, list, "Userlist not found")
			require.Equal(t, "second", list.Name, "Userlist name %v returned, expected %v", list.Name, "second")

			require.Equal(t, 2, len(list.Groups), "Userlist %v, %d groups returned, expected 2", list.Name, len(list.Groups))
			_, ok = list.Groups["one"]
			require.True(t, ok, "Userlist %v, missing user one", list.Name)

			require.Equal(t, 3, len(list.Users), "Userlist %v, %d users returned, expected 2", list.Name, len(list.Users))
			_, ok = list.Users["neo"]
			require.True(t, ok, "Userlist %v, missing user neo", list.Name)

			// fetch test, nonexisting userlist
			_, list, err = c.GetStructuredUserList("fake", "")
			require.Error(t, err, "No error thrown for a non existing userlist")
			require.Nil(t, list, "Non existing userlist found")

			// add tests
			add := models.Userlist{UserlistBase: models.UserlistBase{Name: "third"}}
			require.NoError(t, c.CreateStructuredUserList(&add, "", 1), "Adding a new userlist failed")
			add = models.Userlist{UserlistBase: models.UserlistBase{Name: "fourth"}}
			require.NoError(t, c.CreateStructuredUserList(&add, "", 2), "Adding a new userlist failed")
			require.Error(t, c.CreateStructuredUserList(&add, "", 2), "Duplicated list creation allowed")

			// delete tests
			require.NoError(t, c.DeleteUserList("first", "", 3), "Deleting an existing userlist failed")
			require.NoError(t, c.DeleteUserList("second", "", 4), "Deleting an existing userlist failed")
			require.Error(t, c.DeleteUserList("fake", "", 5), "Deleting an non existing userlist succeeded")
		})
	}
}
