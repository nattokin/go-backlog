package sharedfile_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestProjectSharedFileService_List(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		expectError bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success-project-id": {
			projectIDOrKey: "1234",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/files", spath)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},

		"success-project-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/files", spath)
				return mock.NewJSONResponse(fixture.SharedFile.ListJSON), nil
			},
		},

		"error-project-id-or-key-empty": {
			projectIDOrKey: "",
			expectError:    true,
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			projectIDOrKey: "1234",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey: "1234",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := sharedfile.NewProjectService(method)

			files, err := s.List(context.Background(), tc.projectIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, files)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, files)
			assert.Len(t, files, len(fixture.SharedFile.List))

			for i, v := range fixture.SharedFile.List {
				assert.Equal(t, v.ID, files[i].ID)
				assert.Equal(t, v.Type, files[i].Type)
				assert.Equal(t, v.Dir, files[i].Dir)
				assert.Equal(t, v.Name, files[i].Name)
				assert.Equal(t, v.Size, files[i].Size)
			}
		})
	}
}
