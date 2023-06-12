package simplejsonschemalike

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type Person struct {
	Name string
	Age  int
}

type User struct {
	ID        int
	FirstName string
	LastName  string
}

type Event struct {
	ID         int
	Title      string
	Date       time.Time
	Guests     []Person
	Details    map[string]string
	PtrDate    *time.Time
	unexported string
}

func TestParse(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "bool",
			input:    true,
			expected: "bool",
		},
		{
			name:     "string",
			input:    "hello",
			expected: "string",
		},
		{
			name:     "int",
			input:    42,
			expected: "int",
		},
		{
			name: "struct",
			input: Person{
				Name: "Alice",
				Age:  30,
			},
			expected: `{"Age":"int","Name":"string"}`,
		},
		{
			name: "struct with one time.Time field filled and the other nil",
			input: Event{
				ID:      1,
				Title:   "Party",
				Date:    time.Now(),
				Guests:  []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 40}},
				Details: map[string]string{"location": "New York", "host": "John"},
			},
			expected: `{"Date":"date-time","Details": {"location": "string", "host": "string"},"Guests":"[{"Age":"int","Name":"string"}, {"Age":"int","Name":"string"}]","ID":"int","Title":"string","PtrDate":"date-time"}`,
		},
		{
			name:     "slice",
			input:    []int{1, 2, 3},
			expected: `["int","int","int"]`,
		},
		{
			name:     "slice with interface",
			input:    []interface{}{1, "hello", true},
			expected: `["int","string","bool"]`,
		},
		{
			name:     "empty slice",
			input:    []int{},
			expected: `["int"]`,
		},
		{
			name:     "map",
			input:    map[string]int{"one": 1, "two": 2},
			expected: `{"one":"int","two":"int"}`,
		},
		{
			name:     "empty map",
			input:    map[string]int{},
			expected: `{}`,
		},
		{
			name: "pointer to struct",
			input: &User{
				ID:        1,
				FirstName: "John",
				LastName:  "Doe",
			},
			expected: `{"FirstName":"string","ID":"int","LastName":"string"}`,
		},
		{
			name:     "nil pointer",
			input:    (*User)(nil),
			expected: `{"FirstName":"string","ID":"int","LastName":"string"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, _ := Parse(tt.input)

			expectedMap := make(map[string]interface{})
			expectedBytes, _ := json.Marshal(actual)
			_ = json.Unmarshal(expectedBytes, &expectedMap)

			actualMap := make(map[string]interface{})
			actualBytes, _ := json.Marshal(actual)
			_ = json.Unmarshal(actualBytes, &actualMap)

			if !reflect.DeepEqual(actualMap, expectedMap) {
				t.Errorf("Expected %v, but got %v", tt.expected, actual)
			}
		})
	}
}
