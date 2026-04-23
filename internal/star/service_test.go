package star_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/star"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func newNoContentResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
	}
}

func TestStarService_Add(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts       []core.RequestOption
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErr    bool
	}{
		// --- Success cases ------------------------------------------------------------
		"success-with-issueID": {
			opts: []core.RequestOption{o.WithIssueID(1)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "stars", spath)
				assert.Equal(t, "1", form.Get("issueId"))
				return newNoContentResponse(), nil
			},
		},
		"success-with-commentID": {
			opts: []core.RequestOption{o.WithCommentID(5)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "stars", spath)
				assert.Equal(t, "5", form.Get("commentId"))
				return newNoContentResponse(), nil
			},
		},
		"success-with-wikiID": {
			opts: []core.RequestOption{o.WithWikiID(10)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "stars", spath)
				assert.Equal(t, "10", form.Get("wikiId"))
				return newNoContentResponse(), nil
			},
		},
		"success-with-pullRequestID": {
			opts: []core.RequestOption{o.WithPullRequestID(3)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "stars", spath)
				assert.Equal(t, "3", form.Get("pullRequestId"))
				return newNoContentResponse(), nil
			},
		},
		"success-with-pullRequestCommentID": {
			opts: []core.RequestOption{o.WithPullRequestCommentID(7)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "stars", spath)
				assert.Equal(t, "7", form.Get("pullRequestCommentId"))
				return newNoContentResponse(), nil
			},
		},

		// --- Error cases --------------------------------------------------------------
		"error-no-required-option": {
			wantErr: true,
		},
		"error-invalid-option-type": {
			opts:    []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErr: true,
		},
		"error-invalid-option-value": {
			opts:    []core.RequestOption{o.WithIssueID(0)},
			wantErr: true,
		},
		"error-client-network": {
			opts: []core.RequestOption{o.WithIssueID(1)},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}

			s := star.NewService(method)

			err := s.Add(context.Background(), tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}

func TestStarService_Remove(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Delete = func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		assert.Equal(t, "stars", spath)
		assert.Equal(t, "42", form.Get("id"))
		return newNoContentResponse(), nil
	}
	s := star.NewService(method)

	err := s.Remove(context.Background(), 42)
	require.NoError(t, err)
}

func TestStarService_Remove_invalidID(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	s := star.NewService(method)

	err := s.Remove(context.Background(), 0)
	require.Error(t, err)
}

func TestStarService_Remove_clientError(t *testing.T) {
	t.Parallel()

	method := mock.NewMethod(t)
	method.Delete = func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
		return nil, errors.New("network error")
	}
	s := star.NewService(method)

	err := s.Remove(context.Background(), 1)
	require.Error(t, err)
}

func TestStarService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"StarService.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			o := &core.OptionService{}
			s := star.NewService(m)
			s.Add(ctx, o.WithIssueID(1)) //nolint:errcheck
		}},
		{"StarService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := star.NewService(m)
			s.Remove(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
