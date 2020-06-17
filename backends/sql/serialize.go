package sqlcrm

import (
	"encoding/json"
	"fmt"
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
		switch value.(type) {
		case []string, map[string]interface{}:
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

// TODO: Add deserialize function
