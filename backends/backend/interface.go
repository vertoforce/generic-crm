package backend

import "context"

// Item is a generic item from the crm
type Item struct {
	// Internal item fields, usually left untouched
	Internal interface{}
	Fields   map[string]interface{}
}

// Backend is the interface that a crm needs to comply to
type Backend interface {
	GetItems(ctx context.Context) ([]*Item, error)
	RemoveItem(ctx context.Context, i *Item) error
	CreateItem(ctx context.Context, i *Item) error
	UpdateItem(ctx context.Context, i *Item, updateFields map[string]interface{}) error
}
