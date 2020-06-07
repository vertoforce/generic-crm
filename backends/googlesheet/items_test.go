package googlesheet

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func getTestingClient() (*Client, error) {
	return New(context.Background(), &Config{
		GoogleClientSecretFile: os.Getenv("GoogleClientSecretFile"),
		SpreadsheetURL:         os.Getenv("TESTING_SPREADSHEET_URL"),
		SheetName:              "Sheet1",
	})
}

func TestItems(t *testing.T) {
	c, err := getTestingClient()
	if err != nil {
		t.Error(err)
		return
	}

	// Insert three items
	for i := 0; i < 3; i++ {
		err = c.CreateItemFromValues([]string{fmt.Sprintf("Name %d", i), fmt.Sprintf("Item %d", i)})
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Get the items
	items := c.GetItemsInternal(context.Background())
	count := 0
	var lastItem *Item
	toRemove := Items{}
	for _, item := range items {
		if item.ToMap()["Name"] != fmt.Sprintf("Name %d", count) {
			t.Errorf("Incorrect name for item %d", count)
			return
		}
		if count == 1 || count == 2 {
			toRemove = append(toRemove, item)
		}
		count++
		lastItem = item
	}
	if count != 3 {
		t.Errorf("Did not find correct number of items")
	}

	// Try updating the first item
	err = lastItem.UpdateInternal(map[string]interface{}{
		"Name": "New name",
	})
	if err != nil {
		t.Error(err)
	}

	// Check to make sure it updated
	_, err = c.GetItem(context.Background(), map[string]string{
		"Name": "New name",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// -- Test removing the middle two items --

	// Add new item
	err = c.CreateItemFromValues([]string{"Name last", "Item last"})
	if err != nil {
		t.Error(err)
		return
	}

	// Purposely make out of order
	toRemove.Swap(0, 1)
	err = c.RemoveItems(toRemove)
	if err != nil {
		t.Error(err)
	}

	// Fetch all remaining items and remove them
	toRemove = Items{}
	for _, item := range c.GetItemsInternal(context.Background()) {
		toRemove = append(toRemove, item)
	}

	err = c.RemoveItems(toRemove)
	if err != nil {
		t.Error(err)
	}
}
