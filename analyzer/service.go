package analyzer

import (
	"github.com/whitesmith/unplugg-of-warcraft"
	"gopkg.in/mgo.v2/bson"
)

const AuctionCollection = "auctions"

// Service represents a service for interacting with the dump files.
type Service struct {
	client *Client
}

func (s *Service) AnalyzeDumps(lastTimestamp interface{}, newAuctions []warcraft.Auction) {
	query := bson.M{"timestamp": lastTimestamp}

	lastAuctions, err := s.client.DatabaseService.Find(AuctionCollection, query)
	if err != nil {
		return
	}

	if err := s.client.AddAuctionsThatEndedToBuyoutsCollection(lastAuctions, newAuctions); err != nil {
		return
	}

	return
}
