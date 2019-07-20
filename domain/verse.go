package domain

import (
	// TODO choose anything else ?
	"database/sql"
	"log"

	uuid "github.com/satori/go.uuid"
)

// Verse ... tbd
type Verse struct {
	ID      uuid.NullUUID `sql:",pk,type:uuid default uuid_generate_v4()" json:"-"`
	Ordinal int           `json:"ord"`
	Data    string        `json:"data,omitempty"`
	SongID  uuid.NullUUID `sql:",type:uuid" json:"-"`
}

// saveVerse ... tbd
func (v *Verse) saveVerse(db *sql.DB) error {
	err := db.QueryRow(`INSERT INTO verses(ordinal, data, song_id) 
											VALUES($1, $2, $3) 
											RETURNING id`,
		v.Ordinal, v.Data, v.SongID).Scan(&v.ID)

	if err != nil {
		return err
	}

	log.Printf("Persisted verse id -> %s", v.ID.UUID.String())

	return nil
}
