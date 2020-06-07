package airtable

import (
	"context"

	"github.com/vertoforce/generic-crm/backends/crm"
)

// CreateItem in the crm
func (c *Client) CreateItem(ctx context.Context, i *crm.Item) error {
	request := struct {
		Fields map[string]interface{} `json:"fields"`
	}{
		Fields: i.Fields,
	}
	return c.client.CreateRecord(c.tableName, request)
}
