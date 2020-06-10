package airtablecrm

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// RemoveItems from crm
func (c *Client) RemoveItems(ctx context.Context, items ...crm.Item) error {
	for _, item := range items {
		airtableItem, ok := item.(*Item)
		if !ok {
			return fmt.Errorf("Invalid item")
		}
		err := c.client.DestroyRecord(c.tableName, airtableItem.AirtableID)
		if err != nil {
			return err
		}
	}
	return nil
}
