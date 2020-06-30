package sqlcrm

import (
	"reflect"
	"testing"
)

func TestSerialize(t *testing.T) {
	test := map[string]interface{}{
		"String": "I am a string",
		// Currently not supported
		// "Bytes":  []byte{0, 1, 2},
		"Array": []interface{}{"One", "two", "three"},
		"Struct": struct {
			One string
			Two []byte
		}{One: "One", Two: []byte{0, 1, 2}},
		"Number": 5,
		"float":  float64(5.3),
	}

	// Try to serialize and deserialize and make sure the result is the same as the original
	processed := deserializeFields(serializeFields(test))
	if !reflect.DeepEqual(processed, test) {
		t.Errorf("test failed")
	}

}
