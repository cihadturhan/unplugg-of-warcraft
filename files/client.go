package files

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft"
	"io/ioutil"
	"os"
	"strconv"
)

type Client struct {
	// package logger.
	logger *log.Entry

	// default configs.
	Realm  string
	Locale string
	Key    string

	// service for interacting with the dump files.
	service Service

	// database service.
	DatabaseService warcraft.DatabaseService

	// blizzard service.
	BlizzardService warcraft.BlizzardService
}

// NewClient returns a new configuration client.
func NewClient(realm, locale, key string) *Client {
	c := &Client{
		logger: log.WithFields(log.Fields{"package": "files"}),
		Realm:  realm,
		Locale: locale,
		Key:    key,
	}

	c.service.client = c
	return c
}

// IsDumpFile checks if the file is a dump or not
func isDumpFile(filename string) bool {
	if filename[0:1] != "1" {
		return false
	}
	return true
}

// GetFilenames returns all the filenames present in the root directory
func (c *Client) GetFilenames(files []os.FileInfo) []string {
	dumpFiles := make([]string, 0)

	for _, file := range files {
		filename := file.Name()

		if isDumpFile(filename) {
			dumpFiles = append(dumpFiles, filename)
		}
	}

	return dumpFiles
}

// Read reads a dump file
func (c *Client) Read(filename string) (*warcraft.APIDump, error) {
	dump := warcraft.APIDump{}

	// get file binary data
	rawFile, err := ioutil.ReadFile(filename)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "filename": filename}).Error(errFailedRead)
		return nil, err
	}

	// unmarshal data
	err = json.Unmarshal(rawFile, &dump)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "filename": filename}).Error(errFailedUnmarshal)
		return nil, err
	}

	// set dump timestamp
	dump.Timestamp, err = strconv.ParseInt(filename, 10, 64)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err, "filename": filename}).Error(errFailedStringToInt)
		return nil, err
	}

	return &dump, nil
}

// Remove removes a dump file
func (c *Client) Remove(filename string) error {
	path := "./" + filename

	if err := os.Remove(path); err != nil {
		c.logger.WithFields(log.Fields{"error": err, "path": path}).Error(errFailedRemoveFile)
		return err
	}

	c.logger.WithFields(log.Fields{"filename": filename}).Info("File removed")
	return nil
}

func (c *Client) LoadFileIntoDatabase(filename string) error {

	// read dump file
	dump, err := c.Read(filename)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errFailedRead)
		return err
	}

	// filter dump file
	auctions, err := c.BlizzardService.ValidateAuctions(dump)
	if err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errFailedDumpFilter)
		return err
	}

	// save dump.
	if err := c.DatabaseService.InsertAuctions(auctions); err != nil {
		c.logger.WithFields(log.Fields{"error": err}).Error(errFailedDatabaseSave)
		return err
	}

	return nil
}

// Service returns the service for interacting with the saved dump files
func (c *Client) Service() warcraft.FilesService { return &c.service }
