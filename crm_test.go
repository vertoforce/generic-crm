package crm

import (
	"context"
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

	c := backend.Backend(co)

	err = c.CreateItem(context.Background(), &backend.Item{
		Fields: map[string]interface{}{
			"Name": "test",
			"Item": "test2",
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

	err = c.RemoveItem(context.Background(), items[0])
	if err != nil {
		t.Error(err)
		return
	}
}
