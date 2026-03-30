package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectService_All(t *testing.T) {
	o := newProjectOptionService()

	cases := map[string]struct {
		options []*QueryOption

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantIDs     []int
		wantNames   []string
		wantErrType error
	}{
		"success-without-option": {
			options: []*QueryOption{},

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				require.NotNil(t, query)
				assert.Empty(t, query.Values)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
				}, nil
			},

			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			wantErrType: nil,
		},
		"success-with-valid-option": {
			options: []*QueryOption{
				o.WithQueryAll(false),
				o.WithQueryArchived(true),
			},

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.Equal(t, "false", query.Get("all"))
				assert.Equal(t, "true", query.Get("archived"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
				}, nil
			},

			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			wantErrType: nil,
		},
		"error-with-option-handler": {
			options: []*QueryOption{{
				queryAll,
				nil,
				func(p *QueryParams) error { return errors.New("error") },
			}},

			wantErrType: errors.New(""),
		},
		"error-with-invalid-option": {
			options: []*QueryOption{{
				0,
				nil,
				func(p *QueryParams) error { return nil },
			}},

			wantErrType: &InvalidOptionError[queryType]{},
		},
		"error-client-failure": {
			options: []*QueryOption{},

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-invalid-json": {
			options: []*QueryOption{},

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			projects, err := s.All(tc.options...)

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

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantErrType error
	}{
		"success-with-projectKey": {
			projectIDOrKey: "TEST",

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"success-with-projectID": {
			projectIDOrKey: strconv.Itoa(6),

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/6", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"error-with-empty-projectIDOrKey": {
			projectIDOrKey: "",

			wantErrType: &ValidationError{},
		},
		"error-client-failure": {
			projectIDOrKey: "TEST",

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			project, err := s.One(tc.projectIDOrKey)

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
	o := newProjectOptionService()

	cases := map[string]struct {
		key     string
		name    string
		options []*FormOption

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
	}{
		"success-basic-create": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.NotNil(t, form)
				assert.Equal(t, "TEST", form.Get("key"))
				assert.Equal(t, "test", form.Get("name"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"success-without-option": {
			key:     "TEST",
			name:    "test",
			options: []*FormOption{},

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "", form.Get("chartEnabled"))
				assert.Equal(t, "", form.Get("subtaskingEnabled"))
				assert.Equal(t, "", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "", form.Get("textFormattingRule"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"success-with-valid-option": {
			key:  "TEST",
			name: "test",

			options: []*FormOption{
				o.WithFormChartEnabled(true),
				o.WithFormSubtaskingEnabled(true),
				o.WithFormProjectLeaderCanEditProjectLeader(true),
				o.WithFormTextFormattingRule(FormatBacklog),
			},

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"error-empty-key": {
			key:  "",
			name: "test",

			wantErrType: &ValidationError{},
		},

		"error-empty-name": {
			key:  "TEST",
			name: "",

			wantErrType: &ValidationError{},
		},

		"error-option-validation": {
			key:  "TEST",
			name: "test",

			options: []*FormOption{
				o.WithFormTextFormattingRule("invalid"),
			},

			wantErrType: &ValidationError{},
		},

		"error-invalid-option": {
			key:  "TEST",
			name: "test",

			options: []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},

			wantErrType: &InvalidOptionError[formType]{},
		},

		"error-client-failure": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},

		"error-invalid-json": {
			key:  "TEST",
			name: "test",

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()

			// default: unexpected API call
			s.method.Post = newUnexpectedPostFn(t)

			if tc.mockPostFn != nil {
				s.method.Post = tc.mockPostFn
			}

			project, err := s.Create(tc.key, tc.name, tc.options...)

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
	o := newProjectOptionService()

	cases := map[string]struct {
		projectIDOrKey string
		options        []*FormOption

		mockPatchFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
	}{
		"success-basic": {
			projectIDOrKey: "TEST",

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.NotNil(t, form)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"success-project-id": {
			projectIDOrKey: "1234",

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"error-empty-key": {
			projectIDOrKey: "",

			wantErrType: &ValidationError{},
		},

		"error-zero-id": {
			projectIDOrKey: "0",

			wantErrType: &ValidationError{},
		},

		"success-with-options": {
			projectIDOrKey: "TEST",

			options: []*FormOption{
				o.WithFormKey("TEST1"),
				o.WithFormName("test1"),
				o.WithFormChartEnabled(true),
				o.WithFormSubtaskingEnabled(true),
				o.WithFormProjectLeaderCanEditProjectLeader(true),
				o.WithFormTextFormattingRule(FormatBacklog),
				o.WithFormArchived(true),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "TEST1", form.Get("key"))
				assert.Equal(t, "test1", form.Get("name"))
				assert.Equal(t, "true", form.Get("chartEnabled"))
				assert.Equal(t, "true", form.Get("subtaskingEnabled"))
				assert.Equal(t, "true", form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, "backlog", form.Get("textFormattingRule"))
				assert.Equal(t, "true", form.Get("archived"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},

		"error-option-validation": {
			projectIDOrKey: "TEST",

			options: []*FormOption{
				o.WithFormTextFormattingRule("invalid"),
			},

			wantErrType: &ValidationError{},
		},

		"error-invalid-option": {
			projectIDOrKey: "TEST",

			options: []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},

			wantErrType: &InvalidOptionError[formType]{},
		},

		"error-client-failure": {
			projectIDOrKey: "TEST",

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},

		"error-invalid-json": {
			projectIDOrKey: "TEST",

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()

			// default: unexpected API call
			s.method.Patch = newUnexpectedPatchFn(t)

			if tc.mockPatchFn != nil {
				s.method.Patch = tc.mockPatchFn
			}

			project, err := s.Update(tc.projectIDOrKey, tc.options...)

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

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
	}{
		"success-project-key": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.NotNil(t, form)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"success-project-id": {
			projectIDOrKey: "1234",

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"error-empty-key": {
			projectIDOrKey: "",

			wantErrType: &ValidationError{},
		},
		"error-zero-id": {
			projectIDOrKey: "0",

			wantErrType: &ValidationError{},
		},
		"error-client-failure": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()

			// default: unexpected API call
			s.method.Delete = newUnexpectedDeleteFn(t)

			if tc.mockDeleteFn != nil {
				s.method.Delete = tc.mockDeleteFn
			}

			project, err := s.Delete(tc.projectIDOrKey)

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
