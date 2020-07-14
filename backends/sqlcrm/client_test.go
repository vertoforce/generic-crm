package sqlcrm

import (
	"context"
	"fmt"
	"testing"

	crm "github.com/vertoforce/generic-crm"
)

// To perform this test a database must be set up with a table with the following columns
// name varchar
// item varchar
// test varchar
func TestClient(t *testing.T) {
	c, err := NewCRM("root:pass@tcp(127.0.0.1:3306)/db", "test")
	if err != nil {
		t.Error(err)
		return
	}

	// Try creating a new column
	err = c.UpdateColumns(context.Background(), &crm.DefaultItem{Fields: map[string]interface{}{
		"TestColumn":  "TEST",
		"TestColumn2": "TEST",
		"name":        "name",
	}})
	if err != nil {
		t.Error(err)
		return
	}

	err = c.CreateItem(context.Background(), &crm.DefaultItem{
		Fields: map[string]interface{}{
			"name": "Name 1",
			"item": "item",
			"test": []string{"1", "2"},
		},
	})
	if err != nil {
		t.Error(err)
		return
	}

	// Try to get that specific item
	item, err := c.GetItem(context.Background(), map[string]interface{}{
		"name": "Name 1",
	})
	if err != nil {
		t.Error(err)
		return
	}
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

	items, err := c.GetItems(context.Background())
	if err != nil {
		t.Error(err)
		return
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
