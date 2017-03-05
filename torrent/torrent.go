package torrent

import (
	"log"
	"strings"

	"github.com/doenietzomoeilijk/showfetcher/config"
	"github.com/doenietzomoeilijk/showfetcher/episode"
	"github.com/odwrtw/transmission"
)

// Setup sets up a connection.
func Setup(address string, tmpdir string) (err error) {
	conf := transmission.Config{Address: address}
	config.Torrent, err = transmission.New(conf)
	config.Torrent.Session.Set(transmission.SetSessionArgs{
		IncompleteDir:        tmpdir,
		IncompleteDirEnabled: true,
	})
	if err != nil {
		log.Fatalln("Stuff went wrong:", err)
	}

	log.Println(config.Torrent)

	return
}

// Cleanup cleans up torrents that are done.
func Cleanup() (e []*episode.Episode) {
	tmap, err := config.Torrent.GetTorrentMap()
	if err != nil {
		log.Fatalln("Error while fetching torrents from Transmission:", err)
	}

	var done []*transmission.Torrent

	for hash, tor := range tmap {
		hash = strings.ToLower(hash)
		log.Printf("h=%s id=%2d up=%01.3f fin=%5v stat=%d pct=%3d\n", tor.HashString, tor.ID, tor.UploadRatio, tor.IsFinished, tor.Status, int(100*tor.PercentDone))
		if tor.UploadRatio > 1 || tor.IsFinished || (tor.Status == transmission.StatusSeeding && tor.PercentDone == 1) {
			log.Println("This torrent is done")
			tor.Stop()
			done = append(done, tor)
			e = append(e, &episode.Episode{Hash: hash, Status: "done"})
		}
	}

	err = config.Torrent.RemoveTorrents(done, false)
	if err != nil {
		log.Println("Error while removing torrents:", err)
	}

	return
}

// Add new episodes to our torrent client.
func Add(eps []*episode.Episode) {
	for _, ep := range eps {
		config.Torrent.Session.Set(transmission.SetSessionArgs{
			DownloadDir: ep.Show.Location,
		})

		tor, err := config.Torrent.Add(ep.Magnet)
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

// List all currently active torrents in Transmission.
func List() (t []*transmission.Torrent) {
	t, err := config.Torrent.GetTorrents()
	if err != nil {
		log.Println("Error while fetching torrents:", err)
	}

	return
}
