package googlesheet

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/opentracing/opentracing-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmot"
)

func getTestingClient() (*Client, error) {
	return New(context.Background(), &Config{
		GoogleClientSecretFile: os.Getenv("GoogleClientSecretFile"),
		SpreadsheetURL:         os.Getenv("TESTING_SPREADSHEET_URL"),
		SheetName:              "Sheet1",
	})
}

func TestItems(t *testing.T) {
	ctx := context.Background()

	opentracing.SetGlobalTracer(apmot.New())
	defer apm.DefaultTracer.Flush(nil)
	c, err := getTestingClient()
	if err != nil {
		t.Error(err)
		return
	}

	// Insert three items
	for i := 0; i < 3; i++ {
		err = c.CreateItemFromValues(ctx, []string{fmt.Sprintf("Name %d", i), fmt.Sprintf("Item %d", i)})
		if err != nil {
			t.Error(err)
			return
		}
	}

	// Get the items
	items := make(chan crm.Item)
	go func() {
		defer close(items)
		err := c.GetItems(context.Background(), items)
		require.NoError(t, err)
	}()
	count := 0
	var lastItem crm.Item
	toRemove := []crm.Item{}
	for item := range items {
		if item.GetFields()["Name"] != fmt.Sprintf("Name %d", count) {
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
	err = c.UpdateItem(ctx, lastItem, map[string]interface{}{
		"Name": "New name",
	})
	if err != nil {
		t.Error(err)
	}

	// Check to make sure it updated
	_, err = c.GetItem(context.Background(), map[string]interface{}{
		"Name": "New name",
	})
	if err != nil {
		t.Error(err)
		return
	}

	// -- Test removing the middle two items --

	// Add new item
	err = c.CreateItemFromValues(ctx, []string{"Name last", "Item last"})
	if err != nil {
		t.Error(err)
		return
	}

	// Purposely make out of order
	toRemove[0], toRemove[1] = toRemove[1], toRemove[0]
	for _, toRemoveI := range toRemove {
		err = c.RemoveItems(ctx, toRemoveI)
		if err != nil {
			t.Error(err)
		}
	}

	// Fetch all remaining items and remove them
	toRemove = []crm.Item{}
	items = make(chan crm.Item)
	go func() {
		defer close(items)
		err := c.GetItems(context.Background(), items)
		require.NoError(t, err)
	}()
	for item := range items {
		toRemove = append(toRemove, item)
	}

	err = c.RemoveItems(ctx, toRemove...)
	assert.NoError(t, err)

	// Make sure there are no items
	items = make(chan crm.Item)
	go func() {
		defer close(items)
		err := c.GetItems(context.Background(), items)
		require.NoError(t, err)
	}()
	itemCount := 0
	for range items {
		itemCount++
	}
	assert.NoError(t, err)
	assert.Equal(t, itemCount, 0)
}
