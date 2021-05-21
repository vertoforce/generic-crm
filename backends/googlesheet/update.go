package googlesheet

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	crm "github.com/vertoforce/generic-crm"
)

// UpdateItem Updates an item's fields
func (c *Client) UpdateItem(ctx context.Context, i crm.Item, fields map[string]interface{}) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "UpdateItem")
	defer span.Finish()
	c.Lock()
	defer c.Unlock()

	// convert to google sheet item
	googleSheetItem, ok := i.(*Item)
	if !ok {
		return fmt.Errorf("invalid item")
	}
	for key, value := range fields {
		// Check the column number of this field
		columnNumber := c.getHeaderIndex(key)
		if columnNumber == -1 {
			continue
		}
		// Update it
		updateCell(c.Sheet, googleSheetItem.RowNumber, columnNumber, fmt.Sprintf("%v", value))
	}
	if c.WaitToSynchronize {
		return nil
	}
	return c.Synchronize()
}

// UpdateItemFromStruct Updates an item's fields using the struct names as headers and values as values
func (c *Client) UpdateItemFromStruct(ctx context.Context, i *Item, v interface{}) error {
	return c.UpdateItem(ctx, i, structToMap(v))
}
