package airtablecrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/vertoforce/airtable-go"
	crm "github.com/vertoforce/generic-crm"
	"golang.org/x/sync/errgroup"
)

// GetItem from CRM
func (c *Client) GetItem(ctx context.Context, searchFields map[string]interface{}) (crm.Item, error) {
	// Get items with this filter
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	items := make(chan crm.Item)

	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		defer close(items)
		return c.GetItems(ctx, items, searchFields)
	})

	// Return first found item
	for item := range items {
		return item, nil
	}

	err := errGroup.Wait()
	if err != nil {
		return nil, err
	}

	return nil, crm.ErrItemNotFound
}

// GetItems gets all items from this airtable crm
func (c *Client) GetItems(ctx context.Context, items chan crm.Item, searchFields ...map[string]interface{}) error {
	itemsI := []*Item{}
	var err error
	if len(searchFields) > 0 && searchFields[0] != nil {
		// Get rows with a certain filter
		searchFilters := []string{}
		for key, value := range searchFields[0] {
			searchFilters = append(searchFilters, fmt.Sprintf("%s='%s'", key, strings.ReplaceAll(fmt.Sprintf("%v", value), "'", "\\'")))
		}
		filterFormula := strings.Join(searchFilters, " AND ")
		err = c.client.ListRecords(c.tableName, &itemsI, airtable.ListParameters{FilterByFormula: filterFormula})
	} else {
		err = c.client.ListRecords(c.tableName, &itemsI)
	}
	if err != nil {
		return err
	}

	// Send each item on
	for _, item := range itemsI {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case items <- item:
		}
	}

	return nil
}
