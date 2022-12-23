package sqlcrm

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
)

// To perform this test a database must be set up with a table with the following columns
// name varchar
// item varchar
// test varchar
func TestClient(t *testing.T) {
	c, err := NewCRM("root:pass@tcp(10.0.0.99:3306)/db", "test")
	if err != nil {
		t.Error(err)
		return
	}

	// Try creating a new column
	err = c.UpdateColumns(context.Background(), &crm.DefaultItem{Fields: map[string]interface{}{
		"TestColumn":  "TEST",
		"TestColumn2": "TEST",
		"name":        "name",
		"item":        "item",
		"test":        "",
		"date":        "1/2/2022 10:00:00",
	}})
	require.NoError(t, err)

	err = c.CreateItem(context.Background(), &crm.DefaultItem{
		Fields: map[string]interface{}{
			"name": "Name 1",
			"item": "item",
			"test": []string{"1", "2"},
			"date": "1/2/2022 10:00:00 Z",
		},
	})
	require.NoError(t, err)

	num, err := c.Len(context.Background())
	require.NoError(t, err)
	require.GreaterOrEqual(t, num, uint64(1))

	// Try to get that specific item
	item, err := c.GetItem(context.Background(), map[string]interface{}{
		"name": "Name 1",
	})
	require.NoError(t, err)
	itemFields := item.GetFields()
	require.Equal(t, "Name 1", itemFields["name"])
	require.Equal(t, time.Date(2022, 01, 02, 10, 00, 0, 0, time.UTC), itemFields["date"])
	if item.GetFields()["name"].(string) != "Name 1" {
		t.Errorf("Did not get expected item")
		return
	}

	// Update that item
	err = c.UpdateItem(context.Background(), item, map[string]interface{}{"name": "new name"})
	if err != nil {
		t.Error(err)
		return
	}

	itemsChan := make(chan crm.Item)
	go func() {
		defer close(itemsChan)
		err := c.GetItems(context.Background(), itemsChan)
		require.NoError(t, err)
	}()
	items := []crm.Item{}
	for item := range itemsChan {
		items = append(items, item)
	}

	if len(items) == 0 {
		t.Errorf("Not enough items")
	}
	// Make sure we updated this item
	if items[0].GetFields()["name"].(string) != "new name" {
		t.Errorf("Update did not work")
		return
	}
	for _, item := range items {
		for key, value := range item.GetFields() {
			fmt.Printf("%s:%s ", key, value)
		}
		fmt.Println()
	}

	// Remove all items
	err = c.RemoveItems(context.Background(), items...)
	if err != nil {
		t.Error(err)
		return
	}
}
