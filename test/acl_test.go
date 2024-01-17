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
	_ "embed"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	parser "github.com/haproxytech/config-parser/v5"
	"github.com/stretchr/testify/require"

	"github.com/haproxytech/client-native/v5/configuration"
	"github.com/haproxytech/client-native/v5/models"
)

func aclExpectation() map[string]models.Acls {
	initStructuredExpected()
	res := StructuredToACLMap()
	// Add individual entries
	for k, vs := range res {
		for _, v := range vs {
			key := fmt.Sprintf("%s/%d", k, *v.Index)
			keyName := fmt.Sprintf("%s/%s", k, v.ACLName)
			if _, ok := res[keyName]; !ok {
				res[keyName] = models.Acls{}
			}
			res[keyName] = append(res[keyName], v)
			res[key] = models.Acls{v}
		}
	}
	return res
}

func TestClient_GetACLs(t *testing.T) {
	type args struct {
		parentType    string
		parentName    string
		transactionID string
		aclName       []string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "Should return acl list in one fcgi-app",
			args:    args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: ""},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Should return acl list in one fcgi-app by acl name",
			args:    args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: "", aclName: []string{"invalid_src"}},
			want:    1,
			wantErr: false,
		},
		{
			name: "Should return acl list in one fcgi-app by acl name 2",
			args: args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: "", aclName: []string{"local_dst"}},
			want: 1,

			wantErr: false,
		},
		{
			name:    "Should return acl list in one frontend",
			args:    args{parentType: string(parser.Frontends), parentName: "test", transactionID: ""},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Should return acl list in one frontend by acl name",
			args:    args{parentType: string(parser.Frontends), parentName: "test", transactionID: "", aclName: []string{"invalid_src"}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Should return acl list in one frontend by acl name 2",
			args:    args{parentType: string(parser.Frontends), parentName: "test", transactionID: "", aclName: []string{"local_dst"}},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Should return empty slice when no acl in given section",
			args:    args{parentType: string(parser.Backends), parentName: "test_2", transactionID: ""},
			want:    1,
			wantErr: false,
		},
		{
			name:    "Should return an error when parentName doesn't exist",
			args:    args{parentType: string(parser.Backends), parentName: "not_exists", transactionID: ""},
			want:    1,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			macls := make(map[string]models.Acls)
			got, got1, err := clientTest.GetACLs(tt.args.parentType, tt.args.parentName, tt.args.transactionID, tt.args.aclName...)
			key := fmt.Sprintf("%s/%s", tt.args.parentType, tt.args.parentName)
			if len(tt.args.aclName) > 0 {
				key += fmt.Sprintf("/%s", strings.Join(tt.args.aclName, ","))
			}
			if !tt.wantErr {
				macls[key] = got1
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("clientTest.GetACLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("clientTest.GetACLs() got = %v, want %v", got, tt.want)
			}
			checksACLs(t, macls)
		})
	}
}

func TestGetACL(t *testing.T) {
	macl := make(map[string]models.Acls)

	v, acl, err := clientTest.GetACL(0, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	macl["frontend/test/0"] = []*models.ACL{acl}
	checksACLs(t, macl)

	_, err = acl.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetACL(3, configuration.BackendParentName, "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent ACL")
	}
	_, _, err = clientTest.GetACL(100, configuration.FrontendParentName, "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent frontend and ACL")
	}

	_, _, err = clientTest.GetACL(100, configuration.BackendParentName, "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent backend and ACL")
	}

	_, _, err = clientTest.GetACL(100, configuration.FCGIAppParentName, "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent backend and ACL")
	}
}

func checksACLs(t *testing.T, got map[string]models.Acls) {
	exp := aclExpectation()
	for k, v := range got {
		want, ok := exp[k]
		require.True(t, ok, "k=%s", k)
		require.Equal(t, len(want), len(v), "k=%s", k)
		for _, g := range v {
			for _, w := range want {
				if *g.Index == *w.Index {
					require.True(t, g.Equal(*w), "k=%s - diff %v", k, cmp.Diff(*g, *w))
					break
				}
			}
		}
	}
}

func TestCreateEditDeleteACL(t *testing.T) {
	id := int64(1)

	// TestCreateACL
	r := &models.ACL{
		Index:     &id,
		ACLName:   "site_dead",
		Criterion: "nbsrv(dynamic)",
		Value:     "lt 2",
	}

	err := clientTest.CreateACL(configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetACL(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Created ACL rule: %v\n", ondiskR)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Created ACL rule not equal to given ACL rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditACL
	r = &models.ACL{
		Index:     &id,
		ACLName:   "site_dead",
		Criterion: "nbsrv(static)",
		Value:     "lt 4",
	}

	err = clientTest.EditACL(1, configuration.FrontendParentName, "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetACL(1, configuration.FrontendParentName, "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(ondiskR, r) {
		fmt.Printf("Edited ACL rule: %v\n", ondiskR)
		fmt.Printf("Given ACL rule: %v\n", r)
		t.Error("Edited ACL rule not equal to given ACL rule")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteACL
	err = clientTest.DeleteACL(4, configuration.FrontendParentName, "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetACL(4, configuration.FrontendParentName, "test", "")
	if err == nil {
		t.Error("DeleteACL failed, ACL Rule 4 still exists")
	}

	err = clientTest.DeleteACL(2, configuration.BackendParentName, "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent ACL Rule")
		version++
	}
}
