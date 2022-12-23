package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
	"golang.org/x/sync/errgroup"
)

// GetItems gets all items from this sql crm
//
// Note that this will deserialize special types stored in the database.
// So if you store a map or []string, it will deserialize it back in to the appropriate type
func (c *Client) GetItems(ctx context.Context, items chan crm.Item, searchFields ...map[string]interface{}) error {
	var rows *sqlx.Rows
	var err error
	if len(searchFields) > 0 && searchFields[0] != nil && len(searchFields[0]) > 0 {
		// This is a where search, generate where query
		serializedFields, err := c.serializeFields(ctx, searchFields[0])
		if err != nil {
			return fmt.Errorf("error serializing fields: %w", err)
		}
		whereQuery, whereValues := fieldsToSQLWhere(serializedFields)
		rows, err = c.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s",
			strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", ""),
			whereQuery,
		), whereValues...)
	} else {
		// Just query for all items
		rows, err = c.DB.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", "")))
	}
	if err != nil {
		return fmt.Errorf("error performing query: %w", err)
	}

	defer rows.Close()
	// Get each item and add to array
	for rows.Next() {
		row := map[string]interface{}{}
		err = rows.MapScan(row)
		if err != nil {
			return fmt.Errorf("error scanning row: %w", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case items <- &Item{Fields: row}:
		}
	}

	return nil
}

// GetItem gets a single item from this sql crm
func (c *Client) GetItem(ctx context.Context, searchValues map[string]interface{}) (crm.Item, error) {
	items := make(chan crm.Item)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errGroup, ctx := errgroup.WithContext(ctx)
	errGroup.Go(func() error {
		defer close(items)
		return c.GetItems(ctx, items, searchValues)
	})

	// Return the first found item
	for item := range items {
		return item, nil
	}

	err := errGroup.Wait()
	if err != nil {
		return nil, err
	}

	// Return not found
	return nil, crm.ErrItemNotFound
}

func (c *Client) Len(ctx context.Context) (uint64, error) {
	// Just query for all items
	row := c.DB.QueryRowxContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", "")))
	if err := row.Err(); err != nil {
		return 0, fmt.Errorf("error running query: %w", err)
	}

	var count uint64
	err := row.Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error scanning query: %w", err)
	}

	return count, nil
}

// fieldToSQLWhere Converts a list of fields to a SQL WHERE query (just the part after the WHERE).
// You should pass in pre-serialized fields.
//
// EX: map[string]interface{}{"name": "test", "item": "item"} -> name="?" AND item="?"
func fieldsToSQLWhere(fields map[string]interface{}) (query string, values []interface{}) {
	whereQueries := []string{}
	whereValues := []interface{}{}
	for key, value := range fields {
		whereQueries = append(whereQueries, fmt.Sprintf("%s=?", key))
		whereValues = append(whereValues, value)
	}
	return strings.Join(whereQueries, " AND "), whereValues
}
