package analyzer

import (
	"github.com/whitesmith/unplugg-of-warcraft"
)

const AuctionCollection = "auctions"

// Service represents a service for interacting with the dump files.
type Service struct {
	client *Client
}

func (s *Service) AnalyzeDumps(lastTimestamp interface{}, newAuctions []warcraft.Auction) {
	lastAuctions, err := s.client.DatabaseService.Find(AuctionCollection, lastTimestamp)
	if err != nil {
		return
	}

	if err := s.client.AddAuctionsThatEndedToBuyoutsCollection(lastAuctions, newAuctions); err != nil {
		return
	}

	return
}
