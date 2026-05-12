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

// DiskUsageProject represents disk usage for a specific project.
type DiskUsageProject struct {
	DiskUsageBase
	ProjectID int `json:"projectId,omitempty"`
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
