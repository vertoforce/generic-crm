package crm

import (
	"context"
	"fmt"
	"reflect"

	"github.com/opentracing/opentracing-go"
)

// SearchFunction is a function that given an item, return the search fields to find that unique item.
// It should return a search for some unique ID of that item
type SearchFunction func(i Item) map[string]interface{}

// DefaultSearchFunction just returns a search for the ID field of the item
var DefaultSearchFunction = func(i Item) map[string]interface{} {
	return map[string]interface{}{
		"ID": i.GetFields()["ID"],
	}
}

// SyncMachine is a tool making it easy to update a CRM based on a stream of new items.
//
// It's basically the equivalent of deleting everything and adding it again.
//
// It loops through the channel of new items and
//   * Updates the old item in the crm (If it exists)
//   * Creates the new item if it does not exist
//
// At the end it will delete any items in the CRM that were not updated
type SyncMachine struct {
	deleteUntouchedItems bool // At the end should we delete items that were not updated
	crms                 []CRM
	searchFunction       SearchFunction
}

// NewSyncMachine creates a new sync machine with the default search function, no crms, and deleteUntouchedItems to true
func NewSyncMachine() *SyncMachine {
	return &SyncMachine{
		deleteUntouchedItems: true,
		crms:                 []CRM{},
		searchFunction:       DefaultSearchFunction,
	}
}

// WithSearchFunction Set the search function to find a unique item in this sync machine.
// For more check the description of a SearchFunction
func (s *SyncMachine) WithSearchFunction(SearchFunction SearchFunction) *SyncMachine {
	s.searchFunction = SearchFunction
	return s
}

// WithCRMs adds CRMs to the sync machine.
// Each added CRM will be synced with the incoming list of new items
func (s *SyncMachine) WithCRMs(crms ...CRM) *SyncMachine {
	s.crms = append(s.crms, crms...)
	return s
}

// SetDeleteUntouchedItems If set true, at the end of the sync it will delete all items that weren't deleted or created
//
// This allows us to make the CRM directly reflect the incoming stream of new items
func (s *SyncMachine) SetDeleteUntouchedItems(deleteUntouchedItems bool) *SyncMachine {
	s.deleteUntouchedItems = deleteUntouchedItems
	return s
}

// Sync Performs the actual sync task, see the description of SyncMachine
func (s *SyncMachine) Sync(ctx context.Context, items chan Item) error {
	var span opentracing.Span
	span, ctx = opentracing.StartSpanFromContext(ctx, "Sync")
	defer span.Finish()

	// safeItems stores the items that were either created or updated, and therefore should not be removed
	// It marks the item safe by storing the search function used to find it
	safeItems := map[*map[string]interface{}]bool{}

	for newItem := range items {
		markedSafe := false
		// Update each crm
		for _, crm := range s.crms {
			select {
			case <-ctx.Done():
				return ctx.Err()
			default:
			}

			processItemSpan, processItemCtx := opentracing.StartSpanFromContext(ctx, "ProcessItem")
			processItemSpan.SetTag("CRM", reflect.TypeOf(crm).String())
			// Check if ths CRM contains this item
			SearchSpan, _ := opentracing.StartSpanFromContext(processItemCtx, "SearchFunction")
			newItemSearch := s.searchFunction(newItem)
			SearchSpan.Finish()
			if oldItem, err := crm.GetItem(processItemCtx, newItemSearch); err == nil && oldItem != nil {
				// We found the item, update it
				err = crm.UpdateItem(processItemCtx, oldItem, newItem.GetFields())
				if err != nil {
					processItemSpan.Finish()
					return err
				}
				if !markedSafe {
					safeItems[&newItemSearch] = true
					markedSafe = true
				}
				continue
			}

			// Create the item
			err := crm.CreateItem(processItemCtx, newItem)
			if err != nil {
				processItemSpan.Finish()
				return err
			}
			if !markedSafe {
				safeItems[&newItemSearch] = true
				markedSafe = true
			}
			processItemSpan.Finish()
		}
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	// Now we need to delete any items that were not updates or created
	if s.deleteUntouchedItems {
		deleteSpan, deleteCtx := opentracing.StartSpanFromContext(ctx, "DeleteUntouchedItems")
		for _, crm := range s.crms {
			toRemove := []Item{}

			// Go through each item to check if it's safe
			items, err := crm.GetItems(deleteCtx)
			if err != nil {
				return fmt.Errorf("failed to get items to delete old items: %s", err)
			}

		itemLoop:
			for item := range items {
				select {
				case <-ctx.Done():
					return ctx.Err()
				default:
				}

				// Check if this item is safe
			safeItemSearchLoop:
				for safeItemSearch := range safeItems {
					for safeItemKey, safeItemValue := range *safeItemSearch {
						if itemValue, ok := item.GetFields()[safeItemKey]; !ok || !ForgivingEqual(itemValue, safeItemValue) {
							// This item does not match this safe item, try next one
							continue safeItemSearchLoop
						}
					}
					// This item matches this safe item, it is safe
					continue itemLoop
				}

				// This item is NOT safe, mark to be removed
				toRemove = append(toRemove, item)
			}

			// Remove all items marked for deletion
			err = crm.RemoveItems(deleteCtx, toRemove...)
			if err != nil {
				return fmt.Errorf("failed to delete item: %s", err)
			}
		}
		deleteSpan.Finish()
	}

	return nil
}

// forgivingEqual compares two values with some forgiveness in making sure they are equal.
// For example, if a is a float, and b is an int, but they are the same value, it will return true.
func ForgivingEqual(a, b interface{}) bool {
	if a == b {
		return true
	}

	// Do number comparisson
	// Convert both to floats then compare
	aValue := float64(0)
	switch a.(type) {
	case float64:
		aValue = a.(float64)
	case int64:
		aValue = float64(a.(int64))
	case int32:
		aValue = float64(a.(int32))
	default:
		// If it is not number, we don't use this comparison
		return false
	}

	bValue := float64(0)
	switch b.(type) {
	case float64:
		bValue = b.(float64)
	case int64:
		bValue = float64(b.(int64))
	case int32:
		bValue = float64(b.(int32))
	default:
		return false
	}

	if aValue == bValue {
		return true
	}

	return false
}
