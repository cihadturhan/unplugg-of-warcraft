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

// NewRequest makes a request to the Blizzard API for a dump url.
func NewRequest(config *warcraft.Config) (*warcraft.Request, error) {
	// build request url.
	u, err := url.Parse(Host + config.Realm)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to load host url")
		return nil, err
	}

	// build request query.
	q := u.Query()
	q.Set("realm", config.Realm)
	q.Set("locale", config.Locale)
	q.Set("apikey", config.Key)
	u.RawQuery = q.Encode()

	// make request.
	response, err := http.Get(u.String())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to make request")
		return nil, err
	}
	defer response.Body.Close()

	// decode response.
	r := warcraft.ApiRequest{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&r); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to parse request response")
		return nil, err
	}
	return &r.Requests[0], nil
}
