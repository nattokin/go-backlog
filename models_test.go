package backlog_test

import (
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestFormat_String(t *testing.T) {
	cases := map[string]struct {
		format backlog.Format
		want   string
	}{
		"Markdown": {
			format: backlog.FormatMarkdown,
			want:   "Markdown",
		},
		"Backlog": {
			format: backlog.FormatBacklog,
			want:   "Backlog",
		},
		"unknown": {
			format: backlog.Format("test"),
			want:   "unknown Format type test",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.format.String(), tc.want)
		})

	}
}

func TestOrder_String(t *testing.T) {
	cases := map[string]struct {
		format backlog.Order
		want   string
	}{
		"Markdown": {
			format: backlog.OrderAsc,
			want:   "Asc",
		},
		"Backlog": {
			format: backlog.OrderDesc,
			want:   "Desc",
		},
		"unknown": {
			format: backlog.Order("test"),
			want:   "unknown Order type test",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.format.String(), tc.want)
		})

	}
}

func TestRole_String(t *testing.T) {
	cases := map[string]struct {
		roleType backlog.Role
		want     string
	}{
		"Administrator": {
			roleType: backlog.RoleAdministrator,
			want:     "Administrator",
		},
		"NormalUser": {
			roleType: backlog.RoleNormalUser,
			want:     "NormalUser",
		},
		"Reporter": {
			roleType: backlog.RoleReporter,
			want:     "Reporter",
		},
		"Viewer": {
			roleType: backlog.RoleViewer,
			want:     "Viewer",
		},
		"GuestReporter": {
			roleType: backlog.RoleGuestReporter,
			want:     "GuestReporter",
		},
		"GuestViewer": {
			roleType: backlog.RoleGuestViewer,
			want:     "GuestViewer",
		},
		"unknown": {
			roleType: backlog.Role(0),
			want:     "unknown Role type 0",
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.roleType.String(), tc.want)
		})

	}
}
