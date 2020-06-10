package googlesheet

// Item is the base element of an item for each row in the google sheet
type Item struct {
	RowNumber int
	Fields    []string
	client    *Client
}

// GetFields Returns a map of header -> value
func (i *Item) GetFields() map[string]interface{} {
	ret := map[string]interface{}{}
	for c, value := range i.Fields {
		// Break if this row has too many values
		if c >= len(i.client.Headers) {
			break
		}

		ret[i.client.Headers[c]] = value
	}

	return ret
}
