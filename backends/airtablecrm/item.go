package airtablecrm

// Item is the internal item to keep track of the airtable ID
type Item struct {
	airtableID string
	fields     map[string]interface{}
}

// GetFields of this item
func (i *Item) GetFields() map[string]interface{} {
	return i.fields
}
