package main

import (
	"log"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	"github.com/doenietzomoeilijk/showfetcher/feed"
	"github.com/doenietzomoeilijk/showfetcher/storage"
	"github.com/doenietzomoeilijk/showfetcher/torrent"
)

var shows []episode.Show

func main() {
	conf := config.Load()

	storage.Setup()
	defer storage.Close()

	torrent.Setup(
		conf.Transmission,
		conf.IncompleteDir,
	)

	// Find and clean the torrents we already had going
	episodes, err := torrent.Cleanup()
	if err != nil {
		log.Fatalln("Couldn't run torrent cleanup:", err)
		return
	}

	doneEps := len(episodes)
	storage.MarkDone(episodes)

	// Fill the db with show info, based on the RSS feed.
	episodes = feed.Parse(conf.FeedURL)
	feedEps := len(episodes)
	storage.Store(episodes)

	// Now fetch anything with status new and add it to the tracker
	episodes = storage.Get()
	newEps := len(episodes)
	torrent.Add(episodes)

	ts := torrent.List()
	for _, t := range ts {
		log.Println(t.Status, t.Name, t.PercentDone, t.IsFinished, t.RateUpload)
	}
	torEps := len(ts)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf(
		"Done! %d were done, %d found in feed, %d were new, %d currently being fetched.\n",
		doneEps,
		feedEps,
		newEps,
		torEps,
	)
}
