package model

import "time"

// PullRequest represents pull request of Backlog git.
type PullRequest struct {
	ID           int           `json:"id,omitempty"`
	ProjectID    int           `json:"projectId,omitempty"`
	RepositoryID int           `json:"repositoryId,omitempty"`
	Number       int           `json:"number,omitempty"`
	Summary      string        `json:"summary,omitempty"`
	Description  string        `json:"description,omitempty"`
	Base         string        `json:"base,omitempty"`
	Branch       string        `json:"branch,omitempty"`
	Status       *Status       `json:"status,omitempty"`
	Assignee     *User         `json:"assignee,omitempty"`
	Issue        *Issue        `json:"issue,omitempty"`
	BaseCommit   interface{}   `json:"baseCommit,omitempty"`
	BranchCommit interface{}   `json:"branchCommit,omitempty"`
	CloseAt      time.Time     `json:"closeAt,omitempty"`
	MergeAt      time.Time     `json:"mergeAt,omitempty"`
	CreatedUser  *User         `json:"createdUser,omitempty"`
	Created      time.Time     `json:"created,omitempty"`
	UpdatedUser  *User         `json:"updatedUser,omitempty"`
	Updated      time.Time     `json:"updated,omitempty"`
	Attachments  []*Attachment `json:"attachments,omitempty"`
	Stars        []*Star       `json:"stars,omitempty"`
}

// Repository represents repository of Backlog git.
type Repository struct {
	ID           int       `json:"id,omitempty"`
	ProjectID    int       `json:"projectId,omitempty"`
	Name         string    `json:"name,omitempty"`
	Description  string    `json:"description,omitempty"`
	HookURL      string    `json:"hookUrl,omitempty"`
	HTTPURL      string    `json:"httpUrl,omitempty"`
	SSHURL       string    `json:"sshUrl,omitempty"`
	DisplayOrder int       `json:"displayOrder,omitempty"`
	PushedAt     time.Time `json:"pushedAt,omitempty"`
	CreatedUser  *User     `json:"createdUser,omitempty"`
	Created      time.Time `json:"created,omitempty"`
	UpdatedUser  *User     `json:"updatedUser,omitempty"`
	Updated      time.Time `json:"updated,omitempty"`
}
