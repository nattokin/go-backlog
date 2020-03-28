package backlog_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestWikiService_All(t *testing.T) {
	projectIDOrKey := "test"
	want := struct {
		spath          string
		projectIDOrKey string
		keyword        string
	}{
		spath:          "wikis",
		projectIDOrKey: projectIDOrKey,
		keyword:        "",
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, params.Get("projectIDOrKey"))
			assert.Equal(t, want.keyword, params.Get("keyword"))
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	ws.All(projectIDOrKey)
}

func TestWikiService_Search(t *testing.T) {
	projectIDOrKey := "103"
	keyword := "test"
	bj, err := os.Open("testdata/json/wiki/get-wiki-page-list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath          string
		projectIDOrKey string
		keyword        string
		idList         []int
		nameList       []string
	}{
		spath:          "wikis",
		projectIDOrKey: projectIDOrKey,
		keyword:        keyword,
		idList:         []int{112, 115},
		nameList:       []string{"test1", "test2"},
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, params.Get("projectIDOrKey"))
			assert.Equal(t, want.keyword, params.Get("keyword"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	wikis, err := ws.Search(projectIDOrKey, keyword)
	assert.Nil(t, err)
	count := len(wikis)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], wikis[i].ID)
		assert.Equal(t, want.nameList[i], wikis[i].Name)
	}
}

func TestWikiService_Search_clientError(t *testing.T) {
	projectIDOrKey := "test"
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.All(projectIDOrKey)
	assert.Error(t, err)
}

func TestWikiService_Count(t *testing.T) {
	projectIDOrKey := "test"
	body := ioutil.NopCloser(strings.NewReader(`{"count":5}`))
	want := struct {
		spath          string
		projectIDOrKey string
		count          int
	}{
		spath:          "wikis/count",
		projectIDOrKey: projectIDOrKey,
		count:          5,
	}
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, params.Get("projectIDOrKey"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	count, err := ws.Count(projectIDOrKey)
	assert.Nil(t, err)
	assert.Equal(t, want.count, count)
}

func TestWikiService_Count_clientError(t *testing.T) {
	projectIDOrKey := "test"
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.Count(projectIDOrKey)
	assert.Error(t, err)
}

func TestWikiService_One(t *testing.T) {
	wikiID := 112
	bj, err := os.Open("testdata/json/wiki/get-wiki-page.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath string
		id    int
		name  string
	}{
		spath: "wikis/" + strconv.Itoa(wikiID),
		id:    wikiID,
		name:  "test1",
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
	ws := backlog.ExportNewWikiService(cm)
	wiki, err := ws.One(wikiID)
	assert.Nil(t, err)
	assert.Equal(t, want.id, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
}

func TestWikiService_One_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.One(1)
	assert.Error(t, err)
}

func TestWikiService_Create(t *testing.T) {
	projectIDOrKey := 1
	name := "Home"
	content := "test"
	mailNotify := true
	bj, err := os.Open("testdata/json/wiki/add-wiki-page.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		projectIDOrKey int
		spath          string
		name           string
		content        string
		mailNotify     string
	}{
		projectIDOrKey: projectIDOrKey,
		spath:          "wikis",
		name:           name,
		content:        content,
		mailNotify:     "true",
	}
	cm := &backlog.ExportClientMethod{
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
	}
	ws := backlog.ExportNewWikiService(cm)
	wiki, err := ws.Create(projectIDOrKey, name, content, mailNotify)
	assert.Nil(t, err)
	assert.Equal(t, want.projectIDOrKey, wiki.ID)
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
			projectID:  1,
			name:       "Test",
			content:    "test",
			mailNotify: false,
			wantError:  false,
		},
		"no_error_2": {
			projectID:  100,
			name:       "Test Name",
			content:    "test content",
			mailNotify: true,
			wantError:  false,
		},
		"projectId_zero": {
			projectID:  0,
			name:       "Test",
			content:    "test",
			mailNotify: false,
			wantError:  true,
		},
		"name_empty": {
			projectID:  1,
			name:       "",
			content:    "test",
			mailNotify: false,
			wantError:  true,
		},
		"content_empty": {
			projectID:  1,
			name:       "Test",
			content:    "",
			mailNotify: false,
			wantError:  true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki/add-wiki-page.json")
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
			ws := backlog.ExportNewWikiService(cm)

			if _, err := ws.Create(tc.projectID, tc.name, tc.content, tc.mailNotify); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWikiService_Create_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.Create(1, "name", "test", false)
	assert.Error(t, err)
}

func TestWikiService_Update(t *testing.T) {
	id := 1
	name := "Home"
	content := "test"
	mailNotify := true
	bj, err := os.Open("testdata/json/wiki/update-wiki-page.json")
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
	cm := &backlog.ExportClientMethod{
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
	}
	ws := backlog.ExportNewWikiService(cm)
	wiki, err := ws.Update(id, name, content, mailNotify)
	assert.Nil(t, err)
	assert.Equal(t, want.id, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
	assert.Equal(t, want.content, wiki.Content)
}

func TestWikiService_Update_param(t *testing.T) {
	cases := map[string]struct {
		wikiID     int
		name       string
		content    string
		mailNotify bool
		wantError  bool
	}{
		"no_error_1": {
			wikiID:     1,
			name:       "Test",
			content:    "test",
			mailNotify: false,
			wantError:  false,
		},
		"no_error_2": {
			wikiID:     100,
			name:       "Test Name",
			content:    "test content",
			mailNotify: true,
			wantError:  false,
		},
		"wikiId_zero": {
			wikiID:     0,
			name:       "Test",
			content:    "test",
			mailNotify: false,
			wantError:  true,
		},
		"name_empty": {
			wikiID:     1,
			name:       "",
			content:    "test",
			mailNotify: false,
			wantError:  false,
		},
		"content_empty": {
			wikiID:     1,
			name:       "Test",
			content:    "",
			mailNotify: false,
			wantError:  false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki/update-wiki-page.json")
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
			ws := backlog.ExportNewWikiService(cm)

			if _, err := ws.Update(tc.wikiID, tc.name, tc.content, tc.mailNotify); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWikiService_Update_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.Update(1, "name", "test", false)
	assert.Error(t, err)
}

func TestWikiService_Delete(t *testing.T) {
	id := 1
	mailNotify := true
	bj, err := os.Open("testdata/json/wiki/delete-wiki-page.json")
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
	cm := &backlog.ExportClientMethod{
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
	}
	ws := backlog.ExportNewWikiService(cm)
	wiki, err := ws.Delete(id, mailNotify)
	assert.Nil(t, err)
	assert.Equal(t, want.id, wiki.ID)
}

func TestWikiService_Delete_param(t *testing.T) {
	cases := map[string]struct {
		wikiID     int
		mailNotify bool
		wantError  bool
	}{
		"mailNotify_false": {
			wikiID:     1,
			mailNotify: false,
			wantError:  false,
		},
		"mailNotify_true": {
			wikiID:     100,
			mailNotify: true,
			wantError:  false,
		},
		"wikiId_zero": {
			wikiID:     0,
			mailNotify: false,
			wantError:  true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki/delete-wiki-page.json")
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
			ws := backlog.ExportNewWikiService(cm)

			if _, err := ws.Delete(tc.wikiID, tc.mailNotify); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

func TestWikiService_Delete_clientError(t *testing.T) {
	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			return nil, errors.New("error")
		},
	}
	ws := backlog.ExportNewWikiService(cm)
	_, err := ws.Delete(1, false)
	assert.Error(t, err)
}
