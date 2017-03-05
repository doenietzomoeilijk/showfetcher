package episode

import (
	"time"
)

// Episodes holds a group of episodes, mapped to the hash.
type Episodes map[string]*Episode

// Episode that works for the feed as wel as for the database.
type Episode struct {
	Hash      string // xt=urn:btih:<HASH>
	Show      *Show
	Episode   string // 2x10
	Published *time.Time
	Status    string
	File      string
	Magnet    string // Magnet link
}
