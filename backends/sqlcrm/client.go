package sqlcrm

import (
	"context"
	"fmt"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/segmentio/agecache"
)

type Client struct {
	DB           *sqlx.DB
	Table        string
	columnsCache *agecache.Cache[string, map[string]string]
}

type Config struct {
	ConnectionURL          string
	Table                  string
	CreateTableIfNotExists bool
}

var ErrTableNotFound = fmt.Errorf("table not found")

// NewCRM Creates a new sql crm
//
// Connection string should look like `user:password@tcp(127.0.0.1:3306)/hello`
func NewCRM(ctx context.Context, config Config) (*Client, error) {
	db, err := sqlx.Open("mysql", config.ConnectionURL)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	c := &Client{
		DB:    db,
		Table: config.Table,
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
		if !config.CreateTableIfNotExists {
			return nil, ErrTableNotFound
		}
		// Try to create table
		_, err := db.ExecContext(ctx, `
			CREATE TABLE `+strings.ReplaceAll(pq.QuoteIdentifier(c.Table), "\"", "")+` (
				sqlid int auto_increment NOT NULL,
				CONSTRAINT Test2_PK PRIMARY KEY (sqlid)
			)`)
		if err != nil {
			return nil, fmt.Errorf("error creating table since it does not exist: %w", err)
		}
		// Reset column cache
		c.columnsCache.Clear()
	}

	return c, nil
}

// Close the database
func (c *Client) Close() error {
	return c.DB.Close()
}
