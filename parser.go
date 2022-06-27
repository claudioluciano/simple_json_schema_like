package simplejsonschemalike

import (
	"fmt"
	"reflect"
	"unicode"
)

const timeStringKind = "time.Time"

/*Parse receive anything and try to parse to a JSON where with the type of the field.
Something like this
{
    "Int": "int",
    "Map": {
        "Key1": "string",
        "Key2": "int",
        "Key3": "bool"
    },
    "SliceOfString": "[string]",
    "String": "string",
    "Struct": {
        "Int": "int",
        "Map": {
            "Key1": "string",
            "Key2": "int",
            "Key3": "bool"
        },
        "SliceOfString": "[string]",
        "String": "string"
    },
    "StructPtr": {
        "Int": "int",
        "Map": {
            "Key1": "string",
            "Key2": "int",
            "Key3": "bool"
        },
        "SliceOfString": "[string]",
        "String": "string"
    }
}
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
	case reflect.Slice:
		return fmt.Sprintf("[%v]", parseSlice(v))
	default:
		vv := v.Type().String()

		return vv
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
		fieldName := structfield.Name

		if !exportedField(structfield.Name) {
			continue
		}

		m[fieldName] = parse(field)
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
	return v.Type().Elem()
}

func isTime(v reflect.Value) bool {
	t := v.Type().String()
	return t == timeStringKind || t == fmt.Sprintf("*%v", timeStringKind)
}

func exportedField(name string) bool {
	r := []rune(name)

	return unicode.IsUpper(r[0]) && unicode.IsLetter(r[0])
}
