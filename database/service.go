package database

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"time"
)

// Service represents a service for interacting with the database.
type Service struct {
	client *Client
}

// InsertAuctions insets a slice of auctions into the database.
func (s *Service) InsertAuctions(collectionName string, auctions []warcraft.Auction) error {
	start := time.Now()

	// convert auctions to interface.
	as := make([]interface{}, 0)
	for _, a := range auctions {
		as = append(as, a)
	}

	// batch auctions.
	for i := 0; i < len(as); i = i + 1000 {
		end := i + 1000
		if end > len(as) {
			end = len(as) - 1
		}
		if err := s.client.InsertAuctions(collectionName, as[i:end]); err != nil {
			return nil
		}
	}

	s.client.logger.WithFields(log.Fields{"count": len(as), "time": time.Since(start)}).Info("auctions inserted")
	return nil
}

// GetAuctions returns all the auctions
func (s *Service) GetAuctions(collectionName string) ([]warcraft.Auction, error) {
	auctions, err := s.client.GetAuctions(collectionName)
	if err != nil {
		return nil, err
	}

	return auctions, nil
}

// GetAuctionsInTimeStamp returns all the auctions present in the timestamp provided
func (s *Service) GetAuctionsInTimeStamp(collectionName string, timestamp int64) ([]warcraft.Auction, error) {
	auctions, err := s.client.GetAuctionsInTimeStamp(collectionName, timestamp)
	if err != nil {
		return nil, err
	}

	return auctions, nil
}
