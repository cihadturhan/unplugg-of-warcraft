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

// AnalyzeDumps compares the last dump with the new one
func (s *Service) AnalyzeDumps(lastTimestamp int64, newAuctions []warcraft.Auction) {
	var lastAuctions []warcraft.Auction
	var err error

	if lastAuctions, err = s.client.DatabaseService.GetAuctionsInTimeStamp("auctions", lastTimestamp); err != nil {
		return
	}
	occurencesHash := s.client.NumberOfOcurrences(lastAuctions, newAuctions)

	if err = s.client.AddAuctionsThatEndedToBuyoutsCollection(occurencesHash, lastAuctions); err != nil {
		return
	}

	return
}
