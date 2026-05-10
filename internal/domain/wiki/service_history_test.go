package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestWikiHistoryService_List(t *testing.T) {
	cases := map[string]struct {
		wikiID int

		expectError bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/history", spath)

				return mock.NewJSONResponse(fixture.WikiHistory.ListJSON), nil
			},
		},

		"error-wikiID-zero": {
			wikiID:      0,
			expectError: true,
			mockGetFn:   mock.NewUnexpectedGetFn(t),
		},

		"error-wikiID-negative": {
			wikiID:      -1,
			expectError: true,
			mockGetFn:   mock.NewUnexpectedGetFn(t),
		},

		"error-client": {
			wikiID:      1234,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},

		"error-invalid-json": {
			wikiID:      1234,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewHistoryService(method)

			entries, err := s.List(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, entries)
				return
			}

			assert.NoError(t, err)
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
