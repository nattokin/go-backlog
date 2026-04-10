package model

import "fmt"

// Format defines the text formatting rule for the Backlog wiki.
const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
)

// Order defines the sort order (ascending or descending).
const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Role defines the type of user role within a project.
const (
	_ Role = iota
	RoleAdministrator
	RoleNormalUser
	RoleReporter
	RoleViewer
	RoleGuestReporter
	RoleGuestViewer
)

// Format defines the text formatting rule for the Backlog wiki.
type Format string

func (f Format) String() string {
	switch f {
	case FormatMarkdown:
		return "Markdown"
	case FormatBacklog:
		return "Backlog"
	default:
		return fmt.Sprintf("unknown Format type %s", string(f))
	}
}

// Order defines the sort order (ascending or descending).
type Order string

func (o Order) String() string {
	switch o {
	case OrderAsc:
		return "Asc"
	case OrderDesc:
		return "Desc"
	default:
		return fmt.Sprintf("unknown Order type %s", string(o))
	}
}

// Role defines the type of user role within a project.
type Role int

func (r Role) String() string {
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
		return "GuestReporter"
	case RoleGuestViewer:
		return "GuestViewer"
	default:
		return fmt.Sprintf("unknown Role type %d", r)
	}
}
