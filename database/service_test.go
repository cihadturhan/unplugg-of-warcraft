package database_test

import (
	"github.com/whitesmith/unplugg-of-warcraft"
	"testing"
	"time"
)

// TestService_InsertAuctions tests inserting auctions to the database.
func TestService_InsertAuctions(t *testing.T) {
	c := MustOpenClient()
	defer c.Close()
	s := c.Service()

	a := []warcraft.Auction{
		{
			ID:       01,
			Item:     34,
			Realm:    "grim-batol",
			Buyout:   234,
			Quantity: 32,
			Timeleft: "LONG",
		},
		{
			ID:       02,
			Item:     34,
			Realm:    "grim-batol",
			Buyout:   234,
			Quantity: 32,
			Timeleft: "SHORT",
		},
		{
			ID:       03,
			Item:     34,
			Realm:    "grim-batol",
			Buyout:   234,
			Quantity: 32,
			Timeleft: "LONG",
		},
		{
			ID:       04,
			Item:     34,
			Realm:    "grim-batol",
			Buyout:   234,
			Quantity: 32,
			Timeleft: "LONG",
		},
	}

	if err := s.InsertAuctions(a, time.Now().Unix()); err != nil {
		t.Fatalf("failed to inster auctions: %v", err)
	}
}
