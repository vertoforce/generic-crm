package googlesheet

import (
	"reflect"
	"regexp"

	"github.com/vertoforce/regexgrouphelp"
)

var spreadsheetIDFromURLRegex = regexp.MustCompile(`\/d\/(?P<id>.*?)\/`)

// GetSpreadsheetID Given a google spreadsheet URL, get the ID
func GetSpreadsheetID(url string) string {
	groups := regexgrouphelp.FindRegexGroups(spreadsheetIDFromURLRegex, url)
	if group, ok := groups["id"]; ok {
		if len(group) == 0 {
			return ""
		}
		return group[0]
	}

	return ""
}

func structToMap(v interface{}) map[string]interface{} {
	value := reflect.ValueOf(v).Elem()
	vType := reflect.TypeOf(v).Elem()

	updateMap := map[string]interface{}{}
	for i := 0; i < value.NumField(); i++ {
		updateMap[vType.Field(i).Name] = value.Field(i)
	}

	return updateMap
}
