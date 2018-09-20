package misc

import (
	"reflect"
)

func ObjInArray(value string, y []interface{}, identifier string) bool {
    for _, b := range(y) {
    	objValue := reflect.ValueOf(b).Elem().FieldByName(identifier).String()
        if objValue == value {
            return true
        }
    }
    return false
}

func GetObjByField(l []interface{}, identifier string, value string) interface{} {
	for _, b := range(l) {
    	objValue := reflect.ValueOf(b).Elem().FieldByName(identifier).String()
        if objValue == value {
            return b
        }
    } 
    return nil
}
