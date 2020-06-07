// Package googlesheet makes it easy to use a spreadsheet as a CRM
// It breaks each row into an "item" with distinct fields
package googlesheet

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/vertoforce/generic-crm/backends/googlesheet/quotatrack"

	"golang.org/x/oauth2/google"
	"gopkg.in/Iwark/spreadsheet.v2"
)

// Google API limits
// This version of the Google Sheets API has a limit of 500 requests per 100 seconds per project,
// and 100 requests per 100 seconds per user. Limits for reads and writes are tracked separately. There is no daily usage limit
const (
	GoogleSheetUsageLimit     = 100
	GoogleSheetUsageLimitTime = time.Second * 105 // Set to 105 to be safe
)

// Client is a session with a google sheet
type Client struct {
	Service           *spreadsheet.Service     // Authenticated google sheets api service
	Spreadsheet       *spreadsheet.Spreadsheet // Spreadsheet we are working with
	Sheet             *spreadsheet.Sheet       // Sheet we are working with
	Headers           []string                 // Header row of column names.  If this is blank, no headers for this sheet
	WaitToSynchronize bool                     // Don't synchronize the sheet after every request, wait for Synchronize to be called
	quota             *quotatrack.Quota        // Quota to track our usage to see if we need to slow down
}

// Config to create a new client
type Config struct {
	GoogleClientSecretFile string
	SpreadsheetURL         string
	// SheetName in the name of the sheet we will use
	SheetName string
	// Don't synchronize the sheet after every request, wait for Synchronize to be called
	WaitToSynchronize bool
}

// New creates a new client
func New(ctx context.Context, config *Config) (*Client, error) {
	// Connect
	bytes, err := ioutil.ReadFile(config.GoogleClientSecretFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open client secret file")
	}
	conf, err := google.JWTConfigFromJSON(bytes, spreadsheet.Scope)
	if err != nil {
		return nil, fmt.Errorf("failed to read client secret: %v", err)
	}
	googleClient := conf.Client(ctx)

	client := &Client{
		Service:           spreadsheet.NewServiceWithClient(googleClient),
		WaitToSynchronize: config.WaitToSynchronize,
		quota:             quotatrack.New(GoogleSheetUsageLimitTime),
	}

	// Load spreadsheet
	spreadsheet, err := client.Service.FetchSpreadsheet(GetSpreadsheetID(config.SpreadsheetURL))
	if err != nil {
		return nil, fmt.Errorf("failed to load spreadsheet: %s", err)
	}
	client.Spreadsheet = &spreadsheet

	// Get sheet
	sheet, err := spreadsheet.SheetByTitle(config.SheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to load sheet: %s", err)
	}
	client.Sheet = sheet

	// Get headers
	client.LoadHeaders()

	return client, nil
}

func updateCell(sheet *spreadsheet.Sheet, row int, col int, value string) {
	sheet.Update(row, col, value)
}

// Synchronize - If the client is set to waitToSynchronize, this function synchronizes the sheet after a series of operations
func (c *Client) Synchronize() error {
	c.consumeQuota()
	return c.Sheet.Synchronize()
}

// consumeQuota by waiting for one to be available, and then consuming it
func (c *Client) consumeQuota() {
	if quotaUsage := c.quota.Usage(); quotaUsage >= GoogleSheetUsageLimit {
		sleepTime := c.quota.TimeUntilQuotaAvailable(GoogleSheetUsageLimit, 1)
		time.Sleep(sleepTime)
	}
	c.quota.Consume(1)
}
