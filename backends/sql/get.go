package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// GetItems gets all items from this sql crm
func (c *Client) GetItems(ctx context.Context) ([]crm.Item, error) {
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", "")))
	if err != nil {
		return nil, err
	}

	ret := []crm.Item{}
	for rows.Next() {
		row := map[string]interface{}{}
		err = rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &Item{Fields: row})
	}

	return ret, nil
}

// GetItem gets a single item from this sql crm
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	whereQueries := []string{}
	for key, value := range serializeFields(searchValues) {
		whereQueries = append(whereQueries, fmt.Sprintf("%s=\"%s\"", key, value))

	}
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s", strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", ""), strings.Join(whereQueries, " AND ")))
	if err != nil {
		return nil, err
	}

	ret := []crm.Item{}
	for rows.Next() {
		row := map[string]interface{}{}
		err = rows.MapScan(row)
		if err != nil {
			return nil, err
		}
		ret = append(ret, &Item{Fields: row})
	}
	if len(ret) == 0 {
		return nil, crm.ErrItemNotFound
	}

	return ret[0], nil
}
