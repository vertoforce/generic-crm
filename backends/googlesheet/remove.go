package googlesheet

import (
	"context"
	"crm/backends/backend"
	"sort"
)

// RemoveItem removes a single item
func (c *Client) RemoveItem(ctx context.Context, i *backend.Item) error {
	return c.RemoveItemInternal(i.Internal.(*Item))
}

// RemoveItemInternal removes a single item
func (c *Client) RemoveItemInternal(item *Item) error {
	return c.RemoveItems(Items{item})
}

// RemoveItems from the CRM, NOTE - YOU MUST fetch the items again after removing items because the row numbers will change
func (c *Client) RemoveItems(items Items) error {
	// First sort to be in order of row numbers
	sort.Sort(items)
	offset := 0
	for _, item := range items {
		// Set the row to be blank, and delete that row
		c.consumeQuota()
		err := c.Service.DeleteRows(c.Sheet, item.RowNumber+offset, item.RowNumber+offset+1)
		if err != nil {
			return err
		}
		offset--
	}

	return nil
}
