package backlog

import (
	"bytes"
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
		options     []*QueryOption
		expectError bool
		wantIDs     []int
		wantNames   []string
		mockGetFn   func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success-without-option": {
			options:     []*QueryOption{},
			expectError: false,
			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			mockGetFn: newMockGetFn(t, "projects", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
			}),
		},
		"success-with-valid-option": {
			options: []*QueryOption{
				o.WithQueryAll(false),
				o.WithQueryArchived(true),
			},
			expectError: false,
			wantIDs:     []int{1, 2, 3},
			wantNames:   []string{"test", "test2", "test3"},
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects", spath)
				assert.Equal(t, "false", query.Get("all"))
				assert.Equal(t, "true", query.Get("archived"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
				}, nil
			},
		},
		"error-with-option-handler": {
			options: []*QueryOption{{
				queryAll,
				nil,
				func(p *QueryParams) error { return errors.New("error") },
			}},
			expectError: true,
			mockGetFn:   newUnexpectedGetFn(t),
		},
		"error-with-invalid-option": {
			options: []*QueryOption{{
				0,
				nil,
				func(p *QueryParams) error { return nil },
			}},
			expectError: true,
			mockGetFn:   newUnexpectedGetFn(t),
		},
		"error-client-failure": {
			options:     []*QueryOption{},
			expectError: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			options:     []*QueryOption{},
			expectError: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Get = tc.mockGetFn

			projects, err := s.All(tc.options...)
			if tc.expectError {
				assert.Error(t, err)
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
		expectError    bool
		mockGetFn      func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success-with-projectKey": {
			projectIDOrKey: "TEST",
			expectError:    false,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},
		},
		"success-with-projectID": {
			projectIDOrKey: strconv.Itoa(6),
			expectError:    false,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/6", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}, nil
			},
		},
		"error-with-empty-projectIDOrKey": {
			projectIDOrKey: "",
			expectError:    true,
			mockGetFn:      newUnexpectedGetFn(t),
		},
		"error-client-failure": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Get = tc.mockGetFn

			project, err := s.One(tc.projectIDOrKey)
			if tc.expectError {
				assert.Error(t, err)
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
		key         string
		name        string
		options     []*FormOption
		expectError bool
		mockPostFn  func(spath string, form *FormParams) (*http.Response, error)
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
		},

		"error-empty-key": {
			key:         "",
			name:        "test",
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t),
		},

		"error-empty-name": {
			key:         "TEST",
			name:        "",
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t),
		},

		"error-option-handler": {
			key:  "TEST",
			name: "test",
			options: []*FormOption{
				o.WithFormTextFormattingRule("invalid"),
			},
			expectError: true,
			mockPostFn: newUnexpectedPostFn(
				t,
			),
		},

		"error-invalid-option": {
			key:         "TEST",
			name:        "test",
			options:     []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t),
		},

		"error-client-failure": {
			key:  "TEST",
			name: "test",
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
			expectError: true,
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
			expectError: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Post = tc.mockPostFn

			project, err := s.Create(tc.key, tc.name, tc.options...)

			if tc.expectError {
				assert.Error(t, err)
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
		expectError    bool
		mockPatchFn    func(spath string, form *FormParams) (*http.Response, error)
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
		},

		"error-empty-key": {
			projectIDOrKey: "",
			expectError:    true,
			mockPatchFn:    newUnexpectedPatchFn(t),
		},

		"error-zero-id": {
			projectIDOrKey: "0",
			expectError:    true,
			mockPatchFn:    newUnexpectedPatchFn(t),
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
		},

		"error-option-validation": {
			projectIDOrKey: "TEST",
			options: []*FormOption{
				o.WithFormTextFormattingRule("invalid"),
			},
			expectError: true,
			mockPatchFn: newUnexpectedPatchFn(t),
		},

		"error-invalid-option": {
			projectIDOrKey: "TEST",
			options:        []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},
			expectError:    true,
			mockPatchFn:    newUnexpectedPatchFn(t),
		},

		"error-client-failure": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Patch = tc.mockPatchFn

			project, err := s.Update(tc.projectIDOrKey, tc.options...)

			if tc.expectError {
				assert.Error(t, err)
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
		expectError    bool
		mockDeleteFn   func(spath string, form *FormParams) (*http.Response, error)
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
		},

		"error-empty-key": {
			projectIDOrKey: "",
			expectError:    true,
			mockDeleteFn:   newUnexpectedDeleteFn(t),
		},

		"error-zero-id": {
			projectIDOrKey: "0",
			expectError:    true,
			mockDeleteFn:   newUnexpectedDeleteFn(t),
		},

		"error-client-failure": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Delete = tc.mockDeleteFn

			project, err := s.Delete(tc.projectIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, project)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, project)

			assert.Equal(t, "TEST", project.ProjectKey)
		})
	}
}
