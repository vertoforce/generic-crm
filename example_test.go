package crm

import (
	"context"
	"fmt"
	"os"

	"github.com/vertoforce/generic-crm/backends/airtable"

	"github.com/vertoforce/generic-crm/backends/crm"
)

func Example() {
	// This code does not check for errors

	// Create a object that implements the interface
	a, _ := airtable.New(os.Getenv("AIRTABLE_API_KEY"), os.Getenv("AIRTABLE_BASE_ID"), "Testing")

	// Cast it to the interface so we can drop in another frontend at any time
	c := crm.CRM(a)

	// Use it as a crm!
	// Note that the fields must already bet set up in your connected CRM.
	// So in this example, you must have the first row of your google sheet contain the header "Name"
	// Or for airtable, you must have a column named "Name"
	c.CreateItem(context.Background(), &crm.Item{
		Fields: map[string]interface{}{
			"Name": "test",
		},
	})

	items, _ := c.GetItems(context.Background())
	for _, item := range items {
		fmt.Println(item.Fields["Name"])
	}

	// Output: test
}
