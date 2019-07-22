package domain

// Verse is a struct that describe verses fields
type Verse struct {
	Ordinal int    `json:"ord"`
	Data    string `json:"data,omitempty"`
}
