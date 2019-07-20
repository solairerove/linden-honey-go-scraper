package domain

import (
	// TODO choose anything else ?
	"database/sql"
	"log"

	uuid "github.com/satori/go.uuid"
)

// Song ... tbd
type Song struct {
	ID     uuid.NullUUID `sql:",pk,type:uuid default uuid_generate_v4()" json:"-"`
	Title  string        `json:"title,omitempty"`
	Link   string        `json:"link,omitempty"`
	Author string        `json:"author,omitempty"`
	Album  string        `json:"album,omitempty"`
	Verses []Verse       `json:"verses,omitempty"`
}

// SaveSong ... tbd
func (s *Song) SaveSong(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO songs(title, link, author, album) 
											VALUES($1, $2, $3, $4) 
											RETURNING id`,
		s.Title, s.Link, s.Author, s.Album).Scan(&s.ID)

	if err != nil {
		return err
	}

	log.Printf("Persisted song id -> %s", s.ID.UUID.String())

	// TODO pff
	for _, v := range s.Verses {
		v.SongID = s.ID

		v.saveVerse(db)
	}

	return nil
}
