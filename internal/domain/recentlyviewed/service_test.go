package recentlyviewed_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/recentlyviewed"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_ListIssues(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts      []*option.APIParamOption
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr   bool
		wantLen   int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself/recentlyViewedIssues", spath)
				return mock.NewResponse(`[{"id":1},{"id":2}]`), nil
			},
			wantLen: 2,
		},
		"success-with-count": {
			opts: []*option.APIParamOption{o.WithCount(5)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself/recentlyViewedIssues", spath)
				assert.Equal(t, "5", query.Get("count"))
				return mock.NewResponse(`[{"id":1}]`), nil
			},
			wantLen: 1,
		},
		"success-with-offset": {
			opts: []*option.APIParamOption{o.WithOffset(10)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "10", query.Get("offset"))
				return mock.NewResponse(`[]`), nil
			},
			wantLen: 0,
		},
		"error-invalid-option": {
			opts:    []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErr: true,
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := recentlyviewed.NewService(method)
			got, err := s.ListIssues(context.Background(), tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestService_AddIssue(t *testing.T) {
	cases := map[string]struct {
		issueID    int
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErr    bool
		wantID     int
	}{
		"success": {
			issueID: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1/recentlyViewedIssues", spath)
				return mock.NewResponse(`{"id":1,"summary":"test issue"}`), nil
			},
			wantID: 1,
		},
		"error-invalid-issueID": {
			issueID: 0,
			wantErr: true,
		},
		"error-client-network": {
			issueID: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			issueID: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			s := recentlyviewed.NewService(method)
			got, err := s.AddIssue(context.Background(), tc.issueID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_ListProjects(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts      []*option.APIParamOption
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr   bool
		wantLen   int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself/recentlyViewedProjects", spath)
				return mock.NewResponse(`[{"id":1},{"id":2}]`), nil
			},
			wantLen: 2,
		},
		"success-with-count": {
			opts: []*option.APIParamOption{o.WithCount(3)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "3", query.Get("count"))
				return mock.NewResponse(`[{"id":1}]`), nil
			},
			wantLen: 1,
		},
		"error-invalid-option": {
			opts:    []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErr: true,
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := recentlyviewed.NewService(method)
			got, err := s.ListProjects(context.Background(), tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestService_ListWikis(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts      []*option.APIParamOption
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr   bool
		wantLen   int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself/recentlyViewedWikis", spath)
				return mock.NewResponse(`[{"id":1},{"id":2}]`), nil
			},
			wantLen: 2,
		},
		"success-with-order": {
			opts: []*option.APIParamOption{o.WithOrder("asc")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewResponse(`[{"id":1}]`), nil
			},
			wantLen: 1,
		},
		"error-invalid-option": {
			opts:    []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErr: true,
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := recentlyviewed.NewService(method)
			got, err := s.ListWikis(context.Background(), tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestService_AddWiki(t *testing.T) {
	cases := map[string]struct {
		wikiID     int
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErr    bool
		wantID     int
	}{
		"success": {
			wikiID: 10,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/10/recentlyViewedWikis", spath)
				return mock.NewResponse(`{"id":10,"name":"TestWiki"}`), nil
			},
			wantID: 10,
		},
		"error-invalid-wikiID": {
			wikiID:  0,
			wantErr: true,
		},
		"error-client-network": {
			wikiID: 10,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			wikiID: 10,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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

			s := recentlyviewed.NewService(method)
			got, err := s.AddWiki(context.Background(), tc.wikiID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeGetFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	makePostFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *client.Method)
	}{
		{"ListIssues", func(t *testing.T, m *client.Method) {
			m.Get = makeGetFn(t)
			s := recentlyviewed.NewService(m)
			s.ListIssues(ctx) //nolint:errcheck
		}},
		{"AddIssue", func(t *testing.T, m *client.Method) {
			m.Post = makePostFn(t)
			s := recentlyviewed.NewService(m)
			s.AddIssue(ctx, 1) //nolint:errcheck
		}},
		{"ListProjects", func(t *testing.T, m *client.Method) {
			m.Get = makeGetFn(t)
			s := recentlyviewed.NewService(m)
			s.ListProjects(ctx) //nolint:errcheck
		}},
		{"ListWikis", func(t *testing.T, m *client.Method) {
			m.Get = makeGetFn(t)
			s := recentlyviewed.NewService(m)
			s.ListWikis(ctx) //nolint:errcheck
		}},
		{"AddWiki", func(t *testing.T, m *client.Method) {
			m.Post = makePostFn(t)
			s := recentlyviewed.NewService(m)
			s.AddWiki(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &client.Method{})
		})
	}
}
