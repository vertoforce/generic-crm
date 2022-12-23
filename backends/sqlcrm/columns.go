package sqlcrm

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	crm "github.com/vertoforce/generic-crm"
)

type TimeDefaultNow time.Time

// UpdateColumns would add new columns based on the example item.
// Currently not supported
func (c *Client) UpdateColumns(ctx context.Context, exampleItem crm.Item) error {
	columns, err := c.getColumns(ctx)
	if err != nil {
		return err
	}

exampleItemLoop:
	for key, value := range exampleItem.GetFields() {
		for columnName := range columns {
			if strings.ToLower(key) == strings.ToLower(columnName) {
				// We already have this
				// TODO: check type
				continue exampleItemLoop
			}
		}

		// We didn't find it, create the column
		// Create fieldType (default is text)
		fieldType := "TEXT(500)"
		switch reflect.TypeOf(value).Kind() {
		case reflect.String:
			// Check if it can be parsed as a date
			_, err := dateparse.ParseAny(value.(string))
			if err == nil {
				fieldType = "Datetime"
			}
		case reflect.Int, reflect.Int32, reflect.Int64:
			fieldType = "INT(11)"
		case reflect.Float32, reflect.Float64:
			fieldType = "FLOAT(11)"
		case reflect.TypeOf(time.Time{}).Kind():
			fieldType = "Datetime"
		case reflect.TypeOf(TimeDefaultNow{}).Kind():
			fieldType = "Datetime DEFAULT now()"
		}
		switch value.(type) {
		case int64, int32, int:
		case float64, float32:
		}
		a, err := c.DB.QueryxContext(ctx, fmt.Sprintf("ALTER TABLE `%s` ADD %s %s NULL DEFAULT NULL; ", c.Table, key, fieldType))
		if err != nil {
			return err
		}
		a.Close()
	}

	return nil
}

// getColumns returns a map of column name to it's type
func (c *Client) getColumns(ctx context.Context) (map[string]string, error) {
	columns, found := c.columnsCache.Get("columns")
	if found {
		return columns, nil
	}

	// Get table columns
	rows, err := c.DB.QueryxContext(ctx, fmt.Sprintf("SELECT COLUMN_NAME,COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'%s'", c.Table))
	if err != nil {
		return nil, err
	}

	columns = map[string]string{}
	for rows.Next() {
		row := map[string]interface{}{}
		rows.MapScan(row)
		columnName, ok := row["COLUMN_NAME"]
		if !ok {
			continue
		}
		columnType, ok := row["COLUMN_TYPE"]
		if !ok {
			continue
		}
		columns[fmt.Sprintf("%s", columnName)] = fmt.Sprintf("%s", columnType)
	}

	c.columnsCache.Set("columns", columns)

	return columns, nil
}
