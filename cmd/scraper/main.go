package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/config"
	"gopkg.in/mgo.v2"
	"time"
)

func main() {
	// loads user flags.
	var (
		realm      = flag.String("realm", "", "target realm to fetch data")
		locale     = flag.String("locale", "", "server locale information")
		apikey     = flag.String("apikey", "", "api authentication key")
		mongoUrl   = flag.String("mongoUrl", "", "mongoDB url")
		configPath = flag.String("config", ".env", "config file path")
		debug      = flag.Bool("debug", false, "enable debug level logs")
	)
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// builds configuration.
	c, err := config.NewConfig(*realm, *locale, *apikey, *mongoUrl, *configPath)

	if err != nil {
		panic(err)
	}

	log.WithFields(log.Fields{"configPath": *configPath}).Info("starting crawler")

	lastDump := 0
	for {
		n, err := getDump(c, lastDump)

		if err == nil {
			lastDump = n
		}
		time.Sleep(30 * time.Minute)
	}
}

// auctionIsValid cheks if an auction is valid
func auctionIsValid(auction warcraft.Auction) bool {
	if auction.Timeleft == "SHORT" {
		return false
	}

	return true
}

// buildValidAuctionsSlice takes the auctions array and returns an slice with the valid auctions to be inserted into the DB
func buildValidAuctionsSlice(allAuctions []warcraft.Auction) []interface{} {
	validAuctions := make([]interface{}, 0)

	for _, auction := range allAuctions {
		if auctionIsValid(auction) {
			validAuctions = append(validAuctions, auction)
		}
	}

	return validAuctions
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

	// open Mongo database
	db, err := openDatabase(c.MongoUrl)

	if err != nil {
		return 0, err
	}

	// get collection and insert valid auctions into the database
	collection := db.C("auctions")
	validAuctions := buildValidAuctionsSlice(d.Auctions)

	err = collection.Insert(validAuctions...)
	if err != nil {
		return 0, err
	}

	log.WithFields(log.Fields{"dump": r.Modified}).Info("new dump created")
	return r.Modified, nil
}

// Connects to MongoDB, establishes a session and returns the database
func openDatabase(url string) (*mgo.Database, error) {
	log.WithFields(log.Fields{"mongoUrl": url}).Info("Opening mongodb session")
	session, err := mgo.Dial(url)

	if nil != err {
		log.WithFields(log.Fields{"error": err, "url": url}).Error("failed to connect to MongoDB")
		return nil, err
	}

	session.EnsureSafe(&mgo.Safe{FSync: true, J: true})
	log.WithFields(log.Fields{"database": config.MongoDBDatabase}).Info("Opening database")

	return session.DB(config.MongoDBDatabase), nil
}
