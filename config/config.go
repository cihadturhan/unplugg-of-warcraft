package config

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"os"
)

// Package is the name of current package.
const Package string = "config"

// MongoDBDatabase is the name of the mongo database
const MongoDBDatabase string = "unplugg-of-warcraft"

// setConfigFlags sets the configuration flags
func setConfigFlags(c warcraft.Config, realm, locale, key, mongoUrl, path string) {
	if realm != "" {
		c.Realm = realm
	}
	if locale != "" {
		c.Locale = locale
	}
	if key != "" {
		c.Key = key
	}
	if mongoUrl != "" {
		c.MongoUrl = mongoUrl
	}
}

// NewConfig builds a new configuration from the dotenv and the command flags.
func NewConfig(realm, locale, key, mongoUrl, path string) (*warcraft.Config, error) {
	// open config file.
	file, err := os.Open(path)

	if err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "path": path}).Error(errInvalidFile)
		return nil, err
	}
	defer file.Close()

	// decode config file.
	c := warcraft.Config{}
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&c); err != nil {
		log.WithFields(log.Fields{"package": Package, "error": err, "path": path}).Error(errFailedParse)
		return nil, err
	}

	// set configuration flags
	setConfigFlags(c, realm, locale, key, mongoUrl, path)

	log.WithFields(log.Fields{"package": Package, "config": c}).Debug(errFailedParse)
	return &c, nil
}
