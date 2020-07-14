package airtablecrm

import (
	"fmt"

	crm "github.com/vertoforce/generic-crm"
	"golang.org/x/net/context"
)

// UpdateColumns would add new columns to the airtable based on the example item.
// Currently not supported
func (c *Client) UpdateColumns(ctx context.Context, exampleItem crm.Item) error {
	return fmt.Errorf("not supported")
}
