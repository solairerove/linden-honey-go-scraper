package domain

// Song is a struct that describe necessary song fields
type Song struct {
	Title  string  `json:"title,omitempty"`
	Link   string  `json:"link,omitempty"`
	Author string  `json:"author,omitempty"`
	Album  string  `json:"album,omitempty"`
	Verses []Verse `json:"verses,omitempty"`
}
