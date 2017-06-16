package blizzard

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"net/http"
	"net/url"
)

// NewDump makes a request to the Blizzard API for a HA dump.
func NewDump(path string) (*warcraft.ApiDump, error) {
	// build request url.
	u, err := url.Parse(path)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to load host url")
		return nil, err
	}

	// make request.
	response, err := http.Get(u.String())
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to make request")
		return nil, err
	}
	defer response.Body.Close()

	// decode response.
	r := warcraft.ApiDump{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&r); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to parse dump response")
		return nil, err
	}
	return &r, nil
}
