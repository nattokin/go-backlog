package core_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestNewClient_validation(t *testing.T) {
	cases := map[string]struct {
		baseURL   string
		token     string
		wantError bool
		errMsg    string
	}{
		"missing-token": {
			baseURL:   "https://example.com",
			token:     "",
			wantError: true,
			errMsg:    "missing token",
		},
		"missing-baseURL": {
			baseURL:   "",
			token:     "token",
			wantError: true,
			errMsg:    "missing baseURL",
		},

		"invalid-baseURL": {
			baseURL:   "://invalid-url",
			token:     "token",
			wantError: true,
		},
		"valid-input": {
			baseURL:   "https://example.com",
			token:     "token",
			wantError: false,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c, err := core.NewClient(tc.baseURL, tc.token)

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
			assert.Equal(t, tc.baseURL, c.BaseURL.String())
			assert.Equal(t, tc.token, c.Token)
		})
	}
}

func TestNewClient_initialization(t *testing.T) {
	baseURL := "https://example.com"
	token := "token"

	t.Run("with-Doer", func(t *testing.T) {
		t.Parallel()

		mockDoer := &mock.Doer{T: t,
			DoFunc: func(_ *http.Request) (*http.Response, error) { return nil, errors.New("mockDoer error") },
		}
		c, err := core.NewClient(baseURL, token, core.WithDoer(mockDoer))
		require.NoError(t, err)

		{
			req, _ := c.NewRequest(context.Background(), http.MethodGet, "test")
			res, err := c.Doer.Do(req)
			require.Error(t, err)
			assert.Nil(t, res)
			assert.Equal(t, "mockDoer error", err.Error())
		}

	})

	t.Run("without-Doer", func(t *testing.T) {
		t.Parallel()

		c, err := core.NewClient(baseURL, token)
		require.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, baseURL, c.BaseURL.String())
		assert.Equal(t, token, c.Token)
		assert.Equal(t, http.DefaultClient, c.Doer)
		assert.IsType(t, &core.DefaultWrapper{}, c.Wrapper)
		assert.IsType(t, &core.Method{}, c.Method)
	})
}

func TestClient_Do(t *testing.T) {
	user := &model.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example",
	}
	now := time.Now()

	wantWiki := model.Wiki{
		ID:          1,
		ProjectID:   1,
		Name:        "Home",
		Content:     "test",
		Tags:        []*model.Tag{},
		Attachments: []*model.Attachment{},
		SharedFiles: []*model.SharedFile{},
		Stars:       []*model.Star{},
		CreatedUser: user,
		Created:     now,
		UpdatedUser: user,
		Updated:     now,
	}

	wikiJSON, _ := json.Marshal(wantWiki)

	apiErrors := &core.APIResponseError{
		Errors: []*core.Error{
			{
				Message: "No project.",
				Code:    6,
			},
		},
	}
	errJSON, _ := json.Marshal(apiErrors)

	header := http.Header{}
	header.Set("Content-Type", "application/json;charset=utf-8")

	cases := map[string]struct {
		doFunc  func(*http.Request) (*http.Response, error)
		wantErr bool
	}{
		"success": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Header:     header,
					Body:       io.NopCloser(bytes.NewReader(wikiJSON)),
				}, nil
			},
		},
		"http-client-error": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("http client error")
			},
			wantErr: true,
		},
		"api-error-response": {
			doFunc: func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: 404,
					Header:     header,
					Body:       io.NopCloser(bytes.NewReader(errJSON)),
				}, nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := mock.NewClient(t, tc.doFunc)

			res, err := c.Do(
				context.Background(),
				http.MethodGet,
				"test",
			)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, res)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, res)

			var wiki model.Wiki
			require.NoError(t, json.NewDecoder(res.Body).Decode(&wiki))

			assert.Equal(t, wantWiki.ID, wiki.ID)
			assert.Equal(t, wantWiki.Name, wiki.Name)
			assert.Equal(t, wantWiki.CreatedUser.Name, wiki.CreatedUser.Name)
		})
	}
}

