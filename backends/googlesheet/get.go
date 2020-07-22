package googlesheet

import (
	"context"
	"reflect"

	crm "github.com/vertoforce/generic-crm"
)

// GetItem searches for an item based on field values, will return first item that matches
// It loops through all items searching for it (it's in memory anyway)
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	subContext, cancel := context.WithCancel(ctx)
	defer cancel()
	items, err := c.GetItems(subContext, searchValues)
	if err != nil {
		return nil, err
	}
	for item := range items {
		return item, nil
	}

	return nil, crm.ErrItemNotFound
}

// GetItems returns all the items in the sheet
func (c *Client) GetItems(ctx context.Context, searchFields ...map[string]interface{}) (chan crm.Item, error) {
	// TODO: Reload sheet if we haven't in some time (to make sure we got the latest updates)

	items := make(chan crm.Item)

	go func() {
		defer close(items)

	itemLoop:
		for r := 1; r < c.NumItems()+1; r++ {
			select {
			case <-ctx.Done():
				return
			default:
			}

			row := c.Sheet.Rows[r]

			// Build item to send
			item := &Item{RowNumber: r, Fields: []string{}, client: c}
			// Fill in fields
			for i, col := range row {
				if i >= len(c.Headers) {
					// Fields beyond the headers, ignore these
					break
				}
				item.Fields = append(item.Fields, col.Value)
			}

			if len(searchFields) > 0 {
				itemMap := item.GetFields()
				// Check if this item matches
				for searchKey, searchValue := range searchFields[0] {
					if foundValue, ok := itemMap[searchKey]; !ok || !reflect.DeepEqual(foundValue, searchValue) {
						// This didn't match
						continue itemLoop
					}
				}
			}

			select {
			case <-ctx.Done():
				return
			case items <- item:
			}
		}
	}()

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
