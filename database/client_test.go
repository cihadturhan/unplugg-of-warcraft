package database_test

import (
	"github.com/whitesmith/unplugg-of-warcraft/database"
	"time"
)

// Now is the mocked current time for testing.
var Now = time.Now().Round(time.Second)

// Client is a test wrapper for database.Client.
type Client struct {
	*database.Client
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	// create client wrapper.
	c := &Client{
		Client: database.NewClient("localhost/test"),
	}
	c.Now = func() time.Time { return Now }

	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}

	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	return c.Client.Close()
}
