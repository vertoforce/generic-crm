package airtablecrm

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// RemoveItem from crm
func (c *Client) RemoveItem(ctx context.Context, i crm.Item) error {
	airtableItem, ok := i.(*Item)
	if !ok {
		return fmt.Errorf("Invalid item")
	}
	return c.client.DestroyRecord(c.tableName, airtableItem.AirtableID)
}
