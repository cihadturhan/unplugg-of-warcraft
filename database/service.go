package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"time"
)

// AuctionCollection is the default name for the auction collection.
const AuctionCollection = "auctions"

// Service represents a service for interacting with the database.
type Service struct {
	client *Client
}

// InsertAuctions insets a slice of auctions into the database.
func (s *Service) InsertAuctions(auctions []warcraft.Auction) error {
	start := time.Now()

	// filter auctions.
	as := make([]interface{}, 0)
	for _, a := range auctions {
		if a.Timeleft != "SHORT" && a.Buyout != 0 {
			as = append(as, a)
		}
	}

	// insert auctions.
	for i := 0; i < len(as); i = i + 1000 {
		end := i + 1000
		if end > len(as) {
			end = len(as) - 1
		}
		if err := s.client.InsertAuctions(as[i:end]); err != nil {
			return nil
		}
	}

	s.client.logger.WithFields(log.Fields{"count": len(as), "time": time.Since(start)}).Info("auctions inserted")
	return nil
}
