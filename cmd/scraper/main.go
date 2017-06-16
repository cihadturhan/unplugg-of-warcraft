package main

import (
	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/config"
	"io/ioutil"
	"strconv"
	"time"
)

func main() {
	// loads user flags.
	var (
		realm  = flag.String("realm", "", "target realm to fetch data")
		locale = flag.String("locale", "", "server locale information")
		apikey = flag.String("apikey", "", "api authentication key")
		path   = flag.String("config", ".env", "config file path")
	)
	flag.Parse()

	// builds configuration.
	c, err := config.NewConfig(*realm, *locale, *apikey, *path)
	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{path: path}).Info("starting crawler")

	lastDump := 0
	for {
		n, err := getDump(c, lastDump)
		if err != nil {
			lastDump = n
		}
		time.Sleep(30 * time.Minute)
	}
}

func getDump(c *warcraft.Config, last int) (int, error) {
	// makes request to get dump url.
	r, err := blizzard.NewRequest(c)
	if err != nil {
		return 0, err
	}

	// dump already exists.
	if r.Modified == last {
		log.WithFields(log.Fields{"dump": last}).Warn("dump already exists")
		return last, nil
	}

	// gets the AH dump.
	d, err := blizzard.NewDump(r.URL)
	if err != nil {
		return 0, err
	}

	// serialize AH dump.
	dumpfile, err := json.Marshal(d)
	if err != nil {
		return 0, err
	}

	// save AH dump.
	s := strconv.Itoa(r.Modified)
	err = ioutil.WriteFile(s, dumpfile, 0644)
	if err != nil {
		log.WithFields(log.Fields{"dump": r.Modified}).Error("failed to create file")
		return 0, err
	}

	log.WithFields(log.Fields{"dump": r.Modified}).Info("new dump created")
	return r.Modified, nil
}
