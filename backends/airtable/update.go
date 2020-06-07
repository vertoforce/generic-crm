package airtable

import (
	"context"
	"crm/backends/backend"
	"fmt"
)

// UpdateItem in the crm
func (c *Client) UpdateItem(ctx context.Context, i *backend.Item, updateFields map[string]interface{}) error {
	internalItem, ok := i.Internal.(*Item)
	if !ok {
		return fmt.Errorf("bad item")
	}
	return c.client.UpdateRecord(c.tableName, internalItem.airtableID, updateFields, nil)
}
