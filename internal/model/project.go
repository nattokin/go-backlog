package model

import "time"

// Project represents a Backlog project.
type Project struct {
	ID                                int    `json:"id,omitempty"`
	ProjectKey                        string `json:"projectKey,omitempty"`
	Name                              string `json:"name,omitempty"`
	ChartEnabled                      bool   `json:"chartEnabled,omitempty"`
	SubtaskingEnabled                 bool   `json:"subtaskingEnabled,omitempty"`
	ProjectLeaderCanEditProjectLeader bool   `json:"projectLeaderCanEditProjectLeader,omitempty"`
	TextFormattingRule                string `json:"textFormattingRule,omitempty"`
	Archived                          bool   `json:"archived,omitempty"`
	UseResolvedForChart               bool   `json:"useResolvedForChart,omitempty"`
	UseWiki                           bool   `json:"useWiki,omitempty"`
	UseFileSharing                    bool   `json:"useFileSharing,omitempty"`
	UseWikiTreeView                   bool   `json:"useWikiTreeView,omitempty"`
	UseSubversion                     bool   `json:"useSubversion,omitempty"`
	UseGit                            bool   `json:"useGit,omitempty"`
	UseOriginalImageSizeAtWiki        bool   `json:"useOriginalImageSizeAtWiki,omitempty"`
	DisplayOrder                      int    `json:"displayOrder,omitempty"`
	UseDevAttributes                  bool   `json:"useDevAttributes,omitempty"`
}

// Category represents a project category.
type Category struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// CustomField represents a custom field defined in the project.
type CustomField struct {
	ID                     int                `json:"id,omitempty"`
	TypeID                 int                `json:"typeId,omitempty"`
	Name                   string             `json:"name,omitempty"`
	Description            string             `json:"description,omitempty"`
	Required               bool               `json:"required,omitempty"`
	ApplicableIssueTypeIDs []int              `json:"applicableIssueTypes,omitempty"`
	AllowAddItem           bool               `json:"allowAddItem,omitempty"`
	Items                  []*CustomFieldItem `json:"items,omitempty"`
}

// CustomFieldItem represents one selectable item in a List type CustomField.
type CustomFieldItem struct {
	ID           int    `json:"id,omitempty"`
	Name         string `json:"name,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// DiskUsageProject represents disk usage for a specific project.
type DiskUsageProject struct {
	DiskUsageBase
	ProjectID int `json:"projectId,omitempty"`
}

// Status represents a project status that can be assigned to issues.
type Status struct {
	ID           int    `json:"id,omitempty"`
	ProjectID    int    `json:"projectId,omitempty"`
	Name         string `json:"name,omitempty"`
	Color        string `json:"color,omitempty"`
	DisplayOrder int    `json:"displayOrder,omitempty"`
}

// Version represents a project version (milestone).
type Version struct {
	ID             int    `json:"id,omitempty"`
	ProjectID      int    `json:"projectId,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	StartDate      string `json:"startDate,omitempty"`
	ReleaseDueDate string `json:"releaseDueDate,omitempty"`
	Archived       bool   `json:"archived,omitempty"`
	DisplayOrder   int    `json:"displayOrder,omitempty"`
}

// Webhook represents a Backlog webhook.
type Webhook struct {
	ID              int       `json:"id,omitempty"`
	Name            string    `json:"name,omitempty"`
	Description     string    `json:"description,omitempty"`
	HookURL         string    `json:"hookUrl,omitempty"`
	AllEvent        bool      `json:"allEvent,omitempty"`
	ActivityTypeIDs []int     `json:"activityTypeIds,omitempty"`
	CreatedUser     *User     `json:"createdUser,omitempty"`
	Created         time.Time `json:"created,omitempty"`
	UpdatedUser     *User     `json:"updatedUser,omitempty"`
	Updated         time.Time `json:"updated,omitempty"`
}
