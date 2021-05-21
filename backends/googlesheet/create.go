package googlesheet

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	crm "github.com/vertoforce/generic-crm"
)

// CreateItem Creates an item from the generic backend type
func (c *Client) CreateItem(ctx context.Context, i crm.Item) error {
	return c.CreateItemFromMap(ctx, i.GetFields())
}

// CreateItemFromValues Creates an item by creating a new row
func (c *Client) CreateItemFromValues(ctx context.Context, values []string) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "CreateItemGoogleSheet")
	defer span.Finish()

	c.Lock()
	defer c.Unlock()
	rowNumberToPlaceAt := c.NumItems(ctx) + 1

	// Insert the new values
	for i, value := range values {
		updateCell(c.Sheet, rowNumberToPlaceAt, i, value)
	}
	if c.WaitToSynchronize {
		return nil
	}
	return c.Synchronize()
}

// CreateItemFromStruct Creates an items using the field names as header values
func (c *Client) CreateItemFromStruct(ctx context.Context, v interface{}) error {
	return c.CreateItemFromMap(ctx, structToMap(v))
}

// CreateItemFromMap Creates an item using the map of headers to value
func (c *Client) CreateItemFromMap(ctx context.Context, m map[string]interface{}) error {
	values := []string{}
	// Loop through every header to try and find it's value in the struct
	for _, header := range c.Headers {
		value := ""
		if foundValue, ok := m[header]; ok {
			value = fmt.Sprintf("%v", foundValue)
		}
		values = append(values, value)
	}

	return c.CreateItemFromValues(ctx, values)
}
