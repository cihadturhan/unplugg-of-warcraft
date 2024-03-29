package warcraft

import (
	"gopkg.in/mgo.v2/bson"
)

const AuctionCollection = "auctions"
const BuyoutCollection = "buyouts"

// BlizzardService handles interaction with the blizzard API.
type BlizzardService interface {
	GetAPIDump(realm, locale, key string, last int64) (*APIDump, error)
	ValidateAuctions(dump *APIDump) ([]Auction, error)
}

// DatabaseService handles interaction with the database.
type DatabaseService interface {
	Insert(collectionName string, records []interface{}) error
	Find(collectionName string, options bson.M) ([]Auction, error)
	GetLastRecord(collectionName string) (Auction, error)
}

// FilesService handles interaction with the saved dump files.
type FilesService interface {
	LoadFilesIntoDatabase(path string) error
}

// AnalyzerService handles interaction with the dump analyzer
type AnalyzerService interface {
	AnalyzeDumps(lastTimestamp interface{}, newAuctions []Auction)
}

// Request stores the url and timestamp of the requested dump.
type Request struct {
	URL      string `json:"url"`
	Modified int64  `json:"lastModified"`
}

// APIDump stores the dump requested.
type APIDump struct {
	Realms    []Realm   `json:"realms"`
	Auctions  []Auction `json:"auctions"`
	Timestamp int64     `json:"timestamp,omitempty"`
}

// Realm stores a realm metadata.
type Realm struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Auction stores an auction metadata.
type Auction struct {
	ID        int    `json:"auc" bson:"auc"`
	Item      int    `json:"item" bson:"item"`
	Realm     string `json:"ownerRealm" bson:"ownerRealm"`
	Buyout    int    `json:"buyout" bson:"buyout"`
	Quantity  int    `json:"quantity" bson:"quantity"`
	Timeleft  string `json:"timeLeft" bson:"timeLeft"`
	Timestamp int64  `json:"timestamp,omitempty" bson:"timestamp"`
}

// Buyout stores an auction that ended.
type Buyout struct {
	ID        int   `json:"auc" bson:"auc"`
	Item      int   `json:"item" bson:"item"`
	Buyout    int   `json:"buyout" bson:"buyout"`
	Quantity  int   `json:"quantity" bson:"quantity"`
	Timestamp int64 `json:"timestamp,omitempty" bson:"timestamp"`
}
