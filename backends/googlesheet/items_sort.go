package googlesheet

// Items is a list of item
type Items []*Item

// Len Length of items
func (a Items) Len() int {
	return len(a)
}

// Swap two items
func (a Items) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

// Returns if the item has a lower row number than the later
func (a Items) Less(i, j int) bool {
	return a[i].RowNumber < a[j].RowNumber
}
