package backlog

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat_String(t *testing.T) {
	cases := map[string]struct {
		format Format
		want   string
	}{
		"Markdown": {
			format: FormatMarkdown,
			want:   "Markdown",
		},
		"Backlog": {
			format: FormatBacklog,
			want:   "Backlog",
		},
		"Unknown": {
			format: Format("test"),
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
		format Order
		want   string
	}{
		"Markdown": {
			format: OrderAsc,
			want:   "Asc",
		},
		"Backlog": {
			format: OrderDesc,
			want:   "Desc",
		},
		"Unknown": {
			format: Order("test"),
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
		roleType Role
		want     string
	}{
		"Administrator": {
			roleType: RoleAdministrator,
			want:     "Administrator",
		},
		"NormalUser": {
			roleType: RoleNormalUser,
			want:     "NormalUser",
		},
		"Reporter": {
			roleType: RoleReporter,
			want:     "Reporter",
		},
		"Viewer": {
			roleType: RoleViewer,
			want:     "Viewer",
		},
		"GuestReporter": {
			roleType: RoleGuestReporter,
			want:     "GuestReporter",
		},
		"GuestViewer": {
			roleType: RoleGuestViewer,
			want:     "GuestViewer",
		},
		"Unknown": {
			roleType: Role(0),
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
