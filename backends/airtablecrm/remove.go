package airtablecrm

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// RemoveItem from crm
func (c *Client) RemoveItem(ctx context.Context, i *crm.Item) error {
	internalItem, ok := i.Internal.(*Item)
	if !ok {
		return fmt.Errorf("bad item")
	}
	return c.client.DestroyRecord(c.tableName, internalItem.airtableID)
}
