package config

import (
	"os"
	log "github.com/Sirupsen/logrus"
	"encoding/json"
	"github.com/whitesmith/unplugg-of-warcraft"
)

// NewConfig builds a new configuration from the dotenv and the command flags.
func NewConfig(realm, locale, key, path string) (*warcraft.Config, error) {
	// open config file.
	file, err := os.Open(path)
	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to find config file")
		return nil, err
	}
	defer file.Close()

	// decode config file.
	c := warcraft.Config{}
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		log.WithFields(log.Fields{"error": err}).Error("failed to load config file")
		return nil, err
	}

	// check flags.
	if realm != "" {
		c.Realm = realm
	}
	if locale != "" {
		c.Locale = locale
	}
	if key != "" {
		c.Key = key
	}
	return &c, nil
}
