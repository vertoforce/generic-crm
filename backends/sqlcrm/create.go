package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	crm "github.com/vertoforce/generic-crm"
)

// CreateItem in the crm
//
// Note that this will serialize special types stored in the database.
// So if you store a map or []string, it will serialize it to store it as JSON
func (c *Client) CreateItem(ctx context.Context, i crm.Item) error {
	query, values := c.generateCreateQueryFromItem(i)
	_, err := c.db.QueryxContext(ctx, query, values...)
	if err != nil {
		return err
	}

	return nil
}

// generateCreateQueryFromItem converts the crm item to the sql query to insert it
func (c *Client) generateCreateQueryFromItem(i crm.Item) (query string, values []interface{}) {
	fields := serializeFields(i.GetFields())

	fieldNames := []string{}
	values = []interface{}{}
	for key, value := range fields {
		fieldNames = append(fieldNames, fmt.Sprintf("%s", key))
		values = append(values, value)
	}

	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.table, strings.Join(fieldNames, ","), strings.Join(repeat("?", len(fieldNames)), ","))

	return query, values
}

func repeat(s string, count int) []string {
	ret := make([]string, count)
	for i := 0; i < count; i++ {
		ret[i] = s
	}
	return ret
}
