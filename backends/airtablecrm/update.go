package airtablecrm

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// UpdateItem in the crm
func (c *Client) UpdateItem(ctx context.Context, i crm.Item, updateFields map[string]interface{}) error {
	airtableItem, ok := i.(*Item)
	if !ok {
		return fmt.Errorf("Invalid item")
	}
	return c.client.UpdateRecord(c.tableName, airtableItem.AirtableID, updateFields, nil, true)
}
