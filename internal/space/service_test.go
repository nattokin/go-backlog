package space_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestSpaceService_One(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType  error
		wantSpaceKey string
		wantName     string
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.SpaceJSON))),
				}, nil
			},
			wantSpaceKey: "nulab",
			wantName:     "Nulab Inc.",
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := space.NewService(method)

			got, err := s.One(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantSpaceKey, got.SpaceKey)
			assert.Equal(t, tc.wantName, got.Name)
		})
	}
}

func TestSpaceService_DiskUsage(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType   error
		wantCapacity  int
		wantIssue     int
		wantDetailLen int
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/diskUsage", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.DiskUsageJSON))),
				}, nil
			},
			wantCapacity:  1073741824,
			wantIssue:     119511,
			wantDetailLen: 1,
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/diskUsage", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/diskUsage", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := space.NewService(method)

			got, err := s.DiskUsage(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantCapacity, got.Capacity)
			assert.Equal(t, tc.wantIssue, got.Issue)
			assert.Len(t, got.Details, tc.wantDetailLen)
		})
	}
}

func TestSpaceService_Notification(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantContent string
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.NotificationJSON))),
				}, nil
			},
			wantContent: "Backlog is a project management tool.",
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := space.NewService(method)

			got, err := s.Notification(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantContent, got.Content)
		})
	}
}

func TestSpaceService_UpdateNotification(t *testing.T) {
	cases := map[string]struct {
		content string

		mockPutFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantContent string
	}{
		"success": {
			content: "Backlog is a project management tool.",
			mockPutFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				assert.Equal(t, "Backlog is a project management tool.", form.Get("content"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Space.NotificationJSON))),
				}, nil
			},
			wantContent: "Backlog is a project management tool.",
		},
		"error-validation-content-empty": {
			content:     "",
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			content: "some content",
			mockPutFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				assert.Equal(t, "some content", form.Get("content"))
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			content: "some content",
			mockPutFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "space/notification", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPutFn != nil {
				method.Put = tc.mockPutFn
			}

			s := space.NewService(method)

			got, err := s.UpdateNotification(context.Background(), tc.content)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantContent, got.Content)
		})
	}
}

// TestSpaceService_contextPropagation verifies that the context passed to each
// SpaceService method is correctly relayed to the underlying method call.
// A sentinel value is embedded in the context and its pointer identity is
// asserted inside the mock to catch any ctx substitution.
func TestSpaceService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeGetMock := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	makePutMock := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"One", func(t *testing.T, m *core.Method) {
			m.Get = makeGetMock(t)
			s := space.NewService(m)
			s.One(ctx) //nolint:errcheck
		}},
		{"DiskUsage", func(t *testing.T, m *core.Method) {
			m.Get = makeGetMock(t)
			s := space.NewService(m)
			s.DiskUsage(ctx) //nolint:errcheck
		}},
		{"Notification", func(t *testing.T, m *core.Method) {
			m.Get = makeGetMock(t)
			s := space.NewService(m)
			s.Notification(ctx) //nolint:errcheck
		}},
		{"UpdateNotification", func(t *testing.T, m *core.Method) {
			m.Put = makePutMock(t)
			s := space.NewService(m)
			s.UpdateNotification(ctx, "content") //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
