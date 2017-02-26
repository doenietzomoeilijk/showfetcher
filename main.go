package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"os/exec"

	"github.com/mmcdole/gofeed"
)

const showPattern = "(.+) ([0-9]+)x([0-9]+) (.+)"
const transPattern = "^\\s*(?P<id>\\d+)  +(?P<done>\\d+)%  +(?P<have>[0-9.]+ (kB|MB|GB|TB))  +(?P<eta>.*?)  +(?P<up>[0-9.]+)  +(?P<down>[0-9.]+)  +(?P<ratio>[0-9.]+)  +(?P<status>.*?)  +(?P<name>.*)$"

var (
	showRe  *regexp.Regexp
	transRe *regexp.Regexp
	config  Configuration
)

func main() {
	showRe = regexp.MustCompile(showPattern)
	transRe = regexp.MustCompile(transPattern)

	config, err := loadConfig("config.json")
	if err != nil {
		panic("Config can't be read")
	}

	log.Printf("%#v\n", config)

	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(config.FeedURL)

	log.Println(feed.Title)

	for _, item := range feed.Items {
		epi := Episode{}
		epi.fill(item)

		// log.Printf("Item: %#v\n", epi)
		show, ok := config.findShow(item.Title)
		if !ok {
			continue
		}
		log.Printf("  Entry '%s' matches show '%s'\n", item.Title, show.Title)
		// log.Println("  ", epi.Published.Format(time.RFC3339))
		log.Println("")
	}

	transmatch(config.Transmission)
}

// Configuration holds the entire JSON config
type Configuration struct {
	FeedURL         string `json:"feed_url"`
	BaseTorrentPath string `json:"base_torrent_path"`
	Transmission    string `json:"transmission_daemon"`
	Shows           []Show `json:"shows"`
}

// Show holds one singular show entry
type Show struct {
	Title        string `json:"title"`
	SearchString string `json:"search_string"`
	Location     string `json:"location"`
}

func (c *Configuration) findShow(str string) (s Show, ok bool) {
	for _, show := range c.Shows {
		matches := showRe.FindStringSubmatch(str)
		if len(matches) < 1 {
			continue
		}

		name := matches[1]
		// season := matches[2]
		// episode := matches[3]
		// title := matches[4]

		if m, _ := regexp.Match(show.SearchString, []byte(name)); m {
			return show, true
		}
	}

	return Show{}, false
}

func loadConfig(filename string) (Configuration, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return Configuration{}, err
	}

	var c Configuration
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return Configuration{}, err
	}

	return c, nil
}

// Episode holds a singular episode from the feed.
type Episode struct {
	Name      string
	Season    int
	Episode   int
	Title     string
	File      string
	GUID      string
	Published *time.Time
	Status    string
}

func (e *Episode) fill(item *gofeed.Item) {
	matches := showRe.FindStringSubmatch(item.Title)
	e.Name = matches[1]
	e.Season, _ = strconv.Atoi(matches[2])
	e.Episode, _ = strconv.Atoi(matches[3])
	e.Title = matches[4]
	e.GUID = item.GUID
	v, _ := url.ParseQuery(strings.Split(item.Link, "?")[1])
	e.File = strings.Replace(v["dn"][0], " ", ".", -1)
	e.Published = item.PublishedParsed
}

type TransmissionStatus struct {
	Lines []TransmissionLine
}
type TransmissionLine struct {
	ID     int
	Done   string
	ETA    string
	Name   string
	Status string
}

func transmatch(s string) {
	output, _ := exec.Command("/usr/bin/transmission-remote", s, "--list").Output()
	lines := strings.Split(string(output), "\n")[1:]

	m2 := transRe.FindStringSubmatch(lines[0])
	ma := map[string]string{}

	for i, fld := range transRe.SubexpNames() {
		if fld == "" {
			continue
		}

		ma[fld] = m2[i]
	}

	id, _ := strconv.Atoi(ma["id"])
	tl := TransmissionLine{
		ID:     id,
		Done:   ma["done"],
		ETA:    ma["eta"],
		Name:   ma["name"],
		Status: ma["status"],
	}
	fmt.Printf("%#v\n", tl)
	return

}
