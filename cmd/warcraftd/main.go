package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/whitesmith/unplugg-of-warcraft/analyzer"
	"github.com/whitesmith/unplugg-of-warcraft/blizzard"
	"github.com/whitesmith/unplugg-of-warcraft/database"
	"github.com/whitesmith/unplugg-of-warcraft/files"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// setup graceful shutdown.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	log.Info("Starting Up")

	// parse flags.
	var (
		debug         = flag.Bool("debug", false, "enable debug level logs")
		loadDumpFiles = flag.Bool("loadDumpFiles", false, "load dump files")
		realm         = flag.String("realm", "grim-batol", "target realm to fetch data")
		locale        = flag.String("locale", "en_GB", "server locale information")
		apiKey        = flag.String("apikey", "", "api authentication key")
		timer         = flag.Int("timer", 2700, "time between api requests")
		mongoURL      = flag.String("mongoUrl", "localhost/warcraft", "mongoDB url")
	)
	flag.Parse()

	// set log level.
	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	// create database client and service.
	dbClient := database.NewClient(*mongoURL)
	if err := dbClient.Open(); err != nil {
		panic(err)
	}
	dbService := dbClient.Service()

	// create analyzer client and service
	analyzerClient := analyzer.NewClient()
	analyzerClient.DatabaseService = dbService
	analyzerService := analyzerClient.Service()

	// create scraper client and service.
	scraperClient := blizzard.NewClient(*timer, *realm, *locale, *apiKey)
	scraperClient.DatabaseService = dbService
	scraperService := scraperClient.Service()

	// load files to the database.
	if *loadDumpFiles {
		// load all the auctions from the files
		filesClient := files.NewClient()
		filesClient.DatabaseService = dbService
		filesClient.BlizzardService = scraperService
		filesService := filesClient.Service()
		filesService.LoadFilesIntoDatabase("./")

		// load buyouts
		auctions, _ := dbService.GetAuctions("auctions")
		analyzerService.AnalyzeDumpsFirstTime(auctions)
		return
	} else {
		if err := scraperClient.Open(); err != nil {
			panic(err)
		}
	}

	// graceful shutdown.
	<-sigs
	log.Info("Shutting Down")
}
