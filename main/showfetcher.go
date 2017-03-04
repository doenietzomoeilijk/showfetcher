package main

import (
	"database/sql"
	"log"

	"github.com/doenietzomoeilijk/showfetcher"
	_ "github.com/mattn/go-sqlite3" // Use sqlite3
	"github.com/odwrtw/transmission"
)

var (
	db     *sql.DB
	tr     *transmission.Client
	config *showfetcher.Configuration
)

func main() {
	config = showfetcher.LoadConfig()
	db, err := showfetcher.SetupDb()
	if err != nil {
		log.Panicln("DB error:", err)
		panic("DBConfig can't be read:")
	}
	defer db.Close()

	// Find and clean the torrents we already had going
	tr = showfetcher.SetupTransmission(config)
	showfetcher.CleanTorrents(tr, db)

	// Fill the db with show info, based on the RSS feed.
	showfetcher.ParseFeeds(config, db)

	// Now fetch anything with status new and add it to the tracker
	showfetcher.RunNewTorrents(config, tr, db)

	ts := showfetcher.ListTorrents(tr)
	for _, t := range ts {
		log.Println(t.Status, t.Name, t.PercentDone, t.IsFinished, t.RateUpload)
	}
}
