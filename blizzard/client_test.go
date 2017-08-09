package blizzard_test

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/mock"
)

// Client is a test wrapper.
type Client struct {
	*blizzard.Client
	DatabaseService *mock.DatabaseService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)
	// create client wrapper.
	c := &Client{
		Client:          blizzard.NewClient(1, "grim-batol", "en_GB", ""),
		DatabaseService: &mock.DatabaseService{},
	}
	c.Client.DatabaseService = c.DatabaseService
	return c
}
