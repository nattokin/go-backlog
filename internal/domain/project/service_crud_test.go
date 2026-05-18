package project_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_One(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: strconv.Itoa(6),

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",

			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",

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

			project, err := s.One(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, project)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, project)

			assert.Equal(t, 6, project.ID)
			assert.Equal(t, "TEST", project.ProjectKey)
			assert.Equal(t, "test", project.Name)
		})
	}
}

func TestService_Create(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		key  string
		name string
		opts []core.RequestOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-key-name": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.NotNil(t, form)
				assert.Equal(t, "TEST", form.Get("key"))
				assert.Equal(t, "test", form.Get("name"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"success-without-option": {
			key:  "TEST",
			name: "test",
			opts: []core.RequestOption{},

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "", form.Get("chartEnabled"))
				assert.Equal(t, "", form.Get("subtaskingEnabled"))
				assert.Equal(t, "", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "", form.Get("textFormattingRule"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"success-with-valid-option": {
			key:  "TEST",
			name: "test",

			opts: []core.RequestOption{
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule("backlog"),
			},

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"error-validation-key-empty": {
			key:  "",
			name: "test",

			wantErrType: &core.ValidationError{},
		},

		"error-validation-name-empty": {
			key:  "TEST",
			name: "",

			wantErrType: &core.ValidationError{},
		},

		"error-option-invalid-value": {
			key:  "TEST",
			name: "test",

			opts: []core.RequestOption{
				o.WithTextFormattingRule("invalid"),
			},

			wantErrType: &core.ValidationError{},
		},

		"error-option-invalid-type": {
			key:  "TEST",
			name: "test",

			opts: []core.RequestOption{mock.NewInvalidTypeOption()},

			wantErrType: &core.InvalidOptionKeyError{},
		},

		"error-client-network": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := project.NewService(method)

			project, err := s.Create(context.Background(), tc.key, tc.name, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, project)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, project)

			assert.Equal(t, tc.key, project.ProjectKey)
			assert.Equal(t, tc.name, project.Name)
		})
	}
}

func TestService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"success-projectIDOrKey-id": {
			projectIDOrKey: "1234",
			option:         o.WithName("test"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			option:         o.WithName("test"),

			wantErrType: &core.ValidationError{},
		},

		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey: "0",
			option:         o.WithName("test"),

			wantErrType: &core.ValidationError{},
		},

		"success-with-options": {
			projectIDOrKey: "TEST",
			option:         o.WithKey("TEST1"),

			opts: []core.RequestOption{
				o.WithName("test1"),
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule("backlog"),
				o.WithArchived(true),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "TEST1", form.Get("key"))
				assert.Equal(t, "test1", form.Get("name"))
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))
				assert.Equal(t, "true", form.Get("archived"))
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"error-option-invalid-value": {
			projectIDOrKey: "TEST",
			option:         o.WithTextFormattingRule("invalid"),

			wantErrType: &core.ValidationError{},
		},

		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			option:         mock.NewInvalidTypeOption(),

			wantErrType: &core.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			option:         o.WithName("test"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}
			s := project.NewService(method)

			project, err := s.Update(context.Background(), tc.projectIDOrKey, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, project)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, project)
		})
	}
}

func TestService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "1234",

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey: "0",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := project.NewService(method)

			project, err := s.Delete(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, project)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, project)

			assert.Equal(t, "TEST", project.ProjectKey)
		})
	}
}
