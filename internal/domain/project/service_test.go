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
		opts []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantIDs                []int
		wantNames              []string
		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success-without-option": {
			opts: []*core.APIParamOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				return mock.NewResponse(fixture.Project.ListJSON), nil
			},
			wantIDs:   []int{1, 2, 3},
			wantNames: []string{"test", "test2", "test3"},
		},
		"success-with-option": {
			opts: []*core.APIParamOption{o.WithAll(false), o.WithArchived(true)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "false", query.Get("all"))
				assert.Equal(t, "true", query.Get("archived"))
				return mock.NewResponse(fixture.Project.ListJSON), nil
			},
			wantIDs:   []int{1, 2, 3},
			wantNames: []string{"test", "test2", "test3"},
		},

		// WithTextFormattingRule is not in List's valid types (all, archived only)
		// → InvalidOptionKeyError, not ValidationErrors
		"error-validation-opt-single": {
			opts:        []*core.APIParamOption{o.WithTextFormattingRule("invalid")},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- fail-fast: nil option among valid values ---
		"error-nil-option-with-valid-values": {
			opts:                   []*core.APIParamOption{o.WithAll(true), nil},
			wantInvalidOptionError: true,
		},
		// WithTextFormattingRule hits invalid key before nil is checked
		// → InvalidOptionKeyError
		"error-nil-option-with-invalid-values": {
			opts:        []*core.APIParamOption{o.WithTextFormattingRule("invalid"), nil},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type-with-valid-values": {
			opts:        []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-type-with-invalid-values": {
			opts:        []*core.APIParamOption{o.WithTextFormattingRule("invalid"), mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-option-set-failed": {
			opts:        []*core.APIParamOption{mock.NewFailingSetOption(core.ParamAll)},
			wantErrType: errors.New(""),
		},
		"error-client-network": {
			opts: []*core.APIParamOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			opts: []*core.APIParamOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, projects)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, projects)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantErrType            error
		wantValidationErrCount int
		wantProjectID          int
		wantIssue              int
	}{
		"success": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/diskUsage", spath)
				return mock.NewResponse(fixture.Project.DiskUsageJSON), nil
			},
			wantProjectID: 1,
			wantIssue:     11931,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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

		wantErrType            error
		wantValidationErrCount int
		wantFilename           string
		wantContentType        string
	}{
		"success-key": {
			projectIDOrKey: "TEST",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/image", spath)
				return mock.NewBinaryResponse("icon.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "icon.png",
			wantContentType: "image/png",
		},
		"success-id": {
			projectIDOrKey: "123",
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/123/image", spath)
				return mock.NewBinaryResponse("icon.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "icon.png",
			wantContentType: "image/png",
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			wantValidationErrCount: 1,
		},

		// --- other errors ---
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

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
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
