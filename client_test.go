package backlog_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewClientError(t *testing.T) {
	t.Parallel()

	msg := "test error message"
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
			t.Parallel()

			c, err := backlog.NewClient(tc.url, tc.token)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, c)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, c)
				assert.Equal(t, tc.url, c.ExportURL().String())
				assert.Equal(t, tc.token, c.ExportToken())
			}
		})

	}
}

func TestNewClient_project(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataProjectJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	key := "TEST"
	project, err := c.Project.Update(key, c.Project.Option.WithFormName("test"))
	assert.NoError(t, err)
	require.NotNil(t, project)
	assert.Equal(t, key, project.ProjectKey)
}

func TestNewClient_projectUser(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	projectKey := "TEST"
	userID := 1
	user, err := c.Project.User.Delete(projectKey, userID)
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
}

func TestNewClient_projectActivity(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataActivityListJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	projectKey := "SUB"
	activities, err := c.Project.Activity.List(projectKey, c.Project.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	require.NotNil(t, activities)
	assert.Equal(t, projectKey, activities[0].Project.ProjectKey)
}

func TestNewClient_spaceActivity(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataActivityListJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	projectKey := "SUB"
	activities, err := c.Space.Activity.List(c.Space.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	require.NotNil(t, activities)
	assert.Equal(t, projectKey, activities[0].Project.ProjectKey)
}

func TestNewClient_spaceAttachment(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataAttachmentUploadJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	fname := "test.txt"
	attachment, err := c.Space.Attachment.Upload(fname, f)
	assert.NoError(t, err)
	require.NotNil(t, attachment)
	assert.Equal(t, fname, attachment.Name)
}

func TestNewClient_user(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	userID := 1
	userName := "admin"
	user, err := c.User.Update(userID, c.User.Option.WithFormName(userName))
	assert.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, userName, user.Name)
}

func TestNewClient_userActivity(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataActivityListJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	userID := 1
	activities, err := c.User.Activity.List(userID, c.User.Activity.Option.WithQueryCount(1))
	assert.NoError(t, err)
	require.NotNil(t, activities)
	assert.Equal(t, userID, activities[0].CreatedUser.ID)
}

func TestNewClient_wiki(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       io.NopCloser(bytes.NewReader([]byte(testdataWikiMinimumJSON))),
		}

		return resp, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	projectID := 1
	name := "Minimum Wiki Page"
	content := "This is a minimal wiki page."
	wiki, err := c.Wiki.Create(projectID, name, content, c.Wiki.Option.WithFormMailNotify(false))
	assert.NoError(t, err)
	require.NotNil(t, wiki)
	assert.Equal(t, name, wiki.Name)
}

func TestClient_NewRequest(t *testing.T) {
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
		"method-error": {
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
			t.Parallel()

			c, _ := backlog.NewClient("https://test.backlog.com", "test")
			request, err := backlog.ExportClientNewRequest(c, tc.method, tc.spath, tc.header, tc.body, tc.query)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, request)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, request)
			}
		})

	}

}

func TestClient_Do(t *testing.T) {
	t.Parallel()

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

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
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

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	require.NoError(t, err)

	wiki := backlog.Wiki{}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&wiki))
	assert.Equal(t, want.ID, wiki.ID)
	assert.Equal(t, want.Name, wiki.Name)
	assert.Equal(t, want.CreatedUser.Name, wiki.CreatedUser.Name)
}

func TestClient_Do_httpClientError(t *testing.T) {
	t.Parallel()

	emsg := "http client error"
	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(emsg)
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	resp, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Nil(t, resp)
	assert.Error(t, err, emsg)

}

