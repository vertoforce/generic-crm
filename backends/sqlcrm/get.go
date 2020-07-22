package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// GetItems gets all items from this sql crm
//
// Note that this will deserialize special types stored in the database.
// So if you store a map or []string, it will deserialize it back in to the appropriate type
func (c *Client) GetItems(ctx context.Context, searchFields ...map[string]interface{}) (chan crm.Item, error) {
	var rows *sqlx.Rows
	var err error
	if len(searchFields) > 0 && searchFields[0] != nil {
		// This is a where search, generate where query
		whereQuery, whereValues := fieldsToSQLWhere(serializeFields(searchFields[0]))
		rows, err = c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s",
			strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", ""),
			whereQuery,
		), whereValues...)
	} else {
		// Just query for all items
		rows, err = c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", "")))
	}
	if err != nil {
		return nil, err
	}

	ret := make(chan crm.Item)

	go func() {
		defer close(ret)
		defer rows.Close()
		// Get each item and add to array
		for rows.Next() {
			row := map[string]interface{}{}
			err = rows.MapScan(row)
			if err != nil {
				return
			}
			select {
			case <-ctx.Done():
				return
			case ret <- &Item{Fields: row}:
			}
		}
	}()

	return ret, nil
}

// GetItem gets a single item from this sql crm
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	oneTimeCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	items, err := c.GetItems(oneTimeCtx, searchValues)
	if err != nil {
		return nil, err
	}
	// Return the first found item
	for item := range items {
		return item, nil
	}

	// Return not found
	return nil, crm.ErrItemNotFound
}

// fieldToSQLWhere Converts a list of fields to a SQL WHERE query (just the part after the WHERE)
//
// EX: map[string]interface{}{"name": "test", "item": "item"} -> name="?" AND item="?"
func fieldsToSQLWhere(fields map[string]interface{}) (query string, values []interface{}) {
	whereQueries := []string{}
	whereValues := []interface{}{}
	for key, value := range serializeFields(fields) {
		whereQueries = append(whereQueries, fmt.Sprintf("%s=?", key))
		whereValues = append(whereValues, value)
	}
	return strings.Join(whereQueries, " AND "), whereValues
}
