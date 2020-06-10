package googlesheet

import (
	"context"
	"fmt"
	"reflect"

	crm "github.com/vertoforce/generic-crm"
)

// GetItem searches for an item based on field values, will return first item that matches
// It loops through all items searching for it (it's in memory anyway)
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	subContext, cancel := context.WithCancel(ctx)
	defer cancel()
	items, err := c.GetItems(subContext)
	if err != nil {
		return nil, err
	}
itemLoop:
	for _, item := range items {
		itemMap := item.GetFields()
		// Check if this item matches
		for searchKey, searchValue := range searchValues {
			if foundValue, ok := itemMap[searchKey]; !ok || !reflect.DeepEqual(foundValue, searchValue) {
				// This didn't match
				continue itemLoop
			}
		}
		// Found it!
		return item, nil
	}

	return nil, fmt.Errorf("item not found")
}

// GetItems returns all the items in the sheet
func (c *Client) GetItems(ctx context.Context) ([]crm.Item, error) {
	items := []crm.Item{}
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

	return items, nil
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
