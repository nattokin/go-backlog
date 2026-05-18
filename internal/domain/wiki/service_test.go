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

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/wiki"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectIDOrKey-id": {
			projectIDOrKey: "103",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(fixture.Wiki.ListJSON), nil
			},
		},

		"success-projectIDOrKey-key-with-options": {
			projectIDOrKey: "PRJ_KEY",
			opts: []core.RequestOption{
				o.WithKeyword("test"),
			},

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				assert.Equal(t, "test", query.Get("keyword"))
				return mock.NewJSONResponse(fixture.Wiki.ListJSON), nil
			},
		},

		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},

		"error-option-invalid-type": {
			projectIDOrKey: "invalid",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		"error-option-set-failed": {
			projectIDOrKey: "PRJ",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamKeyword)},
			wantErrType:    errors.New(""),
		},

		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},

		"error-response-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			s := wiki.NewService(method)

			wikis, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, wikis)

			assert.Len(t, wikis, 2)
			assert.Equal(t, 112, wikis[0].ID)
			assert.Equal(t, "test1", wikis[0].Name)
		})
	}
}

func TestService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-projectIDOrKey-id": {
			projectIDOrKey: "103",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "103", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(`{"count": 34}`), nil
			},

			wantCount: 34,
		},
		"success-projectIDOrKey-key": {
			projectIDOrKey: "PRJ_KEY",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "PRJ_KEY", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(`{"count": 10}`), nil
			},

			wantCount: 10,
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return nil, errors.New("network error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "1",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				assert.Equal(t, "1", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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

			s := wiki.NewService(method)

			count, err := s.Count(context.Background(), tc.projectIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}
