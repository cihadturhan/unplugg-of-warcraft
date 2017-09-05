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

// Search searchs for an auction ID in an array
func (c *Client) Search(searchID int, auctions []warcraft.Auction) bool {
	for _, auction := range auctions {
		if auction.ID == searchID {
			return true
		}
	}

	return false
}

// CreateAuctionsMap returns an map with auctionID has key and a pointer to auction has the value
func (c *Client) CreateAuctionsMap(auctions []warcraft.Auction) map[int]warcraft.Auction {
	auctionsHash := make(map[int]warcraft.Auction, 0)

	for _, auction := range auctions {
		auctionsHash[auction.ID] = auction
	}

	return auctionsHash
}

// AuctionsThatEnded returns all the auctions that not present in the next dump (ended)
func (c *Client) AuctionsThatEnded(prevAuctions map[int]warcraft.Auction, auctions []warcraft.Auction) []warcraft.Buyout {
	results := make([]warcraft.Buyout, 0)

	for key, value := range prevAuctions {
		if c.Search(key, auctions) {
			results = append(results, c.CreateBuyout(value))
		}
	}

	return results
}

// AddAuctionsThatEndedToBuyoutsCollection checks which auctions have ended and adds the ones that ended to the buyout collection
func (c *Client) AddAuctionsThatEndedToBuyoutsCollection(prevAuctions []warcraft.Auction, auctions []warcraft.Auction) error {
	prevAuctionsMap := c.CreateAuctionsMap(prevAuctions)
	buyouts := c.AuctionsThatEnded(prevAuctionsMap, auctions)

	if err := c.DatabaseService.Insert(BuyoutsCollection, nil, buyouts); err != nil {
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
