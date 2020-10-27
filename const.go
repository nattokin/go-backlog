package backlog

// Order type
const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

// Fomat of Backlog wiki
const (
	FormatMarkdown format = "markdown"
	FormatBacklog  format = "backlog"
)

// Role type
const (
	_ role = iota
	RoleAdministrator
	RoleNormalUser
	RoleReporter
	RoleViewer
	RoleGuestReporter
	RoleGuestViewer
)
