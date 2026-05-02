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
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_All(t *testing.T) {
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

			projects, err := s.All(context.Background(), tc.opts...)

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
				o.WithTextFormattingRule(model.FormatBacklog),
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
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"success-projectIDOrKey-id": {
			projectIDOrKey: "1234",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)
				return mock.NewJSONResponse(fixture.Project.SingleJSON), nil
			},

			wantErrType: nil,
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",

			wantErrType: &core.ValidationError{},
		},

		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey: "0",

			wantErrType: &core.ValidationError{},
		},

		"success-with-options": {
			projectIDOrKey: "TEST",

			opts: []core.RequestOption{
				o.WithKey("TEST1"),
				o.WithName("test1"),
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule(model.FormatBacklog),
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

			opts: []core.RequestOption{
				o.WithTextFormattingRule("invalid"),
			},

			wantErrType: &core.ValidationError{},
		},

		"error-option-invalid-type": {
			projectIDOrKey: "TEST",

			opts: []core.RequestOption{mock.NewInvalidTypeOption()},

			wantErrType: &core.InvalidOptionKeyError{},
		},

		"error-client-network": {
			projectIDOrKey: "TEST",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
			projectIDOrKey: "TEST",

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

			project, err := s.Update(context.Background(), tc.projectIDOrKey, tc.opts...)

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

func TestCategoryService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Category.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/categories", spath)
				return mock.NewJSONResponse(fixture.Category.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewCategoryService(method)

			categories, err := s.All(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, categories)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, categories)
			assert.Len(t, categories, tc.wantLen)
		})
	}
}

func TestCategoryService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Bug",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Bug",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Bug",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Bug",

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
			s := project.NewCategoryService(method)

			category, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
			assert.Equal(t, "Bug", category.Name)
		})
	}
}

func TestCategoryService_Update(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		categoryID     int
		name           string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				assert.Equal(t, "Bug Fixed", form.Get("name"))
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			categoryID:     12,
			name:           "Bug Fixed",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey: "TEST",
			categoryID:     0,
			name:           "Bug Fixed",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			categoryID:     12,
			name:           "Bug Fixed",

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
			s := project.NewCategoryService(method)

			category, err := s.Update(context.Background(), tc.projectIDOrKey, tc.categoryID, tc.name)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
		})
	}
}

func TestCategoryService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		categoryID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			categoryID:     12,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/categories/12", spath)
				assert.NotNil(t, form)
				return mock.NewJSONResponse(fixture.Category.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			categoryID:     12,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-categoryID-zero": {
			projectIDOrKey: "TEST",
			categoryID:     0,
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			categoryID:     12,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			categoryID:     12,

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
			s := project.NewCategoryService(method)

			category, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.categoryID)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, category)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, category)
			assert.Equal(t, 12, category.ID)
			assert.Equal(t, "Bug", category.Name)
		})
	}
}

func TestIssueTypeService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.IssueType.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/issueTypes", spath)
				return mock.NewJSONResponse(fixture.IssueType.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewIssueTypeService(method)

			issueTypes, err := s.All(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueTypes)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueTypes)
			assert.Len(t, issueTypes, tc.wantLen)
		})
	}
}

func TestIssueTypeService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string
		color          string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes", spath)
				assert.Equal(t, "Bug", form.Get("name"))
				assert.Equal(t, "#e30000", form.Get("color"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Bug",
			color:          "#e30000",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			color:          "#e30000",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-color-empty": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Bug",
			color:          "#e30000",

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
			s := project.NewIssueTypeService(method)

			issueType, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name, tc.color)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
			assert.Equal(t, "Bug", issueType.Name)
			assert.Equal(t, "#e30000", issueType.Color)
		})
	}
}

func TestIssueTypeService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		issueTypeID    int
		option         core.RequestOption
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			opts:           []core.RequestOption{o.WithColor("#990000")},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "Bug Updated", form.Get("name"))
				assert.Equal(t, "#990000", form.Get("color"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-option-only": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "Bug Updated", form.Get("name"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey: "TEST",
			issueTypeID:    0,
			option:         o.WithName("Bug Updated"),
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         mock.NewInvalidTypeOption(),
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			issueTypeID:    1,
			option:         o.WithName("Bug Updated"),

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
			s := project.NewIssueTypeService(method)

			issueType, err := s.Update(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}

func TestIssueTypeService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey        string
		issueTypeID           int
		substituteIssueTypeID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/issueTypes/1", spath)
				assert.Equal(t, "2", form.Get("substituteIssueTypeId"))
				return mock.NewJSONResponse(fixture.IssueType.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:        "",
			issueTypeID:           1,
			substituteIssueTypeID: 2,
			wantErrType:           &core.ValidationError{},
		},
		"error-validation-issueTypeID-zero": {
			projectIDOrKey:        "TEST",
			issueTypeID:           0,
			substituteIssueTypeID: 2,
			wantErrType:           &core.ValidationError{},
		},
		"error-validation-substituteIssueTypeID-zero": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 0,
			wantErrType:           &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey:        "TEST",
			issueTypeID:           1,
			substituteIssueTypeID: 2,

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
			s := project.NewIssueTypeService(method)

			issueType, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.issueTypeID, tc.substituteIssueTypeID)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, issueType)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issueType)
			assert.Equal(t, 1, issueType.ID)
		})
	}
}

