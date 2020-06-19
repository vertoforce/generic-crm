package sqlcrm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// serializeFields Converts an item's fields to how they will be stored in the crm.
// It basically serialize hard to store fields to something sql can store
//
// Ex: converting []string to csv, etc
func serializeFields(fields map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for key, value := range fields {
		// Convert the value if it needs to be changed
		var valueP interface{}
		switch reflect.TypeOf(value).Kind() {
		case reflect.TypeOf(time.Time{}).Kind():
			valueP = value
		case reflect.Slice, reflect.Array, reflect.Map, reflect.Struct:
			// Convert to json
			json, _ := json.Marshal(value)
			valueP = fmt.Sprintf("\"%s\"", json)
		default:
			valueP = value
		}

		ret[key] = valueP
	}

	return ret
}

// deserializeFields will look for JSON fields and unmarshal them
func deserializeFields(fields map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for key, value := range fields {
		switch value.(type) {
		case []byte:
			// Try to unmarshal
			var newValue interface{}
			err := json.Unmarshal(value.([]byte), &newValue)
			if err != nil {
				// Just convert this to a string
				ret[key] = fmt.Sprintf("%s", value)
				continue
			}
			// We did it, use this new value
			ret[key] = value
		default:
			// Leave as is
			ret[key] = value
		}
	}
	return ret
}
