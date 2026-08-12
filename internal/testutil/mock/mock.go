package mock

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
)

// ──────────────────────────────────────────────────────────────
//  Doer mock
// ──────────────────────────────────────────────────────────────

type Doer struct {
	T      *testing.T
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *Doer) Do(req *http.Request) (*http.Response, error) {
	assert.NotNil(m.T, req)
	return m.DoFunc(req)
}

// ──────────────────────────────────────────────────────────────
//  Wrapper mock
// ──────────────────────────────────────────────────────────────

type Wrapper struct {
	CreateErr error
	CopyErr   error
	CloseErr  error
}

func (w Wrapper) NewMultipartWriter(_ io.Writer) core.MultipartWriter {
	return &multipartWriter{wrapper: w}
}

func (w Wrapper) Copy(_ io.Writer, _ io.Reader) error {
	return w.CopyErr
}

type multipartWriter struct {
	wrapper Wrapper
}

func (mw *multipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	if mw.wrapper.CreateErr != nil {
		return nil, mw.wrapper.CreateErr
	}
	return io.Discard, nil
}
func (mw *multipartWriter) FormDataContentType() string { return "mock/type" }
func (mw *multipartWriter) Close() error                { return mw.wrapper.CloseErr }

// ──────────────────────────────────────────────────────────────
//  RequestOption mock
// ──────────────────────────────────────────────────────────────

// NewFailingCheckOption returns a RequestOption whose check function always fails.
func NewFailingCheckOption(t core.APIParamOptionType) *core.APIParamOption {
	return &core.APIParamOption{
		Type: t,
		CheckFunc: func() *core.ValidationError {
			return core.NewValidationError("test", "check error")
		},
		SetFunc: func(_ url.Values) error { return nil },
	}
}

// NewFailingSetOption returns a RequestOption whose set function always fails.
func NewFailingSetOption(t core.APIParamOptionType) *core.APIParamOption {
	return &core.APIParamOption{
		Type:      t,
		CheckFunc: func() *core.ValidationError { return nil },
		SetFunc: func(_ url.Values) error {
			return errors.New("set error")
		},
	}
}

// NewInvalidTypeOption returns a RequestOption with an invalid type.
func NewInvalidTypeOption() *core.APIParamOption {
	return &core.APIParamOption{
		Type:      "invalid",
		CheckFunc: func() *core.ValidationError { return nil },
		SetFunc:   func(_ url.Values) error { return nil },
	}
}

// ──────────────────────────────────────────────────────────────
//  HTTP response helpers
// ──────────────────────────────────────────────────────────────

// NewResponse returns an HTTP 200 OK response with the given JSON string as body.
// It allocates a fresh reader on each call so the body can only be consumed once,
// matching the behaviour of a real HTTP response.
func NewResponse(json string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(strings.NewReader(json)),
	}
}

// NewCreatedResponse returns an HTTP 201 Created response with the given JSON string as body.
// It allocates a fresh reader on each call so the body can only be consumed once,
// matching the behaviour of a real HTTP response.
func NewCreatedResponse(json string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusCreated,
		Body:       io.NopCloser(strings.NewReader(json)),
	}
}

// NewNoContentResponse returns an HTTP 204 No Content response with an empty body.
func NewNoContentResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
	}
}

// NewBinaryResponse returns an HTTP 200 OK response simulating a binary file download.
// filename is used to construct the Content-Disposition header.
// contentType is set as the Content-Type header.
// body is the raw bytes of the file content.
func NewBinaryResponse(filename, contentType string, body []byte) *http.Response {
	header := http.Header{}
	header.Set("Content-Disposition", "attachment; filename="+filename)
	header.Set("Content-Type", contentType)
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     header,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

// NewErrorResponse returns an HTTP response with the given status code and JSON error body.
// It allocates a fresh reader on each call so the body can only be consumed once,
// matching the behaviour of a real HTTP response.
func NewErrorResponse(statusCode int, json string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(json)),
	}
}

