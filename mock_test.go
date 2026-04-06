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

// ──────────────────────────────────────────────────────────────
//  Doer mock
// ──────────────────────────────────────────────────────────────

type mockDoer struct {
	t      *testing.T
	doFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	assert.NotNil(m.t, req)
	return m.doFunc(req)
}

// ──────────────────────────────────────────────────────────────
//  NewClient mock
// ──────────────────────────────────────────────────────────────

// newClientMock creates and returns a test Client instance initialized with the given Doer.
func newClientMock(t *testing.T, baseURL, token string, doer Doer) *Client {
	t.Helper()

	if doer == nil {
		doer = &mockDoer{
			t: t,
			doFunc: func(_ *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Status:     http.StatusText(http.StatusOK),
					Body:       io.NopCloser(bytes.NewBufferString(`{}`)),
					Header:     make(http.Header),
				}, nil
			},
		}
	}

	c, err := NewClient(baseURL, token, doer)
	require.NoError(t, err)

	return c
}

// ──────────────────────────────────────────────────────────────
//  RequestOption mock
// ──────────────────────────────────────────────────────────────

// newFailingCheckOption returns a RequestOption whose check function always fails.
func newFailingCheckOption(t apiParamOptionType) *apiParamOption {
	return &apiParamOption{
		t: t,
		checkFunc: func() error {
			return errors.New("check error")
		},
		setFunc: func(_ url.Values) error { return nil },
	}
}

// newFailingSetOption returns a RequestOption whose set function always fails.
func newFailingSetOption(t apiParamOptionType) *apiParamOption {
	return &apiParamOption{
		t:         t,
		checkFunc: func() error { return nil },
		setFunc: func(_ url.Values) error {
			return errors.New("set error")
		},
	}
}

// newInvalidTypeOption returns a RequestOption whose has invalid type.
func newInvalidTypeOption() *apiParamOption {
	return &apiParamOption{
		t:         "invalid",
		checkFunc: func() error { return nil },
		setFunc:   func(_ url.Values) error { return nil },
	}
}
