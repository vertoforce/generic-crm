package airtablecrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/vertoforce/airtable-go"
	crm "github.com/vertoforce/generic-crm"
)

// GetItem from CRM
func (c *Client) GetItem(ctx context.Context, searchFields map[string]interface{}) (crm.Item, error) {
	items := []*Item{}
	searchFilters := []string{}
	for key, value := range searchFields {
		searchFilters = append(searchFilters, fmt.Sprintf("%s='%v'", key, value))
	}
	filterFormula := strings.Join(searchFilters, " AND ")
	c.client.ListRecords(c.tableName, &items, airtable.ListParameters{FilterByFormula: filterFormula})

	if len(items) == 0 {
		return nil, crm.ErrItemNotFound
	}

	return items[0], nil
}

// GetItems gets all items from this airtable crm
func (c *Client) GetItems(ctx context.Context) ([]crm.Item, error) {
	items := []*Item{}
	err := c.client.ListRecords(c.tableName, &items)
	if err != nil {
		return nil, err
	}

	ret := []crm.Item{}
	for _, item := range items {
		ret = append(ret, item)
	}

	return ret, nil
}
