package v2v3

import (
	"cmp"
	"encoding/json"
	"errors"
	"reflect"
	"slices"

	"github.com/haproxytech/client-native/v6/misc"
)

// V2Tov3 converts a Structured type from v2 to v3
// For example
// v3f, err := V2Tov3[v2.Frontend, models.Frontend](&v2s.Frontend)
// The conversion allows to have an evolving type v3:
// - extra fields would be ignored in the conversion, keeping the v2 fields converted
func V2Tov3[TV2, TV3 any](resource *TV2, skip ...string) (*TV3, error) {
	original, err := json.Marshal(resource)
	if err != nil {
		return nil, err
	}
	originalMap := make(map[string]interface{})
	err = json.Unmarshal(original, &originalMap)
	if err != nil {
		return nil, err
	}
	for _, s := range skip {
		delete(originalMap, s)
	}
	versionj, err := json.Marshal(originalMap)
	if err != nil {
		return nil, err
	}
	var otherversion TV3
	err = json.Unmarshal(versionj, &otherversion)
	if err != nil {
		return nil, err
	}
	return &otherversion, nil
}

func ListV2ToV3[TV2, TV3 any](listV2 []*TV2) ([]*TV3, error) {
	SortListByIndex(listV2)
	listV3 := make([]*TV3, len(listV2))
	for i, v := range listV2 {
		resourceV3, err := V2Tov3[TV2, TV3](v)
		if err != nil {
			return nil, err
		}
		listV3[i] = resourceV3
	}
	return listV3, nil
}

func NamedResourceArrayToMap[T any](namedResource []*T) (map[string]T, error) {
	return NamedResourceArrayToMapWithKey[T](namedResource, "Name")
}

//nolint:nilnil
func NamedResourceArrayToMapWithKey[T any](namedResource []*T, key string) (map[string]T, error) {
	if len(namedResource) == 0 {
		return nil, nil
	}
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

// getKey returns the value of the 'Name' field from any struct or pointer to struct using reflection.
// Constraint: the struct must have an exportable 'Name' field
func getKey(obj interface{}, keyName string) (string, error) {
	value := reflect.ValueOf(obj)
	// If Pointer, first get the pointed value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return "", errors.New("object is not a struct")
	}
	nameField := value.FieldByName(keyName)
	if !nameField.IsValid() || !nameField.CanInterface() {
		return "", errors.New("object does not have an exportable 'Name' field")
	}
	name := nameField.Interface().(string)
	return name, nil
}

func SortListByIndex[T any](list []*T) {
	slices.SortFunc(list,
		func(a, b *T) int {
			ia, _ := getIndex(a)
			ib, _ := getIndex(b)
			if ia == nil || ib == nil {
				return -1
			}
			return cmp.Compare(*ia, *ib)
		})
}

func getIndex(obj interface{}) (*int64, error) {
	value := reflect.ValueOf(obj)
	// If Pointer, first get the pointed value
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return misc.Ptr[int64](-1), errors.New("object is not a struct")
	}
	nameField := value.FieldByName("Index")
	if !nameField.IsValid() || !nameField.CanInterface() {
		return misc.Ptr[int64](-1), errors.New("object does not have an exportable 'Index' field")
	}
	index := nameField.Interface().(*int64)
	return index, nil
}

func NamedResourceArrayV2ToMapV3[TV2, TV3 any](namedResource []*TV2, options ...optionV2ToV3) (map[string]TV3, error) {
	var res map[string]TV3
	if len(namedResource) == 0 {
		return res, nil
	}
	keyName := "Name"
	if len(options) > 0 {
		keyName = options[0].keyName
	}
	res = make(map[string]TV3)
	for _, r := range namedResource {
		name, err := getKey(r, keyName)
		if err != nil {
			return nil, err
		}

		// version change
		otherversion, err := V2Tov3[TV2, TV3](r)
		if err != nil {
			return nil, err
		}

		res[name] = *otherversion
	}
	return res, nil
}

type optionV2ToV3 struct {
	keyName string
}
