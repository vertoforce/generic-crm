package sqlcrm

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	crm "github.com/vertoforce/generic-crm"
)

// UpdateItem in the crm
func (c *Client) UpdateItem(ctx context.Context, i crm.Item, updateFields map[string]interface{}) error {
	whereQuery, whereValues := fieldsToSQLWhere(serializeFields(i.GetFields()))

	// Create set instructions
	sets := []string{}
	setValues := []interface{}{}
	for key, value := range serializeFields(updateFields) {
		sets = append(sets, fmt.Sprintf("%s=?", key))
		setValues = append(setValues, value)
	}
	setQuery := strings.Join(sets, ",")

	_, err := c.db.QueryxContext(ctx, fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", ""),
		setQuery,
		whereQuery,
	), append(setValues, whereValues...)...)
	if err != nil {
		return err
	}

	return nil
}
