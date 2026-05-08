package project_test

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
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestActivityService_List(t *testing.T) {
	o := &core.OptionService{}

	type want struct {
		spath          string
		activityTypeID []string
		minID          string
		maxID          string
		count          string
		order          string
	}
	cases := map[string]struct {
		projectIDOrKey string
		opts           []core.RequestOption
		mockGetFn      func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErrType    error
		want           want
	}{
		"success-no-option": {
			projectIDOrKey: "TEST",
			opts:           []core.RequestOption{},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/activities", spath)
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{
				spath:          "projects/TEST/activities",
				activityTypeID: nil,
				minID:          "",
				maxID:          "",
				count:          "",
				order:          "",
			},
		},
		"success-withActivityTypeIDs": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1}),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"1"}, query["activityTypeId[]"])
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{activityTypeID: []string{"1"}},
		},
		"success-withMinID": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithMinID(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("minId"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{minID: "1"},
		},
		"success-withMaxID": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithMaxID(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("maxId"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{maxID: "1"},
		},
		"success-withCount": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithCount(1),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "1", query.Get("count"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{count: "1"},
		},
		"success-withOrder": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithOrder(model.OrderAsc),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{order: "asc"},
		},
		"success-multiple-options": {
			projectIDOrKey: "TEST",
			opts: []core.RequestOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(1),
				o.WithMaxID(26),
				o.WithCount(20),
				o.WithOrder(model.OrderAsc),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/activities", spath)
				assert.Equal(t, []string{"1", "2"}, query["activityTypeId[]"])
				assert.Equal(t, "1", query.Get("minId"))
				assert.Equal(t, "26", query.Get("maxId"))
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Activity.ListJSON), nil
			},
			want: want{
				spath:          "projects/TEST/activities",
				activityTypeID: []string{"1", "2"},
				minID:          "1",
				maxID:          "26",
				count:          "20",
				order:          "asc",
			},
		},
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey: "",
			wantErrType:    &core.ValidationError{},
		},
		"error-client-network": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/activities", spath)
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
		"error-option-invalid-value": {
			projectIDOrKey: "TEST",
			opts:           []core.RequestOption{o.WithCount(0)},
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "TEST",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-option-set-failed": {
			projectIDOrKey: "TEST",
			opts:           []core.RequestOption{mock.NewFailingSetOption(core.ParamCount)},
			wantErrType:    errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := project.NewActivityService(method)

			got, err := s.List(context.Background(), tc.projectIDOrKey, tc.opts...)

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
