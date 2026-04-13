package mock

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
)

// ──────────────────────────────────────────────────────────────
//  Doer mock
// ──────────────────────────────────────────────────────────────

type MockDoer struct {
	T      *testing.T
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *MockDoer) Do(req *http.Request) (*http.Response, error) {
	assert.NotNil(m.T, req)
	return m.DoFunc(req)
}

// ──────────────────────────────────────────────────────────────
//  Wrapper mock
// ──────────────────────────────────────────────────────────────

type MockWrapper struct {
	CreateErr error
	CopyErr   error
	CloseErr  error
}

func (w MockWrapper) NewMultipartWriter(_ io.Writer) core.MultipartWriter {
	return &MockMultipartWriter{wrapper: w}
}

func (w MockWrapper) Copy(_ io.Writer, _ io.Reader) error {
	return w.CopyErr
}

type MockMultipartWriter struct {
	wrapper MockWrapper
}

func (mw *MockMultipartWriter) CreateFormFile(fieldname, filename string) (io.Writer, error) {
	if mw.wrapper.CreateErr != nil {
		return nil, mw.wrapper.CreateErr
	}
	return io.Discard, nil
}
func (mw *MockMultipartWriter) FormDataContentType() string { return "mock/type" }
func (mw *MockMultipartWriter) Close() error                { return mw.wrapper.CloseErr }

// ──────────────────────────────────────────────────────────────
//  RequestOption mock
// ──────────────────────────────────────────────────────────────

// NewFailingCheckOption returns a RequestOption whose check function always fails.
func NewFailingCheckOption(t core.APIParamOptionType) *core.APIParamOption {
	return &core.APIParamOption{
		Type: t,
		CheckFunc: func() error {
			return errors.New("check error")
		},
		SetFunc: func(_ url.Values) error { return nil },
	}
}

// NewFailingSetOption returns a RequestOption whose set function always fails.
func NewFailingSetOption(t core.APIParamOptionType) *core.APIParamOption {
	return &core.APIParamOption{
		Type:      t,
		CheckFunc: func() error { return nil },
		SetFunc: func(_ url.Values) error {
			return errors.New("set error")
		},
	}
}

// NewInvalidTypeOption returns a RequestOption whose has invalid type.
func NewInvalidTypeOption() *core.APIParamOption {
	return &core.APIParamOption{
		Type:      "invalid",
		CheckFunc: func() error { return nil },
		SetFunc:   func(_ url.Values) error { return nil },
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
