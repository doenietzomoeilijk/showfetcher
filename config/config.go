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
	// Config holds our struct
	Config *Configuration

	// Torrent holds a Transmission Client
	Torrent *transmission.Client

	// DB holds what you'd think it holds
	DB *sql.DB

	// Showmap is a convenience mapping of hash=>Show
	Showmap map[string]*episode.Show
)

// Configuration holds the entire JSON config
type Configuration struct {
	FeedURL       string         `json:"feed_url"`
	IncompleteDir string         `json:"incomplete_dir"`
	Transmission  string         `json:"transmission_rpc_url"`
	Shows         []episode.Show `json:"shows"`
}

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

	Showmap = make(map[string]*episode.Show)
	for _, show := range Config.Shows {
		Showmap[show.Title] = &show
	}
}
