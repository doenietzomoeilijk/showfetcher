package config

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"

	"github.com/doenietzomoeilijk/showfetcher/episode"
	_ "github.com/mattn/go-sqlite3" // Use sqlite3
	"github.com/odwrtw/transmission"
)

const filename string = "config.json"

var (
	Config  *Configuration
	Torrent *transmission.Client
	DB      *sql.DB
	Showmap map[string]*episode.Show
)

// Configuration holds the entire JSON config
type Configuration struct {
	FeedURL       string         `json:"feed_url"`
	IncompleteDir string         `json:"incomplete_dir"`
	Transmission  string         `json:"transmission_rpc_url"`
	Shows         []episode.Show `json:"shows"`
}

func init() {
	Showmap = make(map[string]*episode.Show)
}

// Find a show by title.
/*
func (c *Configuration) findShow(str string) (s episode.Show, ok bool) {
	for _, show := range c.Shows {

		matches := ShowRe.FindStringSubmatch(str)
		if len(matches) < 1 {
			log.Println("No matches")
			continue
		}
		name := matches[1]
		log.Println("Trying searchstring=", show.SearchString, "against match name=", name)

		if m, _ := regexp.Match(show.SearchString, []byte(name)); m {
			return show, true
		}
	}

	return episode.Show{}, false
}

func (c *Configuration) showByTitle(str string) (s episode.Show, ok bool) {
	log.Println("Finding show by title", str)
	for _, show := range c.Shows {
		if show.Title == str {
			return show, true
		}
	}
	return episode.Show{}, false
}
*/

// Load our config.
func Load() {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	Config = &Configuration{}
	err = json.Unmarshal(bytes, &Config)
	if err != nil {
		panic(err)
	}

	// Config.Showmap = episode.Showmap{}
	for _, show := range Config.Shows {
		Showmap[show.Title] = &show
	}
}
