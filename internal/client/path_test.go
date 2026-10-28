package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/client"
)

func TestPathJoin(t *testing.T) {
	cases := map[string]struct {
		segments []string
		want     string
	}{
		"single-segment": {
			segments: []string{"issues"},
			want:     "issues",
		},
		"multiple-segments": {
			segments: []string{"projects", "TEST", "versions"},
			want:     "projects/TEST/versions",
		},
		"dot-segments-not-resolved": {
			segments: []string{"issues", "..", "projects"},
			want:     "issues/../projects",
		},
		"reserved-characters-not-escaped": {
			segments: []string{"projects", "TEST", "files", "metadata", "資料 a?b#c%d"},
			want:     "projects/TEST/files/metadata/資料 a?b#c%d",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, tc.want, client.PathJoin(tc.segments...))
		})
	}
}
