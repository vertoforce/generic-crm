package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	crm "github.com/vertoforce/generic-crm"
)

// CreateItem in the crm
func (c *Client) CreateItem(ctx context.Context, i crm.Item) error {
	query, values := c.generateCreateQueryFromItem(i)
	statement, err := c.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = statement.ExecContext(ctx, values...)
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
