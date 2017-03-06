package feed

import (
	"regexp"
	"strings"

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
	feed, _ := gofeed.NewParser().ParseURL(url)
	for _, item := range feed.Items {
		tv := item.Extensions["tv"]
		showname := tv["show_name"][0].Value
		hash := strings.ToLower(tv["info_hash"][0].Value)
		show := episode.Shows[showname]
		matches := re.FindStringSubmatch(item.Title)

		eps = append(eps, &episode.Episode{
			Show:    show,
			Hash:    hash,
			Magnet:  item.Link,
			Episode: matches[2],
			File:    strings.Replace(tv["raw_title"][0].Value, " ", ".", -1),
		})
	}

	return
}
