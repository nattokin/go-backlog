package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestProjectService_Joined(t *testing.T) {
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
		all:      "false",
		archived: "",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.Joined()
	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_All(t *testing.T) {
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
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	projects, err := s.All()
	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_Archived(t *testing.T) {
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
		all:      "false",
		archived: "true",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	projects, err := s.Archived()

	assert.Nil(t, projects)
	assert.Error(t, err)
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
		all:      "true",
		archived: "true",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	projects, err := s.AllArchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_Unarchived(t *testing.T) {
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
		all:      "false",
		archived: "false",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	projects, err := s.Unarchived()

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
		all:      "true",
		archived: "false",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.all, params.Get("all"))
			assert.Equal(t, want.archived, params.Get("archived"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	projects, err := s.AllUnarchived()

	assert.Nil(t, projects)
	assert.Error(t, err)
}

func TestProjectService_GetList(t *testing.T) {
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

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}

	s := backlog.ExportNewProjectService(cm)
	projects, err := s.Joined()
	assert.Nil(t, err)
	count := len(projects)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], projects[i].ID)
		assert.Equal(t, want.nameList[i], projects[i].Name)
	}
}

func TestProjectService_GetList_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewProjectService(cm)
	_, err := s.Joined()
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
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.One(backlog.ProjectKey(projectKey))
	assert.Nil(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_id(t *testing.T) {
	projectID := 1
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
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.One(backlog.ProjectID(projectID))
	assert.Nil(t, err)
	assert.Equal(t, want.id, project.ID)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_One_key_error(t *testing.T) {
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.One(backlog.ProjectKey(""))
	assert.Nil(t, project)
	assert.Error(t, err)
}

func TestProjectService_One_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewProjectService(cm)
	_, err := s.One(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_One_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
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
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			assert.Equal(t, want.key, params.Get("key"))
			assert.Equal(t, want.name, params.Get("name"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.Create(key, name)
	assert.Nil(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
	assert.Equal(t, want.name, project.Name)
}

func TestProjectService_Create_param(t *testing.T) {
	ops := &backlog.ProjectOptionService{}
	type options struct {
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                string
	}
	cases := map[string]struct {
		key       string
		name      string
		options   []backlog.ProjectOption
		wantError bool
		want      options
	}{
		"no-option": {
			key:       "TEST",
			name:      "test",
			options:   []backlog.ProjectOption{},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"key_empty": {
			key:       "",
			name:      "test",
			options:   []backlog.ProjectOption{},
			wantError: true,
			want:      options{},
		},
		"name_empty": {
			key:       "TEST",
			name:      "",
			options:   []backlog.ProjectOption{},
			wantError: true,
			want:      options{},
		},
		"option-chartEnabled_true": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(true),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "true",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-chartEnabled_false": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(false),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-subtaskingEnabled_true": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithSubtaskingEnabled(true),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "true",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-subtaskingEnabled_false": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithSubtaskingEnabled(false),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-projectLeaderCanEditProjectLeader_true": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithProjectLeaderCanEditProjectLeader(true),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "true",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-projectLeaderCanEditProjectLeader_false": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithProjectLeaderCanEditProjectLeader(false),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-textFormattingRule_backlog": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule(backlog.FormatBacklog),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatBacklog,
			},
		},
		"option-textFormattingRule_markdown": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule(backlog.FormatMarkdown),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"option-textFormattingRule_error": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule("error"),
			},
			wantError: true,
			want:      options{},
		},
		"multi-option-1": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(true),
				ops.WithSubtaskingEnabled(true),
				ops.WithProjectLeaderCanEditProjectLeader(true),
				ops.WithTextFormattingRule(backlog.FormatBacklog),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "true",
				subtaskingEnabled:                 "true",
				projectLeaderCanEditProjectLeader: "true",
				textFormattingRule:                backlog.FormatBacklog,
			},
		},
		"multi-option-2": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(false),
				ops.WithSubtaskingEnabled(false),
				ops.WithProjectLeaderCanEditProjectLeader(false),
				ops.WithTextFormattingRule(backlog.FormatMarkdown),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"invalid-option-key": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithKey("OPTION"),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"invalid-option-name": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithName("option"),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"invalid-option-key_empty": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithKey(""),
			},
			wantError: true,
			want:      options{},
		},
		"invalid-option-name_empty": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithName(""),
			},
			wantError: true,
			want:      options{},
		},
		"invalid-option-archived_true": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithArchived(true),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
		},
		"invalid-option-archived_false": {
			key:  "TEST",
			name: "test",
			options: []backlog.ProjectOption{
				ops.WithArchived(false),
			},
			wantError: false,
			want: options{
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
			},
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

			cm := &backlog.ExportClientMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					// Check options.
					assert.Equal(t, tc.want.chartEnabled, params.Get("chartEnabled"))
					assert.Equal(t, tc.want.subtaskingEnabled, params.Get("subtaskingEnabled"))
					assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, params.Get("projectLeaderCanEditProjectLeader"))
					assert.Equal(t, tc.want.textFormattingRule, params.Get("textFormattingRule"))

					// Check that invalid options are disabled.
					assert.Equal(t, tc.key, params.Get("key"))
					assert.Equal(t, tc.name, params.Get("name"))
					assert.Empty(t, params.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			}
			s := backlog.ExportNewProjectService(cm)

			if _, err := s.Create(tc.key, tc.name, tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestProjectService_Create_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewProjectService(cm)
	_, err := s.Create("TEST", "test")
	assert.Error(t, err)
}

