package model

// User represents a Backlog user.
type User struct {
	ID          int    `json:"id,omitempty"`
	UserID      string `json:"userId,omitempty"`
	Name        string `json:"name,omitempty"`
	RoleType    Role   `json:"roleType,omitempty"`
	Lang        string `json:"lang,omitempty"`
	MailAddress string `json:"mailAddress,omitempty"`
}
