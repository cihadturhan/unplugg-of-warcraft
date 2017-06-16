package main

import (
	"encoding/json"
	"flag"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/config"
	"io/ioutil"
	"strconv"
)

func main() {
	// loads user flags.
	var (
		realm  = flag.String("realm", "", "target realm to fetch data")
		locale = flag.String("locale", "", "server locale information")
		apikey = flag.String("apikey", "", "api authentication key")
		path   = flag.String("config", "/tmp/.env", "config file path")
	)
	flag.Parse()

	// builds configuration.
	c, err := config.NewConfig(*realm, *locale, *apikey, *path)
	if err != nil {
		panic(err)
	}

	lastDump := 0

	// makes request to get dump url.
	r, err := blizzard.NewRequest(c)
	if err != nil {
		panic(err)
	}

	// gets the HA dump.
	d, err := blizzard.NewDump(r.Url)
	if err != nil {
		panic(err)
	}

	dumpfile, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}

	s := strconv.Itoa(lastDump)
	err = ioutil.WriteFile(s, dumpfile, 0644)
	if err != nil {
		panic(err)
	}
}
