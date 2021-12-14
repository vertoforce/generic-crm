package crm_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	crm "github.com/vertoforce/generic-crm"
)

const (
	testSyncItemCount = 3
)

func TestSync(t *testing.T) {
	testCRMs, err := getTestCRMs()
	if err != nil {
		t.Error(err)
		return
	}

	// Create some items
	ctx := context.Background()
	for i := 0; i < testSyncItemCount; i++ {
		for _, testCRM := range testCRMs {
			testCRM.CreateItem(ctx, &crm.DefaultItem{
				Fields: map[string]interface{}{
					"Name": fmt.Sprintf("Name %d", i),
				},
			})
		}
	}

	// Build stream of updates
	// These updates will update Name 1, and Name 2, and create a Name 3
	newItems := make(chan crm.Item)
	go func() {
		for i := 1; i < testSyncItemCount+1; i++ {
			newItems <- &crm.DefaultItem{
				Fields: map[string]interface{}{
					"Name": fmt.Sprintf("Name %d", i),
					"Item": "Updated content",
				},
			}
		}
		close(newItems)
	}()

	// Build sync machine
	syncMachine := crm.NewSyncMachine().
		SetDeleteUntouchedItems(true).
		WithCRMs(testCRMs...).
		WithSearchFunction(func(i crm.Item) map[string]interface{} {
			return map[string]interface{}{
				"Name": i.GetFields()["Name"],
			}
		})

	err = syncMachine.Sync(ctx, newItems)
	if err != nil {
		t.Error(err)
		return
	}

	// Check if the CRMs are in the state we'd expect
	for _, testCRM := range testCRMs {
		items, err := testCRM.GetItems(ctx)
		if err != nil {
			t.Error(err)
			return
		}
		foundNames := map[string]bool{}
		toDelete := []crm.Item{}
		for item := range items {
			toDelete = append(toDelete, item)
			foundNames[item.GetFields()["Name"].(string)] = true
		}
		if len(foundNames) > testSyncItemCount {
			t.Errorf("too many items in CRM")
		}
		for i := 1; i < testSyncItemCount+1; i++ {
			if _, ok := foundNames[fmt.Sprintf("Name %d", i)]; !ok {
				t.Errorf("CRM does not have the expected values")
			}
		}

		// Delete all items
		testCRM.RemoveItems(ctx, toDelete...)
	}

}

func TestForgivingEqual(t *testing.T) {
	tests := []struct {
		A     interface{}
		B     interface{}
		Equal bool
	}{
		{A: float64(1), B: int64(1), Equal: true},
		{A: int64(1), B: float64(1), Equal: true},
		{A: "1", B: int64(1), Equal: false},
		{A: 1, B: 2, Equal: false},
		{A: "1", B: "1", Equal: true},
		{A: "1", B: "2", Equal: false},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			result := crm.ForgivingEqual(test.A, test.B)
			require.Equal(t, test.Equal, result)
		})
	}
}
