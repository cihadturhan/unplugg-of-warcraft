package blizzard_test

import (
	"fmt"
	"github.com/whitesmith/unplugg-of-warcraft"
	"reflect"
	"testing"
	"time"
)

// request path
var path string
var dump *warcraft.APIDump

// TestClient_GetDumpUrl tests retrieving an API dump request.
func TestClient_GetDumpUrl(t *testing.T) {
	c := NewClient()
	r, err := c.GetDumpUrl("grim-batol", "en_GB", "")
	if err != nil {
		t.Fatalf("failed to get dumo url: %v", err)
	}
	path = r.URL
}

// TestClient_GetDump tests retrieving an API dump.
func TestClient_GetDump(t *testing.T) {
	c := NewClient()
	d, err := c.GetDump(path)
	if err != nil {
		t.Fatalf("failed to get dumo url: %v", err)
	}
	dump = d
}

// TestService_GetDump tests retrieving an API dump.
func TestService_GetDump(t *testing.T) {
	c := NewClient()
	d, err := c.Service().GetAPIDump("grim-batol", "en_GB", "", time.Now().Unix())
	if err != nil && reflect.DeepEqual(dump, d) {
		t.Fatalf("failed to get dumo url: %v", err)
	}
}

// TestClient_Daemon tests api request daemon.
func TestClient_Daemon(t *testing.T) {
	c := NewClient()

	// mock database service.
	c.DatabaseService.InsertAuctionsFn = func(auctions []warcraft.Auction) error {
		fmt.Println(len(auctions))
		return nil
	}

	// start client.
	if err := c.Open(); err != nil {
		t.Fatalf("failed to start client: %v", err)
	}
	defer c.Close()
	time.Sleep(5 * time.Minute)
}
