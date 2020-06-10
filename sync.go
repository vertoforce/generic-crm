package crm

import (
	"context"
	"reflect"
	"time"
)

const (
	// graceTime is a small amount of time to be more graceful in request to a CRMc
	graceTime = time.Millisecond * 50
)

// CompareFunction Takes two items and returns true if they are the same unique object (with potentially different fields).
//
// It should compare some unique id that would remain the same despite field values updating.
//
// This is used to identify an item in the CRM when the field values have changed.
// For example in synchronization when checking if an item should be updated with new incoming items.
//
// See the DefaultCompareFunction for an example
type CompareFunction func(a Item, b Item) bool

// DefaultCompareFunction compares the "ID" field of each item
var DefaultCompareFunction = func(a Item, b Item) bool {
	IDA, ok := a.GetFields()["ID"]
	if !ok {
		return false
	}
	IDB, ok := b.GetFields()["ID"]
	if !ok {
		return false
	}
	return reflect.DeepEqual(IDA, IDB)
}

// Synchronize Updates a list of CRMs given a stream of new items
//
// It loops through the provided channel of new items, and for each CRM
// 	* updates the old item (if it exists using the CompareFunction)
// 	* creates the new item (since it does not exist in the crm)
func Synchronize(ctx context.Context, items chan Item, compareFunction CompareFunction, crms ...CRM) error {
	// First cache the current contents of the crms so we don't need to fetch it each time
	CrmsOldItems := [][]Item{}
	for _, crm := range crms {
		oldItems, err := crm.GetItems(ctx)
		if err != nil {
			return err
		}
		CrmsOldItems = append(CrmsOldItems, oldItems)
	}

	// Then loop through each new item and update the old or create the new item
	for newItem := range items {
	crmLoop:
		for crmI, crm := range crms {
			select {
			case <-time.After(graceTime):
			case <-ctx.Done():
				return ctx.Err()
			}

			// Search for the new item in our old items
			for _, oldItem := range CrmsOldItems[crmI] {
				if compareFunction(oldItem, newItem) {
					// We have this item, update it
					err := crm.UpdateItem(ctx, oldItem, newItem.GetFields())
					if err != nil {
						return err
					}
					break crmLoop
				}
			}

			// The item doesn't exist, create it
			err := crm.CreateItem(ctx, newItem)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
