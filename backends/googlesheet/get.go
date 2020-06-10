package googlesheet

import (
	"context"
	"fmt"

	crm "github.com/vertoforce/generic-crm"
)

// GetItem searches for an item based on field values, will return first item that matches
func (c *Client) GetItem(ctx context.Context, searchValues map[string]string) (*Item, error) {
	subContext, cancel := context.WithCancel(ctx)
itemLoop:
	for _, item := range c.GetItemsInternal(subContext) {
		itemMap := item.ToMap()
		// Check if this item matches
		for key, value := range searchValues {
			if val, ok := itemMap[key]; !ok || val != value {
				// This didn't match
				continue itemLoop
			}
		}
		// Found it!
		cancel()
		return item, nil
	}

	cancel()
	return nil, fmt.Errorf("item not found")
}

// GetItems gets items in this crm converted from the internal item type
func (c *Client) GetItems(ctx context.Context) ([]*crm.Item, error) {
	items := c.GetItemsInternal(ctx)

	ret := []*crm.Item{}
	for _, item := range items {
		newMap := map[string]interface{}{}
		for key, value := range item.ToMap() {
			newMap[key] = value
		}
		ret = append(ret, &crm.Item{
			Fields:   newMap,
			Internal: item,
		})
	}

	return ret, nil
}

// GetItemsInternal returns a channel of all the items in the sheet
func (c *Client) GetItemsInternal(ctx context.Context) []*Item {
	items := []*Item{}
	for r := 1; r < c.NumItems()+1; r++ {
		row := c.Sheet.Rows[r]

		// Build item to send
		item := &Item{
			RowNumber: r,
			Fields:    []string{},
			client:    c,
		}
		// Fill in fields
		for i, col := range row {
			if i >= len(c.Headers) {
				// Fields beyond the headers, ignore these
				break
			}
			item.Fields = append(item.Fields, col.Value)
		}

		items = append(items, item)
	}

	return items
}

// NumItems Gets the number of items in the sheet.  If headers are enabled, it does NOT count the first row
// It finds the first empty row to mark the end of the items
func (c *Client) NumItems() int {
	if len(c.Sheet.Rows) <= 1 {
		return 0
	}

rowLoop:
	for r := 1; r < len(c.Sheet.Rows); r++ {
		row := c.Sheet.Rows[r]
		if len(row) == 0 {
			// Row is empty for sure
			return r - 1
		}
		// Check if every value is empty
		for _, col := range row {
			if col.Value != "" {
				// Row has data
				continue rowLoop
			}
		}
		// Empty row, return
		return r - 1
	}

	// Reached end of sheet, so it's just the length of the sheet (minus headers)
	return len(c.Sheet.Rows) - 1
}
