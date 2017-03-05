package feed

import (
	"log"
	"regexp"
	"strings"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	"github.com/mmcdole/gofeed"
)

const pt = "(.+) ([0-9]+x[0-9]+) ?(.*)"

var re *regexp.Regexp

func init() {
	re = regexp.MustCompile(pt)
}

// Parse grabs and parses the show feed.
func Parse(url string) (eps []*episode.Episode) {
	log.Println("Parse Feeds")
	feed, _ := gofeed.NewParser().ParseURL(url)

	for _, item := range feed.Items {
		log.Println("Feed item:", item.Title, re)
		ex := item.Extensions["tv"]
		showname := ex["show_name"][0].Value
		hash := strings.ToLower(ex["info_hash"][0].Value)
		show := config.Showmap[showname]
		matches := re.FindStringSubmatch(item.Title)

		eps = append(eps, &episode.Episode{
			Show:    show,
			Hash:    hash,
			Magnet:  item.Link,
			Episode: matches[2],
			File:    strings.Replace(ex["raw_title"][0].Value, " ", ".", -1),
		})
	}

	return
}
