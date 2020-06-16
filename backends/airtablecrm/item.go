package airtablecrm

// Item is the internal item to keep track of the airtable ID
type Item struct {
	AirtableID string `json:"ID"`
	Fields     map[string]interface{}
}

// Attachment is a field type to upload an attachment
//
// So if you want to upload an attachment, set this as on of the fields when calling CreateRecord
type Attachment struct {
	URL string `json:"url"`
}

// GetFields of this item
func (i *Item) GetFields() map[string]interface{} {
	return i.Fields
}
