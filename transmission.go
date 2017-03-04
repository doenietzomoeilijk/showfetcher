package showfetcher

import (
	"database/sql"
	"log"
	"strings"

	"github.com/odwrtw/transmission"
)

// SetupTransmission sets up a connection.
func SetupTransmission(config *Configuration) *transmission.Client {
	conf := transmission.Config{Address: config.Transmission}
	tr, err := transmission.New(conf)
	if err != nil {
		log.Fatalln("Stuff went wrong:", err)
	}

	return tr
}

// CleanTorrents cleans up torrents that are done.
func CleanTorrents(tr *transmission.Client, db *sql.DB) {
	log.Println("CleanTorrents")
	tmap, err := tr.GetTorrentMap()
	if err != nil {
		log.Fatalln("Stuff went wrong:", err)
	}

	var done []*transmission.Torrent

	for hash, tor := range tmap {
		log.Printf("h=%s id=%2d up=%01.3f fin=%5v stat=%d pct=%3d\n", tor.HashString, tor.ID, tor.UploadRatio, tor.IsFinished, tor.Status, int(100*tor.PercentDone))
		if tor.UploadRatio > 1 || tor.IsFinished || (tor.Status == transmission.StatusSeeding && tor.PercentDone == 1) {
			log.Println("This torrent is done")
			_, err := db.Exec("UPDATE shows SET status = 'done' WHERE hash = ?", strings.ToLower(hash))
			if err != nil {
				log.Fatalln("DB error:", err)
				continue
			}
			tor.Stop()
			done = append(done, tor)
		}
	}

	err = tr.RemoveTorrents(done, false)
	if err != nil {
		log.Println("Error while removing torrents:", err)
	}
}

func RunNewTorrents(conf *Configuration, tr *transmission.Client, db *sql.DB) {
	log.Println("RunNewTorrents")
	rows, err := db.Query("SELECT magnet, show, episode FROM shows WHERE status == 'new'")
	defer rows.Close()
	if err != nil {
		log.Fatalln("Error fetching new data from db:", err)
		return
	}

	for rows.Next() {
		var m string
		var s string
		var e string
		rows.Scan(&m, &s, &e)
		log.Printf("Adding episode %s for show %s\n", e, s)
		show, ok := conf.showByTitle(s)
		if !ok {
			log.Println("Did not find show", s)
			continue
		}

		tr.Session.Set(transmission.SetSessionArgs{
			DownloadDir:          show.Location,
			IncompleteDir:        conf.IncompleteDir,
			IncompleteDirEnabled: true,
		})

		tor, err := tr.Add(m)
		if err != nil {
			log.Println("Error while adding torrent:", err)
			continue
		}

		err = tor.Start()
		if err != nil {
			log.Println("Error while starting torrent:", err)
		}
	}
}

func ListTorrents(tr *transmission.Client) (t []*transmission.Torrent) {
	t, err := tr.GetTorrents()
	if err != nil {
		log.Println("Error while fetching torrents:", err)
	}

	return
}
