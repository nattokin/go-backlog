package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_All(t *testing.T) {
	bj, err := os.Open("testdata/json/project_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		idList   []int
		nameList []string
	}{
		idList:   []int{1, 2, 3},
		nameList: []string{"test", "test2", "test3"},
	}

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.All()
	assert.NoError(t, err)
	count := len(projects)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], projects[i].ID)
		assert.Equal(t, want.nameList[i], projects[i].Name)
	}
}

func TestProjectService_All_option(t *testing.T) {
	o := &backlog.ProjectOptionService{}
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*backlog.ProjectOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.ProjectOption{},
			wantError: false,
			want: options{
				all:      "",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*backlog.ProjectOption{
				o.WithAll(false),
				o.WithArchived(true),
			},
			wantError: false,
			want: options{
				all:      "false",
				archived: "true",
			},
		},
		"OptionError": {
			options: []*backlog.ProjectOption{
				backlog.ExportNewProjectOption(backlog.ExportOptionAll, func(p *backlog.ExportRequestParams) error {
					return errors.New("error")
				}),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.ProjectOption{
				o.WithName("test"),
			},
			wantError: true,
			want:      options{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/project_list.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.want.all, params.Get("all"))
					assert.Equal(t, tc.want.archived, params.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.All(tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_AdminAll(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.AdminAll()
	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_AdminAll_option(t *testing.T) {
	o := &backlog.ProjectOptionService{}
	type options struct {
		all      string
		archived string
	}
	cases := map[string]struct {
		options   []*backlog.ProjectOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.ProjectOption{},
			wantError: false,
			want: options{
				all:      "true",
				archived: "",
			},
		},
		"ValidOption": {
			options: []*backlog.ProjectOption{
				o.WithArchived(true),
			},
			wantError: false,
			want: options{
				all:      "true",
				archived: "true",
			},
		},
		"OptionError": {
			options: []*backlog.ProjectOption{
				backlog.ExportNewProjectOption(backlog.ExportOptionArchived, func(p *backlog.ExportRequestParams) error {
					return errors.New("error")
				}),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.ProjectOption{
				o.WithAll(true),
			},
			wantError: true,
			want:      options{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/project_list.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.want.all, params.Get("all"))
					assert.Equal(t, tc.want.archived, params.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.AdminAll(tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_AllArchived(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		archived: "true",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.AllArchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_AdminAllArchived(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "true",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.AdminAllArchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_AllUnarchived(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		archived: "false",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.AllUnarchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_AdminAllUnarchived(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath    string
		all      string
		archived string
	}{
		spath:    "projects",
		all:      "true",
		archived: "false",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	projects, err := s.AdminAllUnarchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_All_clientError(t *testing.T) {
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.All()
	assert.Error(t, err)
}

func TestProjectService_One_key(t *testing.T) {
	projectKey := "TEST"
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath string
		key   string
		name  string
	}{
		spath: "projects/" + projectKey,
		key:   projectKey,
		name:  "test",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.One(backlog.ProjectKey(projectKey))
	assert.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_id(t *testing.T) {
	projectID := 6
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath string
		id    int
		name  string
	}{
		spath: "projects/" + strconv.Itoa(projectID),
		id:    projectID,
		name:  "test",
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.One(backlog.ProjectID(projectID))
	assert.NoError(t, err)
	assert.Equal(t, want.id, project.ID)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_key_error(t *testing.T) {
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.One(backlog.ProjectKey(""))
	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_One_clientError(t *testing.T) {
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.One(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_One_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.One(backlog.ProjectKey("TEST"))

	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_Create(t *testing.T) {
	key := "TEST"
	name := "test"
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath string
		key   string
		name  string
	}{
		spath: "projects",
		key:   key,
		name:  name,
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			assert.Equal(t, want.key, params.Get("key"))
			assert.Equal(t, want.name, params.Get("name"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Create(key, name)
	assert.NoError(t, err)
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
			bj, err := os.Open("testdata/json/project.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.key, params.Get("key"))
					assert.Equal(t, tc.name, params.Get("name"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.Create(tc.key, tc.name); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_Create_option(t *testing.T) {
	o := &backlog.ProjectOptionService{}
	type options struct {
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                string
	}
	cases := map[string]struct {
		options   []*backlog.ProjectOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.ProjectOption{},
			wantError: false,
			want: options{
				chartEnabled:                      "",
				subtaskingEnabled:                 "",
				projectLeaderCanEditProjectLeader: "",
				textFormattingRule:                "",
			},
		},
		"ValidOption": {
			options: []*backlog.ProjectOption{
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule(backlog.FormatBacklog),
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
			options: []*backlog.ProjectOption{
				o.WithTextFormattingRule("invalid"),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.ProjectOption{
				backlog.ExportNewProjectOption(backlog.ExportOptionType(0), func(p *backlog.ExportRequestParams) error {
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
			bj, err := os.Open("testdata/json/project.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.want.chartEnabled, params.Get("chartEnabled"))
					assert.Equal(t, tc.want.subtaskingEnabled, params.Get("subtaskingEnabled"))
					assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, params.Get("projectLeaderCanEditProjectLeader"))
					assert.Equal(t, string(tc.want.textFormattingRule), params.Get("textFormattingRule"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.Create("TEST", "test", tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_Create_clientError(t *testing.T) {
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Create("TEST", "test")
	assert.Error(t, err)
}

func TestProjectService_Create_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Create("TEST", "test")

	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_Update(t *testing.T) {
	projectKey := "TEST"
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath      string
		projectKey string
	}{
		spath:      "projects/" + projectKey,
		projectKey: projectKey,
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Update(backlog.ProjectKey((projectKey)))
	assert.NoError(t, err)
	assert.Equal(t, want.projectKey, project.ProjectKey)
}

func TestProjectService_Update_param(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey backlog.ProjectIDOrKeyGetter
		wantError      bool
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			wantError:      false,
		},
		"projectIDOrKey_number": {
			projectIDOrKey: backlog.ProjectID(1234),
			wantError:      false,
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: backlog.ProjectKey(""),
			wantError:      true,
		},
		"projectIDOrKey_zero": {
			projectIDOrKey: backlog.ProjectID(0),
			wantError:      true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/project.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.Update(tc.projectIDOrKey); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_Update_option(t *testing.T) {
	o := &backlog.ProjectOptionService{}
	type options struct {
		key                               string
		name                              string
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                backlog.ExportFormat
		archived                          string
	}
	cases := map[string]struct {
		options   []*backlog.ProjectOption
		wantError bool
		want      options
	}{
		"ValidOption": {
			options: []*backlog.ProjectOption{
				o.WithKey("TEST1"),
				o.WithName("test1"),
				o.WithChartEnabled(true),
				o.WithSubtaskingEnabled(true),
				o.WithProjectLeaderCanEditProjectLeader(true),
				o.WithTextFormattingRule(backlog.FormatBacklog),
				o.WithArchived(true),
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
			options: []*backlog.ProjectOption{
				o.WithKey(""),
				o.WithName(""),
				o.WithChartEnabled(false),
				o.WithSubtaskingEnabled(false),
				o.WithProjectLeaderCanEditProjectLeader(false),
				o.WithTextFormattingRule("invalid"),
				o.WithArchived(false),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.ProjectOption{
				backlog.ExportNewProjectOption(backlog.ExportOptionType(0), func(p *backlog.ExportRequestParams) error {
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
			bj, err := os.Open("testdata/json/project.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					assert.Equal(t, tc.want.key, params.Get("key"))
					assert.Equal(t, tc.want.name, params.Get("name"))
					assert.Equal(t, tc.want.chartEnabled, params.Get("chartEnabled"))
					assert.Equal(t, tc.want.subtaskingEnabled, params.Get("subtaskingEnabled"))
					assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, params.Get("projectLeaderCanEditProjectLeader"))
					assert.Equal(t, string(tc.want.textFormattingRule), params.Get("textFormattingRule"))
					assert.Equal(t, tc.want.archived, params.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.Update(backlog.ProjectKey("TEST"), tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_Update_clientError(t *testing.T) {
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Update(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_Update_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Update(backlog.ProjectKey("TEST"))

	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_Delete_param(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey backlog.ProjectIDOrKeyGetter
		wantError      bool
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			wantError:      false,
		},
		"projectIDOrKey_number": {
			projectIDOrKey: backlog.ProjectID(1234),
			wantError:      false,
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: backlog.ProjectKey(""),
			wantError:      true,
		},
		"projectIDOrKey_zero": {
			projectIDOrKey: backlog.ProjectID(0),
			wantError:      true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/project.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.ProjectService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if _, err := s.Delete(tc.projectIDOrKey); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestProjectService_Delete(t *testing.T) {
	projectKey := "TEST"
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath string
		key   string
	}{
		spath: "projects/" + projectKey,
		key:   projectKey,
	}
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Delete(backlog.ProjectKey(projectKey))
	assert.NoError(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
}

func TestProjectService_Delete_clientError(t *testing.T) {
	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Delete(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_Delete_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.ProjectService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	project, err := s.Delete(backlog.ProjectKey("TEST"))

	assert.Nil(t, project)
	assert.Error(t, err)
}
