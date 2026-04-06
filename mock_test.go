package backlog

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// newMockResponse creates an *http.Response object with the given status code and body content.
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
func newClientMethodMock() *method {
	return &method{
		Get: func(spath string, query url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Post: func(spath string, form url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Patch: func(spath string, form url.Values) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Delete: func(spath string, form url.Values) (*http.Response, error) {
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

// newFailingCheckOption returns a RequestOption whose check function always fails.
func newFailingCheckOption(t apiParamOptionType) RequestOption {
	return &apiOption{
		t: t,
		checkFunc: func() error {
			return errors.New("check error")
		},
		setFunc: func(_ url.Values) error { return nil },
	}
}

// newFailingSetOption returns a RequestOption whose set function always fails.
func newFailingSetOption(t apiParamOptionType) RequestOption {
	return &apiOption{
		t:         t,
		checkFunc: func() error { return nil },
		setFunc: func(_ url.Values) error {
			return errors.New("set error")
		},
	}
}

// newInvalidTypeOption returns a RequestOption whose has invalid type.
func newInvalidTypeOption() RequestOption {
	return &apiOption{
		t:         "invalid",
		checkFunc: func() error { return nil },
		setFunc:   func(_ url.Values) error { return nil },
	}
}
