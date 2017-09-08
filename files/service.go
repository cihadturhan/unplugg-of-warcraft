package files

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"strconv"
)

// Service represents a service for interacting with the dump files.
type Service struct {
	client *Client
}

// LoadFilesIntoDatabase loads the Blizzard API dump files into the DB
func (s *Service) LoadFilesIntoDatabase(path string) error {
	// get all the files in the directory
	f, err := ioutil.ReadDir(path)
	if err != nil {
		s.client.logger.WithFields(log.Fields{"error": err}).Error("Failed to get directory files")
		return err
	}

	filenames := s.client.GetFilenames(f)

	// load files to database
	var lastTimestamp, timestamp int64
	var auctions []warcraft.Auction

	lastTimestamp, _ = strconv.ParseInt(filenames[0], 10, 64)
	if err := s.client.LoadFileIntoDatabase(filenames[0]); err != nil {
		return err
	}

	for _, filename := range filenames[1 : len(filenames)-1] {
		// load file into database
		err := s.client.LoadFileIntoDatabase(filename)

		if err == nil {
			// get timestamp from filename
			timestamp, err = strconv.ParseInt(filename, 10, 64)
			if err != nil {
				s.client.logger.WithFields(log.Fields{"error": err}).Error("Failed to parse filename to timestamp")
				continue
			}

			// build query
			query := bson.M{"timestamp": timestamp}

			// get auctions from timestamp
			auctions, err = s.client.DatabaseService.Find(AuctionCollection, query)
			if err != nil {
				continue
			}

			// add auctions that ended to buyouts collection
			s.client.AnalyzerService.AnalyzeDumps(lastTimestamp, auctions)
			if err != nil {
				continue
			}

			// update last timestamp
			lastTimestamp, err = strconv.ParseInt(filename, 10, 64)
			if err != nil {
				s.client.logger.WithFields(log.Fields{"error": err}).Error("Failed to parse filename to timestamp")
			}
		}
	}

	return nil
}
