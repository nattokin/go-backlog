package backlog_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
			}
			return resp, nil
		},
	})

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
	option := backlog.ExportNewProjectOptionService()
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*backlog.QueryOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.QueryOption{},
			wantError: false,
			want: options{
				all:      "",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*backlog.QueryOption{
				option.WithQueryAll(false),
				option.WithQueryArchived(true),
			},
			wantError: false,
			want: options{
				all:      "false",
				archived: "true",
			},
		},
		"OptionError": {
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(backlog.ExportQueryAll, func(p *backlog.QueryParams) error {
					return errors.New("error")
				}),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(0, func(p *backlog.QueryParams) error {
					return nil
				}),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					assert.Equal(t, tc.want.all, query.Get("all"))
					assert.Equal(t, tc.want.archived, query.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
					}
					return resp, nil
				},
			})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, query.Get("all"))
			assert.Equal(t, want.archived, query.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	projects, err := s.AdminAll()
	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_AdminAll_option(t *testing.T) {
	option := backlog.ExportNewProjectOptionService()
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*backlog.QueryOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.QueryOption{},
			wantError: false,
			want: options{
				all:      "true",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*backlog.QueryOption{
				option.WithQueryArchived(true),
			},
			wantError: false,
			want: options{
				all:      "true",
				archived: "true",
			},
		},
		"OptionError": {
			options: []*backlog.QueryOption{
				backlog.ExportNewQueryOption(backlog.ExportQueryArchived, func(p *backlog.QueryParams) error {
					return errors.New("error")
				}),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.QueryOption{
				option.WithQueryAll(true),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					assert.Equal(t, tc.want.all, query.Get("all"))
					assert.Equal(t, tc.want.archived, query.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectListJSON))),
					}
					return resp, nil
				},
			})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.archived, query.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, query.Get("all"))
			assert.Equal(t, want.archived, query.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.archived, query.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, query.Get("all"))
			assert.Equal(t, want.archived, query.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	projects, err := s.AdminAllUnarchived()

	assert.Error(t, err)
	assert.Nil(t, projects)
}

func TestProjectService_All_clientError(t *testing.T) {
	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

	project, err := s.One(strconv.Itoa(projectID))
	require.NoError(t, err)
	assert.Equal(t, want.id, project.ID)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_key_error(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

	project, err := s.One("")
	assert.Error(t, err)
	assert.Nil(t, project)
}

func TestProjectService_One_clientError(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	resp, err := s.One("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_One_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			assert.Equal(t, want.key, form.Get("key"))
			assert.Equal(t, want.name, form.Get("name"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

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

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.key, form.Get("key"))
					assert.Equal(t, tc.name, form.Get("name"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
					}
					return resp, nil
				},
			})

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

	option := backlog.ExportNewProjectOptionService()
	cases := map[string]struct {
		options   []*backlog.FormOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.FormOption{},
			wantError: false,
			want: options{
				chartEnabled:                      "",
				subtaskingEnabled:                 "",
				projectLeaderCanEditProjectLeader: "",
				textFormattingRule:                "",
			},
		},
		"ValidOption": {
			options: []*backlog.FormOption{
				option.WithFormChartEnabled(true),
				option.WithFormSubtaskingEnabled(true),
				option.WithFormProjectLeaderCanEditProjectLeader(true),
				option.WithFormTextFormattingRule(backlog.FormatBacklog),
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
			options: []*backlog.FormOption{
				option.WithFormTextFormattingRule("invalid"),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.FormOption{
				backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
					return nil
				}),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.want.chartEnabled, form.Get("chartEnabled"))
					assert.Equal(t, tc.want.subtaskingEnabled, form.Get("subtaskingEnabled"))
					assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, form.Get("projectLeaderCanEditProjectLeader"))
					assert.Equal(t, string(tc.want.textFormattingRule), form.Get("textFormattingRule"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
					}
					return resp, nil
				},
			})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	resp, err := s.Create("TEST", "test")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Create_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

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

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
					}
					return resp, nil
				},
			})

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
		textFormattingRule                backlog.Format
		archived                          string
	}

	option := backlog.ExportNewProjectOptionService()
	cases := map[string]struct {
		options   []*backlog.FormOption
		wantError bool
		want      options
	}{
		"ValidOption": {
			options: []*backlog.FormOption{
				option.WithFormKey("TEST1"),
				option.WithFormName("test1"),
				option.WithFormChartEnabled(true),
				option.WithFormSubtaskingEnabled(true),
				option.WithFormProjectLeaderCanEditProjectLeader(true),
				option.WithFormTextFormattingRule(backlog.FormatBacklog),
				option.WithFormArchived(true),
			},
			wantError: false,
			want: options{
				key:                               "TEST1",
				name:                              "test1",
				chartEnabled:                      "true",
				subtaskingEnabled:                 "true",
				projectLeaderCanEditProjectLeader: "true",
				textFormattingRule:                backlog.FormatBacklog,
				archived:                          "true",
			},
		},
		"OptionError": {
			options: []*backlog.FormOption{
				option.WithFormKey(""),
				option.WithFormName(""),
				option.WithFormChartEnabled(false),
				option.WithFormSubtaskingEnabled(false),
				option.WithFormProjectLeaderCanEditProjectLeader(false),
				option.WithFormTextFormattingRule("invalid"),
				option.WithFormArchived(false),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.FormOption{
				backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
					return nil
				}),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
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
				},
			})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	resp, err := s.Update("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Update_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

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

			s := backlog.ExportNewProjectService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
					}
					return resp, nil
				},
			})

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

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
			}
			return resp, nil
		},
	})

	project, err := s.Delete(projectKey)
	require.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
}

func TestProjectService_Delete_clientError(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})

	resp, err := s.Delete("TEST")
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestProjectService_Delete_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewProjectService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	project, err := s.Delete("TEST")
	assert.Error(t, err)
	assert.Nil(t, project)
}
