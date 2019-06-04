
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

package misc

import (
	"reflect"
)

// ObjInArray returns true if struct in list y has field named identifier with value value
func ObjInArray(value string, y []interface{}, identifier string) bool {
	for _, b := range y {
		objValue := reflect.ValueOf(b).Elem().FieldByName(identifier).String()
		if objValue == value {
			return true
		}
	}
	return false
}

// GetObjByField returns struct from list l if it has field named identifier with value value
func GetObjByField(l []interface{}, identifier string, value string) interface{} {
	for _, b := range l {
		objValue := reflect.ValueOf(b).Elem().FieldByName(identifier).String()
		if objValue == value {
			return b
		}
	}
	return nil
}
