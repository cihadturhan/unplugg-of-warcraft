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
	"time"
)

func main() {
	// loads user flags.
	var (
		realm         = flag.String("realm", "", "target realm to fetch data")
		locale        = flag.String("locale", "", "server locale information")
		apikey        = flag.String("apikey", "", "api authentication key")
		mongoUrl      = flag.String("mongoUrl", "", "mongoDB url")
		configPath    = flag.String("config", ".env", "config file path")
		debug         = flag.Bool("debug", false, "enable debug level logs")
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
		return nil, err
	}

	// unmarshal data
	err = json.Unmarshal(rawFile, &dump)
	if err != nil {
		return nil, err
	}

	return dump.Auctions, nil
}

func loadFileIntoDatabase(filename string, db *mgo.Database) error {
	collection := db.C("auctions")

	auctions, _ := readFile(filename)
	validAuctions := buildValidAuctionsSlice(auctions)

	if err := collection.Insert(validAuctions...); err != nil {
		log.WithFields(log.Fields{"error": err, "filename": filename}).Error("Failed to load file to database")
		return err
	}

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
		if err := loadFileIntoDatabase(filename, db); err != nil {
			return err
		}

		log.WithFields(log.Fields{"dump": filename}).Info("new dump created")
	}

	return nil
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
