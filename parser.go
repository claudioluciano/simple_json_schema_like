package simplejsonschemalike

import (
	"fmt"
	"reflect"

	"github.com/fatih/structtag"
)

const timeStringKind = "time.Time"

// Parse receives anything and tries to parse it to a JSON schema-like representation
// that has the same type as the input.
//
// For example, if the input is a struct, the output will be a map with the same fields
// and nested types as the input struct.
//
// If the input is a slice or array, the output will be a slice of the same length
// containing the JSON schema-like representation of each element in the input.
//
// If the input is a map, the output will be a map with the same keys and nested types
// as the input map.
//
// For primitive types (e.g. bool, string, int), the output will be the type name as a string.
//
// The function returns an interface{} value that can be asserted to the appropriate type
// after the function call.
func Parse(in interface{}) (interface{}, error) {
	v := reflect.ValueOf(in)
	return parse(v)
}

func parse(v reflect.Value) (interface{}, error) {
	switch v.Kind() {
	case reflect.Struct:
		if ok := isTime(v); ok {
			return "date-time", nil
		}

		return parseStruct(v)
	case reflect.Interface:
		return parseInterface(v)
	case reflect.Ptr:
		return parsePrt(v)
	case reflect.Map:
		return parseMap(v)
	case reflect.Slice, reflect.Array:
		return parseSlice(v)
	default:
		return v.Type().String(), nil
	}
}

func parseNil(v reflect.Value) interface{} {
	if ok := isTime(v); ok {
		return "date-time"
	}

	return v.Type().String()
}

func parseStruct(v reflect.Value) (map[string]interface{}, error) {
	m := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)

		if !field.IsExported() {
			continue
		}

		name := field.Name
		tagName, err := getNameFromTag(field)
		if err != nil {
			return nil, err
		}
		if tagName != "" {
			name = tagName
		}

		v, err := parse(v.Field(i))
		if err != nil {
			return nil, err
		}

		m[name] = v
	}

	return m, nil
}

func getNameFromTag(field reflect.StructField) (string, error) {
	tags, err := structtag.Parse(string(field.Tag))
	if err != nil {
		return "", err
	}

	jsonTag, err := tags.Get("json")
	if err != nil {
		// json tag does not exist
		return "", nil
	}

	return jsonTag.Name, nil
}

func parseInterface(v reflect.Value) (interface{}, error) {
	field := v.Elem()
	return parse(field)
}

func parsePrt(v reflect.Value) (interface{}, error) {
	if v.IsNil() {
		return parseNil(v), nil
	}

	field := v.Elem()
	return parse(field)
}

func parseMap(v reflect.Value) (interface{}, error) {
	m := make(map[string]interface{})

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)
		v, err := parse(value)
		if err != nil {
			return nil, err
		}

		m[key.String()] = v
	}

	return m, nil
}

func parseSlice(v reflect.Value) (interface{}, error) {
	if v.Len() == 0 {
		el := v.Type().Elem()
		elType := el.String()

		if el.Kind() == reflect.Interface {
			elType = "any"
		}

		return fmt.Sprintf("[%v]", elType), nil
	}

	var m []interface{}
	for i := 0; i < v.Len(); i++ {
		field := v.Index(i)
		t, err := parse(field)
		if err != nil {
			return nil, err
		}

		m = append(m, t)
	}

	return m, nil
}

func isTime(v reflect.Value) bool {
	t := v.Type().String()
	return t == timeStringKind || t == fmt.Sprintf("*%v", timeStringKind)
}
