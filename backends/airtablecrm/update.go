package airtablecrm

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// UpdateItem in the crm
func (c *Client) UpdateItem(ctx context.Context, i *crm.Item, updateFields map[string]interface{}) error {
	internalItem, ok := i.Internal.(*Item)
	if !ok {
		return fmt.Errorf("bad item")
	}
	return c.client.UpdateRecord(c.tableName, internalItem.airtableID, updateFields, nil, true)
}
