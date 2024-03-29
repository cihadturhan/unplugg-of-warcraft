package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

// InsertAuctions inserts a slice of auctions into the database.
func (c *Client) Insert(collectionName string, auctions []interface{}) error {
	// connect to collection.
	session := c.Session.Copy()
	defer session.Close()
	col := session.DB("warcraft").C(collectionName)

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

// GetAuctions returns all the auctions
func (c *Client) Find(collectionName string, options bson.M) ([]warcraft.Auction, error) {
	// connect to collection.
	session := c.Session.Copy()
	defer session.Close()
	col := session.DB("warcraft").C(collectionName)

	// get auctions.
	var auctions []warcraft.Auction

	if err := col.Find(options).All(&auctions); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errDatabaseQuery)
		return nil, err
	}

	return auctions, nil
}

// GetLastRecord returns the last record present in a collection
func (c *Client) GetLastRecord(collectionName string) (warcraft.Auction, error) {
	var lastRecord warcraft.Auction

	// connect to collection.
	session := c.Session.Copy()
	defer session.Close()
	col := session.DB("warcraft").C(collectionName)

	if err := col.Find(nil).Sort("-timestamp").One(&lastRecord); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errDatabaseQuery)
		return lastRecord, err
	}

	return lastRecord, nil
}

// Service returns the service associated with the client.
func (c *Client) Service() warcraft.DatabaseService { return &c.service }
