package wiki_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestWikiHistoryService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
	}{
		"success": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/history", spath)
				return mock.NewResponse(fixture.WikiHistory.ListJSON), nil
			},
		},

		"error-validation-wikiID-zero": {
			wikiID:                 0,
			wantValidationErrCount: 1,
		},
		"error-validation-wikiID-negative": {
			wikiID:                 -1,
			wantValidationErrCount: 1,
		},

		"error-client-network": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
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
			s := wiki.NewHistoryService(method)

			entries, err := s.List(context.Background(), tc.wikiID)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, entries)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, entries)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, entries)
			assert.Len(t, entries, len(fixture.WikiHistory.List))

			for i, w := range fixture.WikiHistory.List {
				assert.Equal(t, w.PageID, entries[i].PageID)
				assert.Equal(t, w.Version, entries[i].Version)
				assert.Equal(t, w.Name, entries[i].Name)
				assert.Equal(t, w.Content, entries[i].Content)
			}
		})
	}
}
