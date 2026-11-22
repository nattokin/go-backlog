package activity_test

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
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/shared/activity"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func TestService_ApplyOptions(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts []*option.APIParamOption

		wantQuery              url.Values
		wantErrType            error
		wantValidationErrCount int
	}{
		"success-no-options": {
			wantQuery: url.Values{},
		},
		"success-all-options": {
			opts: []*option.APIParamOption{
				o.WithActivityTypeIDs([]int{1, 2}),
				o.WithMinID(10),
				o.WithMaxID(20),
				o.WithCount(30),
				o.WithOrder("asc"),
			},
			wantQuery: url.Values{
				"activityTypeId[]": {"1", "2"},
				"minId":            {"10"},
				"maxId":            {"20"},
				"count":            {"30"},
				"order":            {"asc"},
			},
		},
		"error-validation-count": {
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-option-not-allowed": {
			opts:        []*option.APIParamOption{o.WithName("name")},
			wantErrType: &option.InvalidOptionKeyError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := activity.NewService(mock.NewMethod(t))
			query := url.Values{}
			err := s.ApplyOptions(query, tc.opts...)

			if tc.wantValidationErrCount > 0 {
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantQuery, query)
		})
	}
}

func TestService_Fetch(t *testing.T) {
	query := url.Values{"count": {"10"}}

	cases := map[string]struct {
		mockGetFn func(context.Context, string, url.Values) (*http.Response, error)

		wantFirstID int
		wantErrType error
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, q url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/activities", spath)
				assert.Equal(t, query, q)
				return mock.NewResponse(fixture.Activity.ListJSON), nil
			},
			wantFirstID: 3153,
		},
		"error-client-network": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, errors.New("network")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return nil, &client.APIResponseError{}
			},
			wantErrType: &client.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(context.Context, string, url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			m := mock.NewMethod(t)
			m.Get = tc.mockGetFn
			s := activity.NewService(m)
			got, err := s.Fetch(context.Background(), "projects/TEST/activities", query)

			if tc.wantErrType != nil {
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotEmpty(t, got)
			assert.Equal(t, tc.wantFirstID, got[0].ID)
		})
	}
}
