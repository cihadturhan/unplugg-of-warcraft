package analyzer

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
)

const BuyoutsCollection = "buyouts"

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
func (c *Client) CreateHash(auctions []warcraft.Auction) map[int]int {
	hash := make(map[int]int)

	for _, auction := range auctions {
		hash[auction.ID] = 0
	}

	return hash
}

// NumberOfOcurrences checks if the auction is present in the next dump
func (c *Client) NumberOfOcurrences(auctions []warcraft.Auction, nextAuctions []warcraft.Auction) map[int]int {
	hash := c.CreateHash(auctions)

	for _, auction := range nextAuctions {
		hash[auction.ID]++
	}

	return hash
}

// AuctionIsPresentInNextDump checks if the auction is present in the next dump
func (c *Client) AuctionIsPresentInNextDump(auctionID int, occurencesHash map[int]int) bool {
	if occurencesHash[auctionID] == 0 {
		return false
	}
	return true
}

// TODO optimize code
// Search for the auction with a given an ID
func (c *Client) Search(auctions []warcraft.Auction, search int) warcraft.Auction {
	var ac warcraft.Auction

	for _, auction := range auctions {
		if auction.ID == search {
			ac = auction
			return ac
		}
	}

	return ac
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

// AddAuctionsThatEndedToBuyoutsCollection checks which auctions have ended and adds the ones that ended to the buyout collection
func (c *Client) AddAuctionsThatEndedToBuyoutsCollection(occurencesHash map[int]int, auctions []warcraft.Auction) error {
	buyoutAuctions := make([]warcraft.Buyout, 0)

	for key, value := range occurencesHash {
		if value == 0 {
			auction := c.Search(auctions, key)
			buyout := c.CreateBuyout(auction)
			buyoutAuctions = append(buyoutAuctions, buyout)
		}
	}
	if err := c.DatabaseService.Insert(BuyoutsCollection, nil, buyoutAuctions); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error("Failed to insert buyout")
		return err
	}

	return nil
}

// GetChunkOfAuctions returns a chunk of auctions from the same timestamp
func (c *Client) GetChunkOfAuctions(index int, auctions []warcraft.Auction) ([]warcraft.Auction, int) {
	// get the first timestamp,
	firstTimestamp := auctions[index].Timestamp
	chunk := make([]warcraft.Auction, 0)
	chunk = append(chunk, auctions[index])
	i := index + 1
	if i >= len(auctions) {
		return chunk, i
	}

	for i < len(auctions) {
		if auctions[i].Timestamp == firstTimestamp {
			chunk = append(chunk, auctions[i])
			i++
		} else {
			break
		}
	}

	return chunk, i
}

// Service returns the service for interacting with the dump analyzer
func (c *Client) Service() warcraft.AnalyzerService { return &c.service }
