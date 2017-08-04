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

// which collection should we dump data into?
const collectionForDump string = "auctions"
const sleepTime time.Duration = 30 * time.Minute

func main() {
	// loads user flags.
	var (
		realm      = flag.String("realm", "", "target realm to fetch data")
		locale     = flag.String("locale", "", "server locale information")
		apikey     = flag.String("apikey", "", "api authentication key")
		mongoUrl   = flag.String("mongoUrl", "", "mongoDB url")
		configPath = flag.String("config", ".env", "config file path")
		debug      = flag.Bool("debug", true, "enable debug level logs")
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

	log.WithFields(log.Fields{"configPath": *configPath}).
		Info("👀 starting crawler 👀")

	// TODO: get last dump timestamp from DB or somthing
	lastDumpTs := 0
	for {
		n, err := getDump(c, lastDumpTs)

		if err == nil {
			lastDumpTs = n
		}

		log.WithFields(log.Fields{"sleepTime": sleepTime}).
			Info("so tired... going to rest a bit 😴")
		time.Sleep(sleepTime)
	}
}

func getDump(cfg *warcraft.Config, previousTimestamp int) (int, error) {
	// makes request to get dump url.
	request, err := blizzard.NewRequest(cfg)
	if err != nil {
		return 0, err
	}
	timestamp := request.Modified

	// dump already exists.
	if timestamp == previousTimestamp {
		log.WithFields(log.Fields{"dump": previousTimestamp}).
			Warn("dump already exists")
		return previousTimestamp, nil
	}

	// gets the AH dump.
	dump, err := blizzard.NewDump(request.URL)
	if err != nil {
		return 0, err
	}

	// get valid auctions from this dump
	log.WithFields(log.Fields{"dump": timestamp}).
		Info("validating auctions")
	validAuctions := buildValidAuctionsSlice(dump.Auctions, timestamp)

	// open Mongo database and store them
	log.WithFields(
		log.Fields{"dump": timestamp, "collection": collectionForDump},
	).Info("storing into collection")
	db, err := openDatabase(cfg.MongoUrl)
	if err != nil {
		return 0, err
	}
	err = db.C(collectionForDump).Insert(validAuctions...)
	if err != nil {
		return 0, err
	}

	log.WithFields(log.Fields{"dump": timestamp}).Info("new dump created")
	return timestamp, nil
}

// buildValidAuctionsSlice takes the auctions array and returns an array
// with only the valid ones
func buildValidAuctionsSlice(
	allAuctions []warcraft.Auction, timestamp int,
) []interface{} {
	validAuctions := make([]interface{}, 0)
	for _, auction := range allAuctions {
		if auctionIsValid(auction) {
			auction.Timestamp = timestamp
			// TODO: we should check auction.ID to prevent duplicates from going into
			// the database
			validAuctions = append(validAuctions, auction)
		}
	}
	return validAuctions
}

// auctionIsValid checks if an auction is valid
func auctionIsValid(auction warcraft.Auction) bool {
	if auction.Timeleft == "SHORT" {
		log.WithFields(log.Fields{"auction": auction}).Debug("invalid TimeLeft")
		return false
	}
	if auction.Buyout == 0 {
		log.WithFields(log.Fields{"auction": auction}).Debug("invalid Buyout")
		return false
	}
	return true
}

// Connects to MongoDB, establishes a session and returns the database
func openDatabase(url string) (*mgo.Database, error) {
	log.WithFields(log.Fields{"mongoUrl": url}).Info("Opening mongodb session")
	session, err := mgo.Dial(url)

	if nil != err {
		log.WithFields(log.Fields{"error": err, "url": url}).
			Error("failed to connect to MongoDB")
		return nil, err
	}

	session.EnsureSafe(&mgo.Safe{FSync: true, J: true})
	log.WithFields(log.Fields{"database": config.MongoDBDatabase}).
		Info("Opening database")

	return session.DB(config.MongoDBDatabase), nil
}
