package data

import "time"

type Movie struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"-"`
	Title     string    `json:"title"`
	Year      int32     `json:"year,omitempty"`    // release year
	Runtime   int32     `json:"runtime,omitempty"` // movie run time (in minutes)
	Genres    []string  `json:"genres,omitempty"`
	Version   int32     `json:"version"` // the version number starts at 1 and will be incremented each time the movie information is updated
}
