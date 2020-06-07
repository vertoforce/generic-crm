package googlesheet

// LoadHeaders loads headers from the first row of the sheet into the Client object
func (c *Client) LoadHeaders() {
	if len(c.Sheet.Rows) >= 1 {
		for _, cell := range c.Sheet.Rows[0] {
			c.Headers = append(c.Headers, cell.Value)
		}
	}
}

// SetHeaders Sets the first row of the sheet.  Note headers MUST be enabled, otherwise nothing will happen
func (c *Client) SetHeaders(headers []string) error {
	for i, header := range headers {
		updateCell(c.Sheet, 0, i, header)
	}
	if c.WaitToSynchronize {
		return nil
	}
	return c.Synchronize()
}

// getHeaderIndex Finds the index of a header in c.Headers
func (c *Client) getHeaderIndex(header string) int {
	for i, h := range c.Headers {
		if h == header {
			return i
		}
	}
	return -1
}
