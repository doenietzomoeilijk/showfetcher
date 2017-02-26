package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3" // Use sqlite3
	"github.com/mmcdole/gofeed"
)

const showPattern = "(.+) ([0-9]+)x([0-9]+) ?(.*)"
const transPattern = "^\\s*(?P<id>\\d+)  +(?P<done>\\d+)%  +(?P<have>[0-9.]+ (kB|MB|GB|TB))  +(?P<eta>.*?)  +(?P<up>[0-9.]+)  +(?P<down>[0-9.]+)  +(?P<ratio>[0-9.]+)  +(?P<status>.*?)  +(?P<name>.*)$"

var (
	db      *sql.DB
	showRe  *regexp.Regexp
	transRe *regexp.Regexp
	config  *Configuration
)

// Create our database if it doesn't exist, make sure it's open.
func setupDb() {
	var err error

	db, err = sql.Open("sqlite3", "shows.db")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS shows (
        guid VARCHAR(50) NOT NULL PRIMARY KEY,
        show VARCHAR(50) NOT NULL,
        episode VARCHAR(5) NOT NULL,
        published DATETIME NULL,
        status VARCHAR(10) NOT NULL DEFAULT 'new',
        filename VARCHAR(100) NOT NULL
    )`)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS st ON shows(status)`)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS sh ON shows(show)`)
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS fn ON shows(filename)`)
}

// Go go go.
func main() {
	config = &Configuration{}
	err := config.load("config.json")
	if err != nil {
		panic("Config can't be read")
	}

	showRe = regexp.MustCompile(showPattern)
	transRe = regexp.MustCompile(transPattern)
	setupDb()

	parseFeeds()
	transmatch()
}

func parseFeeds() {
	feed, _ := gofeed.NewParser().ParseURL(config.FeedURL)

	for _, item := range feed.Items {
		show, ok := config.findShow(item.Title)

		if !ok {
			log.Println("I cannot haz show")
			continue
		}
		log.Printf("Entry '%s' matches show '%s'\n", item.Title, show.Title)

		epi := Episode{}
		epi.fill(&show, item)
		ok = epi.store()
		if !ok {
			log.Println("Could not store show.")
			continue
		}
	}
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

// Load our config
func (c *Configuration) load(filename string) error {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return err
	}

	return nil
}

// Find a show by title.
func (c *Configuration) findShow(str string) (s Show, ok bool) {
	log.Println("Trying to find show", str, c)
	for _, show := range c.Shows {
		log.Println("Trying show", show)

		matches := showRe.FindStringSubmatch(str)
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

// Episode holds a singular episode from the feed.
type Episode struct {
	Show      *Show
	Name      string
	Season    int
	Episode   int
	Title     string
	File      string
	GUID      string
	Published *time.Time
	Status    string
}

func (e *Episode) fill(show *Show, item *gofeed.Item) {
	matches := showRe.FindStringSubmatch(item.Title)
	if len(matches) < 1 {
		log.Println("Didn't find anything!", matches, item.Title)
		return
	}
	e.Show = show
	e.Name = matches[1]
	e.Season, _ = strconv.Atoi(matches[2])
	e.Episode, _ = strconv.Atoi(matches[3])
	e.Title = matches[4]
	e.GUID = item.GUID
	v, _ := url.ParseQuery(strings.Split(item.Link, "?")[1])
	e.File = strings.Replace(v["dn"][0], " ", ".", -1)
	e.Published = item.PublishedParsed
}

func (e *Episode) store() bool {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO shows
		(guid, show, episode, published, filename)
		VALUES (?, ?, ?, ?, ?)`,
		e.GUID,
		e.Show.Title,
		fmt.Sprintf("%dx%d", e.Season, e.Episode),
		e.Published,
		e.File)

	if err == nil {
		log.Println("Stored show", e.Show.Title)
		return true
	}

	log.Panicln("Error while inserting show:", err)
	return false
}

type TransmissionStatus struct {
	Lines []TransmissionLine
}

func (t *TransmissionStatus) fetch() {
	output, _ := exec.Command("/usr/bin/transmission-remote", config.Transmission, "--list").Output()
	lines := strings.Split(string(output), "\n")[1:]

	for _, line := range lines {
		fields := transRe.FindStringSubmatch(line)
		if len(fields) < 1 {
			continue
		}
		ma := map[string]string{}
		for i, fld := range transRe.SubexpNames() {
			if fld == "" {
				continue
			}
			ma[fld] = fields[i]
		}
		id, _ := strconv.Atoi(ma["id"])
		tl := TransmissionLine{
			ID:     id,
			Done:   ma["done"],
			ETA:    ma["eta"],
			Name:   ma["name"],
			Status: ma["status"],
		}

		t.Lines = append(t.Lines, tl)
	}
}

// TransmissionLine holds a single line from the Transmission process.
type TransmissionLine struct {
	ID     int
	Done   string
	ETA    string
	Name   string
	Status string
}

// Match against regex to find SxEx part
func (tl *TransmissionLine) Episode() string {
	return ""
}

func transmatch() {
	t := &TransmissionStatus{}
	t.fetch()

	for _, tl := range t.Lines {
		if tl.Status == "Seeding" || tl.Status == "Idle" {

		}
	}
}
