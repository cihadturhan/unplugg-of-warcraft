package analyzer

import (
	"github.com/whitesmith/unplugg-of-warcraft"
	"gopkg.in/mgo.v2/bson"
)

// Service represents a service for interacting with the dump files.
type Service struct {
	client *Client
}

// AnalyzeDumps iterates the previous auctions and checks if they still exist in the new dump
// if they don't exist it adds them to they buyouts collection
func (s *Service) AnalyzeDumps(lastTimestamp interface{}, newAuctions []warcraft.Auction) {
	query := bson.M{"timestamp": lastTimestamp}

	lastAuctions, err := s.client.DatabaseService.Find(warcraft.AuctionCollection, query)
	if err != nil {
		return
	}

	if err := s.client.AddAuctionsThatEndedToBuyoutsCollection(lastAuctions, newAuctions); err != nil {
		return
	}

	return
}
