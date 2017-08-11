package warcraft

// BlizzardService handles interaction with the blizzard API.
type BlizzardService interface {
	GetAPIDump(realm, locale, key string, last int64) (*APIDump, error)
	ValidateAuctions(dump *APIDump) ([]Auction, error)
}

// DatabaseService handles interaction with the database.
type DatabaseService interface {
	InsertAuctions(auctions []Auction) error
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
