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
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestUserStarService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		userID    int
		opts      []core.RequestOption
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr   bool
		wantLen   int
	}{
		"success-no-options": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/stars", spath)
				return newJSONResponse(`[{"id":10},{"id":20}]`), nil
			},
			wantLen: 2,
		},
		"success-with-count": {
			userID: 2,
			opts:   []core.RequestOption{o.WithCount(5)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/2/stars", spath)
				assert.Equal(t, "5", query.Get("count"))
				return newJSONResponse(`[{"id":1}]`), nil
			},
			wantLen: 1,
		},
		"error-invalid-userID": {
			userID:  0,
			wantErr: true,
		},
		"error-invalid-option": {
			userID:  1,
			opts:    []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErr: true,
		},
		"error-client-network": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return newJSONResponse(fixture.InvalidJSON), nil
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

			s := star.NewUserService(method)
			got, err := s.List(context.Background(), tc.userID, tc.opts...)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestUserStarService_Count(t *testing.T) {
	cases := map[string]struct {
		userID    int
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr   bool
		wantCount int
	}{
		"success": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/stars/count", spath)
				assert.Nil(t, query)
				return newJSONResponse(`{"count":42}`), nil
			},
			wantCount: 42,
		},
		"error-invalid-userID": {
			userID:  0,
			wantErr: true,
		},
		"error-client-network": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-json-decode": {
			userID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return newJSONResponse(fixture.InvalidJSON), nil
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

			s := star.NewUserService(method)
			got, err := s.Count(context.Background(), tc.userID)

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, got)
		})
	}
}
