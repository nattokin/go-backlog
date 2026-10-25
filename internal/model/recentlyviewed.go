package model

import "time"

// RecentlyViewedIssue represents an entry in the recently viewed issues list.
type RecentlyViewedIssue struct {
	Issue   *Issue    `json:"issue,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

// RecentlyViewedProject represents an entry in the recently viewed projects list.
type RecentlyViewedProject struct {
	Project *Project  `json:"project,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}

// RecentlyViewedWiki represents an entry in the recently viewed wikis list.
type RecentlyViewedWiki struct {
	Page    *Wiki     `json:"page,omitempty"`
	Updated time.Time `json:"updated,omitempty"`
}
