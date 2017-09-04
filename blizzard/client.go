package blizzard

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"net/http"
	"net/url"
	"sort"
	"time"
)

// Host is the target api location.
const Host string = "https://eu.api.battle.net/wow/auction/data/"

// AuctionCollection is the default name for the auction collection.
const AuctionCollection = "auctions"

// Client represents a client for interacting with the blizzard API.
type Client struct {
	// package logger.
	logger *log.Entry

	// target api location.
	Host string

	// time between requests.
	Timer  int
	ticker *time.Ticker

	// default configs.
	Realm  string
	Locale string
	Key    string

	// last request time.
	Last int64

	// service for interacting with the blizzard API.
	service Service

	// database service.
	DatabaseService warcraft.DatabaseService

	// service to analyze dump data
	AnalyzerService warcraft.AnalyzerService
}

// NewClient returns a new scrapper client.
func NewClient(timer int, realm, locale, key string) *Client {
	c := &Client{
		logger: log.WithFields(log.Fields{"package": "blizzard"}),
		Host:   Host,
		Timer:  timer,
		Realm:  realm,
		Locale: locale,
		Key:    key,
	}
	c.service.client = c
	return c
}

// Open starts the api daemon.
func (c *Client) Open() error {
	// start ticker.
	c.ticker = time.NewTicker(time.Second * time.Duration(c.Timer))
	go c.handleRequests()
	c.logger.WithFields(log.Fields{"host": c.Host}).Info("started the API request daemon")
	return nil
}

// Close stops the connection to the database.
func (c *Client) Close() error {
	c.ticker.Stop()
	c.logger.Info("stopped making api requests")
	return nil
}

// GetDumpURL makes a request to the Blizzard API for a dump url.
func (c *Client) GetDumpURL(realm, locale, key string) (*warcraft.Request, error) {
	// build request url.
	u, err := url.Parse(c.Host + realm)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errInvalidURL)
		return nil, err
	}

	// build request query.
	q := u.Query()
	q.Set("realm", realm)
	q.Set("locale", locale)
	q.Set("apikey", key)
	u.RawQuery = q.Encode()
	c.logger.WithFields(log.Fields{"url": u.String()}).Debug("built request url")

	// make http request.
	response, err := http.Get(u.String())
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errFailedRequest)
		return nil, err
	}
	defer response.Body.Close()

	// check status code.
	if response.StatusCode != 200 {
		c.logger.WithFields(log.Fields{"code": response.StatusCode, "url": u.String()}).Error(errInvalidResponse)
		return nil, errInvalidResponse
	}

	// decode response.
	var r APIRequest
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&r); err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errInvalidResponse)
		return nil, err
	}

	c.logger.WithFields(log.Fields{"url": u.String(), "dump": r.Requests[0].URL, "timestamp": r.Requests[0].Modified}).Info("dump url retrieved")
	return &r.Requests[0], nil
}

// APIRequest stores the response with url for requesting the dump.
type APIRequest struct {
	Requests []warcraft.Request `json:"files"`
}

// GetDump makes a request to the Blizzard API for a auction house dump.
func (c *Client) GetDump(path string) (*warcraft.APIDump, error) {
	// build request url.
	u, err := url.Parse(path)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errInvalidURL)
		return nil, err
	}

	// make http request.
	response, err := http.Get(u.String())
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errFailedRequest)
		return nil, err
	}
	defer response.Body.Close()

	// check status code.
	if response.StatusCode != 200 {
		c.logger.WithFields(log.Fields{"code": response.StatusCode, "url": u.String()}).Error(errInvalidResponse)
		return nil, errInvalidResponse
	}

	// decode response.
	var r warcraft.APIDump
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&r); err != nil {
		c.logger.WithFields(log.Fields{"error": err, "url": u.String()}).Error(errInvalidResponse)
		return nil, err
	}

	c.logger.WithFields(log.Fields{"url": u.String()}).Info("dump retrieved")
	return &r, nil
}

// handleRequests makes periodic requests to the API.
func (c *Client) handleRequests() {
	for range c.ticker.C {
		// get dump.
		d, err := c.Service().GetAPIDump(c.Realm, c.Locale, c.Key, c.Last)
		if err != nil {
			c.logger.WithFields(log.Fields{"error": err}).Warn(errGetDump)
			continue
		}

		// filter dump.
		a, err := c.service.ValidateAuctions(d)
		if err != nil {
			c.logger.WithFields(log.Fields{"error": err}).Warn(errFilterDump)
			continue
		}

		// save dump.
		if err := c.DatabaseService.Insert(AuctionCollection, a, nil); err != nil {
			c.logger.WithFields(log.Fields{"error": err}).Warn(errSaveDump)
			continue
		}

		// TODO fix this
		// save auctions that ended
		if c.Last != 0 {
			c.AnalyzerService.AnalyzeDumps(c.Last, a)

			// update dump timestamp.
			c.Last = d.Timestamp
			c.logger.Info("handled new api dump")
		} else {
			// get all the starting auctions
			auctions, err := c.DatabaseService.GetAuctions(AuctionCollection, nil)
			if err != nil {
				continue
			}

			// sort auctions by timestamp
			sort.Slice(auctions, func(i, j int) bool {
				return auctions[i].Timestamp < auctions[j].Timestamp
			})

			c.Last = auctions[len(auctions)-1].Timestamp
		}
	}
}

// Service returns the service for interacting with the blizzard API.
func (c *Client) Service() warcraft.BlizzardService { return &c.service }
