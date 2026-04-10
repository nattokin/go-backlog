package model

// DiskUsageProject represents project's disk usage.
type DiskUsageProject struct {
	DiskUsageBase
	ProjectID int `json:"projectId,omitempty"`
}

// Project represents a project of Backlog.
type Project struct {
	ID                                int    `json:"id,omitempty"`
	ProjectKey                        string `json:"projectKey,omitempty"`
	Name                              string `json:"name,omitempty"`
	ChartEnabled                      bool   `json:"chartEnabled,omitempty"`
	SubtaskingEnabled                 bool   `json:"subtaskingEnabled,omitempty"`
	ProjectLeaderCanEditProjectLeader bool   `json:"projectLeaderCanEditProjectLeader,omitempty"`
	TextFormattingRule                Format `json:"textFormattingRule,omitempty"`
	Archived                          bool   `json:"archived,omitempty"`
}
