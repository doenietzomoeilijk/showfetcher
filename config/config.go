// Config loads our config.json file.
package config

import (
	"encoding/json"
	"io/ioutil"

	"github.com/doenietzomoeilijk/showfetcher/episode"
	_ "github.com/mattn/go-sqlite3" // Use sqlite3
)

const filename string = "config.json"

// Configuration holds the entire JSON config
type Configuration struct {
	FeedURL       string         `json:"feed_url"`
	IncompleteDir string         `json:"incomplete_dir"`
	Transmission  string         `json:"transmission_rpc_url"`
	Shows         []episode.Show `json:"shows"`
	Logfile       string         `json:"logfile"`
}

// Load our config.
func Load() Configuration {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	conf := Configuration{}
	err = json.Unmarshal(bytes, &conf)
	if err != nil {
		panic(err)
	}

	episode.Shows = make(map[string]episode.Show)
	for _, show := range conf.Shows {
		episode.Shows[show.Title] = show
	}

	return conf
}