// NewUnauthorizedResponse returns an HTTP 401 Unauthorized response with a Backlog
// authentication failure error body.
func NewUnauthorizedResponse() *http.Response {
	return NewErrorResponse(
		http.StatusUnauthorized,
		`{"errors":[{"message":"Authentication failure.","code":11,"moreInfo":""}]}`,
	)
}

// NewNotFoundResponse returns an HTTP 404 Not Found response with a generic Backlog
// not-found error body.
func NewNotFoundResponse() *http.Response {
	return NewErrorResponse(
		http.StatusNotFound,
		`{"errors":[{"message":"No such resource.","code":6,"moreInfo":""}]}`,
	)
}

// NewInternalServerErrorResponse returns an HTTP 500 Internal Server Error response
// with a generic Backlog internal server error body.
func NewInternalServerErrorResponse() *http.Response {
	return NewErrorResponse(
		http.StatusInternalServerError,
		`{"errors":[{"message":"Internal Server Error","code":1,"moreInfo":""}]}`,
	)
}

// ──────────────────────────────────────────────────────────────
//  DoFunc helpers
// ──────────────────────────────────────────────────────────────

// NewDoFunc returns a doFunc that always responds with HTTP 200 and the given JSON body.
func NewDoFunc(json string) func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewResponse(json), nil
	}
}

// NewCreatedDoFunc returns a doFunc that always responds with HTTP 201 and the given JSON body.
func NewCreatedDoFunc(json string) func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewCreatedResponse(json), nil
	}
}

// NewNoContentDoFunc returns a doFunc that always responds with HTTP 204 No Content.
func NewNoContentDoFunc() func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewNoContentResponse(), nil
	}
}

// NewBinaryDoFunc returns a doFunc that always responds with HTTP 200 and a binary
// file download response with the given filename, Content-Type, and body.
func NewBinaryDoFunc(filename, contentType string, body []byte) func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewBinaryResponse(filename, contentType, body), nil
	}
}

// NewUnauthorizedDoFunc returns a doFunc that always responds with HTTP 401 Unauthorized.
func NewUnauthorizedDoFunc() func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewUnauthorizedResponse(), nil
	}
}

// NewNotFoundDoFunc returns a doFunc that always responds with HTTP 404 Not Found.
func NewNotFoundDoFunc() func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewNotFoundResponse(), nil
	}
}

// NewInternalServerErrorDoFunc returns a doFunc that always responds with HTTP 500 Internal Server Error.
func NewInternalServerErrorDoFunc() func(*http.Request) (*http.Response, error) {
	return func(_ *http.Request) (*http.Response, error) {
		return NewInternalServerErrorResponse(), nil
	}
}

// NewUnexpectedDoFunc returns a doFunc for mock.Doer that fails the test if called.
// It mirrors NewUnexpectedGetFn and friends, but for the lower-level Doer used by
// root-package tests (which is not verb-specific, unlike core.Method).
func NewUnexpectedDoFunc(t *testing.T) func(*http.Request) (*http.Response, error) {
	t.Helper()
	return func(*http.Request) (*http.Response, error) {
		t.Helper()
		t.Error("Do must not be called")
		return nil, errors.New("unexpected call")
	}
}

// ──────────────────────────────────────────────────────────────
//  Client helpers
// ──────────────────────────────────────────────────────────────

const (
	testBaseURL = "https://example.com"
	testToken   = "token"
)

// NewClient creates a test Client with the given doFunc as its Doer.
// If doFunc is nil, the client responds to every request with HTTP 200 and an empty JSON object.
// baseURL and token are fixed to well-known test values.
// T is not set on the underlying Doer; use NewCaptureClient or construct Doer directly when
// assertion on the request is needed.
func NewClient(t *testing.T, doFunc func(*http.Request) (*http.Response, error)) *core.Client {
	t.Helper()

	if doFunc == nil {
		doFunc = func(_ *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
				Header:     make(http.Header),
			}, nil
		}
	}

	c, err := core.NewClient(testBaseURL, testToken, core.WithDoer(&Doer{T: t, DoFunc: doFunc}))
	require.NoError(t, err)

	return c
}

