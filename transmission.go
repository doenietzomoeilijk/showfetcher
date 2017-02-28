package showfetcher

import (
	"database/sql"
	"log"

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
	tmap, err := tr.GetTorrentMap()
	if err != nil {
		log.Fatalln("Stuff went wrong:", err)
	}

	for a, tor := range tmap {
		log.Println(a, tor.ID, tor.UploadRatio, tor.IsFinished, tor.Status)
		if tor.UploadRatio > 1 || tor.IsFinished {
			db.Exec("UPDATE shows SET status = 'DONE' WHERE hash = ?")
			tor.Stop()
			tor.Client.RemoveTorrents([]*transmission.Torrent{tor}, false)
		}
	}
}
