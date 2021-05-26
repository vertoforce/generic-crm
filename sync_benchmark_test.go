package crm_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
	"github.com/vertoforce/generic-crm/backends/googlesheet"
)

func BenchmarkSyncGoogleSheet(b *testing.B) {
	ctx := context.Background()
	// Set up testing google sheet client
	googleSheetCRM, err := googlesheet.New(context.Background(), &googlesheet.Config{
		GoogleClientSecretFile: os.Getenv("GoogleClientSecretFile"),
		SpreadsheetURL:         os.Getenv("TESTING_SPREADSHEET_URL"),
		SheetName:              "Sheet1",
	})
	require.NoError(b, err)

	// Disable automatic sync, we are doing the test locally
	googleSheetCRM.WaitToSynchronize = true

	// Insert some data
	for i := 0; i < 1000; i++ {
		googleSheetCRM.CreateItem(ctx, &crm.DefaultItem{Fields: map[string]interface{}{
			"Name": fmt.Sprintf("Name %d", i),
			"Item": fmt.Sprintf("%d", i),
		}})
	}

	// Build set of updates
	updates := []crm.Item{}
	// Add a bunch of updates to items
	for i := 0; i < 1000; i++ {
		updates = append(updates, &crm.DefaultItem{Fields: map[string]interface{}{"Name": fmt.Sprintf("Name %d", i), "Item": "New item title"}})
	}

	searchFunc := func(i crm.Item) map[string]interface{} {
		return map[string]interface{}{"Name": i.GetFields()["Name"]}
	}

	b.Run("GoogleSheetSync", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// set thread to send updated items to channel
			items := make(chan crm.Item)
			go func() {
				for _, update := range updates {
					items <- update
				}
				close(items)
			}()

			err = crm.NewSyncMachine().WithCRMs(googleSheetCRM).SetDeleteUntouchedItems(false).WithSearchFunction(searchFunc).Sync(ctx, items)
			require.NoError(b, err)
		}
	})

	// Clear the sheet
	// Turns out this is too long since it does it one by one
	// items, err = googleSheetCRM.GetItems(ctx)
	// require.NoError(b, err)
	// itemsSlice := []crm.Item{}
	// for item := range items {
	// 	itemsSlice = append(itemsSlice, item)
	// }
	// err = googleSheetCRM.RemoveItems(ctx, itemsSlice...)
	// require.NoError(b, err)
}
