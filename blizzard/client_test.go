package configuration_test

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/brand-digital-box/configuration"
	"github.com/whitesmith/brand-digital-box/mock"
)

// Client is a test wrapper.
type Client struct {
	*configuration.Client
	MQTTService       *mock.MQTTService
	ScreenService     *mock.ScreenService
	DatabaseService   *mock.DatabaseService
	DownloaderService *mock.DownloaderService
	PlayerService     *mock.PlayerService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	log.SetLevel(log.DebugLevel)
	// create client wrapper.
	c := &Client{
		Client:            configuration.NewClient("nuc-aspire"),
		MQTTService:       &mock.MQTTService{},
		ScreenService:     &mock.ScreenService{},
		DatabaseService:   &mock.DatabaseService{},
		DownloaderService: &mock.DownloaderService{},
		PlayerService:     &mock.PlayerService{},
	}
	c.Client.MQTTService = c.MQTTService
	c.Client.ScreenService = c.ScreenService
	c.Client.DatabaseService = c.DatabaseService
	c.Client.DownloaderService = c.DownloaderService
	c.Client.PlayerService = c.PlayerService
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	return c.Client.Close()
}
