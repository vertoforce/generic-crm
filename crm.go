package crm

import (
	"context"
	"fmt"
)

// Standardized errors
var (
	ErrItemNotFound = fmt.Errorf("item not found")
)

// CRM is a set of generic items.
// You can fetch get, remove, or update items in a crm.
//
// The whole point is that a each item in a crm has some unique id, and that other fields can change,
// so a Synchronize function is provided on crms also.
type CRM interface {
	// Get all items with optional searchFields
	GetItems(ctx context.Context, items chan Item, searchFields ...map[string]interface{}) error
	GetItem(ctx context.Context, searchFields map[string]interface{}) (Item, error)
	RemoveItems(ctx context.Context, i ...Item) error
	CreateItem(ctx context.Context, i Item) error
	UpdateItem(ctx context.Context, i Item, updateFields map[string]interface{}) error
	// Create or update current columns using the example item to determine types
	// If a column exists, it will update it's type
	UpdateColumns(ctx context.Context, exampleItem Item) error
}

// Item is a generic item from the crm.
//
// The Fields keys should match whatever fields are configured in the crm, usuall the column names
type Item interface {
	// Get the keys and values of each field in this item
	GetFields() map[string]interface{}
}

// DefaultItem is used for creating items, it just contains the fields of the item
type DefaultItem struct {
	Fields map[string]interface{}
}

// GetFields of the default item
func (d *DefaultItem) GetFields() map[string]interface{} {
	return d.Fields
}
