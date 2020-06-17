package sqlcrm

import (
	"context"
	"fmt"
	"testing"

	crm "github.com/vertoforce/generic-crm"
)

func TestClient(t *testing.T) {
	c, err := NewCRM("root:pass@tcp(127.0.0.1:3306)/db", "test")
	if err != nil {
		t.Error(err)
		return
	}

	err = c.CreateItem(context.Background(), &crm.DefaultItem{
		Fields: map[string]interface{}{
			"name": "Name 1",
			"item": "item",
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
	if fmt.Sprintf("%s", item.GetFields()["name"]) != "Name 1" {
		t.Errorf("Did not get expected item")
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
	for _, item := range items {
		for key, value := range item.GetFields() {
			fmt.Printf("%s:%s ", key, value)
		}
		fmt.Println()
	}
}
