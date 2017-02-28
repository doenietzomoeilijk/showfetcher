package showfetcher

import (
	"log"
	"net/url"
	"strings"

	"database/sql"

	"github.com/mmcdole/gofeed"
)

func ParseFeeds(config *Configuration, db *sql.DB) {
	feed, _ := gofeed.NewParser().ParseURL(config.FeedURL)
	for _, item := range feed.Items {
		show, ok := config.findShow(item.Title)

		if !ok {
			log.Println("I cannot haz show")
			continue
		}
		log.Printf("Entry '%s' matches show '%s'\n", item.Title, show.Title)

		matches := ShowRe.FindStringSubmatch(item.Title)
		if len(matches) < 1 {
			log.Println("Didn't find anything!", matches, item.Title)
			continue
		}

		magnet, _ := url.ParseQuery(strings.Split(item.Link, "?")[1])
		episode := matches[2]
		_, err := db.Exec(
			`INSERT OR IGNORE INTO shows
                (hash, show, episode, published, filename, magnet)
                VALUES (?, ?, ?, ?, ?, ?)`,
			strings.Split(magnet["xt"][0], ":")[2],
			show.Title,
			episode,
			item.PublishedParsed,
			strings.Replace(magnet["dn"][0], " ", ".", -1),
			item.Link,
		)
		if err == nil {
			log.Println("Stored episode", show.Title, episode)
		} else {
			log.Println("Could not store show:", err)
		}
	}
}
