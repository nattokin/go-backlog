package backlog

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
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			token:     "token123",
			wantError: true,
			errMsg:    "missing baseURL",
		},

		"invalid-baseURL": {
			baseURL:   "://invalid-url",
			token:     "token123",
			wantError: true,
		},
		"valid-input": {
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
			assert.Equal(t, tc.baseURL, c.baseURL.String())
			assert.Equal(t, tc.token, c.token)
		})
	}
}

func TestNewClient_initialization(t *testing.T) {
	baseURL := "https://example.com"
	token := "token123"

	t.Run("with-Doer", func(t *testing.T) {
		t.Parallel()

		mockDoer := &mockDoer{t: t,
			doFunc: func(_ *http.Request) (*http.Response, error) { return nil, errors.New("mockDoer error") },
		}
		c, err := NewClient(baseURL, token, WithDoer(mockDoer))
		require.NoError(t, err)

		{
			req, _ := c.newRequest(context.Background(), http.MethodGet, "test")
			res, err := c.doer.Do(req)
			require.Error(t, err)
			assert.Nil(t, res)
			assert.Equal(t, "mockDoer error", err.Error())
		}

	})

	t.Run("without-Doer", func(t *testing.T) {
		t.Parallel()

		c, err := NewClient(baseURL, token)
		require.NoError(t, err)
		assert.NotNil(t, c)

		assert.Equal(t, token, c.token)
		assert.Equal(t, http.DefaultClient, c.doer)
		assert.IsType(t, &defaultWrapper{}, c.wrapper)

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
		assert.NotNil(t, c.Wiki.Option.base)
		assert.NotNil(t, c.Wiki.Option.base)
		assert.Same(t, c.Wiki.Option.base, c.Project.Option.base)
		assert.Same(t, c.Wiki.Option.base, c.Project.Option.base)

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
	})
}

func TestClient_do(t *testing.T) {
	user := &User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    1,
		Lang:        "ja",
		MailAddress: "test@example",
	}
	now := time.Now()

	wantWiki := Wiki{
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

	wikiJSON, _ := json.Marshal(wantWiki)

	apiErrors := &APIResponseError{
		Errors: []*Error{
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

			c := newClientMock(t, "https://test.com", "test", &mockDoer{
				t:      t,
				doFunc: tc.doFunc,
			})

			res, err := c.do(
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

			var wiki Wiki
			require.NoError(t, json.NewDecoder(res.Body).Decode(&wiki))

			assert.Equal(t, wantWiki.ID, wiki.ID)
			assert.Equal(t, wantWiki.Name, wiki.Name)
			assert.Equal(t, wantWiki.CreatedUser.Name, wiki.CreatedUser.Name)
		})
	}
}

