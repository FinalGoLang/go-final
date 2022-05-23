package models

import "time"

type Subscribe struct {
	UserID  int `db:"user_id" json:"user_id,omitempty"`
	FilmdID int `db:"film_id" json:"film_id,omitempty"`
}

type MyFilms struct {
	FilmID    int       `db:"film_id" json:"film_id,omitempty"`
	Name      string    `db:"name" json:"name,omitempty"`
	Price     float32   `db:"price" json:"price,omitempty"`
	Rating    int       `db:"rating" json:"rating,omitempty"`
	UserID    int       `db:"user_id" json:"user_id,omitempty"`
	Expires   time.Time `db:"expires" json:"expires,omitempty"`
	CreatedAt time.Time `db:"created_at" json:"-"`
	UpdatedAt time.Time `db:"updated_at" json:"-"`
}
