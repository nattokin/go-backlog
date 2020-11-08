package backlog_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestWikiService_All(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	projectKey := "test"
	want := struct {
		spath          string
		projectIDOrKey string
		keyword        string
	}{
		spath:          "wikis",
		projectIDOrKey: projectKey,
		keyword:        "",
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, params.Get("projectIdOrKey"))
			assert.Equal(t, want.keyword, params.Get("keyword"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wikis, err := s.All(backlog.ProjectKey(projectKey))
	assert.Nil(t, wikis)
	assert.Error(t, err)
}

func TestWikiService_Search(t *testing.T) {
	projectID := 103
	keyword := "test"
	bj, err := os.Open("testdata/json/wiki_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath     string
		projectID int
		keyword   string
		idList    []int
		nameList  []string
	}{
		spath:     "wikis",
		projectID: projectID,
		keyword:   keyword,
		idList:    []int{112, 115},
		nameList:  []string{"test1", "test2"},
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, strconv.Itoa(want.projectID), params.Get("projectIdOrKey"))
			assert.Equal(t, want.keyword, params.Get("keyword"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wikis, err := s.Search(backlog.ProjectID(projectID), keyword)
	assert.NoError(t, err)
	count := len(wikis)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], wikis[i].ID)
		assert.Equal(t, want.nameList[i], wikis[i].Name)
	}
}
func TestWikiService_Search_param_error(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wikis, err := s.Search(backlog.ProjectID(0), "test")
	assert.Error(t, err)
	assert.Nil(t, wikis)
	wikis, err = s.Search(backlog.ProjectKey(""), "test")
	assert.Error(t, err)
	assert.Nil(t, wikis)
}

func TestWikiService_Search_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.All(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
}

func TestWikiService_Count(t *testing.T) {
	projectKey := "TEST"
	body := ioutil.NopCloser(strings.NewReader(`{"count":5}`))
	want := struct {
		spath      string
		projectKey string
		count      int
	}{
		spath:      "wikis/count",
		projectKey: projectKey,
		count:      5,
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectKey, params.Get("projectIdOrKey"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	count, err := s.Count(backlog.ProjectKey(projectKey))
	assert.NoError(t, err)
	assert.Equal(t, want.count, count)
}
func TestWikiService_Count_param_error(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	count, err := s.Count(backlog.ProjectID(0))
	assert.Error(t, err)
	assert.Zero(t, count)
	count, err = s.Count(backlog.ProjectKey(""))
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_Count_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	count, err := s.Count(backlog.ProjectKey("TEST"))
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_Count_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	count, err := s.Count(backlog.ProjectKey("TEST"))
	assert.Zero(t, count)
	assert.Error(t, err)
}

func TestWikiService_One(t *testing.T) {
	wikiID := 34
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath  string
		wikiID int
		name   string
	}{
		spath:  "wikis/" + strconv.Itoa(wikiID),
		wikiID: wikiID,
		name:   "Maximum Wiki Page",
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.One(wikiID)
	assert.NoError(t, err)
	assert.Equal(t, want.wikiID, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
}
func TestWikiService_One_param(t *testing.T) {
	cases := map[string]struct {
		wikiID    int
		wantError bool
	}{
		"valid": {
			wikiID:    1,
			wantError: false,
		},
		"invalid-1": {
			wikiID:    0,
			wantError: true,
		},
		"invalid-2": {
			wikiID:    -1,
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki_maximum.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{

				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			})

			if wiki, err := s.One(tc.wikiID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 34, wiki.ID)
			}
		})
	}
}

func TestWikiService_One_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	wiki, err := s.One(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_One_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.One(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Create(t *testing.T) {
	projectID := 56
	wikiID := 34
	name := "Minimum Wiki Page"
	content := "This is a minimal wiki page."
	bj, err := os.Open("testdata/json/wiki_minimum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		wikiID     int
		spath      string
		name       string
		content    string
		mailNotify string
	}{
		wikiID:     wikiID,
		spath:      "wikis",
		name:       name,
		content:    content,
		mailNotify: "true",
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			assert.Equal(t, want.name, params.Get("name"))
			assert.Equal(t, want.content, params.Get("content"))
			assert.Equal(t, want.mailNotify, params.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Create(projectID, name, content, s.Option.WithMailNotify(true))
	assert.NoError(t, err)
	assert.Equal(t, want.wikiID, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
	assert.Equal(t, want.content, wiki.Content)
}

func TestWikiService_Create_param(t *testing.T) {
	cases := map[string]struct {
		projectID  int
		name       string
		content    string
		mailNotify bool
		wantError  bool
	}{
		"no_error_1": {
			projectID: 1,
			name:      "Test",
			content:   "test",
			wantError: false,
		},
		"no_error_2": {
			projectID: 100,
			name:      "Test Name",
			content:   "test content",
			wantError: false,
		},
		"projectId_zero": {
			projectID: 0,
			name:      "Test",
			content:   "test",
			wantError: true,
		},
		"name_empty": {
			projectID: 1,
			name:      "",
			content:   "test",
			wantError: true,
		},
		"content_empty": {
			projectID: 1,
			name:      "Test",
			content:   "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki_minimum.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{

				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			})

			if _, err := s.Create(tc.projectID, tc.name, tc.content); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWikiService_Create_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Create(1, "name", "test")
	assert.Error(t, err)
}

func TestWikiService_Create_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Create(1, "name", "test")
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Create_option_error(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_minimum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	error_option := func(p *backlog.ExportRequestParams) error {
		return errors.New("error")
	}
	wiki, err := s.Create(1, "name", "content", error_option)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Update(t *testing.T) {
	id := 34
	name := "Maximum Wiki Page"
	content := "This is a muximal wiki page."
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		id         int
		spath      string
		name       string
		content    string
		mailNotify string
	}{
		id:         id,
		spath:      "wikis/" + strconv.Itoa(id),
		name:       name,
		content:    content,
		mailNotify: "true",
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			assert.Equal(t, want.name, params.Get("name"))
			assert.Equal(t, want.content, params.Get("content"))
			assert.Equal(t, want.mailNotify, params.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	option := s.Option
	wiki, err := s.Update(id, option.WithName(name), option.WithContent(content), option.WithMailNotify(true))
	assert.NoError(t, err)
	assert.Equal(t, want.id, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
	assert.Equal(t, want.content, wiki.Content)
}

func TestWikiService_Update_param(t *testing.T) {
	cases := map[string]struct {
		wikiID    int
		name      string
		content   string
		wantError bool
	}{
		"no_error_1": {
			wikiID:    1,
			name:      "Test",
			content:   "test",
			wantError: false,
		},
		"no_error_2": {
			wikiID:    100,
			name:      "Test Name",
			content:   "test content",
			wantError: false,
		},
		"wikiId_zero": {
			wikiID:    0,
			name:      "Test",
			content:   "test",
			wantError: true,
		},
		"name_empty": {
			wikiID:    1,
			name:      "",
			content:   "test",
			wantError: true,
		},
		"content_empty": {
			wikiID:    1,
			name:      "Test",
			content:   "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki_maximum.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{

				Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			})
			o := s.Option
			if _, err := s.Update(tc.wikiID, o.WithName(tc.name), o.WithContent(tc.content)); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWikiService_Update_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Update(1, s.Option.WithName("name"))
	assert.Error(t, err)
}

func TestWikiService_Update_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Update(1, s.Option.WithName("name"))
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Update_option_required(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Update(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Delete(t *testing.T) {
	id := 34
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		id         int
		spath      string
		mailNotify string
	}{
		id:         id,
		spath:      "wikis/" + strconv.Itoa(id),
		mailNotify: "true",
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, params)
			assert.Equal(t, want.mailNotify, params.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Delete(id, s.Option.WithMailNotify(true))
	assert.NoError(t, err)
	assert.Equal(t, want.id, wiki.ID)
}

func TestWikiService_Delete_param(t *testing.T) {
	cases := map[string]struct {
		wikiID     int
		mailNotify bool
		wantError  bool
	}{
		"mailNotify_false": {
			wikiID:    1,
			wantError: false,
		},
		"mailNotify_true": {
			wikiID:    100,
			wantError: false,
		},
		"wikiId_zero": {
			wikiID:    0,
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki_maximum.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{

				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return backlog.ExportNewResponse(resp), nil
				},
			})

			if _, err := s.Delete(tc.wikiID); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestWikiService_Delete_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Delete(1)
	assert.Error(t, err)
}

func TestWikiService_Delete_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	wiki, err := s.Delete(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Delete_option_error(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{

		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	})
	error_option := func(p *backlog.ExportRequestParams) error {
		return errors.New("error")
	}
	wiki, err := s.Delete(1, error_option)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}
