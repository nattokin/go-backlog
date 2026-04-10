package model

// Activity represents a recent update or change in the project or space.
type Activity struct {
	ID            int              `json:"id,omitempty"`
	Project       *Project         `json:"project,omitempty"`
	Type          int              `json:"type,omitempty"`
	Content       *ActivityContent `json:"content,omitempty"`
	Notifications []*Notification  `json:"notifications,omitempty"`
	CreatedUser   *User            `json:"createdUser,omitempty"`
}

// ActivityContent represents the detailed content of an activity.
type ActivityContent struct {
	ID          int      `json:"id,omitempty"`
	KeyID       int      `json:"key_id,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	Comment     *Comment `json:"comment,omitempty"`
}
