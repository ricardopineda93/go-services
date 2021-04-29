package accountsrv

import (
	"reflect"
	"strings"
)

// Function takes in a Struct and returns a Map of that struct using the tags associated
// to the struct fields to create the keys of the Map
func structToMapByTag(item interface{}, tagName string) map[string]interface{} {

	// Declare the map that will be returned
	res := map[string]interface{}{}
	// If input item is null return the map
	if item == nil {
		return res
	}
	v := reflect.TypeOf(item)
	reflectValue := reflect.ValueOf(item)
	reflectValue = reflect.Indirect(reflectValue)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		tag := v.Field(i).Tag.Get(tagName)

		// remove omitEmpty
		omitEmpty := false
		if strings.HasSuffix(tag, "omitempty") {
			omitEmpty = true
			idx := strings.Index(tag, ",")
			if idx > 0 {
				tag = tag[:idx]
			} else {
				tag = ""
			}
		}

		field := reflectValue.Field(i).Interface()
		if tag != "" && tag != "-" {
			if v.Field(i).Type.Kind() == reflect.Struct {
				res[tag] = structToMapByTag(field, tagName)
			} else {
				if !(omitEmpty && reflectValue.Field(i).IsZero()) {
					res[tag] = field
				}
			}
		}
	}
	return res
}
