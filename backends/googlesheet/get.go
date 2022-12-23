package googlesheet

import (
	"context"
	"reflect"

	"github.com/opentracing/opentracing-go"
	crm "github.com/vertoforce/generic-crm"
	"golang.org/x/sync/errgroup"
)

// GetItem searches for an item based on field values, will return first item that matches
// It loops through all items searching for it (it's in memory anyway)
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetItemGoogleSheet")
	defer span.Finish()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	items := make(chan crm.Item)

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		return c.GetItems(ctx, items, searchValues)
	})

	// Return first found item
	for item := range items {
		return item, nil
	}

	err := errGroup.Wait()
	if err != nil {
		return nil, err
	}

	return nil, crm.ErrItemNotFound
}

// GetItems returns all the items in the sheet
func (c *Client) GetItems(ctx context.Context, items chan crm.Item, searchFields ...map[string]interface{}) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "GetItemsGoogleSheet")
	// TODO: Reload sheet if we haven't in some time (to make sure we got the latest updates)

	// Build map of column number to desired value from the searchFields
	// This is so we can be more efficient in searching each row
	searchValuesRowBased := map[int]interface{}{}
	for _, sF := range searchFields {
		for key, value := range sF {
			// Find column number for this field
			for colNum, header := range c.Headers {
				if header == key {
					searchValuesRowBased[colNum] = value
					break
				}
			}
		}
	}

	defer span.Finish()

	numItems := c.NumItems(ctx)
itemLoop:
	for r := 1; r < numItems+1; r++ {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		row := c.Sheet.Rows[r]

		// Check if this item matches
		for colNum, value := range searchValuesRowBased {
			if colNum >= len(row) {
				// Bad row, doesn't match
				continue itemLoop
			}
			if !reflect.DeepEqual(row[colNum].Value, value) {
				// This didn't match
				continue itemLoop
			}
		}

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

		select {
		case <-ctx.Done():
			return ctx.Err()
		case items <- item:
		}
	}

	return nil
}

// NumItems Gets the number of items in the sheet.  If headers are enabled, it does NOT count the first row
// It finds the first empty row to mark the end of the items
func (c *Client) NumItems(ctx context.Context) int {
	span, _ := opentracing.StartSpanFromContext(ctx, "GetNumItems")
	span.SetTag("result", 0)
	defer span.Finish()
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
	result := len(c.Sheet.Rows) - 1
	span.SetTag("result", result)
	return result
}
