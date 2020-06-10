package crm_test

import (
	"context"
	"os"
	"testing"

	crm "github.com/vertoforce/generic-crm"
	"github.com/vertoforce/generic-crm/backends/airtablecrm"
	"github.com/vertoforce/generic-crm/backends/googlesheet"
)

var crms = []crm.CRM{
	&googlesheet.Client{},
	&airtablecrm.Client{},
}

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

	a, err := airtablecrm.New(os.Getenv("AIRTABLE_API_KEY"), os.Getenv("AIRTABLE_BASE_ID"), "Testing")
	if err != nil {
		t.Error(err)
		return
	}

	// Test each backend individually
	backends := []crm.CRM{
		crm.CRM(co),
		crm.CRM(a),
	}

	for _, b := range backends {
		// Create an item
		err = b.CreateItem(context.Background(), &crm.DefaultItem{
			Fields: map[string]interface{}{
				"Name": "test",
				"Item": "test2",
			},
		})
		if err != nil {
			t.Error(err)
			return
		}

		// Get specific item
		item, err := b.GetItem(context.Background(), map[string]interface{}{"Name": "test"})
		if err != nil {
			t.Error(err)
			return
		}
		if item.GetFields()["Name"] != "test" {
			t.Errorf("wrong item")
			return
		}

		// Get all items
		items, err := b.GetItems(context.Background())
		if err != nil {
			t.Error(err)
			return
		}

		// Delete our created item
		var toDelete crm.Item
		for _, item := range items {
			if item.GetFields()["Name"] == "test" {
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
