package airtablecrm

import "github.com/fabioberger/airtable-go"

// Client to interact with airtable as a crm
type Client struct {
	client    *airtable.Client
	tableName string
}

// New Creates a new airtable CRM
func New(apiKey, baseID string, tableName string) (*Client, error) {
	c := &Client{}

	airtableClient, err := airtable.New(apiKey, baseID)
	if err != nil {
		return nil, err
	}
	c.client = airtableClient
	c.tableName = tableName

	return c, nil

}
