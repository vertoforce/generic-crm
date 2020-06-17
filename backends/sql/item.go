package sqlcrm

type Item struct {
	id     int64 // SQL row id
	Fields map[string]interface{}
}

func (i *Item) GetFields() map[string]interface{} {
	return i.Fields
}