// Capture holds the details of the most recent HTTP request executed by the
// client returned from NewCaptureClient. It is used in tests to inspect the
// outgoing request and verify that the Client constructs it correctly.
//
// Captured fields:
//   - Method: HTTP method used (GET, POST, PATCH, PUT, etc.)
//   - URL:    Full request URL, including query parameters
//   - Header: All headers set on the request
//   - Body:   Raw request body bytes
type Capture struct {
	Method string
	URL    *url.URL
	Header http.Header
	Body   []byte
}

// NewCaptureClient creates a test Client whose Doer records each outgoing
// request into the returned *Capture and responds with HTTP 200 and the given
// responseJSON as the body. Use the *Capture to assert on the request that the
// Client built.
//
// Pass "{}" as responseJSON when the response content is not relevant to the
// test, making it explicit that the return value carries no meaning.
//
// Example:
//
//	client, capture := mock.NewCaptureClient(t, "{}")
//	_, _ = client.Method.Get(ctx, "/wikis", nil)
//	assert.Equal(t, "GET", capture.Method)
//	assert.Equal(t, "/api/v2/wikis", capture.URL.Path)
func NewCaptureClient(t *testing.T, responseJSON string) (*core.Client, *Capture) {
	t.Helper()

	captured := &Capture{}

	c := NewClient(t, func(req *http.Request) (*http.Response, error) {
		var bodyBytes []byte
		if req.Body != nil {
			bodyBytes, _ = io.ReadAll(req.Body)
		}

		captured.Method = req.Method
		captured.URL = req.URL
		captured.Header = req.Header
		captured.Body = bodyBytes

		return NewResponse(responseJSON), nil
	})

	return c, captured
}

// ──────────────────────────────────────────────────────────────
//  Method mock helpers
// ──────────────────────────────────────────────────────────────

// NewMethod returns a *core.Method with all fields initialized to their
// corresponding NewUnexpected*Fn(t) functions. Tests should replace only the
// fields they intend to exercise, so that any accidental call to an unintended
// HTTP method causes an immediate test failure instead of a nil-pointer panic.
func NewMethod(t *testing.T) *core.Method {
	t.Helper()
	return &core.Method{
		Get:      NewUnexpectedGetFn(t),
		Post:     NewUnexpectedPostFn(t),
		Patch:    NewUnexpectedPatchFn(t),
		Put:      NewUnexpectedPutFn(t),
		Delete:   NewUnexpectedDeleteFn(t),
		Upload:   NewUnexpectedUploadFn(t),
		Download: NewUnexpectedDownloadFn(t),
	}
}

// NewUnexpectedGetFn returns a mock function for http GET that fails if called.
func NewUnexpectedGetFn(t *testing.T) func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Get must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedPostFn returns a mock function for http POST that fails if called.
func NewUnexpectedPostFn(t *testing.T) func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Post must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedPatchFn returns a mock function for http PATCH that fails if called.
func NewUnexpectedPatchFn(t *testing.T) func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Patch must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedPutFn returns a mock function for http PUT that fails if called.
func NewUnexpectedPutFn(t *testing.T) func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Put must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedDeleteFn returns a mock function for http DELETE that fails if called.
func NewUnexpectedDeleteFn(t *testing.T) func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Delete must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedUploadFn returns a mock function for http Upload that fails if called.
func NewUnexpectedUploadFn(t *testing.T) func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath, fileName string, r io.Reader) (*http.Response, error) {
		t.Helper()
		t.Error("Upload must not be called")
		return nil, errors.New("unexpected call")
	}
}

// NewUnexpectedDownloadFn returns a mock function for Download that fails if called.
func NewUnexpectedDownloadFn(t *testing.T) func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
	t.Helper()
	return func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		t.Helper()
		t.Error("Download must not be called")
		return nil, errors.New("unexpected call")
	}
}
