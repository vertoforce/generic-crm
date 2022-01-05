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
			// TODO: This needs to be cleaned up, improved, and tested.
			// Try to unmarshal
			var newValue interface{}
			j := fmt.Sprintf("%s", value)
			err := json.Unmarshal([]byte(j), &newValue)
			if err != nil || j[0] != '{' { // If we failed to marshal or this field doesn't start with {
				if newValue != nil {
					// Use whatever we unmarshalled
					// This is not JSON in the field, we unmarshalled to something else...
					// But I am assuming whatever it managed to unmarshal is an accurate representation of what's in the database.
					ret[key] = newValue
					continue
				}
				// Just convert this to a string
				ret[key] = fmt.Sprintf("%s", value)
				continue
			}
			// We did it, extract the value from this
			if p, ok := newValue.(map[string]interface{}); ok {
				if v, ok := p["value"]; ok {
					ret[key] = v
					continue
				}
			}
			// This only happens when this just happened to be serializable
			// Use the raw value
			ret[key] = value
		default:
			// Leave as is
			ret[key] = value
		}
	}
	return ret
}