func TestClient_NewRequest(t *testing.T) {
	cases := map[string]struct {
		method    string
		spath     string
		header    http.Header
		body      io.Reader
		query     url.Values
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
		"method-put": {
			method:    http.MethodPut,
			spath:     "put",
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
			query:     url.Values{},
			wantError: false,
		},
	}

	for n, tc := range cases {
		t.Run(n, func(t *testing.T) {
			t.Parallel()

			c := mock.NewClient(t, nil)
			request, err := c.NewRequest(
				context.Background(),
				tc.method,
				tc.spath,
				core.WithHeader(tc.header),
				core.WithBody(tc.body),
				core.WithQuery(tc.query),
			)

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

func TestClient_method(t *testing.T) {
	cases := map[string]struct {
		call    func(c *core.Client) (*http.Response, error)
		check   func(t *testing.T, captured *mock.Capture)
		wantErr bool
	}{
		"GET": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Get(context.Background(), "/path1", nil)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "GET", captured.Method)
				assert.Equal(t, "/api/v2/path1", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Empty(t, captured.Body)
				assert.Empty(t, captured.Header.Get("Content-Type"))
			},
		},

		"DOWNLOAD": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Download(context.Background(), "/path-download", nil)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "GET", captured.Method)
				assert.Equal(t, "/api/v2/path-download", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.Body)
			},
		},

		"POST": {
			call: func(c *core.Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("k", "v")
				return c.Method.Post(context.Background(), "/path2", form)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "POST", captured.Method)
				assert.Equal(t, "/api/v2/path2", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Equal(t, "application/x-www-form-urlencoded", captured.Header.Get("Content-Type"))
				assert.Contains(t, string(captured.Body), "k=v")
			},
		},

		"PATCH": {
			call: func(c *core.Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("id", "123")
				return c.Method.Patch(context.Background(), "/path3", form)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "PATCH", captured.Method)
				assert.Equal(t, "/api/v2/path3", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Contains(t, string(captured.Body), "id=123")
			},
		},

		"PUT": {
			call: func(c *core.Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("content", "hello")
				return c.Method.Put(context.Background(), "/path4", form)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "PUT", captured.Method)
				assert.Equal(t, "/api/v2/path4", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Equal(t, "application/x-www-form-urlencoded", captured.Header.Get("Content-Type"))
				assert.Contains(t, string(captured.Body), "content=hello")
			},
		},

		"DELETE": {
			call: func(c *core.Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("id", "321")
				return c.Method.Delete(context.Background(), "/path5", form)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "DELETE", captured.Method)
				assert.Equal(t, "/api/v2/path5", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Contains(t, string(captured.Body), "id=321")
			},
		},

		"UPLOAD": {
			call: func(c *core.Client) (*http.Response, error) {
				buf := bytes.NewBufferString("dummyfiledata")
				return c.Method.Upload(context.Background(), "/upload-path", "file.txt", buf)
			},
			check: func(t *testing.T, captured *mock.Capture) {
				assert.Equal(t, "POST", captured.Method)
				assert.Equal(t, "/api/v2/upload-path", captured.URL.Path)
				assert.Equal(t, "Bearer token", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))

				ct := captured.Header.Get("Content-Type")
				assert.Contains(t, ct, "multipart/form-data")
				assert.Contains(t, ct, "boundary=")

				boundary := strings.Split(ct, "boundary=")[1]
				reader := multipart.NewReader(bytes.NewReader(captured.Body), boundary)

				part, err := reader.NextPart()
				require.NoError(t, err)

				assert.Equal(t, "file", part.FormName())
				assert.Equal(t, "file.txt", part.FileName())

				data, err := io.ReadAll(part)
				require.NoError(t, err)
				assert.Equal(t, "dummyfiledata", string(data))

				next, err := reader.NextPart()
				assert.Nil(t, next)
				assert.Equal(t, io.EOF, err)
			},
		},

		// エラーケース
		"GET newRequest error": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Get(context.Background(), "", url.Values{})
			},
			wantErr: true,
		},

		"POST newRequest error": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Post(context.Background(), "", nil)
			},
			wantErr: true,
		},

		"PATCH empty params": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Patch(context.Background(), "spath", nil)
			},
		},

		"PATCH newRequest error": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Patch(context.Background(), "", nil)
			},
			wantErr: true,
		},

		"PUT empty params": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Put(context.Background(), "spath", nil)
			},
		},

		"PUT newRequest error": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Put(context.Background(), "", nil)
			},
			wantErr: true,
		},

		"DELETE empty params": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Delete(context.Background(), "spath", nil)
			},
		},

		"DELETE newRequest error": {
			call: func(c *core.Client) (*http.Response, error) {
				return c.Method.Delete(context.Background(), "", nil)
			},
			wantErr: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, captured := mock.NewCaptureClient(t, "{}")

			resp, err := tc.call(c)

			if tc.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, resp)

			if tc.check != nil {
				tc.check(t, captured)
			}
		})
	}
}

