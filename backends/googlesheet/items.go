package googlesheet

// Item is the base element of an item for each row in the google sheet
type Item struct {
	RowNumber int
	Fields    []string
	client    *Client
}

// ToMap Returns a map of header -> value
func (i *Item) ToMap() map[string]string {
	ret := map[string]string{}
	for c, value := range i.Fields {
		// Break if this row has too many values
		if c >= len(i.client.Headers) {
			break
		}

		ret[i.client.Headers[c]] = value
	}

	return ret
}

// UnmarshalItem Fills an item in to a struct using the struct names -> mapping to -> header names in the sheet
func UnmarshalItem(item *Item, v interface{}) error {
	// TODO:
	return nil
}
