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
	"fmt"
	"reflect"
	"testing"

	parser "github.com/haproxytech/config-parser/v4"
	"github.com/stretchr/testify/assert"

	"github.com/haproxytech/client-native/v4/misc"
	"github.com/haproxytech/client-native/v4/models"
)

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
		want1   models.Acls
		wantErr bool
	}{
		{
			name: "Should return acl list in one fcgi-app",
			args: args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: ""},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "invalid_src", Criterion: "src", Index: misc.Int64P(0), Value: "0.0.0.0/7 224.0.0.0/3"},
				&models.ACL{ACLName: "invalid_src", Criterion: "src_port", Index: misc.Int64P(1), Value: "0:1023"},
				&models.ACL{ACLName: "local_dst", Criterion: "hdr(host)", Index: misc.Int64P(2), Value: "-i localhost"},
			},
			wantErr: false,
		},
		{
			name: "Should return acl list in one fcgi-app by acl name",
			args: args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: "", aclName: []string{"invalid_src"}},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "invalid_src", Criterion: "src", Index: misc.Int64P(0), Value: "0.0.0.0/7 224.0.0.0/3"},
				&models.ACL{ACLName: "invalid_src", Criterion: "src_port", Index: misc.Int64P(1), Value: "0:1023"},
			},
			wantErr: false,
		},
		{
			name: "Should return acl list in one fcgi-app by acl name 2",
			args: args{parentType: string(parser.FCGIApp), parentName: "test", transactionID: "", aclName: []string{"local_dst"}},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "local_dst", Criterion: "hdr(host)", Index: misc.Int64P(2), Value: "-i localhost"},
			},
			wantErr: false,
		},
		{
			name: "Should return acl list in one frontend",
			args: args{parentType: string(parser.Frontends), parentName: "test", transactionID: ""},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "invalid_src", Criterion: "src", Index: misc.Int64P(0), Value: "0.0.0.0/7 224.0.0.0/3"},
				&models.ACL{ACLName: "invalid_src", Criterion: "src_port", Index: misc.Int64P(1), Value: "0:1023"},
				&models.ACL{ACLName: "local_dst", Criterion: "hdr(host)", Index: misc.Int64P(2), Value: "-i localhost"},
			},
			wantErr: false,
		},
		{
			name: "Should return acl list in one frontend by acl name",
			args: args{parentType: string(parser.Frontends), parentName: "test", transactionID: "", aclName: []string{"invalid_src"}},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "invalid_src", Criterion: "src", Index: misc.Int64P(0), Value: "0.0.0.0/7 224.0.0.0/3"},
				&models.ACL{ACLName: "invalid_src", Criterion: "src_port", Index: misc.Int64P(1), Value: "0:1023"},
			},
			wantErr: false,
		},
		{
			name: "Should return acl list in one frontend by acl name 2",
			args: args{parentType: string(parser.Frontends), parentName: "test", transactionID: "", aclName: []string{"local_dst"}},
			want: 1,
			want1: models.Acls{
				&models.ACL{ACLName: "local_dst", Criterion: "hdr(host)", Index: misc.Int64P(2), Value: "-i localhost"},
			},
			wantErr: false,
		},
		{
			name:    "Should return empty slice when no acl in given section",
			args:    args{parentType: string(parser.Backends), parentName: "test_2", transactionID: ""},
			want:    1,
			want1:   models.Acls{},
			wantErr: false,
		},
		{
			name:    "Should return an error when parentName doesn't exist",
			args:    args{parentType: string(parser.Backends), parentName: "not_exists", transactionID: ""},
			want:    1,
			want1:   nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := clientTest.GetACLs(tt.args.parentType, tt.args.parentName, tt.args.transactionID, tt.args.aclName...)
			if (err != nil) != tt.wantErr {
				t.Errorf("clientTest.GetACLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("clientTest.GetACLs() got = %v, want %v", got, tt.want)
			}
			if !assert.EqualValues(t, got1, tt.want1) {
				t.Errorf("clientTest.GetACLs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetACL(t *testing.T) {
	v, acl, err := clientTest.GetACL(0, "frontend", "test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if *acl.Index != 0 {
		t.Errorf("ACL ID not 0, %v found", *acl.Index)
	}
	if acl.ACLName != "invalid_src" {
		t.Errorf("%v: ACLName not invalid_src: %v", *acl.Index, acl.ACLName)
	}
	if acl.Value != "0.0.0.0/7 224.0.0.0/3" {
		t.Errorf("%v: Value not 0.0.0.0/7 224.0.0.0/3: %v", *acl.Index, acl.Value)
	}
	if acl.Criterion != "src" {
		t.Errorf("%v: Criterion not src: %v", *acl.Index, acl.Criterion)
	}

	_, err = acl.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = clientTest.GetACL(3, "backend", "test_2", "")
	if err == nil {
		t.Error("Should throw error, non existent ACL")
	}

	_, _, err = clientTest.GetACL(100, "frontend", "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent frontend and ACL")
	}

	_, _, err = clientTest.GetACL(100, "backend", "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent backend and ACL")
	}

	_, _, err = clientTest.GetACL(100, "fcgi-app", "fake", "")
	if err == nil {
		t.Error("Should throw error, non existent backend and ACL")
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

	err := clientTest.CreateACL("frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err := clientTest.GetACL(1, "frontend", "test", "")
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

	err = clientTest.EditACL(1, "frontend", "test", r, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ondiskR, err = clientTest.GetACL(1, "frontend", "test", "")
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
	err = clientTest.DeleteACL(3, "frontend", "test", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetACL(3, "frontend", "test", "")
	if err == nil {
		t.Error("DeleteACL failed, ACL Rule 3 still exists")
	}

	err = clientTest.DeleteACL(2, "backend", "test_2", "", version)
	if err == nil {
		t.Error("Should throw error, non existent ACL Rule")
		version++
	}
}
