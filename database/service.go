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

// InsertAuctions inserts a slice of auctions into the database.
func (s *Service) Insert(collectionName string, records []interface{}) error {
	start := time.Now()
	as := records

	// batch auctions.
	for i := 0; i < len(as); i = i + 1000 {
		end := i + 1000
		if end > len(as) {
			end = len(as) - 1
		}
		if err := s.client.Insert(collectionName, as[i:end]); err != nil {
			return nil
		}
	}

	s.client.logger.WithFields(log.Fields{"count": len(as), "time": time.Since(start)}).Info("auctions inserted")
	return nil
}

// GetAuctions returns all the auctions
func (s *Service) Find(collectionName string, options interface{}) ([]warcraft.Auction, error) {
	auctions, err := s.client.Find(collectionName, options)
	if err != nil {
		return nil, err
	}

	return auctions, nil
}
