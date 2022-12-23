package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// UpdateItem in the crm
func (c *Client) UpdateItem(ctx context.Context, i crm.Item, updateFields map[string]interface{}) error {
	serializedFields, err := c.serializeFields(ctx, i.GetFields())
	if err != nil {
		return fmt.Errorf("failed to serialize fields: %w", err)
	}
	whereQuery, whereValues := fieldsToSQLWhere(serializedFields)

	// Create set instructions
	sets := []string{}
	setValues := []interface{}{}
	serializedFields, err = c.serializeFields(ctx, updateFields)
	if err != nil {
		return fmt.Errorf("failed to serialize update fields: %w", err)
	}
	for key, value := range serializedFields {
		sets = append(sets, fmt.Sprintf("`%s`=?", key))
		setValues = append(setValues, value)
	}
	setQuery := strings.Join(sets, ",")

	r, err := c.DB.QueryxContext(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", ""),
		setQuery,
		whereQuery,
	), append(setValues, whereValues...)...)
	if err != nil {
		return err
	}
	defer r.Close()

	return nil
}
