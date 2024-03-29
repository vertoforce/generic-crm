package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// RemoveItems from the CRM
func (c *Client) RemoveItems(ctx context.Context, items ...crm.Item) error {
	for i, item := range items {
		err := c.RemoveItem(ctx, item)
		if err != nil {
			return fmt.Errorf("error deleting item number %d: %s", i, err)
		}
	}
	return nil
}

// RemoveItem from the CRM
func (c *Client) RemoveItem(ctx context.Context, item crm.Item) error {
	serializedFields, err := c.serializeFields(ctx, item.GetFields())
	if err != nil {
		return fmt.Errorf("error serializing fields: %w", err)
	}
	whereQuery, whereValues := fieldsToSQLWhere(serializedFields)
	r, err := c.DB.QueryContext(ctx, fmt.Sprintf("DELETE FROM %s WHERE %s",
		strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", ""),
		whereQuery,
	), whereValues...)
	defer r.Close()
	return err
}
