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
	projectIDOrKey := "103"
	bj, err := os.Open("testdata/json/wiki_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	want := struct {
		spath          string
		projectIDOrKey string
		idList         []int
		nameList       []string
	}{
		spath:          "wikis",
		projectIDOrKey: projectIDOrKey,
		idList:         []int{112, 115},
		nameList:       []string{"test1", "test2"},
	}
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectIDOrKey, query.Get("projectIdOrKey"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wikis, err := s.All(projectIDOrKey)
	assert.NoError(t, err)
	count := len(wikis)
	assert.Equal(t, len(want.idList), count)
	for i := 0; i < count; i++ {
		assert.Equal(t, want.idList[i], wikis[i].ID)
		assert.Equal(t, want.nameList[i], wikis[i].Name)
	}
}

func TestWikiService_All_param(t *testing.T) {
	s := &backlog.WikiService{}
	o := s.Option
	cases := map[string]struct {
		projectIdOrKey string
		keywordOption  *backlog.QueryOption
		content        string
		mailNotify     bool
		wantError      bool
	}{
		"valid-1": {
			projectIdOrKey: "1",
			keywordOption:  o.WithQueryKeyword("test"),
			wantError:      false,
		},
		"valid-2": {
			projectIdOrKey: "TEST",
			keywordOption:  o.WithQueryKeyword(""),
			wantError:      false,
		},
		"invalid-ProjectID": {
			projectIdOrKey: "0",
			keywordOption:  o.WithQueryKeyword("test"),
			wantError:      true,
		},
		"invalid-ProjectKey": {
			projectIdOrKey: "",
			keywordOption:  o.WithQueryKeyword("test"),
			wantError:      true,
		},
		"invalid-option": {
			projectIdOrKey: "TEST",
			keywordOption: backlog.ExportNewQueryOption(0, func(p *backlog.QueryParams) error {
				return nil
			}),
			wantError: true,
		},
		"option-error": {
			projectIdOrKey: "TEST",
			keywordOption: backlog.ExportNewQueryOption(backlog.ExportQueryKeyword, func(p *backlog.QueryParams) error {
				return errors.New("error")
			}),
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			bj, err := os.Open("testdata/json/wiki_list.json")
			if err != nil {
				t.Fatal(err)
			}
			defer bj.Close()

			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})

			if wikis, err := s.All(tc.projectIdOrKey, tc.keywordOption); tc.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, wikis, 2)
			}
		})
	}
}

func TestWikiService_All_param_error(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_list.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wikis, err := s.All("0")
	assert.Error(t, err)
	assert.Nil(t, wikis)
	wikis, err = s.All("")
	assert.Error(t, err)
	assert.Nil(t, wikis)
}

func TestWikiService_All_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.All("TEST")
	assert.Error(t, err)
}

func TestWikiService_All_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	_, err = s.All("TEST")
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Equal(t, want.projectKey, query.Get("projectIdOrKey"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       body,
			}
			return resp, nil
		},
	})
	count, err := s.Count(projectKey)
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	count, err := s.Count("0")
	assert.Error(t, err)
	assert.Zero(t, count)
	count, err = s.Count("")
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_Count_clientError(t *testing.T) {
	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	count, err := s.Count("TEST")
	assert.Error(t, err)
	assert.Zero(t, count)
}

func TestWikiService_Count_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	count, err := s.Count("TEST")
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.One(wikiID)
	assert.NoError(t, err)
	assert.Equal(t, want.wikiID, wiki.ID)
	assert.Equal(t, want.name, wiki.Name)
}

func TestWikiService_One_param(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

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
			s := &backlog.WikiService{}
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
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
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	wiki, err := s.One(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_One_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
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
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			assert.Equal(t, want.name, form.Get("name"))
			assert.Equal(t, want.content, form.Get("content"))
			assert.Equal(t, want.mailNotify, form.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.Create(projectID, name, content, s.Option.WithFormMailNotify(true))
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
				Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
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
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Create(1, "name", "test")
	assert.Error(t, err)
}

func TestWikiService_Create_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.Create(1, "name", "test")
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Create_OptionError(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_minimum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	errorOption := backlog.ExportNewFormOption(backlog.ExportFormMailNotify, func(p *backlog.ExportRequestParams) error {
		return errors.New("error")
	})
	wiki, err := s.Create(1, "name", "content", errorOption)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Create_invalidOption(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_minimum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	invalidOption := backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
		return nil
	})
	wiki, err := s.Create(1, "name", "content", invalidOption)
	assert.Nil(t, wiki)
	assert.IsType(t, &backlog.InvalidFormOptionError{}, err)
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
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			assert.Equal(t, want.name, form.Get("name"))
			assert.Equal(t, want.content, form.Get("content"))
			assert.Equal(t, want.mailNotify, form.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	option := s.Option
	wiki, err := s.Update(id, option.WithFormName(name), option.WithFormContent(content), option.WithFormMailNotify(true))
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
				Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
				},
			})
			o := s.Option
			if _, err := s.Update(tc.wikiID, o.WithFormName(tc.name), o.WithFormContent(tc.content)); tc.wantError {
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
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Update(1, s.Option.WithFormName("name"))
	assert.Error(t, err)
}

func TestWikiService_Update_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.Update(1, s.Option.WithFormName("name"))
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
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.Update(1)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Update_invalidOption(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	invalidOption := backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
		return nil
	})
	wiki, err := s.Update(1, invalidOption)
	assert.Nil(t, wiki)
	assert.IsType(t, &backlog.InvalidFormOptionError{}, err)
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.NotNil(t, form)
			assert.Equal(t, want.mailNotify, form.Get("mailNotify"))

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	wiki, err := s.Delete(id, s.Option.WithFormMailNotify(true))
	assert.NoError(t, err)
	assert.Equal(t, want.id, wiki.ID)
}

func TestWikiService_Delete_param(t *testing.T) {
	cases := map[string]struct {
		wikiID     int
		mailNotify bool
		wantError  bool
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
				Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
					resp := &http.Response{
						StatusCode: http.StatusOK,
						Body:       bj,
					}
					return resp, nil
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("error")
		},
	})
	_, err := s.Delete(1)
	assert.Error(t, err)
}

func TestWikiService_Delete_invalidJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalid.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
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
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	errorOption := backlog.ExportNewFormOption(backlog.ExportFormMailNotify, func(p *backlog.ExportRequestParams) error {
		return errors.New("error")
	})
	wiki, err := s.Delete(1, errorOption)
	assert.Nil(t, wiki)
	assert.Error(t, err)
}

func TestWikiService_Delete_invalidOption(t *testing.T) {
	bj, err := os.Open("testdata/json/wiki_maximum.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	s := &backlog.WikiService{}
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.ExportRequestParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return resp, nil
		},
	})
	invalidOption := backlog.ExportNewFormOption(0, func(p *backlog.ExportRequestParams) error {
		return nil
	})
	wiki, err := s.Delete(1, invalidOption)
	assert.Nil(t, wiki)
	assert.IsType(t, &backlog.InvalidFormOptionError{}, err)
}
