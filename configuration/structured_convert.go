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
)

//nolint:nilnil
func namedResourceArrayToMapWithKey[T any](namedResource []*T, key string) (map[string]T, error) {
	res := make(map[string]T)
	for _, r := range namedResource {
		name, err := getKey(r, key)
		if err != nil {
			return nil, err
		}
		res[name] = *r
	}
	return res, nil
}

func namedResourceArrayToMap[T any](namedResource []*T) (map[string]T, error) {
	return namedResourceArrayToMapWithKey[T](namedResource, "Name")
}

// getKey returns the value of the 'Name' field from any struct or pointer to struct using reflection.
// Constraint: the struct must have an exportable 'Name' field
func getKey(obj interface{}, keyName string) (string, error) {
	value := reflect.ValueOf(obj)
	// If Pointer, first get the pointed value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", fmt.Errorf("object is not a struct")
	}
	nameField := value.FieldByName(keyName)
	if !nameField.IsValid() || !nameField.CanInterface() {
		return "", fmt.Errorf("object does not have an exportable 'Name' field")
	}
	name := nameField.Interface().(string)
	return name, nil
}
