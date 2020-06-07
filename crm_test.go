package crm

import (
	"context"
	"crm/backends/airtable"
	"crm/backends/backend"
	"crm/backends/googlesheet"
	"os"
	"testing"
)

func TestCRM(t *testing.T) {
	co, err := googlesheet.New(context.Background(), &googlesheet.Config{
		GoogleClientSecretFile: os.Getenv("GoogleClientSecretFile"),
		SpreadsheetURL:         os.Getenv("TESTING_SPREADSHEET_URL"),
		SheetName:              "Sheet1",
	})
	if err != nil {
		t.Error(err)
		return
	}

	a, err := airtable.New(os.Getenv("AIRTABLE_API_KEY"), os.Getenv("AIRTABLE_BASE_ID"), "Testing")
	if err != nil {
		t.Error(err)
		return
	}

	// Test each backend individually
	backends := []backend.Backend{
		backend.Backend(co),
		backend.Backend(a),
	}

	for _, b := range backends {
		err = b.CreateItem(context.Background(), &backend.Item{
			Fields: map[string]interface{}{
				"Name": "test",
				"Item": "test2",
			},
		})
		if err != nil {
			t.Error(err)
			return
		}

		items, err := b.GetItems(context.Background())
		if err != nil {
			t.Error(err)
			return
		}
		var toDelete *backend.Item
		for _, item := range items {
			if item.Fields["Name"] == "test" {
				toDelete = item
				break
			}
		}

		err = b.RemoveItem(context.Background(), toDelete)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
