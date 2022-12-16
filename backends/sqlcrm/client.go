package sqlcrm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Client struct {
	DB    *sqlx.DB
	table string
}

// NewCRM Creates a new sql crm
//
// Connection string should look like `user:password@tcp(127.0.0.1:3306)/hello`
func NewCRM(connectionURL string, table string) (*Client, error) {
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
		table: table,
	}

	return c, nil
}

// Close the database
func (c *Client) Close() error {
	return c.DB.Close()
}
