package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// GetItems gets all items from this sql crm
//
// Note that this will deserialize special types stored in the database.
// So if you store a map or []string, it will deserialize it back in to the appropriate type
func (c *Client) GetItems(ctx context.Context) ([]crm.Item, error) {
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", "")))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	// TODO: Change to prepared query to avoid sql injection
	whereQuery, whereValues := fieldsToSQLWhere(serializeFields(searchValues))
	rows, err := c.db.QueryxContext(ctx, fmt.Sprintf("SELECT * FROM %s WHERE %s",
		strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", ""),
		whereQuery,
	), whereValues...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