func TestClient_Do_errorResponse(t *testing.T) {
	t.Parallel()

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

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetHTTPClient(httpClient)

	resp, err := backlog.ExportClientDo(
		c, http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Nil(t, resp)
	assert.Error(t, err, apiErrors.Error())
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

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

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := backlog.NewClient(baseURL, apiKey)
	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientGet(c, spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Get_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	resp, err := backlog.ExportClientGet(c, "", backlog.NewQueryParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Post(t *testing.T) {
	t.Parallel()

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

	res, err := backlog.ExportClientPost(c, spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Post_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	resp, err := backlog.ExportClientPost(c, "", backlog.NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Patch(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	want := struct {
		method      string
		url         string
		contentType string
		body        string
	}{
		method:      http.MethodPatch,
		url:         baseURL + "/api/v2/" + spath + "?apiKey=" + apiKey,
		contentType: "application/x-www-form-urlencoded",
		body:        "key=value",
	}

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, want.body, string(content))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c.ExportSetHTTPClient(httpClient)

	form := backlog.NewFormParams()
	form.Set("key", "value")

	res, err := backlog.ExportClientPatch(c, spath, form)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_emptyParams(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.backlog.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := backlog.NewClient(baseURL, apiKey)

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientPatch(c, spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	resp, err := backlog.ExportClientPatch(c, "", backlog.NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()

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

	c, _ := backlog.NewClient(baseURL, apiKey)
	c.ExportSetHTTPClient(httpClient)

	form := backlog.NewFormParams()
	form.Set("key", "value")

	res, err := backlog.ExportClientDelete(c, spath, form)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_emptyParams(t *testing.T) {
	t.Parallel()

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := backlog.NewClient("https://test.backlog.com", "apikey")
	c.ExportSetHTTPClient(httpClient)

	res, err := backlog.ExportClientDelete(c, "spath", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")

	resp, err := backlog.ExportClientDelete(c, "", backlog.NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Upload(t *testing.T) {
	t.Parallel()

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

	httpClient := NewHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())

		contentType := req.Header.Get("Content-Type")
		assert.Regexp(t, "^multipart/form-data; boundary=", contentType)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := backlog.NewClient(baseURL, apiKey)
	c.ExportSetHTTPClient(httpClient)

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	res, err := backlog.ExportClientUpload(c, spath, "filename", f)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Upload_newRequestError(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	resp, err := backlog.ExportClientUpload(c, "", "filename", f)
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Upload_emptyFileName(t *testing.T) {
	t.Parallel()

	f, err := os.Open("testdata/testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	resp, err := backlog.ExportClientUpload(c, "spath", "", f)
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Upload_createFormFileError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetWrapper(&backlog.ExportWrapper{
		CreateFormFile: func(w *multipart.Writer, fname string) (io.Writer, error) {
			return nil, errors.New("")
		},
		Copy: backlog.ExportCopy,
	})

	f := io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON)))
	resp, err := backlog.ExportClientUpload(c, "spath", "filename", f)
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Upload_copyError(t *testing.T) {
	t.Parallel()

	c, _ := backlog.NewClient("https://test.backlog.com", "test")
	c.ExportSetWrapper(&backlog.ExportWrapper{
		CreateFormFile: backlog.ExportCreateFormFile,
		Copy: func(dst io.Writer, src io.Reader) error {
			return errors.New("")
		},
	})

	f := io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON)))
	resp, err := backlog.ExportClientUpload(c, "spath", "filename", f)
	require.Error(t, err)
	assert.Nil(t, resp)
}

func TestCheckResponse(t *testing.T) {
	const apiErrorBody = `{"errors":[{"message": "No project.", "code": 6, "moreInfo": ""}]}`
	const wantErrorStringFormat = "Status Code:%d\nMessage:No project., Code:6"

	cases := map[string]struct {
		statusCode        int
		body              io.ReadCloser
		wantNilResponse   bool
		wantErr           bool
		wantErrStatusCode int
		wantEmptyBodyTest bool // Indicates a test case with a nil response body
	}{
		"Status OK (200)": {
			statusCode:      http.StatusOK,
			body:            io.NopCloser(bytes.NewReader(nil)),
			wantNilResponse: false,
		},
		"Status Created (201)": {
			statusCode:      http.StatusCreated,
			body:            io.NopCloser(bytes.NewReader(nil)),
			wantNilResponse: false,
		},
		// Test 204 No Content handling: should return (nil, nil)
		"Status No Content (204) with nil body": {
			statusCode:      http.StatusNoContent,
			body:            io.NopCloser(bytes.NewReader(nil)),
			wantNilResponse: true,
		},
		// Test 204 No Content handling with body: should return (nil, nil)
		"Status No Content (204) with non-nil body": {
			statusCode:      http.StatusNoContent,
			body:            io.NopCloser(bytes.NewReader([]byte(`{"data":"ignored"}`))),
			wantNilResponse: true,
		},
		// Test 4xx/5xx error handling with valid body
		"Status Bad Request (400)": {
			statusCode:        http.StatusBadRequest,
			body:              io.NopCloser(bytes.NewReader([]byte(apiErrorBody))),
			wantNilResponse:   true,
			wantErr:           true,
			wantErrStatusCode: http.StatusBadRequest,
		},
		// Test 4xx/5xx error handling with nil body (to check defer r.Body.Close() and JSON unmarshal resilience)
		"Status Not Found (404) with nil body": {
			statusCode:        http.StatusNotFound,
			body:              nil,
			wantNilResponse:   true,
			wantErr:           true,
			wantErrStatusCode: http.StatusNotFound,
			wantEmptyBodyTest: true,
		},
		"Status Internal Server Error (500)": {
			statusCode:        http.StatusInternalServerError,
			body:              io.NopCloser(bytes.NewReader([]byte(apiErrorBody))),
			wantNilResponse:   true,
			wantErr:           true,
			wantErrStatusCode: http.StatusInternalServerError,
		},
		"Status Bad Request (400) with invalid JSON": {
			statusCode:        http.StatusBadRequest,
			body:              io.NopCloser(bytes.NewReader([]byte(`{"errors":[{"invalid json...`))),
			wantNilResponse:   true,
			wantErr:           true,
			wantErrStatusCode: http.StatusBadRequest,
			wantEmptyBodyTest: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				StatusCode: tc.statusCode,
				Body:       tc.body,
			}

			// Use the exported function from backlog package
			r, err := backlog.ExportCheckResponse(resp)

			// 1. Validate response pointer
			if tc.wantNilResponse {
				assert.Nil(t, r, "Response should be nil")
			} else {
				assert.NotNil(t, r, "Response should NOT be nil")
			}

			// 2. Validate error
			if tc.wantErr {
				assert.Error(t, err)

				apiErr, ok := err.(*backlog.APIResponseError)
				if assert.True(t, ok, "Error should be *backlog.APIResponseError") {
					// Validate StatusCode
					assert.Equal(t, tc.wantErrStatusCode, apiErr.StatusCode, "StatusCode mismatch")

					// Validate error message only if a body was provided
					if !tc.wantEmptyBodyTest {
						wantMsg := fmt.Sprintf(wantErrorStringFormat, tc.wantErrStatusCode)
						assert.Equal(t, wantMsg, apiErr.Error(), "Error message mismatch")
					}
				}
			} else {
				assert.NoError(t, err)
			}
		})

	}
}
