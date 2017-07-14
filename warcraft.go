package warcraft

// Config stores the crawler configuration.
type Config struct {
	Realm    string `json:"realm"`
	Locale   string `json:"locale"`
	Key      string `json:"apikey"`
	MongoUrl string `json:"mongoUrl"`
}

// APIRequest stores the response with url for requesting the dump.
type APIRequest struct {
	Requests []Request `json:"files"`
}

// Request stores the url and timestamp of the requested dump.
type Request struct {
	URL      string `json:"url"`
	Modified int    `json:"lastModified"`
}

// APIDump stores the dump requested.
type APIDump struct {
	Realms   []Realm   `json:"realms"`
	Auctions []Auction `json:"auctions"`
}

// Realm stores a realm metadata.
type Realm struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// Auction stores an auction metadata.
type Auction struct {
	ID       int    `json:"auc"`
	Item     int    `json:"item"`
	Player   string `json:"owner"`
	Realm    string `json:"ownerRealm"`
	Bid      int    `json:"bid"`
	Buyout   int    `json:"buyout"`
	Quantity int    `json:"quantity"`
	Timeleft string `json:"timeLeft"`
	Rand     int    `json:"rand"`
	Seed     int    `json:"seed"`
	Context  int    `json:"context"`
}
