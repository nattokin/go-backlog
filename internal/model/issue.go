package model

import "time"

// Category represents an issue category.
type Category struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// Issue represents a issue of Backlog.
type Issue struct {
	ID             int            `json:"id,omitempty"`
	ProjectID      int            `json:"projectId,omitempty"`
	IssueKey       string         `json:"issueKey,omitempty"`
	KeyID          int            `json:"keyId,omitempty"`
	IssueType      *IssueType     `json:"issueType,omitempty"`
	Summary        string         `json:"summary,omitempty"`
	Description    string         `json:"description,omitempty"`
	Resolutions    []*Resolution  `json:"resolutions,omitempty"`
	Priority       *Priority      `json:"priority,omitempty"`
	Status         *Status        `json:"status,omitempty"`
	Assignee       *User          `json:"assignee,omitempty"`
	Category       []*Category    `json:"category,omitempty"`
	Versions       *Version       `json:"versions,omitempty"`
	Milestone      *Version       `json:"milestone,omitempty"`
	StartDate      time.Time      `json:"startDate,omitempty"`
	DueDate        time.Time      `json:"dueDate,omitempty"`
	EstimatedHours int            `json:"estimatedHours,omitempty"`
	ActualHours    int            `json:"actualHours,omitempty"`
	ParentIssueID  int            `json:"parentIssueId,omitempty"`
	CreatedUser    *User          `json:"createdUser,omitempty"`
	Created        time.Time      `json:"created,omitempty"`
	UpdatedUser    *User          `json:"updatedUser,omitempty"`
	Updated        time.Time      `json:"updated,omitempty"`
	CustomFields   []*CustomField `json:"customFields,omitempty"`
	Attachments    []*Attachment  `json:"attachments,omitempty"`
	SharedFiles    []*SharedFile  `json:"sharedFiles,omitempty"`
	Stars          []*Star        `json:"stars,omitempty"`
}

// IssueType represents type of Issue.
type IssueType struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// Priority represents a priority.
type Priority struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// Resolution represents a resolution.
type Resolution struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}
