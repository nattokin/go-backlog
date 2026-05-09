package model

import "time"

// DiskUsageSpace represents disk usage for the entire space.
type DiskUsageSpace struct {
	DiskUsageBase
	Capacity int                 `json:"capacity,omitempty"`
	Details  []*DiskUsageProject `json:"details,omitempty"`
}

// Space represents a Backlog space.
type Space struct {
	SpaceKey           string    `json:"spaceKey,omitempty"`
	Name               string    `json:"name,omitempty"`
	OwnerID            int       `json:"ownerId,omitempty"`
	Lang               string    `json:"lang,omitempty"`
	Timezone           string    `json:"timezone,omitempty"`
	ReportSendTime     string    `json:"reportSendTime,omitempty"`
	TextFormattingRule string    `json:"textFormattingRule,omitempty"`
	Created            time.Time `json:"created,omitempty"`
	Updated            time.Time `json:"updated,omitempty"`
}

// SpaceNotification represents the notification message of a Backlog space.
type SpaceNotification struct {
	Content string    `json:"content,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}
