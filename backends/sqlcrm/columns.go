package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	crm "github.com/vertoforce/generic-crm"
)

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
		// Create fieldType (default is varchar)
		fieldType := "TEXT(500)"
		switch value.(type) {
		case int64, int32, int:
			fieldType = "INT(11)"
		case float64, float32:
			fieldType = "FLOAT(11)"
		}
		a, err := c.DB.QueryxContext(ctx, fmt.Sprintf("ALTER TABLE `%s` ADD %s %s NULL DEFAULT NULL; ", c.table, key, fieldType))
		if err != nil {
			return err
		}
		a.Close()
	}

	return nil
}

// getColumns returns a map of column name to it's type
func (c *Client) getColumns(ctx context.Context) (map[string]string, error) {
	// Get table columns
	rows, err := c.DB.QueryxContext(ctx, fmt.Sprintf("SELECT COLUMN_NAME,COLUMN_TYPE FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = N'%s'", c.table))
	if err != nil {
		return nil, err
	}

	columns := map[string]string{}
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

	return columns, nil
}
