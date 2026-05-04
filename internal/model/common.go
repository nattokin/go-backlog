package model

import (
	"io"
	"time"
)

// FileData represents a downloaded binary file with its metadata.
// Body must be closed by the caller after use.
type FileData struct {
	Body        io.ReadCloser
	Filename    string
	ContentType string
}

type Attachment struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	CreatedUser *User     `json:"createdUser"`
	Created     time.Time `json:"created"`
}

type Star struct {
	ID        int       `json:"id"`
	Comment   string    `json:"comment"`
	URL       string    `json:"url"`
	Title     string    `json:"title"`
	Presenter *User     `json:"presenter"`
	Created   time.Time `json:"created"`
}

type Notification struct {
	ID                  int          `json:"id"`
	AlreadyRead         bool         `json:"alreadyRead"`
	Reason              int          `json:"reason"`
	ResourceAlreadyRead bool         `json:"resourceAlreadyRead"`
	Project             *Project     `json:"project"`
	Issue               *Issue       `json:"issue"`
	Comment             *Comment     `json:"comment"`
	PullRequest         *PullRequest `json:"pullRequest"`
	PullRequestComment  *Comment     `json:"pullRequestComment"`
	Sender              *User        `json:"sender"`
	Created             time.Time    `json:"created"`
}

type Comment struct {
	ID            int             `json:"id"`
	Content       string          `json:"content"`
	ChangeLogs    []*ChangeLog    `json:"changeLog"`
	CreatedUser   *User           `json:"createdUser"`
	Created       time.Time       `json:"created"`
	Updated       time.Time       `json:"updated"`
	Stars         []*Star         `json:"stars"`
	Notifications []*Notification `json:"notifications"`
}

type ChangeLog struct {
	Field         string `json:"field"`
	NewValue      string `json:"newValue"`
	OriginalValue string `json:"originalValue"`
}

type SharedFile struct {
	ID          int       `json:"id"`
	Type        string    `json:"type"`
	Dir         string    `json:"dir"`
	Name        string    `json:"name"`
	Size        int       `json:"size"`
	CreatedUser *User     `json:"createdUser"`
	Created     time.Time `json:"created"`
	UpdatedUser *User     `json:"updatedUser"`
	Updated     time.Time `json:"updated"`
}

type Category struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DisplayOrder int    `json:"displayOrder"`
}

type Status struct {
	ID           int    `json:"id"`
	ProjectID    int    `json:"projectId"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	DisplayOrder int    `json:"displayOrder"`
}

type Tag struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Version struct {
	ID             int       `json:"id"`
	ProjectID      int       `json:"projectId"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	StartDate      time.Time `json:"startDate"`
	ReleaseDueDate time.Time `json:"releaseDueDate"`
	Archived       bool      `json:"archived"`
	DisplayOrder   int       `json:"displayOrder"`
}

type WatchingItem struct {
	ID                  int       `json:"id"`
	ResourceAlreadyRead bool      `json:"resourceAlreadyRead"`
	Note                string    `json:"note"`
	Type                string    `json:"type"`
	Issue               *Issue    `json:"issue"`
	LastContentUpdated  time.Time `json:"lastContentUpdated"`
	Created             time.Time `json:"created"`
	Updated             time.Time `json:"updated"`
}
