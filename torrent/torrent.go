// Torrent talks to our Transmission server.
// Torrent talks to our Transmission RPC daemon.
package torrent

import (
	"log"
	"strings"

	"github.com/doenietzomoeilijk/showfetcher/episode"
	"github.com/odwrtw/transmission"
)

var client *transmission.Client

// Setup sets up a connection.
func Setup(address string, tmpdir string) (err error) {
	conf := transmission.Config{Address: address}
	client, err = transmission.New(conf)
	client.Session.Set(transmission.SetSessionArgs{
		IncompleteDir:        tmpdir,
		IncompleteDirEnabled: true,
	})
	if err != nil {
		log.Fatalln("Couldn't set up Transmission:", err)
	}

	return
}

// Cleanup cleans up torrents that are done.
func Cleanup() (e []*episode.Episode, err error) {
	tmap, err := client.GetTorrentMap()
	if err != nil {
		log.Fatalln("Error while fetching torrents from Transmission:", err)
	}

	var done []*transmission.Torrent

	for hash, tor := range tmap {
		hash = strings.ToLower(hash)
		if tor.UploadRatio > 1 || tor.IsFinished || (tor.Status == transmission.StatusSeeding && tor.PercentDone == 1) {
			tor.Stop()
			done = append(done, tor)
			e = append(e, &episode.Episode{Hash: hash, Status: "done"})
		}
	}

	err = client.RemoveTorrents(done, false)
	if err != nil {
		log.Println("Error while removing torrents:", err)
	}

	return
}

// Add new episodes to our torrent client.
func Add(eps []*episode.Episode) {
	for _, ep := range eps {
		log.Println("Adding episode to Transmission:", ep)
		client.Session.Set(transmission.SetSessionArgs{
			DownloadDir: ep.Show.Location,
		})

		tor, err := client.Add(ep.Magnet)
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
	t, err := client.GetTorrents()
	if err != nil {
		log.Println("Error while fetching torrents:", err)
	}

	return
}
