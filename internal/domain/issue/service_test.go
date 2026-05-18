package issue_test

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
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				return mock.NewJSONResponse(`{"count":42}`), nil
			},
			wantCount: 42,
		},
		"success-with-projectIDs": {
			opts: []core.RequestOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return mock.NewJSONResponse(`{"count":5}`), nil
			},
			wantCount: 5,
		},
		"success-with-keyword": {
			opts: []core.RequestOption{o.WithKeyword("bug")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				assert.Equal(t, "bug", query.Get("keyword"))
				return mock.NewJSONResponse(`{"count":3}`), nil
			},
			wantCount: 3,
		},
		"error-option-invalid-type": {
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-projectID": {
			opts:        []core.RequestOption{o.WithProjectIDs([]int{0})},
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
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

			s := issue.NewService(method)

			count, err := s.Count(context.Background(), tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestIssueService_Participants(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
	}{
		"success-by-key": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/participants", spath)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			wantIDs: []int{1, 2, 3, 4},
		},
		"success-by-id": {
			issueIDOrKey: "1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/1/participants", spath)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			wantIDs: []int{1, 2, 3, 4},
		},
		"success-empty-list": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/PRJ-1/participants", spath)
				return mock.NewJSONResponse(`[]`), nil
			},
			wantIDs: []int{},
		},
		"error-empty-issueIDOrKey": {
			issueIDOrKey: "",
			wantErrType:  &core.ValidationError{},
		},
		"error-zero-issueIDOrKey": {
			issueIDOrKey: "0",
			wantErrType:  &core.ValidationError{},
		},
		"error-client-network": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-client-api-error": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, &core.APIResponseError{}
			},
			wantErrType: &core.APIResponseError{},
		},
		"error-response-invalid-json": {
			issueIDOrKey: "PRJ-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
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

			s := issue.NewService(method)

			got, err := s.Participants(context.Background(), tc.issueIDOrKey)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Len(t, got, len(tc.wantIDs))
			for i := range got {
				assert.Equal(t, tc.wantIDs[i], got[i].ID)
			}
		})
	}
}
