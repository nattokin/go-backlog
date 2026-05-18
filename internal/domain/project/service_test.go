package project_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantIDs     []int
		wantNames   []string
		wantErrType error
	}{
		"success-without-option": {
			opts: []core.RequestOption{},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				require.NotNil(t, query)
				assert.Empty(t, query)
				return mock.NewJSONResponse(fixture.Project.ListJSON), nil
			},

			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			wantErrType: nil,
		},
		"success-with-valid-option": {
			opts: []core.RequestOption{
				o.WithAll(false),
				o.WithArchived(true),
			},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.Equal(t, "false", query.Get("all"))
				assert.Equal(t, "true", query.Get("archived"))
				return mock.NewJSONResponse(fixture.Project.ListJSON), nil
			},

			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			wantErrType: nil,
		},
		"error-option-set-failed": {
			opts: []core.RequestOption{mock.NewFailingSetOption(core.ParamAll)},

			wantErrType: errors.New(""),
		},
		"error-option-invalid-type": {
			opts: []core.RequestOption{mock.NewInvalidTypeOption()},

			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			opts: []core.RequestOption{},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			opts: []core.RequestOption{},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := project.NewService(method)

			projects, err := s.List(context.Background(), tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, projects)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, projects)
			assert.Equal(t, len(tc.wantIDs), len(projects))

			for i := range projects {
				assert.Equal(t, tc.wantIDs[i], projects[i].ID)
				assert.Equal(t, tc.wantNames[i], projects[i].Name)
			}
		})
	}
}

func TestService_DiskUsage(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType   error
		wantProjectID int
		wantIssue     int
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/diskUsage", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Project.DiskUsageJSON), nil
			},
			wantProjectID: 1,
			wantIssue:     11931,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/diskUsage", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/diskUsage", spath)
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := project.NewService(method)

			got, err := s.DiskUsage(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantProjectID, got.ProjectID)
			assert.Equal(t, tc.wantIssue, got.Issue)
		})
	}
}

func TestService_Icon(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockDownloadFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/image", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("icon.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "icon.png",
			wantContentType: "image/png",
		},
		"success-project-id": {
			projectIDOrKey: "123",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/123/image", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("icon.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "icon.png",
			wantContentType: "image/png",
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDownloadFn != nil {
				method.Download = tc.mockDownloadFn
			}

			s := project.NewService(method)

			got, err := s.Icon(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantFilename, got.Filename)
			assert.Equal(t, tc.wantContentType, got.ContentType)
			require.NotNil(t, got.Body)
			got.Body.Close()
		})
	}
}
