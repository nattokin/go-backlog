package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type httpCapture struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   []byte
}

func makeClient(t *testing.T) (*Client, *httpCapture) {
	t.Helper()

	captured := &httpCapture{}

	mockTransport := roundTripFunc(func(req *http.Request) (*http.Response, error) {
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
		}
		captured.Method = req.Method
		captured.URL = req.URL
		captured.Header = req.Header
		captured.Body = bodyBytes
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(`{}`))}, nil
	})

	c, err := newClientMock("https://example.com", "token123", mockTransport)
	require.NoError(t, err)

	return c, captured
}

func TestNewClientError(t *testing.T) {
	t.Parallel()

	msg := "test error message"
	e := newInternalClientError(msg)

	assert.Equal(t, msg, e.Error())
}

func TestNewClient_InitAndValidation(t *testing.T) {
	cases := map[string]struct {
		baseURL   string
		token     string
		wantError bool
		errMsg    string
	}{
		"missing token": {
			baseURL:   "https://example.com",
			token:     "",
			wantError: true,
			errMsg:    "missing token",
		},
		"invalid baseURL": {
			baseURL:   "://invalid-url",
			token:     "token123",
			wantError: true,
		},
		"valid input": {
			baseURL:   "https://example.com",
			token:     "token123",
			wantError: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c, err := NewClient(tc.baseURL, tc.token)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, c)
				if tc.errMsg != "" {
					assert.Contains(t, err.Error(), tc.errMsg)
				}
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, c)
			assert.Equal(t, tc.token, c.token)
			assert.Equal(t, http.DefaultClient, c.httpClient)
			assert.IsType(t, &defaultWrapper{}, c.wrapper)
		})
	}
}

func TestNewClient_InitializationStructure(t *testing.T) {
	c, err := NewClient("https://example.com", "token123")
	assert.NoError(t, err)
	assert.NotNil(t, c)

	// Service initialization
	assert.NotNil(t, c.Wiki)
	assert.NotNil(t, c.Project)
	assert.NotNil(t, c.User)
	assert.NotNil(t, c.Issue)
	assert.NotNil(t, c.PullRequest)
	assert.NotNil(t, c.Space)

	// Shared method
	assert.NotNil(t, c.Wiki.method)
	assert.Same(t, c.Wiki.method, c.Project.method)
	assert.Same(t, c.Wiki.method, c.User.method)
	assert.Same(t, c.Wiki.method, c.Space.method)

	// Shared option support
	assert.NotNil(t, c.Wiki.Option.support.query)
	assert.NotNil(t, c.Wiki.Option.support.form)
	assert.Same(t, c.Wiki.Option.support.query, c.Project.Option.support.query)
	assert.Same(t, c.Wiki.Option.support.form, c.Project.Option.support.form)

	// Activity / Attachment presence
	assert.NotNil(t, c.Project.Activity)
	assert.NotNil(t, c.User.Activity)
	assert.NotNil(t, c.Space.Activity)
	assert.NotNil(t, c.Issue.Attachment)
	assert.NotNil(t, c.Wiki.Attachment)
	assert.NotNil(t, c.PullRequest.Attachment)

	// Reflection-based safety check
	clientType := reflect.TypeOf(*c)
	clientValue := reflect.ValueOf(*c)
	for i := 0; i < clientType.NumField(); i++ {
		field := clientType.Field(i)
		value := clientValue.Field(i)
		if field.Type.Kind() == reflect.Ptr && field.Name != "httpClient" {
			assert.Falsef(t, value.IsNil(), "%s should not be nil", field.Name)
		}
	}
}