func TestClient_newRequest(t *testing.T) {
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

			c := newClientMock(t, "https://test.com", "test", nil)
			request, err := c.newRequest(
				context.Background(),
				tc.method,
				tc.spath,
				withHeader(tc.header),
				withBody(tc.body),
				withQuery(tc.query),
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
		call    func(c *Client) (*http.Response, error)
		check   func(t *testing.T, captured *httpCapture)
		wantErr bool
	}{
		"GET": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Get(context.Background(), "/path1", nil)
			},
			check: func(t *testing.T, captured *httpCapture) {
				assert.Equal(t, "GET", captured.Method)
				assert.Equal(t, "/api/v2/path1", captured.URL.Path)
				assert.Equal(t, "Bearer token123", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Empty(t, captured.Body)
				assert.Empty(t, captured.Header.Get("Content-Type"))
			},
		},

		"POST": {
			call: func(c *Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("k", "v")
				return c.method.Post(context.Background(), "/path2", form)
			},
			check: func(t *testing.T, captured *httpCapture) {
				assert.Equal(t, "POST", captured.Method)
				assert.Equal(t, "/api/v2/path2", captured.URL.Path)
				assert.Equal(t, "Bearer token123", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Equal(t, "application/x-www-form-urlencoded", captured.Header.Get("Content-Type"))
				assert.Contains(t, string(captured.Body), "k=v")
			},
		},

		"PATCH": {
			call: func(c *Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("id", "123")
				return c.method.Patch(context.Background(), "/path3", form)
			},
			check: func(t *testing.T, captured *httpCapture) {
				assert.Equal(t, "PATCH", captured.Method)
				assert.Equal(t, "/api/v2/path3", captured.URL.Path)
				assert.Equal(t, "Bearer token123", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Contains(t, string(captured.Body), "id=123")
			},
		},

		"DELETE": {
			call: func(c *Client) (*http.Response, error) {
				form := url.Values{}
				form.Add("id", "321")
				return c.method.Delete(context.Background(), "/path4", form)
			},
			check: func(t *testing.T, captured *httpCapture) {
				assert.Equal(t, "DELETE", captured.Method)
				assert.Equal(t, "/api/v2/path4", captured.URL.Path)
				assert.Equal(t, "Bearer token123", captured.Header.Get("Authorization"))
				assert.Empty(t, captured.URL.Query().Get("apiKey"))
				assert.Contains(t, string(captured.Body), "id=321")
			},
		},

		"UPLOAD": {
			call: func(c *Client) (*http.Response, error) {
				buf := bytes.NewBufferString("dummyfiledata")
				return c.method.Upload(context.Background(), "/upload-path", "file.txt", buf)
			},
			check: func(t *testing.T, captured *httpCapture) {
				assert.Equal(t, "POST", captured.Method)
				assert.Equal(t, "/api/v2/upload-path", captured.URL.Path)
				assert.Equal(t, "Bearer token123", captured.Header.Get("Authorization"))
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

		// エラーケースは変更なし
		"GET newRequest error": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Get(context.Background(), "", url.Values{})
			},
			wantErr: true,
		},

		"POST newRequest error": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Post(context.Background(), "", nil)
			},
			wantErr: true,
		},

		"PATCH empty params": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Patch(context.Background(), "spath", nil)
			},
		},

		"PATCH newRequest error": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Patch(context.Background(), "", nil)
			},
			wantErr: true,
		},

		"DELETE empty params": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Delete(context.Background(), "spath", nil)
			},
		},

		"DELETE newRequest error": {
			call: func(c *Client) (*http.Response, error) {
				return c.method.Delete(context.Background(), "", nil)
			},
			wantErr: true,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c, captured := makeClient(t)

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
		setup    func(c *Client)
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
			setup: func(c *Client) {
				c.wrapper = mockWrapper{createErr: errors.New("mock createFormFile error")}
			},
		},

		"close_error": {
			spath:    "spath",
			fileName: "filename",
			fileData: "dummy",
			setup: func(c *Client) {
				c.wrapper = mockWrapper{closeErr: errors.New("mock close error")}
			},
		},

		"copy_error": {
			spath:    "spath",
			fileName: "filename",
			fileData: "dummy",
			setup: func(c *Client) {
				c.wrapper = mockWrapper{copyErr: errors.New("mock copy error")}
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			c := newClientMock(t, "https://test.com", "test", nil)

			if tc.setup != nil {
				tc.setup(c)
			}

			f := io.NopCloser(bytes.NewBufferString(tc.fileData))

			resp, err := c.method.Upload(context.Background(), tc.spath, tc.fileName, f)

			assert.Error(t, err)
			assert.Nil(t, resp)
		})
	}
}

func Test_checkResponse(t *testing.T) {
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

// ──────────────────────────────────────────────────────────────
//  Test helper types and functions
// ──────────────────────────────────────────────────────────────

// httpCapture holds the details of the most recent HTTP request executed
// by a mock Doer during testing. It is used to inspect the outgoing request
// generated by the Client to verify correctness in unit tests.
//
// Captured fields:
//   - Method: HTTP method used (GET, POST, PATCH, etc.)
//   - URL: Full request URL, including query parameters
//   - Header: All headers set on the request
//   - Body: Raw request body data
type httpCapture struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   []byte
}

// makeClient creates and returns a test Client along with an httpCapture instance
// that records the details of each request made by the client.
//
// This helper is primarily used in unit tests to verify that the Client constructs
// correct HTTP requests without performing actual network calls.
//
// Example:
//
//	client, captured := makeClient(t)
//	_, _ = client.Wiki.All()
//	assert.Equal(t, "GET", captured.Method)
//	assert.Contains(t, captured.URL.Path, "/api/v2/wikis")
func makeClient(t *testing.T) (*Client, *httpCapture) {
	t.Helper()

	captured := &httpCapture{}

	c := newClientMock(t, "https://example.com", "token123", &mockDoer{
		t: t,
		doFunc: func(req *http.Request) (*http.Response, error) {
			var bodyBytes []byte
			if req.Body != nil {
				bodyBytes, _ = io.ReadAll(req.Body)
			}

			captured.Method = req.Method
			captured.URL = req.URL
			captured.Header = req.Header
			captured.Body = bodyBytes

			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(`{}`)),
			}, nil
		},
	})

	return c, captured
}
