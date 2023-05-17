# simple_json_schema_like

## Parse

`Parse` is a Go function that receives anything and tries to parse it to a JSON schema-like representation that has the same type as the input.

### Usage

The function takes any input value and returns a value with the same type as the input, but represented in a JSON schema-like format. The returned value can be of any type, including struct, map, slice, array, pointer to a struct, slice, or map, or any primitive type.

To use the `Parse` function, simply call it with the input value:

```go
v := struct {
    Age  int
    Name string
}{
    Age:  30,
    Name: "John",
}

parsed := Parse(v)
```

The `parsed` variable will now contain a map with the same fields as the input struct:

```go
map[string]interface{}{
    "Age":  "int",
    "Name": "string",
}
```

### Output Format

The output format of the `Parse` function depends on the input type:

- Struct: the output is a map with the same fields as the input struct. The keys of the map are the field names, and the values are the field types represented as strings.

- Map: the output is a map with the same keys as the input map. The values of the map are the corresponding types of the input map values, represented as strings.

- Slice or Array: the output is a slice of the same length as the input slice or array. Each element of the output slice is the corresponding type of the input slice or array elements, represented as strings.

- Primitive types: the output is the type name as a string.

### Examples

Here are some examples of input and output for the `Parse` function:

```go
// Struct
v := struct {
    Age  int
    Name string
}{
    Age:  30,
    Name: "John",
}

parsed := Parse(v)
// parsed = map[string]interface{}{
//     "Age":  "int",
//     "Name": "string",
// }

// Map
v := map[string]int{
    "foo": 42,
    "bar": 69,
}

parsed := Parse(v)
// parsed = map[string]interface{}{
//     "foo": "int",
//     "bar": "int",
// }

// Slice
v := []interface{}{
    "foo",
    42,
}

parsed := Parse(v)
// parsed = []interface{}{
//     "string",
//     "int",
// }


// Primitive types
v := true

parsed := Parse(v)
// parsed = "bool"
```