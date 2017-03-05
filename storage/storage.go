package storage

import (
	"database/sql"
	"log"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	_ "github.com/mattn/go-sqlite3" // Use sqlite3
)

// Setup our database; make sure it exists and open it.
func Setup() (err error) {
	config.DB, err = sql.Open("sqlite3", "shows.db")

	if err != nil {
		log.Fatal(err)
	}

	_, err = config.DB.Exec(`CREATE TABLE IF NOT EXISTS shows (
        hash VARCHAR(50) NOT NULL COLLATE NOCASE PRIMARY KEY,
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
	_, err = config.DB.Exec(`CREATE INDEX IF NOT EXISTS st ON shows(status)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = config.DB.Exec(`CREATE INDEX IF NOT EXISTS sh ON shows(show)`)
	if err != nil {
		log.Fatal(err)
	}
	_, err = config.DB.Exec(`CREATE INDEX IF NOT EXISTS fn ON shows(filename)`)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// Close wraps around the database's closer.
func Close() {
	config.DB.Close()
}

// MarkDone marks a given set of Episodes as 'done' in the database.
// These Episodes should at least have their Hash set.
func MarkDone(episodes []episode.Episode) {
	for _, ep := range episodes {
		_, err := config.DB.Exec("UPDATE shows SET status = 'done' WHERE hash = ?", ep.Hash)
		if err != nil {
			log.Fatalln("DB error:", err)
		}
	}
}

// Get episodes that have status 'new'.
func Get() (episodes []episode.Episode) {
	rows, err := config.DB.Query("SELECT magnet, show, episode FROM shows WHERE status == 'new'")
	defer rows.Close()
	if err != nil {
		log.Fatalln("Error fetching data from db:", err)
		return
	}

	for rows.Next() {
		var ep episode.Episode
		var s string
		rows.Scan(&ep.Magnet, &s, &ep.Episode)
		ep.Show = episode.Shows[s]
		episodes = append(episodes, ep)
	}

	return
}
