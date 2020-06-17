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

	items, err := c.GetItems(context.Background())
	if err != nil {
		t.Error(err)
		return
	}
	if len(items) == 0 {
		t.Errorf("Not enough items")
	}
	for _, item := range items {
		fmt.Printf("%s\n", item.GetFields()["name"])
	}
}
