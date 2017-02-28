package showfetcher

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"

	_ "github.com/mattn/go-sqlite3" // Use sqlite3
)

const filename string = "config.json"
const ShowPattern = "(.+) ([0-9]+x[0-9]+) ?(.*)"

var (
	ShowRe *regexp.Regexp
)

// Configuration holds the entire JSON config
type Configuration struct {
	FeedURL         string `json:"feed_url"`
	BaseTorrentPath string `json:"base_torrent_path"`
	Transmission    string `json:"transmission_rpc_url"`
	Shows           []Show `json:"shows"`
}

// Show holds one singular show entry
type Show struct {
	Title        string `json:"title"`
	SearchString string `json:"search_string"`
	Location     string `json:"location"`
}

// Find a show by title.
func (c *Configuration) findShow(str string) (s Show, ok bool) {
	log.Println("Trying to find show", str, c)
	for _, show := range c.Shows {
		log.Println("Trying show", show)

		matches := ShowRe.FindStringSubmatch(str)
		if len(matches) < 1 {
			log.Println("No matches")
			continue
		}
		name := matches[1]
		log.Println("Trying name", name, "against regex", show.SearchString)

		if m, _ := regexp.Match(show.SearchString, []byte(name)); m {
			return show, true
		}
	}

	return Show{}, false
}

func LoadConfig() (config *Configuration) {
	ShowRe = regexp.MustCompile(ShowPattern)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	config = &Configuration{}
	err = json.Unmarshal(bytes, &config)
	if err != nil {
		panic(err)
	}

	return
}

// Create our database if it doesn't exist, make sure it's open.
func SetupDb() (db *sql.DB, err error) {
	db, err = sql.Open("sqlite3", "shows.db")

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS shows (
        hash VARCHAR(50) NOT NULL PRIMARY KEY,
        show VARCHAR(50) NOT NULL,
        episode VARCHAR(5) NOT NULL,
        published DATETIME NULL,
        status VARCHAR(10) NOT NULL DEFAULT 'new',
        filename VARCHAR(100) NOT NULL,
		magnet TEXT NOT NULL
    )`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS st ON shows(status)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS sh ON shows(show)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS fn ON shows(filename)`)
	if err != nil {
		log.Fatal(err)
	}
	return db, nil
}
