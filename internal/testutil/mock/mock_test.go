package mock_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

// recordingTB is a testing.TB that records Error calls instead of failing the test.
// Calling any other method of testing.TB panics, because the embedded TB is nil.
type recordingTB struct {
	testing.TB
	errors []string
}

func (r *recordingTB) Helper() {}

func (r *recordingTB) Error(args ...any) {
	r.errors = append(r.errors, fmt.Sprint(args...))
}

func TestNewUnexpectedFns(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		call    func(testing.TB) (*http.Response, error)
		wantMsg string
	}{
		"NewUnexpectedGetFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedGetFn(tb)(ctx, "", nil) },
			wantMsg: "Get must not be called",
		},
		"NewUnexpectedPostFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedPostFn(tb)(ctx, "", nil) },
			wantMsg: "Post must not be called",
		},
		"NewUnexpectedPatchFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedPatchFn(tb)(ctx, "", nil) },
			wantMsg: "Patch must not be called",
		},
		"NewUnexpectedPutFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedPutFn(tb)(ctx, "", nil) },
			wantMsg: "Put must not be called",
		},
		"NewUnexpectedDeleteFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedDeleteFn(tb)(ctx, "", nil) },
			wantMsg: "Delete must not be called",
		},
		"NewUnexpectedUploadFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedUploadFn(tb)(ctx, "", "", nil) },
			wantMsg: "Upload must not be called",
		},
		"NewUnexpectedDownloadFn": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedDownloadFn(tb)(ctx, "", nil) },
			wantMsg: "Download must not be called",
		},
		"NewUnexpectedDoFunc": {
			call:    func(tb testing.TB) (*http.Response, error) { return mock.NewUnexpectedDoFunc(tb)(nil) },
			wantMsg: "Do must not be called",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tb := &recordingTB{}
			resp, err := tc.call(tb)

			assert.Nil(t, resp)
			assert.EqualError(t, err, "unexpected call")
			assert.Equal(t, []string{tc.wantMsg}, tb.errors)
		})
	}
}

func TestNewMethod(t *testing.T) {
	ctx := context.Background()

	cases := map[string]struct {
		call    func(*client.Method) (*http.Response, error)
		wantMsg string
	}{
		"Get": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Get(ctx, "", nil) },
			wantMsg: "Get must not be called",
		},
		"Post": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Post(ctx, "", nil) },
			wantMsg: "Post must not be called",
		},
		"Patch": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Patch(ctx, "", nil) },
			wantMsg: "Patch must not be called",
		},
		"Put": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Put(ctx, "", nil) },
			wantMsg: "Put must not be called",
		},
		"Delete": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Delete(ctx, "", nil) },
			wantMsg: "Delete must not be called",
		},
		"Upload": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Upload(ctx, "", "", nil) },
			wantMsg: "Upload must not be called",
		},
		"Download": {
			call:    func(m *client.Method) (*http.Response, error) { return m.Download(ctx, "", nil) },
			wantMsg: "Download must not be called",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			tb := &recordingTB{}
			m := mock.NewMethod(tb)
			assert.Empty(t, tb.errors)

			resp, err := tc.call(m)

			assert.Nil(t, resp)
			assert.EqualError(t, err, "unexpected call")
			assert.Equal(t, []string{tc.wantMsg}, tb.errors)
		})
	}
}
