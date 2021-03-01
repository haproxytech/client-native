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

	"github.com/haproxytech/client-native/v2/misc"
	"github.com/haproxytech/client-native/v2/models"
)

func TestGetResolvers(t *testing.T) {
	v, resolvers, err := client.GetResolvers("")
	if err != nil {
		t.Error(err.Error())
	}

	if len(resolvers) != 1 {
		t.Errorf("%v resolvers returned, expected 1", len(resolvers))
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if resolvers[0].Name != "test" {
		t.Errorf("Expected only test, %v found", resolvers[0].Name)
	}
}

func TestGetResolver(t *testing.T) {
	v, l, err := client.GetResolver("test", "")
	if err != nil {
		t.Error(err.Error())
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	if l.Name != "test" {
		t.Errorf("Expected test resolver, %v found", l.Name)
	}

	_, err = l.MarshalBinary()
	if err != nil {
		t.Error(err.Error())
	}

	_, _, err = client.GetResolver("doesnotexist", "")
	if err == nil {
		t.Error("Should throw error, non existant resolvers section")
	}
}

func TestCreateEditDeleteResolver(t *testing.T) {
	f := &models.Resolver{
		Name:                "created_resolver",
		AcceptedPayloadSize: 4096,
		HoldNx:              misc.Int64P(10),
		HoldObsolete:        misc.Int64P(10),
		HoldOther:           misc.Int64P(10),
		HoldRefused:         misc.Int64P(10),
		HoldTimeout:         misc.Int64P(10),
		HoldValid:           misc.Int64P(100),
		ResolveRetries:      10,
		ParseResolvConf:     true,
		TimeoutResolve:      10,
		TimeoutRetry:        10,
	}
	err := client.CreateResolver(f, "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	v, resolver, err := client.GetResolver("created_resolver", "")
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(resolver, f) {
		fmt.Printf("Created resolver: %v\n", resolver)
		fmt.Printf("Given resolver: %v\n", f)
		t.Error("Created resolver not equal to given resolver")
	}

	if v != version {
		t.Errorf("Version %v returned, expected %v", v, version)
	}

	err = client.CreateResolver(f, "", version)
	if err == nil {
		t.Error("Should throw error resolver already exists")
		version++
	}

	err = client.DeleteResolver("created_resolver", "", version)
	if err != nil {
		t.Error(err.Error())
	} else {
		version++
	}

	if v, _ := client.GetVersion(""); v != version {
		t.Error("Version not incremented")
	}

	err = client.DeleteResolver("created_resolver", "", 999999)
	if err != nil {
		if confErr, ok := err.(*ConfError); ok {
			if confErr.Code() != ErrVersionMismatch {
				t.Error("Should throw ErrVersionMismatch error")
			}
		} else {
			t.Error("Should throw ErrVersionMismatch error")
		}
	}
	_, _, err = client.GetResolver("created_resolver", "")
	if err == nil {
		t.Error("DeleteResolver failed, resolver creatd_resolver still exists")
	}

	err = client.DeleteResolver("doesnotexist", "", version)
	if err == nil {
		t.Error("Should throw error, non existant resolver")
		version++
	}
}
