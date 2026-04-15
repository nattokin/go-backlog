package project_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
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

func TestProjectService_All(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.ListJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.ListJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
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

func TestProjectService_One(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: strconv.Itoa(6),

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/6", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
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

func TestProjectService_Create(t *testing.T) {
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Post: mock.NewUnexpectedPostFn(t),
			}
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

func TestProjectService_Update(t *testing.T) {
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"success-projectIDOrKey-id": {
			projectIDOrKey: "1234",

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Patch: mock.NewUnexpectedPatchFn(t),
			}
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

func TestProjectService_Delete(t *testing.T) {
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

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"success-projectIDOrKey-id": {
			projectIDOrKey: "1234",

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Project.SingleJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Delete: mock.NewUnexpectedDeleteFn(t),
			}
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

// TestProjectService_contextPropagation verifies that the context passed to each
// ProjectService method is correctly relayed to the underlying method call.
func TestProjectService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T)
	}{
		{"All", func(t *testing.T) {
			s := project.NewService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.All(ctx) //nolint:errcheck
		}},
		{"One", func(t *testing.T) {
			s := project.NewService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.One(ctx, "TEST") //nolint:errcheck
		}},
		{"Create", func(t *testing.T) {
			s := project.NewService(&core.Method{
				Post: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Create(ctx, "KEY", "name", o.WithChartEnabled(true)) //nolint:errcheck
		}},
		{"Update", func(t *testing.T) {
			s := project.NewService(&core.Method{
				Patch: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Update(ctx, "TEST") //nolint:errcheck
		}},
		{"Delete", func(t *testing.T) {
			s := project.NewService(&core.Method{
				Delete: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			})
			s.Delete(ctx, "TEST") //nolint:errcheck
		}}}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t)
		})
	}
}
