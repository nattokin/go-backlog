package client_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestPathJoin(t *testing.T) {
	cases := map[string]struct {
		segments []string

		want       string
		wantErrMsg string
	}{
		"single-segment": {
			segments: []string{"issues"},
			want:     "issues",
		},
		"multiple-segments": {
			segments: []string{"projects", "TEST", "versions"},
			want:     "projects/TEST/versions",
		},
		"issue-key": {
			segments: []string{"issues", "PRJ-1"},
			want:     "issues/PRJ-1",
		},
		"dots-inside-segment": {
			segments: []string{"projects", "TEST", "git", "repositories", "app.v2"},
			want:     "projects/TEST/git/repositories/app.v2",
		},
		"non-ascii-and-reserved-characters-not-escaped": {
			segments: []string{"projects", "TEST", "files", "metadata", "資料 a?b#c%d"},
			want:     "projects/TEST/files/metadata/資料 a?b#c%d",
		},

		"error-empty": {
			segments:   []string{"issues", ""},
			wantErrMsg: "invalid path segment: must not be empty",
		},
		"error-dot": {
			segments:   []string{"issues", "."},
			wantErrMsg: `invalid path segment ".": must not be a dot segment`,
		},
		"error-dot-dot": {
			segments:   []string{"issues", ".."},
			wantErrMsg: `invalid path segment "..": must not be a dot segment`,
		},
		"error-traversal": {
			segments:   []string{"issues", "../projects/PRJ"},
			wantErrMsg: `invalid path segment "../projects/PRJ": must not contain path separators`,
		},
		"error-slash": {
			segments:   []string{"issues", "PRJ-1/comments/1"},
			wantErrMsg: `invalid path segment "PRJ-1/comments/1": must not contain path separators`,
		},
		"error-backslash": {
			segments:   []string{"issues", `..\projects`},
			wantErrMsg: `invalid path segment "..\\projects": must not contain path separators`,
		},
		"error-control-character": {
			segments:   []string{"issues", "PRJ-1\n"},
			wantErrMsg: `invalid path segment "PRJ-1\n": must not contain control characters`,
		},
		"error-fixed-segment-with-slash": {
			segments:   []string{"projects/TEST"},
			wantErrMsg: `invalid path segment "projects/TEST": must not contain path separators`,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := client.PathJoin(tc.segments...)

			if tc.wantErrMsg != "" {
				assert.Empty(t, got)
				var ve *validation.Error
				if assert.ErrorAs(t, err, &ve) {
					assert.Equal(t, "path", ve.Target())
					assert.Equal(t, tc.wantErrMsg, ve.Message())
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
