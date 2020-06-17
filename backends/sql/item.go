package sqlcrm

// Item that lives in the sql table
type Item struct {
	Fields map[string]interface{}
}

// GetFields of the sql item
func (i *Item) GetFields() map[string]interface{} {
	return i.Fields
}
