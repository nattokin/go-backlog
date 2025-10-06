package backlog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestNewClientError(t *testing.T) {
	msg := "test error massage"
	e := backlog.ExportNewInternalClientError(msg)

	assert.Equal(t, msg, e.Error())
}

func TestNewClient(t *testing.T) {
	cases := map[string]struct {
		url       string
		token     string
		wantError bool
	}{
		"no-error": {
			url:       "https://test.backlog.com",
			token:     "test",
			wantError: false,
		},
		"url-token-empty": {
			url:       "",
			token:     "",
			wantError: true,
		},
		"url-empty": {
			url:       "",
			token:     "test",
			wantError: true,
		},
		"token-empty": {
			url:       "https://test.backlog.com",
			token:     "",
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			c, err := backlog.NewClient(tc.url, tc.token)

			switch {
			case tc.wantError:
				assert.Error(t, err)
				assert.Nil(t, c)
			case !tc.wantError:
				assert.NoError(t, err)
				assert.NotNil(t, c)
			}

			if c == nil {
				return
			}

			assert.Equal(t, tc.url, c.ExportURL().String())
			assert.Equal(t, tc.token, c.ExportToken())
		})
	}
}

func TestNewClient_project(t *testing.T) {
	key := "TEST"
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/project.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	project, err := c.Project.Update(key, c.Project.Option.WithFormName("test"))
	assert.NoError(t, err)
	assert.NotNil(t, project)
	assert.Equal(t, key, project.ProjectKey)
}

func TestNewClient_projectUser(t *testing.T) {
	projectKey := "TEST"
	userID := 1
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	user, err := c.Project.User.Delete(projectKey, userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
}

func TestNewClient_projectActivity(t *testing.T) {
	projectKey := "SUB"
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/activity_list.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	activities, err := c.Project.Activity.List(projectKey, c.Project.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	assert.NotNil(t, activities)
	assert.Equal(t, projectKey, activities[0].Project.ProjectKey)
}

func TestNewClient_spaceActivity(t *testing.T) {
	projectKey := "SUB"
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/activity_list.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	activities, err := c.Space.Activity.List(c.Space.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	assert.NotNil(t, activities)
	assert.Equal(t, projectKey, activities[0].Project.ProjectKey)
}

func TestNewClient_spaceAttachment(t *testing.T) {
	fname := "test.txt"
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/attachment_upload.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	attachment, err := c.Space.Attachment.Upload(fname, f)
	assert.NoError(t, err)
	assert.NotNil(t, attachment)
	assert.Equal(t, fname, attachment.Name)
}

func TestNewClient_user(t *testing.T) {
	userID := 1
	userName := "admin"
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	user, err := c.User.Update(userID, c.User.Option.WithFormName(userName))
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userName, user.Name)
}

func TestNewClient_userActivity(t *testing.T) {
	userID := 1
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/activity_list.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	activities, err := c.User.Activity.List(userID, c.User.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	assert.NotNil(t, activities)
	assert.Equal(t, userID, activities[0].CreatedUser.ID)
}

func TestNewClient_wiki(t *testing.T) {
	projectID := 1
	name := "Minimum Wiki Page"
	content := "This is a minimal wiki page."
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bj, err := os.Open("testdata/json/wiki_minimum.json")
	if err != nil {
		t.Fatal(err)
	}

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       bj,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	wiki, err := c.Wiki.Create(projectID, name, content, c.Wiki.Option.WithFormMailNotify(false))
	assert.NoError(t, err)
	assert.NotNil(t, wiki)
	assert.Equal(t, name, wiki.Name)
}

func TestClient_NewReqest(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	cases := map[string]struct {
		method    string
		spath     string
		header    http.Header
		body      io.Reader
		query     *backlog.QueryParams
		wantError bool
	}{
		"method-get": {
			method:    http.MethodGet,
			spath:     "get",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"method-post": {
			method:    http.MethodPost,
			spath:     "post",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"method-patch": {
			method:    http.MethodPatch,
			spath:     "patch",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"method-delete": {
			method:    http.MethodDelete,
			spath:     "delete",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"method-empty": {
			method:    "",
			spath:     "nothing",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"method-eroor": {
			method:    "@error",
			spath:     "nothing",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: true,
		},
		"spath-empty": {
			method:    http.MethodGet,
			spath:     "",
			header:    nil,
			body:      nil,
			query:     nil,
			wantError: true,
		},
		"header": {
			method:    http.MethodGet,
			spath:     "test",
			header:    http.Header{},
			body:      nil,
			query:     nil,
			wantError: false,
		},
		"body": {
			method:    http.MethodGet,
			spath:     "test",
			header:    nil,
			body:      bytes.NewReader([]byte("test")),
			query:     nil,
			wantError: false,
		},
		"query": {
			method:    http.MethodGet,
			spath:     "test",
			header:    nil,
			body:      nil,
			query:     backlog.NewQueryParams(),
			wantError: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			request, err := backlog.ExportClientNewReqest(c, tc.method, tc.spath, tc.header, tc.body, tc.query)

			switch {
			case tc.wantError:
				assert.Error(t, err)
				assert.Nil(t, request)
			case !tc.wantError:
				assert.NoError(t, err)
				assert.NotNil(t, request)
			}
		})
	}

}

func TestClient_Do(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	user := &backlog.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example",
	}
	now := time.Now()

	want := backlog.Wiki{
		ID:          1,
		ProjectID:   1,
		Name:        "Home",
		Content:     "test",
		Tags:        []*backlog.Tag{},
		Attachments: []*backlog.Attachment{},
		SharedFiles: []*backlog.SharedFile{},
		Stars:       []*backlog.Star{},
		CreatedUser: user,
		Created:     now,
		UpdatedUser: user,
		Updated:     now,
	}

	bs, _ := json.Marshal(want)
	body := io.NopCloser(bytes.NewReader(bs))

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.NoError(t, err)

	defer res.Body.Close()

	wiki := backlog.Wiki{}
	json.NewDecoder(res.Body).Decode(&wiki)

	assert.Equal(t, want.ID, wiki.ID)
	assert.Equal(t, want.Name, wiki.Name)
	assert.Equal(t, want.CreatedUser.Name, wiki.CreatedUser.Name)
}

func TestClient_Do_httpClientError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	emsg := "http client error"
	// If http.Client.Do return error, Shuld return error.
	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(emsg)
	})
	c.ExportSetHTTPClient(httpClient)

	_, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Error(t, err, emsg)

}

func TestClient_Do_errorResponse(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	apiErrors := &backlog.APIResponseError{
		Errors: []*backlog.Error{
			{
				Message:  "No project.",
				Code:     6,
				MoreInfo: "",
			},
		},
	}

	bs, _ := json.Marshal(apiErrors)
	body := io.NopCloser(bytes.NewReader(bs))

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: 404,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})
	c.ExportSetHTTPClient(httpClient)

	_, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Error(t, err, apiErrors.Error())
}

func TestClient_Get(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodGet,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientGet(c, spath, nil)
	statusCode := res.StatusCode
	assert.Equal(t, http.StatusOK, statusCode)
}

func TestClient_Get_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientGet(c, "", backlog.NewQueryParams())
	assert.Error(t, err)
}