func TestClient_HTTPMethods(t *testing.T) {
	t.Parallel()

	// --- GET ---
	{
		c, captured := makeClient(t)
		_, err := c.method.Get("/path1", nil)
		assert.NoError(t, err)
		assert.Equal(t, "GET", captured.Method)
		assert.Equal(t, "/api/v2/path1", captured.URL.Path)
		assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))
		assert.Empty(t, captured.Body)
		assert.Empty(t, captured.Header.Get("Content-Type"))
	}

	// --- POST ---
	{
		c, captured := makeClient(t)
		form := NewFormParams()
		form.Add("k", "v")
		_, err := c.method.Post("/path2", form)
		assert.NoError(t, err)
		assert.Equal(t, "POST", captured.Method)
		require.NoError(t, err)
		assert.Equal(t, "/api/v2/path2", captured.URL.Path)
		assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))
		assert.Equal(t, "application/x-www-form-urlencoded", captured.Header.Get("Content-Type"))
		assert.Contains(t, string(captured.Body), "k=v")
	}

	// --- PATCH ---
	{
		c, captured := makeClient(t)
		form := NewFormParams()
		form.Add("id", "123")
		_, err := c.method.Patch("/path3", form)
		assert.NoError(t, err)
		assert.Equal(t, "PATCH", captured.Method)
		require.NoError(t, err)
		assert.Equal(t, "/api/v2/path3", captured.URL.Path)
		assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))
		assert.Contains(t, string(captured.Body), "id=123")
	}

	// --- DELETE ---
	{
		c, captured := makeClient(t)
		form := NewFormParams()
		form.Add("id", "321")
		_, err := c.method.Delete("/path4", form)
		assert.NoError(t, err)
		assert.Equal(t, "DELETE", captured.Method)
		require.NoError(t, err)
		assert.Equal(t, "/api/v2/path4", captured.URL.Path)
		assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))
		assert.Contains(t, string(captured.Body), "id=321")
	}

	// --- UPLOAD ---
	{
		c, captured := makeClient(t)
		buf := bytes.NewBufferString("dummyfiledata")
		_, err := c.method.Upload("/path5", "file.txt", buf)
		assert.NoError(t, err)
		assert.Equal(t, "POST", captured.Method)
		assert.Contains(t, captured.URL.String(), "/path5")
		require.NoError(t, err)
		assert.Equal(t, "/api/v2/path5", captured.URL.Path)
		assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))
		assert.Contains(t, captured.Header.Get("Content-Type"), "multipart/form-data")
	}
}

func TestClient_HTTPUpload(t *testing.T) {
	t.Parallel()

	c, captured := makeClient(t)

	content := "dummyfiledata"
	fileName := "file.txt"
	reader := strings.NewReader(content)

	resp, err := c.method.Upload("/upload-path", fileName, reader)
	require.NoError(t, err)
	require.NotNil(t, resp)

	// Verify basic request info
	assert.Equal(t, "POST", captured.Method)
	require.NoError(t, err)
	assert.Equal(t, "/api/v2/upload-path", captured.URL.Path)
	assert.Equal(t, "token123", captured.URL.Query().Get("apiKey"))

	// Verify headers
	ct := captured.Header.Get("Content-Type")
	assert.Contains(t, ct, "multipart/form-data")
	assert.Contains(t, ct, "boundary=")

	// Verify multipart structure
	boundary := strings.Split(ct, "boundary=")[1]
	reader2 := multipart.NewReader(bytes.NewReader(captured.Body), boundary)

	part, err := reader2.NextPart()
	require.NoError(t, err)
	assert.Equal(t, "file", part.FormName())
	assert.Equal(t, "file.txt", part.FileName())

	data, err := io.ReadAll(part)
	require.NoError(t, err)
	assert.Equal(t, content, string(data))

	// Ensure multipart body is closed
	next, err := reader2.NextPart()
	assert.Nil(t, next)
	assert.Equal(t, io.EOF, err)
}

func TestClient_NewRequest(t *testing.T) {
	cases := map[string]struct {
		method    string
		spath     string
		header    http.Header
		body      io.Reader
		query     *QueryParams
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
			query:     NewQueryParams(),
			wantError: false,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			c, _ := NewClient("https://test.com", "test")
			request, err := c.newRequest(tc.method, tc.spath, tc.header, tc.body, tc.query)

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

	user := &User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example",
	}
	now := time.Now()

	want := Wiki{
		ID:          1,
		ProjectID:   1,
		Name:        "Home",
		Content:     "test",
		Tags:        []*Tag{},
		Attachments: []*Attachment{},
		SharedFiles: []*SharedFile{},
		Stars:       []*Star{},
		CreatedUser: user,
		Created:     now,
		UpdatedUser: user,
		Updated:     now,
	}

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")
	bs, _ := json.Marshal(want)
	body := io.NopCloser(bytes.NewReader(bs))

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: http.StatusOK,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})

	c, _ := NewClient("https://test.com", "test")
	c.httpClient = httpClient

	res, err := c.do(
		http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	require.NoError(t, err)

	wiki := Wiki{}
	require.NoError(t, json.NewDecoder(res.Body).Decode(&wiki))
	assert.Equal(t, want.ID, wiki.ID)
	assert.Equal(t, want.Name, wiki.Name)
	assert.Equal(t, want.CreatedUser.Name, wiki.CreatedUser.Name)
}

