package episode

// Shows holds the externally available showmap
var Shows showmap

// Showmap holds all our shows, mapped to the show title.
type showmap map[string]Show

// Show holds one singular show entry.
type Show struct {
	Title        string `json:"title"`
	SearchString string `json:"search_string"`
	Location     string `json:"location"`
}
