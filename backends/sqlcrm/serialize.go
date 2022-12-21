package sqlcrm

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// serializeFields Converts an item's fields to how they will be stored in the crm.
// It basically serializes hard to store fields to something sql can store
func serializeFields(fields map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for key, value := range fields {
		// Convert the value if it needs to be changed
		var valueP interface{}
		switch t := reflect.TypeOf(value); {
		case t == nil:
			continue
		case t.Kind() == reflect.TypeOf(time.Time{}).Kind():
			valueP = value
		case (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) && t.Elem().Kind() == reflect.Uint8:
			// This is a []byte, a special case to just use this raw value
			// Just use this raw value
			valueP = value
		case t.Kind() == reflect.Slice, t.Kind() == reflect.Array, t.Kind() == reflect.Map, t.Kind() == reflect.Struct:
			// Convert to json
			json, _ := json.Marshal(map[string]interface{}{"value": value})
			valueP = fmt.Sprintf("%s", json)
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
		case []byte, string:
			// Use raw value by default
			ret[key] = fmt.Sprintf("%s", value)

			// Try to unmarshal
			var newValue interface{}
			j := fmt.Sprintf("%s", value)
			if len(j) == 0 {
				break
			}
			if j[0] != '{' { // This isn't JSON, just use raw value
				break
			}
			err := json.Unmarshal([]byte(j), &newValue)
			if err != nil {
				// This only happens when this just happened to start with {
				// Use the raw value
				break
			}
			// We did it, extract the value from this
			if p, ok := newValue.(map[string]interface{}); ok {
				if v, ok := p["value"]; ok {
					ret[key] = v
				}
			}
		default:
			// Leave as is
			ret[key] = value
		}
	}
	return ret
}
