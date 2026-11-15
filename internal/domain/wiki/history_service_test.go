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
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestWikiHistoryService_List(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		wikiID int
		opts   []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
	}{
		"success": {
			wikiID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/history", spath)
				assert.Empty(t, query)
				return mock.NewResponse(fixture.WikiHistory.ListJSON), nil
			},
		},
		"success-with-options": {
			wikiID: 1234,
			opts: []*option.APIParamOption{
				o.WithMinID(10),
				o.WithMaxID(20),
				o.WithCount(5),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1234/history", spath)
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "20", query.Get("maxId"))
				assert.Equal(t, "5", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
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
		"error-validation-opt-single": {
			wikiID:                 1234,
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			wikiID:                 0,
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 2,
		},

		"error-nil-option-with-valid-values": {
			wikiID:                 1234,
			opts:                   []*option.APIParamOption{o.WithCount(5), nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			wikiID:                 0,
			opts:                   []*option.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		"error-option-invalid-type": {
			wikiID:      1234,
			opts:        []*option.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType: &option.InvalidOptionKeyError{},
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

			entries, err := s.List(context.Background(), tc.wikiID, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, entries)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, entries)
				var ves validation.Errors
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
