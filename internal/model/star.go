package model

import "time"

// Star represents a star attached to an issue, comment, or wiki page.
type Star struct {
	ID        int       `json:"id,omitempty"`
	Comment   string    `json:"comment,omitempty"`
	URL       string    `json:"url,omitempty"`
	Title     string    `json:"title,omitempty"`
	Presenter *User     `json:"presenter,omitempty"`
	Created   time.Time `json:"created,omitempty"`
}
