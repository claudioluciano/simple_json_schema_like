package simplejsonschemalike

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	type inn struct {
		String        string
		Int           int
		SliceOfString []string
		Map           map[string]interface{}
	}

	type in struct {
		String        string
		Int           int
		SliceOfString []string
		Map           map[string]interface{}
		Time          *time.Time
		Struct        inn
		StructPtr     *inn
	}

	type args struct {
		in interface{}
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "",
			args: args{
				in: in{
					String:        "",
					Int:           0,
					SliceOfString: []string{},
					Map: map[string]interface{}{
						"Key1": "",
						"Key2": 0,
						"Key3": true,
					},
					Struct: inn{
						String:        "",
						Int:           0,
						SliceOfString: []string{},
						Map: map[string]interface{}{
							"Key1": "",
							"Key2": 0,
							"Key3": true,
						},
					},
					StructPtr: &inn{
						String:        "",
						Int:           0,
						SliceOfString: []string{},
						Map: map[string]interface{}{
							"Key1": "",
							"Key2": 0,
							"Key3": true,
						},
					},
				},
			},
			want: mockTest1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Parse(tt.args.in)
			if !equals(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func equals(value interface{}, jsonValue string) bool {
	b, _ := json.Marshal(value)

	fmt.Println(string(b))
	return string(b) == jsonValue
}