func TestClient_Post(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodPost,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientPost(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Post_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientPost(c, "", backlog.NewFormParams())
	assert.Error(t, err)
}

func TestClient_Patch(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodPatch,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content := []byte{}
		req.Body.Read(content)
		assert.Empty(t, content)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	form := backlog.NewFormParams()
	form.Set("key", "value")

	res, _ := backlog.ExportClientPatch(c, spath, form)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_emptyParams(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientPatch(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientPatch(c, "", backlog.NewFormParams())
	assert.Error(t, err)
}

func TestClient_Delete(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method      string
		url         string
		contentType string
	}{
		method:      http.MethodDelete,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()

		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content := []byte{}
		req.Body.Read(content)
		assert.Empty(t, content)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	form := backlog.NewFormParams()
	form.Set("key", "value")

	res, _ := backlog.ExportClientDelete(c, spath, form)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_emptyParams(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	res, _ := backlog.ExportClientDelete(c, spath, nil)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	_, err := backlog.ExportClientDelete(c, "", backlog.NewFormParams())
	assert.Error(t, err)
}

func TestClient_Upload(t *testing.T) {
	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"
	want := struct {
		method string
		url    string
	}{
		method: http.MethodPost,
		url:    baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())

		contentType := req.Header.Get("Content-Type")
		assert.Regexp(t, "^multipart/form-data; boundary=", contentType)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})
	c.ExportSetHTTPClient(httpClient)

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res, err := backlog.ExportClientUpload(c, spath, "filename", f)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Upload_newRequestError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = backlog.ExportClientUpload(c, "", "filename", f)
	assert.Error(t, err)
}

func TestClient_Upload_emptyFileName(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = backlog.ExportClientUpload(c, "spath", "", f)
	assert.Error(t, err)
}

func TestClient_Upload_createFormFileError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetWrapper(&backlog.ExportWrapper{
		CreateFormFile: func(w *multipart.Writer, fname string) (io.Writer, error) {
			return nil, errors.New("")
		},
		Copy: backlog.ExportCopy,
	})

	f, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = backlog.ExportClientUpload(c, "spath", "filename", f)
	assert.Error(t, err)
}

func TestClient_Upload_copyError(t *testing.T) {
	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetWrapper(&backlog.ExportWrapper{
		CreateFormFile: backlog.ExportCreateFormFile,
		Copy: func(dst io.Writer, src io.Reader) error {
			return errors.New("")
		},
	})

	f, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	_, err = backlog.ExportClientUpload(c, "spath", "filename", f)
	assert.Error(t, err)
}

func TestCeckResponse(t *testing.T) {
	cases := map[string]struct {
		statusCode int
		wantError  bool
	}{
		"199": {
			statusCode: 199,
			wantError:  true,
		},
		"200": {
			statusCode: 200,
			wantError:  false,
		},
		"299": {
			statusCode: 299,
			wantError:  false,
		},
		"300": {
			statusCode: 300,
			wantError:  true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			body := io.NopCloser(bytes.NewReader([]byte(`{"errors":[{"message": "No project.","code": 6,"moreInfo": ""}]}`)))

			resp := &http.Response{
				StatusCode: tc.statusCode,
				Body:       body,
			}

			if r, err := backlog.ExportCeckResponse(resp); tc.wantError {
				assert.NotNil(t, err)
			} else {
				assert.Equal(t, resp, r)
			}
		})
	}
}

func TestCeckResponse_emptyBody(t *testing.T) {
	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       nil,
	}
	_, err := backlog.ExportCeckResponse(resp)
	assert.Error(t, err)
}

func TestCeckResponse_StatusBadRequestWithInvalidJSON(t *testing.T) {
	body := io.NopCloser(bytes.NewReader([]byte(`{{invalid json}`)))

	resp := &http.Response{
		StatusCode: http.StatusBadRequest,
		Body:       body,
	}
	want := &backlog.APIResponseError{}

	_, err := backlog.ExportCeckResponse(resp)
	assert.IsType(t, want, err)
}
