package model_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/model"
)

func TestFormat_String(t *testing.T) {
	cases := map[string]struct {
		format model.Format
		want   string
	}{
		"Markdown": {
			format: model.FormatMarkdown,
			want:   "Markdown",
		},
		"Backlog": {
			format: model.FormatBacklog,
			want:   "Backlog",
		},
		"Unknown": {
			format: model.Format("test"),
			want:   "unknown Format type test",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.format.String(), tc.want)
		})

	}
}

func TestOrder_String(t *testing.T) {
	cases := map[string]struct {
		format model.Order
		want   string
	}{
		"Markdown": {
			format: model.OrderAsc,
			want:   "Asc",
		},
		"Backlog": {
			format: model.OrderDesc,
			want:   "Desc",
		},
		"Unknown": {
			format: model.Order("test"),
			want:   "unknown Order type test",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.format.String(), tc.want)
		})

	}
}

func TestRole_String(t *testing.T) {
	cases := map[string]struct {
		roleType model.Role
		want     string
	}{
		"Administrator": {
			roleType: model.RoleAdministrator,
			want:     "Administrator",
		},
		"NormalUser": {
			roleType: model.RoleNormalUser,
			want:     "NormalUser",
		},
		"Reporter": {
			roleType: model.RoleReporter,
			want:     "Reporter",
		},
		"Viewer": {
			roleType: model.RoleViewer,
			want:     "Viewer",
		},
		"GuestReporter": {
			roleType: model.RoleGuestReporter,
			want:     "GuestReporter",
		},
		"GuestViewer": {
			roleType: model.RoleGuestViewer,
			want:     "GuestViewer",
		},
		"Unknown": {
			roleType: model.Role(0),
			want:     "unknown Role type 0",
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			assert.Equal(t, tc.roleType.String(), tc.want)
		})

	}
}
