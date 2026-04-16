package model

import "time"

// DiskUsageProject represents project's disk usage.
type DiskUsageProject struct {
	DiskUsageBase
	ProjectID int `json:"projectId,omitempty"`
}

// DiskUsageSpace represents space's disk usage.
type DiskUsageSpace struct {
	DiskUsageBase
	Capacity int                 `json:"capacity,omitempty"`
	Details  []*DiskUsageProject `json:"details,omitempty"`
}

// Space represents space of Backlog.
type Space struct {
	SpaceKey           string    `json:"spaceKey,omitempty"`
	Name               string    `json:"name,omitempty"`
	OwnerID            int       `json:"ownerId,omitempty"`
	Lang               string    `json:"lang,omitempty"`
	Timezone           string    `json:"timezone,omitempty"`
	ReportSendTime     string    `json:"reportSendTime,omitempty"`
	TextFormattingRule Format    `json:"textFormattingRule,omitempty"`
	Created            time.Time `json:"created,omitempty"`
	Updated            time.Time `json:"updated,omitempty"`
}

// SpaceNotification represents a notification of Space.
type SpaceNotification struct {
	Content string    `json:"content,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}
