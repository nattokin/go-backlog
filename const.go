package backlog

import "github.com/nattokin/go-backlog/internal/model"

// Order defines the sort order (ascending or descending).
const (
	OrderAsc  = model.OrderAsc
	OrderDesc = model.OrderDesc
)

// Format defines the text formatting rule for the Backlog wiki.
const (
	FormatMarkdown = model.FormatMarkdown
	FormatBacklog  = model.FormatBacklog
)

// Role defines the type of user role within a project.
const (
	RoleAdministrator = model.RoleAdministrator
	RoleNormalUser    = model.RoleNormalUser
	RoleReporter      = model.RoleReporter
	RoleViewer        = model.RoleViewer
	RoleGuestReporter = model.RoleGuestReporter
	RoleGuestViewer   = model.RoleGuestViewer
)

const MaxActivityTypeID = 26
