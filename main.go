package main

import (
	"log"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	"github.com/doenietzomoeilijk/showfetcher/feed"
	"github.com/doenietzomoeilijk/showfetcher/storage"
	"github.com/doenietzomoeilijk/showfetcher/torrent"
	_ "github.com/mattn/go-sqlite3" // Use sqlite3
)

var shows []episode.Show

func main() {
	config.Load()

	storage.Setup()
	defer storage.Close()

	torrent.Setup(config.Config.Transmission, config.Config.IncompleteDir)

	// Find and clean the torrents we already had going
	episodes := torrent.Cleanup()
	storage.MarkDone(episodes)

	// Fill the db with show info, based on the RSS feed.
	episodes = feed.Parse(config.Config.FeedURL)
	storage.Store(episodes)

	// Now fetch anything with status new and add it to the tracker
	episodes = storage.Get()
	torrent.Add(episodes)

	ts := torrent.List()
	for _, t := range ts {
		log.Println(t.Status, t.Name, t.PercentDone, t.IsFinished, t.RateUpload)
	}
}
