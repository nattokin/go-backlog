package model

import "time"

// Wiki represents Backlog Wiki.
type Wiki struct {
	ID          int           `json:"id,omitempty"`
	ProjectID   int           `json:"projectId,omitempty"`
	Name        string        `json:"name,omitempty"`
	Content     string        `json:"content,omitempty"`
	Tags        []*Tag        `json:"tags,omitempty"`
	Attachments []*Attachment `json:"attachments,omitempty"`
	SharedFiles []*SharedFile `json:"sharedFiles,omitempty"`
	Stars       []*Star       `json:"stars,omitempty"`
	CreatedUser *User         `json:"createdUser,omitempty"`
	Created     time.Time     `json:"created,omitempty"`
	UpdatedUser *User         `json:"updatedUser,omitempty"`
	Updated     time.Time     `json:"updated,omitempty"`
}

// WikiHistory represents a version history entry for a wiki page.
type WikiHistory struct {
	PageID      int       `json:"pageId,omitempty"`
	Version     int       `json:"version,omitempty"`
	Name        string    `json:"name,omitempty"`
	Content     string    `json:"content,omitempty"`
	CreatedUser *User     `json:"createdUser,omitempty"`
	Created     time.Time `json:"created,omitempty"`
}