func TestClient_Do_httpClientError(t *testing.T) {
	t.Parallel()

	emsg := "http client error"
	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		return nil, errors.New(emsg)
	})

	c, _ := NewClient("https://test.com", "test")
	c.httpClient = httpClient

	resp, err := c.do(
		http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Nil(t, resp)
	assert.Error(t, err, emsg)

}

func TestClient_Do_errorResponse(t *testing.T) {
	t.Parallel()

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	apiErrors := &APIResponseError{
		Errors: []*Error{
			{
				Message:  "No project.",
				Code:     6,
				MoreInfo: "",
			},
		},
	}

	bs, _ := json.Marshal(apiErrors)
	body := io.NopCloser(bytes.NewReader(bs))

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		resp := &http.Response{
			StatusCode: 404,
			Header:     header,
			Body:       body,
		}

		return resp, nil
	})

	c, _ := NewClient("https://test.com", "test")
	c.httpClient = httpClient

	resp, err := c.do(
		http.MethodGet, "test",
		http.Header{}, nil, nil,
	)
	assert.Nil(t, resp)
	assert.Error(t, err, apiErrors.Error())
}

func TestClient_Get(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.com"
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

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := NewClient(baseURL, apiKey)
	c.httpClient = httpClient

	res, err := c.get(spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Get_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := NewClient("https://test.com", "test")

	resp, err := c.get("", NewQueryParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Post(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.com"
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

	c, _ := NewClient(baseURL, apiKey)

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c.httpClient = httpClient

	res, err := c.post(spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Post_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := NewClient("https://test.com", "test")

	resp, err := c.post("", NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Patch(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.com"
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

	c, _ := NewClient(baseURL, apiKey)

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, want.body, string(content))
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c.httpClient = httpClient

	form := NewFormParams()
	form.Set("key", "value")

	res, err := c.patch(spath, form)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_emptyParams(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.com"
	apiKey := "apikey"
	spath := "spath"

	c, _ := NewClient(baseURL, apiKey)

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c.httpClient = httpClient

	res, err := c.patch(spath, nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Patch_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := NewClient("https://test.com", "test")

	resp, err := c.patch("", NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Delete(t *testing.T) {
	t.Parallel()

	baseURL := "https://test.com"
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

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()

		assert.Equal(t, want.method, req.Method)
		assert.Equal(t, want.url, req.URL.String())
		assert.Equal(t, want.contentType, req.Header.Get("Content-Type"))
		content := []byte{}
		req.Body.Read(content)
		assert.Empty(t, content)
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := NewClient(baseURL, apiKey)
	c.httpClient = httpClient

	form := NewFormParams()
	form.Set("key", "value")

	res, err := c.delete(spath, form)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_emptyParams(t *testing.T) {
	t.Parallel()

	httpClient := newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
		defer req.Body.Close()
		return &http.Response{StatusCode: http.StatusOK}, nil
	})

	c, _ := NewClient("https://test.com", "apikey")
	c.httpClient = httpClient

	res, err := c.delete("spath", nil)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
}

func TestClient_Delete_newRequestError(t *testing.T) {
	t.Parallel()

	c, _ := NewClient("https://test.com", "test")

	resp, err := c.delete("", NewFormParams())
	assert.Error(t, err)
	assert.Nil(t, resp)
}

func TestClient_Upload(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		baseURL    string
		apiKey     string
		spath      string
		fileName   string
		fileData   string
		wantMethod string
		wantURL    string
		wantStatus int
		wantError  bool
		setup      func(c *Client)
	}

	cases := map[string]testCase{
		"success": {
			name:       "success",
			baseURL:    "https://test.com",
			apiKey:     "apikey",
			spath:      "spath",
			fileName:   "filename",
			fileData:   "testdata",
			wantMethod: http.MethodPost,
			wantURL:    "https://test.com/api/v2/spath?apiKey=apikey",
			wantStatus: http.StatusOK,
			wantError:  false,
			setup: func(c *Client) {
				c.httpClient = newHTTPClientMock(func(req *http.Request) (*http.Response, error) {
					assert.Equal(t, http.MethodPost, req.Method)
					assert.Equal(t, "https://test.com/api/v2/spath?apiKey=apikey", req.URL.String())
					assert.Regexp(t, "^multipart/form-data; boundary=", req.Header.Get("Content-Type"))
					return &http.Response{StatusCode: http.StatusOK}, nil
				})
			},
		},
		"empty fileName": {
			name:      "empty fileName",
			baseURL:   "https://test.com",
			apiKey:    "test",
			spath:     "spath",
			fileName:  "",
			fileData:  "testdata",
			wantError: true,
		},
		"empty spath": {
			name:      "empty spath",
			baseURL:   "https://test.com",
			apiKey:    "test",
			spath:     "",
			fileName:  "filename",
			fileData:  "dummy",
			wantError: true,
		},
		"createFormFile create error": {
			baseURL:   "https://test.com",
			apiKey:    "test",
			spath:     "spath",
			fileName:  "filename",
			fileData:  "dummy",
			wantError: true,
			setup: func(c *Client) {
				c.wrapper = mockWrapper{createErr: errors.New("mock createFormFile error")}
			},
		},
		"createFormFile close error": {
			name:      "createFormFile error",
			baseURL:   "https://test.com",
			apiKey:    "test",
			spath:     "spath",
			fileName:  "filename",
			fileData:  "invalid",
			wantError: true,
			setup: func(c *Client) {
				c.wrapper = mockWrapper{closeErr: errors.New("mock close error")}
			},
		},
		"copy error": {
			name:      "copy error",
			baseURL:   "https://test.com",
			apiKey:    "test",
			spath:     "spath",
			fileName:  "filename",
			fileData:  "invalid",
			wantError: true,
			setup: func(c *Client) {
				c.wrapper = mockWrapper{copyErr: errors.New("mock copy error")}
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, err := NewClient(tc.baseURL, tc.apiKey)
			require.NoError(t, err)

			if tc.setup != nil {
				tc.setup(c)
			}

			f := io.NopCloser(bytes.NewBufferString(tc.fileData))
			resp, err := c.upload(tc.spath, tc.fileName, f)

			if tc.wantError {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)
			assert.Equal(t, tc.wantStatus, resp.StatusCode)
		})
	}
}

func TestCheckResponse(t *testing.T) {
	const apiErrorBody = `{"errors":[{"message": "No project.", "code": 6, "moreInfo": ""}]}`
	const wantErrorStringFormat = "Status Code:%d\nMessage:No project., Code:6"

	cases := map[string]struct {
		statusCode        int
		body              io.ReadCloser
		wantNilResponse   bool
		wantError         bool
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
			wantError:         true,
			wantErrStatusCode: http.StatusBadRequest,
		},
		// Test 4xx/5xx error handling with nil body (to check defer r.Body.Close() and JSON unmarshal resilience)
		"Status Not Found (404) with nil body": {
			statusCode:        http.StatusNotFound,
			body:              nil,
			wantNilResponse:   true,
			wantError:         true,
			wantErrStatusCode: http.StatusNotFound,
			wantEmptyBodyTest: true,
		},
		"Status Internal Server Error (500)": {
			statusCode:        http.StatusInternalServerError,
			body:              io.NopCloser(bytes.NewReader([]byte(apiErrorBody))),
			wantNilResponse:   true,
			wantError:         true,
			wantErrStatusCode: http.StatusInternalServerError,
		},
		"Status Bad Request (400) with invalid JSON": {
			statusCode:        http.StatusBadRequest,
			body:              io.NopCloser(bytes.NewReader([]byte(`{"errors":[{"invalid json...`))),
			wantNilResponse:   true,
			wantError:         true,
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
			r, err := checkResponse(resp)

			// 1. Validate response pointer
			if tc.wantNilResponse {
				assert.Nil(t, r, "Response should be nil")
			} else {
				assert.NotNil(t, r, "Response should NOT be nil")
			}

			// 2. Validate error
			if tc.wantError {
				assert.Error(t, err)

				apiErr, ok := err.(*APIResponseError)
				if assert.True(t, ok, "Error should be *APIResponseError") {
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
