package airtable

import (
	"context"
	"crm/backends/backend"
	"fmt"
)

// RemoveItem from crm
func (c *Client) RemoveItem(ctx context.Context, i *backend.Item) error {
	internalItem, ok := i.Internal.(*Item)
	if !ok {
		return fmt.Errorf("bad item")
	}
	return c.client.DestroyRecord(c.tableName, internalItem.airtableID)
}
