package blizzard

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"net/http"
	"net/url"
)

// Host is the target location for the request.
const Host string = "https://eu.api.battle.net/wow/auction/data/"

// Package is the name of current package.
const Package string = "blizzard"

// buildRequestQuery builds the params used to query the Blizzard Api
func buildRequestQuery(config *warcraft.Config, u *url.URL, q url.Values) {
	q.Set("realm", config.Realm)
	q.Set("locale", config.Locale)
	q.Set("apikey", config.Key)
	u.RawQuery = q.Encode()
}

// NewRequest makes a request to the Blizzard API for a dump url.
func NewRequest(config *warcraft.Config) (*warcraft.Request, error) {
	// build request url.
	u, err := url.Parse(Host + config.Realm)
	if err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "url": u.String()}).Error(errInvalidURL)
		return nil, err
	}

	// build request query
	q := u.Query()
	buildRequestQuery(config, u, q)
	log.WithFields(log.Fields{"package": Package, "url": u.String()}).Debug("making http request")

	// make request.
	response, err := http.Get(u.String())
	if err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "url": u.String()}).Error(errFailedRequest)
		return nil, err
	}
	defer response.Body.Close()

	// check status code.
	if response.StatusCode != 200 {
		log.WithFields(log.Fields{"package": Package, "code": response.StatusCode, "url": u.String()}).Error(errInvalidResponse)
	}

	// decode response.
	r := warcraft.APIRequest{}
	decoder := json.NewDecoder(response.Body)

	if err := decoder.Decode(&r); err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "url": u.String()}).Error(errInvalidResponse)
		return nil, err
	}

	log.WithFields(log.Fields{"package": Package, "url": u.String(), "dump": r.Requests[0].URL, "timestamp": r.Requests[0].Modified}).Debug("dump url retrieved")
	return &r.Requests[0], nil
}
