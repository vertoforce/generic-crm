package sqlcrm

import (
	"context"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/segmentio/agecache"
)

type Client struct {
	DB           *sqlx.DB
	Table        string
	columnsCache *agecache.Cache[string, map[string]string]
}

var ErrTableNotFound = fmt.Errorf("table not found")

// NewCRM Creates a new sql crm
//
// Connection string should look like `user:password@tcp(127.0.0.1:3306)/hello`
func NewCRM(ctx context.Context, connectionURL string, table string) (*Client, error) {
	db, err := sqlx.Open("mysql", connectionURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c := &Client{
		DB:    db,
		Table: table,
		columnsCache: agecache.New(agecache.Config[string, map[string]string]{
			MaxAge:   time.Minute,
			Capacity: 1,
		}),
	}

	// Make sure table exists
	columns, err := c.getColumns(ctx)
	if err != nil {
		return nil, fmt.Errorf("error listing columns for table: %w", err)
	}
	if len(columns) == 0 {
		return nil, ErrTableNotFound
	}

	return c, nil
}

// Close the database
func (c *Client) Close() error {
	return c.DB.Close()
}
