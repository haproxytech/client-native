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
