package analyzer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
)

// Client represents a client for interacting with the analyzer
type Client struct {
	// package logger.
	logger *log.Entry

	// service for interacting with the analyzer
	service Service

	// database service.
	DatabaseService warcraft.DatabaseService
}

// NewClient returns a new configuration client.
func NewClient() *Client {
	c := &Client{
		logger: log.WithFields(log.Fields{"package": "blizzard"}),
	}

	c.service.client = c
	return c
}

// CreateHash creates the hash to find duplicates
func (c *Client) CreateHash(auctions []warcraft.Auction) map[int]warcraft.Auction {
	auctionsHash := make(map[int]warcraft.Auction, 0)

	for _, auction := range auctions {
		auctionsHash[auction.ID] = auction
	}

	return auctionsHash
}

// CreateBuyout receives an auction and converts it to an buyout
func (c *Client) CreateBuyout(auction warcraft.Auction) warcraft.Buyout {
	return warcraft.Buyout{
		ID:        auction.ID,
		Item:      auction.Item,
		Buyout:    auction.Buyout,
		Quantity:  auction.Quantity,
		Timestamp: auction.Timestamp,
	}
}

// AuctionsThatEnded returns all the auctions that not present in the next dump (ended)
func (c *Client) AuctionsThatEnded(prevAuctions map[int]warcraft.Auction, auctions map[int]warcraft.Auction) []warcraft.Buyout {
	results := make([]warcraft.Buyout, 0)

	for key, value := range prevAuctions {
		if _, present := auctions[key]; !present {
			results = append(results, c.CreateBuyout(value))
		}
	}

	return results
}

// AddAuctionsThatEndedToBuyoutsCollection checks which auctions have ended and adds the ones that ended to the buyout collection
func (c *Client) AddAuctionsThatEndedToBuyoutsCollection(prevAuctions []warcraft.Auction, auctions []warcraft.Auction) error {
	prevAuctionsMap := c.CreateHash(prevAuctions)
	auctionsMap := c.CreateHash(auctions)
	buyouts := c.AuctionsThatEnded(prevAuctionsMap, auctionsMap)

	// convert buyouts to interfaces
	records := make([]interface{}, 0)
	for _, buyout := range buyouts {
		records = append(records, buyout)
	}

	// insert to database
	if err := c.DatabaseService.Insert(warcraft.BuyoutCollection, records); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error("Failed to insert buyout")
		return err
	}

	return nil
}

// Service returns the service for interacting with the dump analyzer
func (c *Client) Service() warcraft.AnalyzerService { return &c.service }
