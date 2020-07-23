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
	// Get items with this filter
	subContext, cancel := context.WithCancel(ctx)
	defer cancel()
	items, err := c.GetItems(subContext, searchFields)
	if err != nil {
		return nil, err
	}

	// Return first found item
	for item := range items {
		return item, nil
	}

	return nil, crm.ErrItemNotFound
}

// GetItems gets all items from this airtable crm
func (c *Client) GetItems(ctx context.Context, searchFields ...map[string]interface{}) (chan crm.Item, error) {
	ret := make(chan crm.Item)

	items := []*Item{}
	var err error
	if len(searchFields) > 0 && searchFields[0] != nil {
		// Get rows with a certain filter
		searchFilters := []string{}
		for key, value := range searchFields[0] {
			searchFilters = append(searchFilters, fmt.Sprintf("%s='%s'", key, strings.ReplaceAll(fmt.Sprintf("%v", value), "'", "\\'")))
		}
		filterFormula := strings.Join(searchFilters, " AND ")
		err = c.client.ListRecords(c.tableName, &items, airtable.ListParameters{FilterByFormula: filterFormula})
	} else {
		err = c.client.ListRecords(c.tableName, &items)
	}
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(ret)
		// Send each item on
		for _, item := range items {
			select {
			case <-ctx.Done():
				return
			case ret <- item:
			}
		}
	}()

	return ret, nil
}
