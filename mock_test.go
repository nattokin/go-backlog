package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockResponse creates an *http.Response object with the given status code
// and body content. It is primarily used in unit tests to simulate API responses.
//
// Example:
//
//	resp := newMockResponse(http.StatusOK, `{"id": 1, "name": "Wiki"}`)
func newMockResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Status:     http.StatusText(statusCode),
		Body:       io.NopCloser(bytes.NewBufferString(body)),
		Header:     make(http.Header),
	}
}

type mockDoer struct {
	t      *testing.T
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	assert.NotNil(m.t, req)
	return m.doFunc(req)
}

// newClientMock creates and returns a test Client instance initialized with the given Doer.
// It simplifies unit testing by allowing complete control over HTTP request–response behavior.
//
// If doer is nil, a default mockDoer is automatically injected that returns
// an empty 200 OK response. This enables quick test setup without requiring
// a custom HTTP mock for every test case.
//
// This helper marks itself as a test helper using t.Helper(), so any failure
// reported inside it will correctly point to the caller's line in test output.
//
// Example:
//
//	client := newClientMock(t, "https://example.com", "dummy-token", &mockDoer{
//		t: t,
//		doFunc: func(req *http.Request) (*http.Response, error) {
//			return newMockResponse(http.StatusOK, `{"id":1}`), nil
//		},
//	})
func newClientMock(t *testing.T, baseURL, token string, doer Doer) *Client {
	t.Helper()

	if doer == nil {
		doer = &mockDoer{
			t: t,
			doFunc: func(_ *http.Request) (*http.Response, error) {
				return newMockResponse(http.StatusOK, `{}`), nil
			},
		}
	}

	c, err := NewClient(baseURL, token, doer)
	require.NoError(t, err)

	return c
}

// newClientMethodMock creates and returns a mock implementation of the `method` struct.
// Each API function (Get, Post, Patch, Delete) returns a default "not implemented" error.
// This helper is typically used in unit tests to validate service-layer behavior
// without performing real HTTP requests.
func newClientMethodMock() *method {
	return &method{
		Get: func(spath string, query *QueryParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Post: func(spath string, form *FormParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Patch: func(spath string, form *FormParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Delete: func(spath string, form *FormParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Upload: func(spath, fileName string, r io.Reader) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
	}
}

type mockWrapper struct {
	createErr error
	copyErr   error
	closeErr  error
}

func (w mockWrapper) NewMultipartWriter(_ io.Writer) multipartWriter {
	return &mockMultipartWriter{wrapper: w}
}

func (w mockWrapper) Copy(_ io.Writer, _ io.Reader) error {
	return w.copyErr
}

type mockMultipartWriter struct {
	wrapper mockWrapper
}

func (mw *mockMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	if mw.wrapper.createErr != nil {
		return nil, mw.wrapper.createErr
	}
	return io.Discard, nil
}
func (mw *mockMultipartWriter) FormDataContentType() string { return "mock/type" }
func (mw *mockMultipartWriter) Close() error                { return mw.wrapper.closeErr }

// newQueryOptionWithCheckError returns a QueryOption whose check function always fails.
func newQueryOptionWithCheckError(t queryType) *QueryOption {
	return &QueryOption{
		t: t,
		checkFunc: func() error {
			return errors.New("check error")
		},
		setFunc: func(_ *QueryParams) error { return nil },
	}
}

// newQueryOptionWithSetError returns a QueryOption whose set function always fails.
func newQueryOptionWithSetError(t queryType) *QueryOption {
	return &QueryOption{
		t:         t,
		checkFunc: func() error { return nil },
		setFunc: func(_ *QueryParams) error {
			return errors.New("set error")
		},
	}
}

// newFormOptionWithCheckError returns a FormOption whose check function always fails.
func newFormOptionWithCheckError(t formType) *FormOption {
	return &FormOption{
		t: t,
		checkFunc: func() error {
			return errors.New("check error")
		},
		setFunc: func(_ *FormParams) error { return nil },
	}
}

// newFormOptionWithSetError returns a FormOption whose set function always fails.
func newFormOptionWithSetError(t formType) *FormOption {
	return &FormOption{
		t:         t,
		checkFunc: func() error { return nil },
		setFunc: func(_ *FormParams) error {
			return errors.New("set error")
		},
	}
}
