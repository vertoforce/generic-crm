package sqlcrm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

var (
	anyNumber           = regexp.MustCompile(`\d+(\.\d+)?`)
	sqlNumberQuantifier = regexp.MustCompile(`\(\d+\)`)
)

// serializeFields Converts an item's fields to how they will be stored in the crm.
// It basically serializes hard to store fields to something sql can store
// It will also reference the existing sql schema to know what datatypes the incoming data should be transformed in to.
func (c *Client) serializeFields(ctx context.Context, fields map[string]interface{}) (map[string]interface{}, error) {
	columns, err := c.GetColumns(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting columns: %w", err)
	}

	ret := map[string]interface{}{}
	for key, value := range fields {
		if key == "" {
			continue
		}
		sqlDataType := strings.ToLower(columns[key])
		sqlDataType = sqlNumberQuantifier.ReplaceAllString(sqlDataType, "")

		// Convert the value if it needs to be changed
		switch t := reflect.TypeOf(value); {
		case t == nil:
			continue
		case t.Kind() == reflect.String:
			switch sqlDataType {
			case "float", "int":
				// Find any number in this string and use it
				anyN := anyNumber.FindString(CleanNumber(value.(string)))
				if anyN == "" {
					ret[key] = 0
				}
				switch sqlDataType {
				case "float":
					ret[key], _ = strconv.ParseFloat(anyN, 64)
				case "int":
					ret[key], _ = strconv.ParseInt(value.(string), 10, 64)
				}
			case "datetime":
				date, err := dateparse.ParseAny(value.(string))
				if err != nil {
					// Cannot convert to date.  Cannot insert this column
					continue
				}
				ret[key] = date
			default:
				ret[key] = value
			}
		case t.Kind() == reflect.TypeOf(time.Time{}).Kind():
			ret[key] = value
		case (t.Kind() == reflect.Slice || t.Kind() == reflect.Array) && t.Elem().Kind() == reflect.Uint8:
			// This is a []byte, a special case to just use this raw value
			// Just use this raw value
			ret[key] = value
		case t.Kind() == reflect.Slice, t.Kind() == reflect.Array, t.Kind() == reflect.Map, t.Kind() == reflect.Struct:
			// Convert to json
			json, _ := json.Marshal(map[string]interface{}{"value": value})
			ret[key] = string(json)
		default:
			ret[key] = value
		}
	}

	return ret, nil
}

// deserializeFields will look for JSON fields and unmarshal them
func deserializeFields(fields map[string]interface{}) map[string]interface{} {
	ret := map[string]interface{}{}
	for key, value := range fields {
		switch value.(type) {
		case []byte, string:
			// Use raw value by default
			ret[key] = fmt.Sprintf("%s", value)

			if time, err := time.Parse("2006-01-02 15:04:05", ret[key].(string)); err == nil {
				ret[key] = time
				continue
			}

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
