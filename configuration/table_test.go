// Copyright 2023 HAProxy Technologies
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

	"github.com/haproxytech/client-native/v5/models"
)

func TestGetTables(t *testing.T) { //nolint:gocognit,gocyclo
	v, tables, err := clientTest.GetTables("mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if len(tables) != 2 {
		t.Errorf("%v tables returned, expected 2", len(tables))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	for _, table := range tables {
		if table.Name != "t1" && table.Name != "t2" {
			t.Errorf("table.Name not t1 or t2: %s", table.Name)
		}
		if table.Type != "string" {
			t.Errorf("%v: table.Type not string: %v", table.Name, table.Type)
		}
		if table.TypeLen == nil {
			t.Errorf("%s: table.TypeLen is nil", table.Name)
		} else if *table.TypeLen != 1000 {
			t.Errorf("%s: table.TypeLen not 1000: %v", table.Name, *table.TypeLen)
		}
		if table.Size == "" {
			t.Errorf("%s: table.Size is nil", table.Name)
		} else if table.Size != "200k" {
			t.Errorf("%v: table.Size not 200k: %v", table.Name, table.Size)
		}
		if table.Expire == nil {
			t.Errorf("%s: table.Expire is nil", table.Name)
		} else if *table.Expire != "5m" {
			t.Errorf("%v: table.Expire not 5m: %v", table.Name, *table.Expire)
		}
		if !table.NoPurge {
			t.Errorf("%s: table.NoPurge is not set", table.Name)
		}
		if table.Store != "gpc0,conn_rate(30s)" && table.Store != "gpc0,gpc1,conn_rate(30s)" {
			t.Errorf("%s: t.Store[0] not gpc0,conn_rate(30s) or gpc0,gpc1,conn_rate(30s): %s", table.Name, table.Store)
		}
	}
}

func TestTable(t *testing.T) {
	v, table, err := clientTest.GetTable("t1", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if table.Name != "t1" {
		t.Errorf("table.Name not t1 : %s", table.Name)
	}
	if table.Type != "string" {
		t.Errorf("%v: table.Type not string: %v", table.Name, table.Type)
	}
	if table.TypeLen == nil {
		t.Errorf("%s: table.TypeLen is nil", table.Name)
	} else if *table.TypeLen != 1000 {
		t.Errorf("%s: table.TypeLen not 1000: %v", table.Name, *table.TypeLen)
	}
	if table.Size == "" {
		t.Errorf("%s: table.Size is nil", table.Name)
	} else if table.Size != "200k" {
		t.Errorf("%v: table.Size not 200k: %v", table.Name, table.Size)
	}
	if table.Expire == nil {
		t.Errorf("%s: table.Expire is nil", table.Name)
	} else if *table.Expire != "5m" {
		t.Errorf("%v: table.Expire not 5m: %v", table.Name, *table.Expire)
	}
	if !table.NoPurge {
		t.Errorf("%s: table.NoPurge is not set", table.Name)
	}
	if table.Store != "gpc0,conn_rate(30s)" {
		t.Errorf("%s: t.Store[0] not gpc0,conn_rate(30s): %s", table.Name, table.Store)
	}
}

func TestCreateEditDeleteTable(t *testing.T) {
	// TestCreateTable
	table := &models.Table{
		Name: "t3",
		Type: "string",
		Size: "200k",
	}

	err := clientTest.CreateTable("mycluster", table, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ta, err := clientTest.GetTable("t3", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(table, ta) {
		fmt.Printf("Created table: %v\n", table)
		fmt.Printf("Given table: %v\n", ta)
		t.Error("Created table not equal to given table")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestEditTable
	table = &models.Table{
		Name:    "t3",
		Type:    "string",
		Size:    "200k",
		NoPurge: true,
		Store:   "gpc0,conn_rate(30s)",
	}

	err = clientTest.EditTable("t3", "mycluster", table, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, ta, err = clientTest.GetTable("t3", "mycluster", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(table, ta) {
		fmt.Printf("Created table: %v\n", table)
		fmt.Printf("Given table: %v\n", ta)
		t.Error("Created table not equal to given table")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	// TestDeleteTable
	err = clientTest.DeleteTable("t3", "mycluster", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := clientTest.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	_, _, err = clientTest.GetTable("t3", "mycluster", "")
	if err == nil {
		t.Error("DeleteTable failed, table t3 still exists")
	}

	err = clientTest.DeleteTable("t4", "mycluster", "", version)
	if err == nil {
		t.Error("Should throw error, non existent stick rule")
		version++
	}
}
