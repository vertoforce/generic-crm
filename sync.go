package crm

import (
	"context"
	"fmt"
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
	// safeItems stores the items that were either created or updated, and therefore should not be removed
	// It marks the item safe by storing the search function used to find it
	safeItems := map[*map[string]interface{}]bool{}

	for newItem := range items {
		markedSafe := false

		// Update each crm
		for _, crm := range s.crms {
			// Check if ths CRM contains this item
			newItemSearch := s.searchFunction(newItem)
			if oldItem, err := crm.GetItem(ctx, newItemSearch); err == nil && oldItem != nil {
				// We found the item, update it
				err = crm.UpdateItem(ctx, oldItem, newItem.GetFields())
				if err != nil {
					return err
				}
				if !markedSafe {
					safeItems[&newItemSearch] = true
					markedSafe = true
				}
				continue
			}

			// Create the item
			err := crm.CreateItem(ctx, newItem)
			if err != nil {
				return err
			}
			if !markedSafe {
				safeItems[&newItemSearch] = true
				markedSafe = true
			}
		}
	}

	// Now we need to delete any items that were not updates or created
	if s.deleteUntouchedItems {
		for _, crm := range s.crms {
			toRemove := []Item{}

			// Go through each item to check if it's safe
			items, err := crm.GetItems(ctx)
			if err != nil {
				return fmt.Errorf("failed to get items to delete old items: %s", err)
			}

		itemLoop:
			for _, item := range items {
				// Check if this item is safe
			safeItemSearchLoop:
				for safeItemSearch := range safeItems {
					for safeItemKey, safeItemValue := range *safeItemSearch {
						if itemValue, ok := item.GetFields()[safeItemKey]; !ok || itemValue != safeItemValue {
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
			err = crm.RemoveItems(ctx, toRemove...)
			if err != nil {
				return fmt.Errorf("failed to delete item: %s", err)
			}
		}
	}

	return nil
}
