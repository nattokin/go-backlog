package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestWikiStarService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID                 int
		mockGetFn              func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErr                bool
		wantValidationErrCount int
		wantLen                int
	}{
		"success": {
			wikiID: 34,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/34/stars", spath)
				assert.Nil(t, query)
				return mock.NewResponse(fixture.Star.ListJSON), nil
			},
			wantLen: 2,
		},
		"error-invalid-wikiID": {
			wikiID:                 0,
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			wikiID: 34,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
		"error-client-api-error": {
			wikiID: 34,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErr: true,
		},
		"error-json-decode": {
			wikiID: 34,
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

			s := wiki.NewStarService(method)
			got, err := s.List(context.Background(), tc.wikiID)

			if tc.wantValidationErrCount > 0 {
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Len(t, got, tc.wantLen)
		})
	}
}
