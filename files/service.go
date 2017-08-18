package files

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
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
	for _, filename := range filenames {
		s.client.LoadFileIntoDatabase(filename)
	}

	return nil
}
