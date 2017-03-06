// Showfetcher is a simple glue layer between ShowRSS.info and your Transmission bittorrent client. See README.md in the Github repo for more information.
package main

import (
	"log"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/feed"
	"github.com/doenietzomoeilijk/showfetcher/storage"
	"github.com/doenietzomoeilijk/showfetcher/torrent"
)

func main() {
	log.Println("Getting config, setting up storage and Transmission client...")
	conf := config.Load()
	storage.Setup()
	defer storage.Close()
	torrent.Setup(
		conf.Transmission,
		conf.IncompleteDir,
	)

	log.Println("Find and clean the torrents we already had going...")
	episodes, err := torrent.Cleanup()
	if err != nil {
		log.Fatalln("Couldn't run torrent cleanup:", err)
		return
	}
	doneEps := len(episodes)
	storage.MarkDone(episodes)

	log.Println("Fill the database with show info, based on the RSS feed...")
	episodes = feed.Parse(conf.FeedURL)
	feedEps := len(episodes)
	storage.Store(episodes)

	log.Println("Figuring out new episodes and adding them to the tracker...")
	episodes = storage.Get()
	newEps := len(episodes)
	torrent.Add(episodes)

	ts := torrent.List()
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
