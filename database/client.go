package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"gopkg.in/mgo.v2"
	"time"
)

// Client represents a client to the mongoDB data store.
type Client struct {
	// package logger.
	logger *log.Entry

	// returns the current time.
	Now func() time.Time

	// location of the mongodb daemon.
	Host string

	// database connection.
	Session *mgo.Session

	// Service for interacting with the database.
	service Service
}

// NewClient creates a new database client.
func NewClient(h string) *Client {
	c := &Client{
		logger: log.WithFields(log.Fields{"package": "database"}),
		Now:    time.Now,
		Host:   h,
	}
	c.service.client = c
	return c
}

// Open opens and initializes the MongoDB database.
func (c *Client) Open() error {
	// connect to the database.
	session, err := mgo.Dial(c.Host)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "host": c.Host}).Error(errDatabaseFailed)
		return err
	}
	c.Session = session

	c.logger.WithFields(log.Fields{"host": c.Host}).Info("connected to the database")
	return nil
}

// Close stops the connection to the database.
func (c *Client) Close() error {
	c.Session.Close()
	return nil
}

// Service returns the service associated with the client.
func (c *Client) Service() warcraft.DatabaseService { return &c.service }

// InsertAuctions inserts a slice of auctions into the database.
func (c *Client) InsertAuctions(auctions []interface{}) error {
	// connect to collection.
	session := c.Session.Copy()
	defer session.Close()
	col := session.DB("").C(AuctionCollection)

	// insert auctions.
	b := col.Bulk()
	b.Unordered()
	b.Insert(auctions...)
	if _, err := b.Run(); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errDatabaseInsert)
		return nil
	}
	return nil
}
