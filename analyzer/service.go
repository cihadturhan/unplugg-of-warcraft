package analyzer

import (
	"github.com/whitesmith/unplugg-of-warcraft"
	"sort"
)

// Service represents a service for interacting with the dump files.
type Service struct {
	client *Client
}

// AnalyzeDumpsFirstTime go through all the dumps and analyzes them
func (s *Service) AnalyzeDumpsFirstTime(auctions []warcraft.Auction) {
	// sort auctions by timestamp
	sort.Slice(auctions, func(i, j int) bool {
		return auctions[i].Timestamp < auctions[j].Timestamp
	})

	index := 0
	var chunk []warcraft.Auction
	var nextChunk []warcraft.Auction

	// Check when index is final in GetChunkOfAuctions
	for index < len(auctions) {
		chunk, index = s.client.GetChunkOfAuctions(index, auctions)

		if index >= len(auctions) {
			return
		}

		nextChunk, index = s.client.GetChunkOfAuctions(index, auctions)
		occurencesHash := s.client.NumberOfOcurrences(chunk, nextChunk)
		if err := s.client.AddAuctionsThatEndedToBuyoutsCollection(occurencesHash, auctions); err != nil {
			continue
		}
	}
}
