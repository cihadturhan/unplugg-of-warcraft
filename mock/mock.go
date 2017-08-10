package mock

import "github.com/whitesmith/unplugg-of-warcraft"

// DatabaseService used to mock this service.
type DatabaseService struct {
	InsertAuctionsFn func(auctions []warcraft.Auction) error
}

// InsertAuctions mocked method.
func (s *DatabaseService) InsertAuctions(auctions []warcraft.Auction, timestamp int64) error {
	return s.InsertAuctionsFn(auctions)
}
