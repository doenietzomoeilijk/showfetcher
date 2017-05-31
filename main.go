// Showfetcher is a simple glue layer between ShowRSS.info and your Transmission bittorrent client. See README.md in the Github repo for more information.
package main

import (
	"log"
	"os"
	"fmt"
	"strings"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/feed"
	"github.com/doenietzomoeilijk/showfetcher/storage"
	"github.com/doenietzomoeilijk/showfetcher/torrent"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	"gopkg.in/telegram-bot-api.v4"
)

var (
	bot *tgbotapi.BotAPI
	rcp []string
)

func main() {
	conf := config.Load()
	rcp = conf.BotRecipients

	f, err := os.OpenFile(conf.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Println("Setting up Telegram bot...")
	bot, err = tgbotapi.NewBotAPI(conf.BotToken)
	if err != nil {
		log.Fatalln("Couldn't connect to Telegram:", err)
	}

	log.Println("Setting up storage and Transmission client...")
	defer storage.Close()
	torrent.Setup(
		conf.Transmission,
		conf.IncompleteDir,
		conf.SeasonFolders,
	)

	log.Println("Find and clean the torrents we already had going...")
	episodes, err := torrent.Cleanup()
	if err != nil {
		log.Fatalln("Couldn't run torrent cleanup:", err)
		return
	}
	doneEpCount := len(episodes)
	if doneEpCount > 0 {
		spam(fmt.Sprintf("Klaar met downloaden van %d episodes", doneEpCount))
	}
	storage.MarkDone(episodes)

	log.Println("Fill the database with show info, based on the RSS feed...")
	episodes = feed.Parse(conf.FeedURL)
	feedEpCount := len(episodes)
	storage.Store(episodes)

	log.Println("Figuring out new episodes and adding them to the tracker...")
	episodes = storage.Get()
	newEpCount := len(episodes)
	torrent.Add(episodes)
	if newEpCount > 0 {
		spam(fmt.Sprintf("Toegevoegd aan downloads: %s", formatEpList(episodes)))
	}

	ts := torrent.List()
	fetchingEpCount := len(ts)

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf(
		"Done! %d done, %d found in feed, %d new, %d currently being fetched.\n",
		doneEpCount,
		feedEpCount,
		newEpCount,
		fetchingEpCount,
	)
}

func formatEpList(eps []*episode.Episode) string {
	epTitles := []string{}
	for _, ep := range eps {
		epTitles = append(epTitles, fmt.Sprintf("%s (%s)", ep.Show.Title, ep.Episode))
	}

	return strings.Join(epTitles, ", ")
}

func spam(msg string) {
	for _, r := range rcp {
		msg := tgbotapi.NewMessageToChannel(r, msg)
		_, err := bot.Send(msg)
		if err != nil {
			log.Println("Could not send Telegram message:", err)
		}
	}
}
