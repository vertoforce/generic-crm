package crm_test

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
	"github.com/vertoforce/generic-crm/backends/airtablecrm"
	"github.com/vertoforce/generic-crm/backends/googlesheet"
	"github.com/vertoforce/generic-crm/backends/sqlcrm"
)

var crms = []crm.CRM{
	&googlesheet.Client{},
	&airtablecrm.Client{},
	&sqlcrm.Client{},
}

func getTestCRMs() ([]crm.CRM, error) {
	co, err := googlesheet.New(context.Background(), &googlesheet.Config{
		GoogleClientSecretFile: os.Getenv("GoogleClientSecretFile"),
		SpreadsheetURL:         os.Getenv("TESTING_SPREADSHEET_URL"),
		SheetName:              "Sheet1",
	})
	if err != nil {
		return nil, err
	}

	a, err := airtablecrm.New(os.Getenv("AIRTABLE_API_KEY"), os.Getenv("AIRTABLE_BASE_ID"), "Testing")
	if err != nil {
		return nil, err
	}

	// Test each backend individually
	backends := []crm.CRM{
		crm.CRM(co),
		crm.CRM(a),
	}
	return backends, nil
}

func TestCRM(t *testing.T) {
	testCRMs, err := getTestCRMs()
	if err != nil {
		t.Error(err)
	}

	for _, testCRM := range testCRMs {
		// Create an item
		err = testCRM.CreateItem(context.Background(), &crm.DefaultItem{
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
		item, err := testCRM.GetItem(context.Background(), map[string]interface{}{"Name": "test"})
		if err != nil {
			t.Error(err)
			return
		}
		if item.GetFields()["Name"] != "test" {
			t.Errorf("wrong item")
			return
		}

		// Get all items
		items := make(chan crm.Item)
		go func() {
			defer close(items)
			err := testCRM.GetItems(context.Background(), items)
			require.NoError(t, err)
		}()

		// Delete our created item
		var toDelete crm.Item
		for item := range items {
			if item.GetFields()["Name"] == "test" {
				toDelete = item
				break
			}
		}
		err = testCRM.RemoveItems(context.Background(), toDelete)
		if err != nil {
			t.Error(err)
			return
		}
	}
}
