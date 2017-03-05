package feed

import (
	"log"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/mmcdole/gofeed"
)

const ShowPattern = "(.+) ([0-9]+x[0-9]+) ?(.*)"

// Parse grabs and parses the show feed.
func Parse(url string) {
	log.Println("ParseFeeds")
	feed, _ := gofeed.NewParser().ParseURL(url)

	for _, item := range feed.Items {
		log.Println("Feed item", item.Title)
		showname := item.Extensions["tv"]["show_name"][0].Value
		show := config.Showmap[showname]
		hash := item.Extensions["tv"]["info_hash"][0].Value

		log.Println("item show name", showname, hash, show)

		// show, ok := episode.Shows[item.Title]

		// if !ok {
		// 	log.Println("I cannot haz show", item.Title)
		// 	continue
		// }
		// log.Printf("Entry '%s' matches show '%s'\n", item.Title, show.Title)

		// matches := ShowRe.FindStringSubmatch(item.Title)
		// if len(matches) < 1 {
		// 	log.Println("Didn't find anything!", matches, item.Title)
		// 	continue
		// }

		// magnet, _ := url.ParseQuery(strings.Split(item.Link, "?")[1])
		// episode := matches[2]
		// hash := strings.Split(magnet["xt"][0], ":")[2]
		// _, err := db.Exec(
		// 	`INSERT OR IGNORE INTO shows
		//         (hash, show, episode, published, filename, magnet)
		//         VALUES (?, ?, ?, ?, ?, ?)`,
		// 	strings.ToLower(hash),
		// 	show.Title,
		// 	episode,
		// 	item.PublishedParsed,
		// 	strings.Replace(magnet["dn"][0], " ", ".", -1),
		// 	item.Link,
		// )
		// if err == nil {
		// 	log.Println("Stored episode", show.Title, episode)
		// } else {
		// 	log.Println("Could not store show:", err)
		// }
	}
}
