package sqlcrm

type Item struct {
	Fields map[string]interface{}
}

func (i *Item) GetFields() map[string]interface{} {
	return i.Fields
}