func TestStatusService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success-projectIDOrKey-key": {
			projectIDOrKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "6",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6/statuses", spath)
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
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
			s := project.NewStatusService(method)

			statuses, err := s.All(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, statuses)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, statuses)
			assert.Len(t, statuses, tc.wantLen)
		})
	}
}

func TestStatusService_Create(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		name           string
		color          string

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses", spath)
				assert.Equal(t, "Open", form.Get("name"))
				assert.Equal(t, "#ed8077", form.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			name:           "Open",
			color:          "#ed8077",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-name-empty": {
			projectIDOrKey: "TEST",
			name:           "",
			color:          "#ed8077",
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-color-empty": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			name:           "Open",
			color:          "#ed8077",

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
			s := project.NewStatusService(method)

			status, err := s.Create(context.Background(), tc.projectIDOrKey, tc.name, tc.color)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
			assert.Equal(t, "Open", status.Name)
			assert.Equal(t, "#ed8077", status.Color)
		})
	}
}

func TestStatusService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		statusID       int
		opts           []core.RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusID:       1,
			opts: []core.RequestOption{
				o.WithName("Open Updated"),
				o.WithColor("#f5ab35"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "Open Updated", form.Get("name"))
				assert.Equal(t, "#f5ab35", form.Get("color"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"success-without-option": {
			projectIDOrKey: "TEST",
			statusID:       1,

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			statusID:       1,
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey: "TEST",
			statusID:       0,
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			statusID:       1,
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			statusID:       1,

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			statusID:       1,

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
			s := project.NewStatusService(method)

			status, err := s.Update(context.Background(), tc.projectIDOrKey, tc.statusID, tc.opts...)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
		})
	}
}

func TestStatusService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey     string
		statusID           int
		substituteStatusID int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/1", spath)
				assert.Equal(t, "2", form.Get("substituteStatusId"))
				return mock.NewJSONResponse(fixture.Status.SingleJSON), nil
			},

			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:     "",
			statusID:           1,
			substituteStatusID: 2,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey:     "TEST",
			statusID:           0,
			substituteStatusID: 2,
			wantErrType:        &core.ValidationError{},
		},
		"error-validation-substituteStatusID-zero": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 0,
			wantErrType:        &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey:     "TEST",
			statusID:           1,
			substituteStatusID: 2,

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
			s := project.NewStatusService(method)

			status, err := s.Delete(context.Background(), tc.projectIDOrKey, tc.statusID, tc.substituteStatusID)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, status)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, status)
			assert.Equal(t, 1, status.ID)
		})
	}
}

func TestStatusService_UpdateOrder(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		statusIDs      []int

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantLen     int
		wantErrType error
	}{
		"success": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{2, 1},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/statuses/updateDisplayOrder", spath)
				assert.Equal(t, []string{"2", "1"}, form["statusId[]"])
				return mock.NewJSONResponse(fixture.Status.ListJSON), nil
			},

			wantLen:     2,
			wantErrType: nil,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			statusIDs:      []int{1, 2},
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusIDs-empty": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{},
			wantErrType:    &core.ValidationError{},
		},
		"error-validation-statusID-zero": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 0},
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 2},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			statusIDs:      []int{1, 2},

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
			s := project.NewStatusService(method)

			statuses, err := s.UpdateOrder(context.Background(), tc.projectIDOrKey, tc.statusIDs)

			if tc.wantErrType != nil {
				require.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, statuses)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, statuses)
			assert.Len(t, statuses, tc.wantLen)
		})
	}
}

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.All(ctx) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.One(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewService(m)
			s.Create(ctx, "KEY", "name", o.WithChartEnabled(true)) //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewService(m)
			s.Update(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewService(m)
			s.Delete(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.DiskUsage", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.DiskUsage(ctx, "TEST") //nolint:errcheck
		}},
		{"CategoryService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"CategoryService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Create(ctx, "TEST", "Bug") //nolint:errcheck
		}},
		{"CategoryService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Update(ctx, "TEST", 12, "Bug Fixed") //nolint:errcheck
		}},
		{"CategoryService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Delete(ctx, "TEST", 12) //nolint:errcheck
		}},
		{"IssueTypeService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"IssueTypeService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Create(ctx, "TEST", "Bug", "#e30000") //nolint:errcheck
		}},
		{"IssueTypeService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Update(ctx, "TEST", 1, o.WithName("Bug Updated")) //nolint:errcheck
		}},
		{"IssueTypeService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Delete(ctx, "TEST", 1, 2) //nolint:errcheck
		}},
		{"StatusService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewStatusService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"StatusService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Create(ctx, "TEST", "Open", "#ed8077") //nolint:errcheck
		}},
		{"StatusService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Update(ctx, "TEST", 1) //nolint:errcheck
		}},
		{"StatusService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Delete(ctx, "TEST", 1, 2) //nolint:errcheck
		}},
		{"StatusService.UpdateOrder", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewStatusService(m)
			s.UpdateOrder(ctx, "TEST", []int{1, 2}) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
