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
	rows, err := c.db.Queryx(fmt.Sprintf("SELECT * FROM %s", strings.ReplaceAll(pq.QuoteIdentifier(c.table), "\"", "")))
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
