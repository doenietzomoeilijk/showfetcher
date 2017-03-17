// Showfetcher is a simple glue layer between ShowRSS.info and your Transmission bittorrent client. See README.md in the Github repo for more information.
package main

import (
	"log"
	"os"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/feed"
	"github.com/doenietzomoeilijk/showfetcher/storage"
	"github.com/doenietzomoeilijk/showfetcher/torrent"
)

func main() {
	conf := config.Load()

	f, err := os.OpenFile(conf.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Setting up storage and Transmission client...")
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
		"Done! %d done, %d found in feed, %d new, %d currently being fetched.\n",
		doneEps,
		feedEps,
		newEps,
		torEps,
	)
}
