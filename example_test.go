package crm_test

import (
	"context"
	"fmt"
	"os"

	crm "github.com/vertoforce/generic-crm"
	"github.com/vertoforce/generic-crm/backends/airtablecrm"
)

func Example() {
	// This code does not check for errors

	// Create a object that implements the interface
	a, _ := airtablecrm.New(os.Getenv("AIRTABLE_API_KEY"), os.Getenv("AIRTABLE_BASE_ID"), "Testing")

	// Cast it to the interface so we can drop in another crm at any time
	c := crm.CRM(a)

	// Use it as a crm!
	// Note that the fields must already bet set up in your connected CRM.
	// So in this example, you must have the first row of your google sheet contain the header "Name"
	// Or for airtable, you must have a column named "Name"
	c.CreateItem(context.Background(), &crm.DefaultItem{
		Fields: map[string]interface{}{
			"Name": "test",
		},
	})

	items := make(chan crm.Item)
	go func() {
		defer close(items)
		err := c.GetItems(context.Background(), items)
		_ = err
	}()
	toRemove := []crm.Item{}
	for item := range items {
		toRemove = append(toRemove, item)
		fmt.Println(item.GetFields()["Name"])
	}
	c.RemoveItems(context.Background(), toRemove...)

	// Output: test
}
