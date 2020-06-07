package googlesheet

import (
	"reflect"
	"testing"
)

func TestStructToMap(t *testing.T) {
	tests := []struct {
		item   interface{}
		wanted map[string]string
	}{
		{
			item: &struct {
				Name  string
				Other string
			}{
				Name:  "My Name",
				Other: "Other Value",
			},
			wanted: map[string]string{
				"Name":  "My Name",
				"Other": "Other Value",
			},
		},
	}

	for i, test := range tests {
		if mapValue := structToMap(test.item); !reflect.DeepEqual(mapValue, test.wanted) {
			t.Errorf("Test %d failed", i)
		}
	}
}
