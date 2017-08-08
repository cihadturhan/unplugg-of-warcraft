package main

import (
	"encoding/json"
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/config"
	"gopkg.in/mgo.v2"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// which collection should we dump data into?
const collectionForDump string = "auctions"
const sleepTime time.Duration = 30 * time.Minute

func main() {
	// loads user flags.
	var (
		realm         = flag.String("realm", "", "target realm to fetch data")
		locale        = flag.String("locale", "", "server locale information")
		apikey        = flag.String("apikey", "", "api authentication key")
		mongoUrl      = flag.String("mongoUrl", "", "mongoDB url")
		configPath    = flag.String("config", ".env", "config file path")
		debug         = flag.Bool("debug", true, "enable debug level logs")
		loadDumpFiles = flag.Bool("loadDumpFiles", false, "load api dump files")
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

	if *loadDumpFiles {
		loadFilesIntoDatabase(c, "./")
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

// isDumpFile checks if the file is a dump or not
func isDumpFile(filename string) bool {
	if filename[0:1] != "1" {
		return false
	}
	return true
}

// buildFilenamesSlice returns a slice with all the api dump files
func buildFilenamesSlice(files []os.FileInfo) []string {
	dumpFiles := make([]string, 0)

	for _, file := range files {
		filename := file.Name()

		if isDumpFile(filename) {
			dumpFiles = append(dumpFiles, filename)
		}
	}

	return dumpFiles
}

// readFile takes a filename reads the file and decodes it to json
func readFile(filename string) ([]warcraft.Auction, error) {
	dump := warcraft.APIDump{}

	// get file binary data
	rawFile, err := ioutil.ReadFile(filename)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "filename": filename}).Error("Failed to read raw file data")
		return nil, err
	}

	// unmarshal data
	err = json.Unmarshal(rawFile, &dump)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "filename": filename}).Error("Failed to unmarshal file binary data")
		return nil, err
	}

	return dump.Auctions, nil
}

// removeFile removes a dump file
func removeFile(filename string) error {
	path := "./" + filename

	if err := os.Remove(path); err != nil {
		log.WithFields(log.Fields{"error": err, "path": path}).Error("Failed to remove file")
		return err
	}

	log.WithFields(log.Fields{"filename": filename}).Info("File removed")
	return nil
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

func insertIntoDatabase(auctions []interface{}, collection *mgo.Collection) error {
	if err := collection.Insert(auctions...); err != nil {
		return err
	}
	return nil
}

// insertIntoDabase splits the auctions slice into smaller slices to insert
func insertBySlices(auctions []interface{}, collection *mgo.Collection) error {
	if len(auctions) > 100 {
		if err := insertIntoDatabase(auctions, collection); err != nil {
			return err
		}

		return nil
	}

	low := 0
	high := 100
	for high <= len(auctions) {
		if err := insertIntoDatabase(auctions[low:high], collection); err != nil {
			return err
		}

		low += 100
		high += 100
	}

	if err := insertIntoDatabase(auctions[low:len(auctions)], collection); err != nil {
		return err
	}

	return nil
}

func loadFileIntoDatabase(filename string, db *mgo.Database) error {
	collection := db.C("auctions")

	timestamp, err := strconv.Atoi(filename)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "filename": filename}).Error("Failed to convert filename to int")
		return err
	}

	auctions, _ := readFile(filename)
	validAuctions := buildValidAuctionsSlice(auctions, timestamp)

	if err := insertBySlices(validAuctions, collection); err != nil {
		log.WithFields(log.Fields{"error": err, "filename": filename}).Error("Failed to load file to database")
		return err
	}

	log.WithFields(log.Fields{"dump": filename}).Info("Dump loaded to database")
	return nil
}

// loadFilesIntoDatabase loads the Blizzard API dump files into the DB
func loadFilesIntoDatabase(c *warcraft.Config, path string) error {
	// get all the files in the directory
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.WithFields(log.Fields{"error": err, "path": path}).Error("Failed to load files in directory")
		return err
	}

	// open Mongo database
	db, err := openDatabase(c.MongoUrl)
	if err != nil {
		return err
	}

	// load dump files
	dumpFiles := buildFilenamesSlice(files)
	for _, filename := range dumpFiles {
		if err := loadFileIntoDatabase(filename, db); err == nil {
			removeFile(filename)
		}
	}

	log.Info("Dump files loaded into database")
	return nil
}

// getDump makes a request to the Blizzard API and loads it into de database
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
