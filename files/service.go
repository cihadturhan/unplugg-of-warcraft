package files

import (
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
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

	// TODO error when reading file does mistake here
	// load files to database
	var lastTimestamp, timestamp int
	var auctions []warcraft.Auction

	lastTimestamp, _ = strconv.Atoi(filenames[0])
	if err := s.client.LoadFileIntoDatabase(filenames[0]); err != nil {
		return err
	}

	for _, filename := range filenames[1 : len(filenames)-1] {
		// load file into database
		err := s.client.LoadFileIntoDatabase(filename)

		if err == nil {
			// get timestamp from filename
			timestamp, err = strconv.Atoi(filename)
			if err != nil {
				s.client.logger.WithFields(log.Fields{"error": err}).Error("Failed to parse filename to timestamp")
				continue
			}

			// get auctions from timestamp
			auctions, err = s.client.DatabaseService.Find(AuctionCollection, timestamp)
			if err != nil {
				continue
			}

			// add auctions that ended to buyouts collection
			s.client.AnalyzerService.AnalyzeDumps(lastTimestamp, auctions)
			if err != nil {
				continue
			}

			// update last timestamp
			lastTimestamp, err = strconv.Atoi(filename)
			if err != nil {
				s.client.logger.WithFields(log.Fields{"error": err}).Error("Failed to parse filename to timestamp")
			}
		}
	}

	return nil
}
