package airtablecrm

import (
	"context"

	"github.com/vertoforce/generic-crm/backends/crm"
)

// GetItems gets all items from this airtable crm
func (c *Client) GetItems(ctx context.Context) ([]*crm.Item, error) {
	items := []struct {
		ID     string
		Fields map[string]interface{}
	}{}
	err := c.client.ListRecords(c.tableName, &items)
	if err != nil {
		return nil, err
	}

	ret := []*crm.Item{}
	for _, item := range items {
		ret = append(ret, &crm.Item{
			Fields: item.Fields,
			Internal: &Item{
				airtableID: item.ID,
			},
		})
	}

	return ret, nil
}