func TestProjectService_Create_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
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
		spath string
		key   string
	}{
		spath: "projects/" + projectKey,
		key:   projectKey,
	}
	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.Update(backlog.ProjectKey((projectKey)))
	assert.Nil(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
}

func TestProjectService_Update_param(t *testing.T) {
	ops := &backlog.ProjectOptionService{}
	type options struct {
		key                               string
		name                              string
		chartEnabled                      string
		subtaskingEnabled                 string
		projectLeaderCanEditProjectLeader string
		textFormattingRule                string
		archived                          string
	}
	cases := map[string]struct {
		projectIDOrKey backlog.ProjectIDOrKeyGetter
		options        []backlog.ProjectOption
		wantError      bool
		want           options
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options:        []backlog.ProjectOption{},
			wantError:      false,
			want:           options{},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: backlog.ProjectID(1234),
			options:        []backlog.ProjectOption{},
			wantError:      false,
			want:           options{},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: backlog.ProjectKey(""),
			options:        []backlog.ProjectOption{},
			wantError:      true,
			want:           options{},
		},
		"projectIDOrKey_zero": {
			projectIDOrKey: backlog.ProjectID(0),
			options:        []backlog.ProjectOption{},
			wantError:      true,
			want:           options{},
		},
		"option-key": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithKey("TEST1"),
			},
			wantError: false,
			want: options{
				key: "TEST1",
			},
		},
		"option-name": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithName("test1"),
			},
			wantError: false,
			want: options{
				name: "test1",
			},
		},
		"option-key_empty": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithKey(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-name_empty": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithName(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-chartEnabled_true": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(true),
			},
			wantError: false,
			want: options{
				chartEnabled: "true",
			},
		},
		"option-chartEnabled_false": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithChartEnabled(false),
			},
			wantError: false,
			want: options{
				chartEnabled: "false",
			},
		},
		"option-subtaskingEnabled_true": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithSubtaskingEnabled(true),
			},
			wantError: false,
			want: options{
				subtaskingEnabled: "true",
			},
		},
		"option-subtaskingEnabled_false": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithSubtaskingEnabled(false),
			},
			wantError: false,
			want: options{
				subtaskingEnabled: "false",
			},
		},
		"option-projectLeaderCanEditProjectLeader_true": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithProjectLeaderCanEditProjectLeader(true),
			},
			wantError: false,
			want: options{
				projectLeaderCanEditProjectLeader: "true",
			},
		},
		"option-projectLeaderCanEditProjectLeader_false": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithProjectLeaderCanEditProjectLeader(false),
			},
			wantError: false,
			want: options{
				projectLeaderCanEditProjectLeader: "false",
			},
		},
		"option-textFormattingRule_backlog": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule(backlog.FormatBacklog),
			},
			wantError: false,
			want: options{
				textFormattingRule: backlog.FormatBacklog,
			},
		},
		"option-textFormattingRule_markdown": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule(backlog.FormatMarkdown),
			},
			wantError: false,
			want: options{
				textFormattingRule: backlog.FormatMarkdown,
			},
		},
		"option-textFormattingRule_error": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithTextFormattingRule("error"),
			},
			wantError: true,
			want:      options{},
		},
		"multi-option-1": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithKey("TEST1"),
				ops.WithName("test1"),
				ops.WithChartEnabled(true),
				ops.WithSubtaskingEnabled(true),
				ops.WithProjectLeaderCanEditProjectLeader(true),
				ops.WithTextFormattingRule(backlog.FormatBacklog),
				ops.WithArchived(true),
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
		"multi-option-2": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithKey("TEST2"),
				ops.WithName("test2"),
				ops.WithChartEnabled(false),
				ops.WithSubtaskingEnabled(false),
				ops.WithProjectLeaderCanEditProjectLeader(false),
				ops.WithTextFormattingRule(backlog.FormatMarkdown),
				ops.WithArchived(false),
			},
			wantError: false,
			want: options{
				key:                               "TEST2",
				name:                              "test2",
				chartEnabled:                      "false",
				subtaskingEnabled:                 "false",
				projectLeaderCanEditProjectLeader: "false",
				textFormattingRule:                backlog.FormatMarkdown,
				archived:                          "false",
			},
		},
		"option-archived_true": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithArchived(true),
			},
			wantError: false,
			want: options{
				archived: "true",
			},
		},
		"option-archived_false": {
			projectIDOrKey: backlog.ProjectKey("TEST"),
			options: []backlog.ProjectOption{
				ops.WithArchived(false),
			},
			wantError: false,
			want: options{
				archived: "false",
			},
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

			cm := &backlog.ExportClientMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					// Check options.
					assert.Equal(t, tc.want.key, params.Get("key"))
					assert.Equal(t, tc.want.name, params.Get("name"))
					assert.Equal(t, tc.want.chartEnabled, params.Get("chartEnabled"))
					assert.Equal(t, tc.want.subtaskingEnabled, params.Get("subtaskingEnabled"))
					assert.Equal(t, tc.want.projectLeaderCanEditProjectLeader, params.Get("projectLeaderCanEditProjectLeader"))
					assert.Equal(t, tc.want.textFormattingRule, params.Get("textFormattingRule"))
					assert.Equal(t, tc.want.archived, params.Get("archived"))

					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			}
			s := backlog.ExportNewProjectService(cm)

			if _, err := s.Update(tc.projectIDOrKey, tc.options...); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestProjectService_Update_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewProjectService(cm)
	_, err := s.Update(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_Update_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
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

			cm := &backlog.ExportClientMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			}
			s := backlog.ExportNewProjectService(cm)

			if _, err := s.Delete(tc.projectIDOrKey); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
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
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.Delete(backlog.ProjectKey(projectKey))
	assert.Nil(t, err)
	assert.Equal(t, want.key, project.ProjectKey)
}

func TestProjectService_Delete_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewProjectService(cm)
	_, err := s.Delete(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestProjectService_Delete_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectService(cm)
	project, err := s.Delete(backlog.ProjectKey("TEST"))

	assert.Nil(t, project)
	assert.Error(t, err)
}
