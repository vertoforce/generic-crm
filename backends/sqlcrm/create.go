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
	query, values, err := c.generateCreateQueryFromItem(ctx, i)
	if err != nil {
		return fmt.Errorf("error generating query: %w", err)
	}
	r, err := c.DB.QueryxContext(ctx, query, values...)
	if err != nil {
		return err
	}
	r.Close()

	return nil
}

// generateCreateQueryFromItem converts the crm item to the sql query to insert it
func (c *Client) generateCreateQueryFromItem(ctx context.Context, i crm.Item) (query string, values []interface{}, err error) {
	fields, err := c.serializeFields(ctx, i.GetFields())
	if err != nil {
		return "", nil, fmt.Errorf("error serializing fields: %w", err)
	}

	fieldNames := []string{}
	values = []interface{}{}
	for key, value := range fields {
		fieldNames = append(fieldNames, fmt.Sprintf("`%s`", key))
		values = append(values, value)
	}

	query = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", c.Table, strings.Join(fieldNames, ","), strings.Join(repeat("?", len(fieldNames)), ","))

	return query, values, nil
}

func repeat(s string, count int) []string {
	ret := make([]string, count)
	for i := 0; i < count; i++ {
		ret[i] = s
	}
	return ret
}
