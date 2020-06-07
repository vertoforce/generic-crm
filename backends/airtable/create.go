package airtable

import (
	"context"
	"crm/backends/backend"
)

// CreateItem in the crm
func (c *Client) CreateItem(ctx context.Context, i *backend.Item) error {
	request := struct {
		Fields map[string]interface{} `json:"fields"`
	}{
		Fields: i.Fields,
	}
	return c.client.CreateRecord(c.tableName, request)
}
