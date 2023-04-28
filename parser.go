package simplejsonschemalike

import (
	"fmt"
	"reflect"
)

const timeStringKind = "time.Time"

/*
Parse receives anything and tries to parse it to a JSON schema-like representation
that has the same type as the input.

For example, if the input is a struct, the output will be a map with the same fields
and nested types as the input struct.

If the input is a slice or array, the output will be a slice of the same length
containing the JSON schema-like representation of each element in the input.

If the input is a map, the output will be a map with the same keys and nested types
as the input map.

If the input is a pointer to a struct, slice, or map, the output will be a pointer to
the corresponding JSON schema-like representation.

For other types (e.g. bool, string, int), the output will be the type name as a string.

The function returns an interface{} value that can be asserted to the appropriate type
after the function call.
*/
func Parse(in interface{}) interface{} {
	v := reflect.ValueOf(in)

	return parse(v)
}

func parse(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Struct:
		if ok := isTime(v); ok {
			return "DateTime"
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
		return v.Type().String()
	}
}

func parseNil(v reflect.Value) interface{} {
	if ok := isTime(v); ok {
		return "DateTime"
	}

	return v.Type().String()
}

func parseStruct(v reflect.Value) map[string]interface{} {
	m := make(map[string]interface{})

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		structfield := v.Type().Field(i)

		if !structfield.IsExported() {
			continue
		}

		m[structfield.Name] = parse(field)
	}

	return m
}

func parseInterface(v reflect.Value) interface{} {
	field := v.Elem()
	return parse(field)
}

func parsePrt(v reflect.Value) interface{} {
	if v.IsNil() {
		return parseNil(v)
	}

	field := v.Elem()
	return parse(field)
}

func parseMap(v reflect.Value) interface{} {
	m := make(map[string]interface{})

	for _, key := range v.MapKeys() {
		value := v.MapIndex(key)

		m[key.String()] = parse(value)
	}

	return m
}

func parseSlice(v reflect.Value) interface{} {
	if v.Len() == 0 {
		el := v.Type().Elem()
		elType := el.String()

		if el.Kind() == reflect.Interface {
			elType = "any"
		}

		return fmt.Sprintf("[%v]", elType)
	}

	m := []interface{}{}
	for i := 0; i < v.Len(); i++ {
		field := v.Index(i)
		t := parse(field)

		m = append(m, t)
	}

	return m
}

func isTime(v reflect.Value) bool {
	t := v.Type().String()
	return t == timeStringKind || t == fmt.Sprintf("*%v", timeStringKind)
}
