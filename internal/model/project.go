package model

import (
	"encoding/json"
	"time"
)

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

	// Type-specific settings. Each appears only for the matching custom field type.
	Min          *CustomFieldBound       `json:"min,omitempty"`
	Max          *CustomFieldBound       `json:"max,omitempty"`
	InitialValue *float64                `json:"initialValue,omitempty"`
	Unit         string                  `json:"unit,omitempty"`
	InitialDate  *CustomFieldInitialDate `json:"initialDate,omitempty"`
	AllowInput   bool                    `json:"allowInput,omitempty"`
}

// CustomFieldBound is the min or max of a CustomField.
// The API returns a number for a Number type custom field and a datetime string
// for a Date type one under the same key.
type CustomFieldBound struct {
	Number *float64
	Date   *time.Time
}

// UnmarshalJSON decodes a JSON string into Date and anything else into Number.
func (b *CustomFieldBound) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && data[0] == '"' {
		return json.Unmarshal(data, &b.Date)
	}
	return json.Unmarshal(data, &b.Number)
}

// CustomFieldInitialDate represents how the initial value of a Date type CustomField is decided.
// ID is the initial value type (1: today, 2: today plus Shift days, 3: Date).
type CustomFieldInitialDate struct {
	ID    int       `json:"id,omitempty"`
	Shift int       `json:"shift,omitempty"`
	Date  time.Time `json:"date,omitempty"`
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
	ID             int       `json:"id,omitempty"`
	ProjectID      int       `json:"projectId,omitempty"`
	Name           string    `json:"name,omitempty"`
	Description    string    `json:"description,omitempty"`
	StartDate      time.Time `json:"startDate,omitempty"`
	ReleaseDueDate time.Time `json:"releaseDueDate,omitempty"`
	Archived       bool      `json:"archived,omitempty"`
	DisplayOrder   int       `json:"displayOrder,omitempty"`
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
