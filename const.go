package backlog

// Order type
const (
	OrderAsc  = "asc"
	OrderDesc = "desc"
)

// Fomat of Backlog wiki
const (
	FormatMarkdown = "markdown"
	FormatBacklog  = "backlog"
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

type role int

func (r role) String() string {
	switch r {
	case RoleAdministrator:
		return "Administrator"
	case RoleNormalUser:
		return "NormalUser"
	case RoleReporter:
		return "Reporter"
	case RoleViewer:
		return "Viewer"
	case RoleGuestReporter:
		return "RoleGuestReporter"
	case RoleGuestViewer:
		return "RoleGuestViewer"
	default:
		return "unknown"
	}
}
