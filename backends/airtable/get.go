package airtable

import (
	"context"
	"crm/backends/backend"
)

// GetItems gets all items from this airtable crm
func (c *Client) GetItems(ctx context.Context) ([]*backend.Item, error) {
	items := []struct {
		ID     string
		Fields map[string]interface{}
	}{}
	err := c.client.ListRecords(c.tableName, &items)
	if err != nil {
		return nil, err
	}

	ret := []*backend.Item{}
	for _, item := range items {
		ret = append(ret, &backend.Item{
			Fields: item.Fields,
			Internal: &Item{
				airtableID: item.ID,
			},
		})
	}

	return ret, nil
}
