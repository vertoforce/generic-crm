package airtablecrm

import (
	"context"

	crm "github.com/vertoforce/generic-crm"
)

// CreateItem in the crm
func (c *Client) CreateItem(ctx context.Context, i *crm.Item) error {
	request := struct {
		Fields map[string]interface{} `json:"fields"`
		// Typecast allows new options to be created with them multiple select field type
		Typecast bool `json:"typecast"`
	}{
		Fields:   i.Fields,
		Typecast: true,
	}
	return c.client.CreateRecord(c.tableName, request)
}
