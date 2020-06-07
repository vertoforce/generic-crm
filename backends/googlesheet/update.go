package googlesheet

import (
	"context"
	"fmt"

	"github.com/vertoforce/generic-crm/backends/crm"
)

// UpdateItem Updates an item in the crm
func (c *Client) UpdateItem(ctx context.Context, i *crm.Item, updateFields map[string]interface{}) error {
	item, ok := i.Internal.(*Item)
	if !ok {
		return fmt.Errorf("bad internal item")
	}
	return item.UpdateInternal(updateFields)
}

// UpdateInternal Updates an item's fields
func (i *Item) UpdateInternal(fields map[string]interface{}) error {
	for key, value := range fields {
		// Check the column number of this field
		columnNumber := i.client.getHeaderIndex(key)
		if columnNumber == -1 {
			continue
		}
		// Update it
		updateCell(i.client.Sheet, i.RowNumber, columnNumber, fmt.Sprintf("%v", value))
	}
	if i.client.WaitToSynchronize {
		return nil
	}
	return i.client.Synchronize()
}

// UpdateFromStruct Updates an item's fields using the struct names as headers and values as values
func (i *Item) UpdateFromStruct(v interface{}) error {
	return i.UpdateInternal(structToMap(v))
}