func TestClient_methodUpload_errors(t *testing.T) {
	type testCase struct {
		spath    string
		fileName string
		fileData string
		setup    func(c *core.Client)
	}

	cases := map[string]testCase{
		"empty_fileName": {
			spath:    "spath",
			fileName: "",
			fileData: "testdata",
		},

		"empty_spath": {
			spath:    "",
			fileName: "filename",
			fileData: "dummy",
		},

		"createFormFile_error": {
			spath:    "spath",
			fileName: "filename",
			fileData: "dummy",
			setup: func(c *core.Client) {
				c.Wrapper = mock.Wrapper{CreateErr: errors.New("mock createFormFile error")}
			},
		},

		"close_error": {
			spath:    "spath",
			fileName: "filename",
			fileData: "dummy",
			setup: func(c *core.Client) {
				c.Wrapper = mock.Wrapper{CloseErr: errors.New("mock close error")}
			},
		},

		"copy_error": {
			spath:    "spath",
			fileName: "filename",
			fileData: "dummy",
			setup: func(c *core.Client) {
				c.Wrapper = mock.Wrapper{CopyErr: errors.New("mock copy error")}
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := mock.NewClient(t, nil)

			if tc.setup != nil {
				tc.setup(c)
			}

			f := io.NopCloser(bytes.NewBufferString(tc.fileData))

			resp, err := c.Method.Upload(context.Background(), tc.spath, tc.fileName, f)

			assert.Error(t, err)
			assert.Nil(t, resp)
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
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				StatusCode: tc.statusCode,
				Body:       tc.body,
			}

			// Use the exported function from backlog package
			r, err := core.CheckResponse(resp)

			// 1. Validate response pointer
			if tc.wantNilResponse {
				assert.Nil(t, r, "Response should be nil")
			} else {
				assert.NotNil(t, r, "Response should NOT be nil")
			}

			// 2. Validate error
			if tc.wantError {
				assert.Error(t, err)

				apiErr, ok := err.(*core.APIResponseError)
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

func TestDecodeResponse(t *testing.T) {
	type target struct {
		Name string `json:"name"`
	}

	cases := map[string]struct {
		body     string
		wantErr  bool
		wantName string
	}{
		"success": {
			body:     `{"name":"test"}`,
			wantName: "test",
		},
		"invalid-json": {
			body:    `{"name":`,
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				Body: io.NopCloser(strings.NewReader(tc.body)),
			}

			var v target
			err := core.DecodeResponse(resp, &v)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantName, v.Name)
		})
	}
}

func TestDownloadResponse(t *testing.T) {
	cases := map[string]struct {
		header http.Header

		wantFilename    string
		wantContentType string
	}{
		"filename-and-content-type": {
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename="file.png"`},
				"Content-Type":        []string{"image/png"},
			},
			wantFilename:    "file.png",
			wantContentType: "image/png",
		},
		"content-type-with-charset": {
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename="doc.txt"`},
				"Content-Type":        []string{"text/plain; charset=utf-8"},
			},
			wantFilename:    "doc.txt",
			wantContentType: "text/plain",
		},
		"missing-headers": {
			header: http.Header{},
		},
		"malformed-content-disposition": {
			header: http.Header{
				"Content-Disposition": []string{`not a valid header;;;`},
				"Content-Type":        []string{"image/png"},
			},
			wantContentType: "image/png",
		},
		"malformed-content-type": {
			header: http.Header{
				"Content-Disposition": []string{`attachment; filename="file.png"`},
				"Content-Type":        []string{`not a valid header;;;`},
			},
			wantFilename: "file.png",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				Header: tc.header,
				Body:   io.NopCloser(bytes.NewReader([]byte("data"))),
			}

			got, err := core.DownloadResponse(resp)

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantFilename, got.Filename)
			assert.Equal(t, tc.wantContentType, got.ContentType)
			require.NotNil(t, got.Body)

			data, err := io.ReadAll(got.Body)
			require.NoError(t, err)
			assert.Equal(t, "data", string(data))
		})
	}
}

func TestError_Error(t *testing.T) {
	e := &core.Error{
		Message:  "No project.",
		Code:     6,
		MoreInfo: "more info",
	}
	want := "Message:No project., Code:6, MoreInfo:more info"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_Error(t *testing.T) {
	e := &core.APIResponseError{
		StatusCode: 404,
		Errors: []*core.Error{
			{
				Message:  "1st error",
				Code:     5,
				MoreInfo: "more info 1",
			},
			{
				Message:  "2nd error",
				Code:     9,
				MoreInfo: "more info 2",
			},
		},
	}
	want := "Status Code:404\nMessage:1st error, Code:5, MoreInfo:more info 1\nMessage:2nd error, Code:9, MoreInfo:more info 2"

	assert.Equal(t, want, e.Error())
}

func TestAPIResponseError_errorsAs(t *testing.T) {
	resp := &http.Response{
		StatusCode: 404,
		Body:       nil,
	}
	_, err := core.CheckResponse(resp)
	require.Error(t, err)

	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.APIResponseError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, 404, target.StatusCode)
}

func TestInternalClientError_errorsAs(t *testing.T) {
	err := core.NewInternalClientError("missing token")
	wrapped := fmt.Errorf("wrap: %w", err)

	var target *core.InternalClientError
	assert.True(t, errors.As(wrapped, &target))
	assert.Equal(t, "missing token", target.Error())
}
