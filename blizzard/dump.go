package blizzard

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"net/http"
	"net/url"
)

// NewDump makes a request to the Blizzard API for a HA dump.
func NewDump(path string) (*warcraft.APIDump, error) {
	// build request url.
	u, err := url.Parse(path)
	if err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "url": u.String()}).Error(errInvalidURL)
		return nil, err
	}

	log.WithFields(log.Fields{"package": Package, "url": u.String()}).Debug("making dump request")

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
	r := warcraft.APIDump{}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&r); err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "url": u.String()}).Error(errInvalidResponse)
		return nil, err
	}

	log.WithFields(log.Fields{"package": Package, "url": u.String()}).Debug("dump retrieved")
	return &r, nil
}
