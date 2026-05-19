package pullrequest_test

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
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_Count(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		opts           []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo1/pullRequests/count", spath)
				return mock.NewJSONResponse(`{"count":5}`), nil
			},
			wantCount: 5,
		},
		"success-with-assigneeIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{o.WithAssigneeIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"10", "20"}, query["assigneeId[]"])
				return mock.NewJSONResponse(`{"count":2}`), nil
			},
			wantCount: 2,
		},
		"error-empty-projectIDOrKey": {
			projectIDOrKey: "",
			repoIDOrName:   "repo1",
			wantErrType:    &core.ValidationError{},
		},
		"error-empty-repoIDOrName": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "",
			wantErrType:    &core.ValidationError{},
		},
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			opts:           []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo1",
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

			s := pullrequest.NewService(method)
			count, err := s.Count(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.opts...)

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
