package user_test

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
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/user"
)

func TestActivityService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		userID    int
		opts      []core.RequestOption
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-no-option": {
			userID: 1234,
			opts:   []core.RequestOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234/activities", spath)
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withActivityTypeIDs": {
			userID: 1234,
			opts:   []core.RequestOption{o.WithActivityTypeIDs([]int{1})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1"}, query["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withMinID": {
			userID: 1234,
			opts:   []core.RequestOption{o.WithMinID(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("minId"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withMaxID": {
			userID: 1234,
			opts:   []core.RequestOption{o.WithMaxID(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("maxId"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withCount": {
			userID: 1234,
			opts:   []core.RequestOption{o.WithCount(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("count"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-withOrder": {
			userID: 1234,
			opts:   []core.RequestOption{o.WithOrder(model.OrderAsc)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"success-multiple-options": {
			userID: 1234,
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder(model.OrderAsc),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234/activities", spath)
				assert.Equal(t, []string{"1", "2"}, query["activityTypeId[]"])
				assert.Equal(t, "1", query.Get("minId"))
				assert.Equal(t, "26", query.Get("maxId"))
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
		},
		"error-validation-userID-zero": {
			userID:      0,
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			userID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234/activities", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			userID: 1234,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
		"error-option-invalid-value": {
			userID:      1234,
			opts:        []core.RequestOption{o.WithCount(0)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-type": {
			userID:      1234,
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-set-failed": {
			userID:      1234,
			opts:        []core.RequestOption{mock.NewFailingSetOption(core.ParamCount)},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewActivityService(method)

			got, err := s.List(context.Background(), tc.userID, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
		})
	}
}
