package googlesheet

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	crm "github.com/vertoforce/generic-crm"
)

// RemoveItems removes items
//
// NOTE that there is a special case with this CRM that after you remove items,
// ALL other items you have cached anywhere become invalid (their row numbers have changed)
// So you must refresh your items
func (c *Client) RemoveItems(ctx context.Context, items ...crm.Item) error {
	// Convert to google sheet item
	internalItems := Items{}
	for _, item := range items {
		googleSheetItem, ok := item.(*Item)
		if !ok {
			return fmt.Errorf("invalid item")
		}
		internalItems = append(internalItems, googleSheetItem)
	}
	return c.RemoveItemsInternal(ctx, internalItems)
}

// RemoveItemsInternal from the CRM, NOTE - YOU MUST fetch the items again after removing items because the row numbers will change
func (c *Client) RemoveItemsInternal(ctx context.Context, items Items) error {
	c.Lock()
	defer c.Unlock()

	// First sort to be in order of row numbers
	sort.Sort(items)
	offset := 0
	for _, item := range items {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		// Set the row to be blank, and delete that row
		c.consumeQuota()
		err := c.Service.DeleteRows(c.Sheet, item.RowNumber+offset, item.RowNumber+offset+1)
		// Keep retrying if we get resource exhausted
		for err != nil && strings.Contains(err.Error(), "RESOURCE_EXHAUSTED") {
			time.Sleep(time.Second * 5)
			c.consumeQuota()
			err = c.Service.DeleteRows(c.Sheet, item.RowNumber+offset, item.RowNumber+offset+1)
		}
		if err != nil {
			return err
		}
		offset--
	}
	// We need to reload the sheet every time after a deletion unfortunately
	err := c.loadSheet()
	if err != nil {
		return fmt.Errorf("error reloading sheet:%s", err)
	}

	return nil
}
