// Episode holds the structs and maps for shows and episodes.
package episode

import (
	"time"
)

// Episodes holds a group of episodes, mapped to the hash.
type Episodes map[string]*Episode

// Episode holds a singular episode (linked with the show).
type Episode struct {
	Hash      string
	Show      Show
	Episode   string // in 2x10 format
	Published *time.Time
	Status    string // new, busy or done
	File      string // file name, with spaces switched for dots
	Magnet    string // magnet link
}
