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
	t.Parallel()

	want := struct {
		idList   []int
		nameList []string
	}{
		idList:   []int{1, 2, 3},
		nameList: []string{"test", "test2", "test3"},
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
		}
		return resp, nil
	}

	projects, err := s.All()
	assert.NoError(t, err)
	require.NotNil(t, projects)
	count := len(projects)
	require.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], projects[i].ID)
		assert.Equal(t, want.nameList[i], projects[i].Name)
	}
}

func TestProjectService_All_option(t *testing.T) {
	o := newProjectOptionService()
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*QueryOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*QueryOption{},
			wantError: false,
			want: options{
				all:      "",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*QueryOption{
				o.WithQueryAll(false),
				o.WithQueryArchived(true),
			},
			wantError: false,
			want: options{
				all:      "false",
				archived: "true",
			},
		},
		"OptionError": {
			options: []*QueryOption{{queryAll, nil, func(p *QueryParams) error {
				return errors.New("error")
			}}},
			wantError: true,
		},
		"InvalidOption": {
			options:   []*QueryOption{{0, nil, func(p *QueryParams) error { return nil }}},
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, tc.want.all, query.Get("all"))
				assert.Equal(t, tc.want.archived, query.Get("archived"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
				}
				return resp, nil
			}

			if resp, err := s.All(tc.options...); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_AdminAll(t *testing.T) {
	t.Parallel()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Equal(t, want.all, query.Get("all"))
		assert.Equal(t, want.archived, query.Get("archived"))

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.AdminAll()
	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_AdminAll_option(t *testing.T) {
	o := newProjectOptionService()
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*QueryOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*QueryOption{},
			wantError: false,
			want: options{
				all:      "true",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*QueryOption{
				o.WithQueryArchived(true),
			},
			wantError: false,
			want: options{
				all:      "true",
				archived: "true",
			},
		},
		"OptionError": {
			options:   []*QueryOption{{queryArchived, nil, func(p *QueryParams) error { return errors.New("error") }}},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*QueryOption{
				o.WithQueryAll(true),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, tc.want.all, query.Get("all"))
				assert.Equal(t, tc.want.archived, query.Get("archived"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
				}
				return resp, nil
			}

			if resp, err := s.AdminAll(tc.options...); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_AllArchived(t *testing.T) {
	t.Parallel()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		archived: "true",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Equal(t, want.archived, query.Get("archived"))

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.AllArchived()
	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_AdminAllArchived(t *testing.T) {
	t.Parallel()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "true",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Equal(t, want.all, query.Get("all"))
		assert.Equal(t, want.archived, query.Get("archived"))

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.AdminAllArchived()

	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_AllUnarchived(t *testing.T) {
	t.Parallel()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		archived: "false",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Equal(t, want.archived, query.Get("archived"))

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.AllUnarchived()

	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_AdminAllUnarchived(t *testing.T) {
	t.Parallel()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "false",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Equal(t, want.all, query.Get("all"))
		assert.Equal(t, want.archived, query.Get("archived"))

		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	projects, err := s.AdminAllUnarchived()

	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_All_clientError(t *testing.T) {
	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		return nil, errors.New("error")
	}

	resp, err := s.All()
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_One_key(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath string
		key   string
		name  string
	}{
		spath: "projects/" + projectKey,
		key:   projectKey,
		name:  "test",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Nil(t, query)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.One(projectKey)
	require.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_id(t *testing.T) {
	t.Parallel()

	projectID := 6

	want := struct {
		spath string
		id    int
		name  string
	}{
		spath: "projects/" + strconv.Itoa(projectID),
		id:    projectID,
		name:  "test",
	}

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.Nil(t, query)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.One(strconv.Itoa(projectID))
	require.NoError(t, err)
	assert.Equal(t, want.id, project.ID)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_key_error(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.One("")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectService_One_clientError(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		return nil, errors.New("error")
	}

	resp, err := s.One("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_One_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	project, err := s.One("TEST")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectService_Create(t *testing.T) {
	t.Parallel()

	key := "TEST"
	name := "test"

	want := struct {
		spath string
		key   string
		name  string
	}{
		spath: "projects",
		key:   key,
		name:  name,
	}

	s := newProjectService()
	s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.NotNil(t, form)
		assert.Equal(t, want.key, form.Get("key"))
		assert.Equal(t, want.name, form.Get("name"))
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.Create(key, name)
	require.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_Create_param(t *testing.T) {
	cases := map[string]struct {
		key       string
		name      string
		wantError bool
	}{
		"WithoutOption": {
			key:       "TEST",
			name:      "test",
			wantError: false,
		},
		"KeyEnpty": {
			key:       "",
			name:      "test",
			wantError: true,
		},
		"NameEmpty": {
			key:       "TEST",
			name:      "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, tc.key, form.Get("key"))
				assert.Equal(t, tc.name, form.Get("name"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}
				return resp, nil
			}

			if resp, err := s.Create(tc.key, tc.name); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_Create_option(t *testing.T) {
	type options struct {
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                string
	}

	o := newProjectOptionService()
	cases := map[string]struct {
		options   []*FormOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*FormOption{},
			wantError: false,
			want: options{
				chartEnabled:                      "",
				subtaskingEnabled:                 "",
				projectLeaderCanEditProjectLeader: "",
				textFormattingRule:                "",
			},
		},
		"ValidOption": {
			options: []*FormOption{
				o.WithFormChartEnabled(true),
				o.WithFormSubtaskingEnabled(true),
				o.WithFormProjectLeaderCanEditProjectLeader(true),
				o.WithFormTextFormattingRule(FormatBacklog),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "true",
				subtaskingEnabled:                 "true",
				projectLeaderCanEditProjectLeader: "true",
				textFormattingRule:                "backlog",
			},
		},
		"OptionError": {
			options: []*FormOption{
				o.WithFormTextFormattingRule("invalid"),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options:   []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, tc.want.chartEnabled, form.Get("chartEnabled"))
				assert.Equal(t, tc.want.subtaskingEnabled, form.Get("subtaskingEnabled"))
				assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, string(tc.want.textFormattingRule), form.Get("textFormattingRule"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}
				return resp, nil
			}

			if resp, err := s.Create("TEST", "test", tc.options...); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_Create_clientError(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
		return nil, errors.New("error")
	}

	resp, err := s.Create("TEST", "test")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Create_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	project, err := s.Create("TEST", "test")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectService_Update(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath      string
		projectKey string
	}{
		spath:      "projects/" + projectKey,
		projectKey: projectKey,
	}

	s := newProjectService()
	s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.NotNil(t, form)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.Update(projectKey)
	require.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, want.projectKey, project.ProjectKey)
}

func TestProjectService_Update_param(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		wantError      bool
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			wantError:      false,
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			wantError:      false,
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			wantError:      true,
		},
		"projectIDOrKey_zero": {
			projectIDOrKey: "0",
			wantError:      true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}
				return resp, nil
			}

			if resp, err := s.Update(tc.projectIDOrKey); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_Update_option(t *testing.T) {
	type options struct {
		key                               string
		name                              string
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                Format
		archived                          string
	}

	o := newProjectOptionService()
	cases := map[string]struct {
		options   []*FormOption
		wantError bool
		want      options
	}{
		"ValidOption": {
			options: []*FormOption{
				o.WithFormKey("TEST1"),
				o.WithFormName("test1"),
				o.WithFormChartEnabled(true),
				o.WithFormSubtaskingEnabled(true),
				o.WithFormProjectLeaderCanEditProjectLeader(true),
				o.WithFormTextFormattingRule(FormatBacklog),
				o.WithFormArchived(true),
			},
			wantError: false,
			want: options{
				key:                               "TEST1",
				name:                              "test1",
				chartEnabled:                      "true",
				subtaskingEnabled:                 "true",
				projectLeaderCanEditProjectLeader: "true",
				textFormattingRule:                FormatBacklog,
				archived:                          "true",
			},
		},
		"OptionError": {
			options: []*FormOption{
				o.WithFormKey(""),
				o.WithFormName(""),
				o.WithFormChartEnabled(false),
				o.WithFormSubtaskingEnabled(false),
				o.WithFormProjectLeaderCanEditProjectLeader(false),
				o.WithFormTextFormattingRule("invalid"),
				o.WithFormArchived(false),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options:   []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, tc.want.key, form.Get("key"))
				assert.Equal(t, tc.want.name, form.Get("name"))
				assert.Equal(t, tc.want.chartEnabled, form.Get("chartEnabled"))
				assert.Equal(t, tc.want.subtaskingEnabled, form.Get("subtaskingEnabled"))
				assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, form.Get("projectLeaderCanEditProjectLeader"))
				assert.Equal(t, string(tc.want.textFormattingRule), form.Get("textFormattingRule"))
				assert.Equal(t, tc.want.archived, form.Get("archived"))

				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}
				return resp, nil
			}

			if resp, err := s.Update("TEST", tc.options...); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_Update_clientError(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
		return nil, errors.New("error")
	}

	resp, err := s.Update("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Update_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Patch = func(spath string, form *FormParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	project, err := s.Update("TEST")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectService_Delete_param(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		wantError      bool
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			wantError:      false,
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			wantError:      false,
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			wantError:      true,
		},
		"projectIDOrKey_zero": {
			projectIDOrKey: "0",
			wantError:      true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := newProjectService()
			s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
				resp := &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
				}
				return resp, nil
			}

			if resp, err := s.Delete(tc.projectIDOrKey); tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
			}
		})

	}
}

func TestProjectService_Delete(t *testing.T) {
	t.Parallel()

	projectKey := "TEST"

	want := struct {
		spath string
		key   string
	}{
		spath: "projects/" + projectKey,
		key:   projectKey,
	}

	s := newProjectService()
	s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
		assert.Equal(t, want.spath, spath)
		assert.NotNil(t, form)
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}
		return resp, nil
	}

	project, err := s.Delete(projectKey)
	require.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
}

func TestProjectService_Delete_clientError(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
		return nil, errors.New("error")
	}

	resp, err := s.Delete("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Delete_invalidJson(t *testing.T) {
	t.Parallel()

	s := newProjectService()
	s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
		}
		return resp, nil
	}

	project, err := s.Delete("TEST")
	assert.Error(t, err)
	assert.Nil(t, project)
}
