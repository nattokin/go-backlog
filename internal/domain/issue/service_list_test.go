package issue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		opts []core.RequestOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
		wantIDs     []int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-projectIDs": {
			opts: []core.RequestOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-keyword": {
			opts: []core.RequestOption{o.WithKeyword("bug")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "bug", query.Get("keyword"))
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-sort-and-order": {
			opts: []core.RequestOption{
				o.WithIssueSort("created"),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "created", query.Get("sort"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-offset": {
			opts: []core.RequestOption{
				o.WithCount(50),
				o.WithOffset(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "50", query.Get("count"))
				assert.Equal(t, "100", query.Get("offset"))
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-date-filters": {
			opts: []core.RequestOption{
				o.WithCreatedSince("2024-01-01"),
				o.WithCreatedUntil("2024-12-31"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "2024-01-01", query.Get("createdSince"))
				assert.Equal(t, "2024-12-31", query.Get("createdUntil"))
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-parentChild": {
			opts: []core.RequestOption{o.WithParentChild(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "1", query.Get("parentChild"))
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"error-option-invalid-type": {
			opts:        []core.RequestOption{mock.NewInvalidTypeOption()},
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-invalid-projectID": {
			opts:        []core.RequestOption{o.WithProjectIDs([]int{0})},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-sort": {
			opts:        []core.RequestOption{o.WithIssueSort("invalid")},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-parentChild": {
			opts:        []core.RequestOption{o.WithParentChild(5)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-count": {
			opts:        []core.RequestOption{o.WithCount(0)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-offset": {
			opts:        []core.RequestOption{o.WithOffset(-1)},
			wantErrType: &core.ValidationError{},
		},
		"error-option-set-failed": {
			opts:        []core.RequestOption{mock.NewFailingSetOption(core.ParamKeyword)},
			wantErrType: errors.New(""),
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

			issues, err := s.List(context.Background(), tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, issues)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, issues)
			assert.Len(t, issues, len(tc.wantIDs))
			for i := range issues {
				assert.Equal(t, tc.wantIDs[i], issues[i].ID)
			}
		})
	}
}

func TestService_All(t *testing.T) {
	ctx := context.Background()

	t.Run("multi-page", func(t *testing.T) {
		t.Parallel()

		var callCount atomic.Int32
		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, query url.Values) (*http.Response, error) {
			n := callCount.Add(1)
			assert.Equal(t, "2", query.Get("count"))
			switch n {
			case 1:
				assert.Equal(t, "0", query.Get("offset"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(fixture.Issue.ListJSON)),
				}, nil
			case 2:
				assert.Equal(t, "2", query.Get("offset"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewBufferString(issueLastPageJSON)),
				}, nil
			default:
				t.Errorf("unexpected request #%d", n)
				return nil, nil
			}
		}

		s := issue.NewService(method)
		seq, err := s.All(ctx, 2)
		require.NoError(t, err)
		var got []int
		for iss, err := range seq {
			require.NoError(t, err)
			got = append(got, iss.ID)
		}

		assert.Equal(t, int32(2), callCount.Load())
		assert.Equal(t, []int{1, 2, 3}, got)
	})

	t.Run("break", func(t *testing.T) {
		t.Parallel()

		var callCount atomic.Int32
		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, _ url.Values) (*http.Response, error) {
			n := callCount.Add(1)
			if n > 1 {
				t.Errorf("unexpected request #%d after break", n)
			}
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString(fixture.Issue.ListJSON)),
			}, nil
		}

		s := issue.NewService(method)
		seq, err := s.All(ctx, 2)
		require.NoError(t, err)
		var got []int
		for iss, err := range seq {
			require.NoError(t, err)
			got = append(got, iss.ID)
			break
		}

		assert.Equal(t, int32(1), callCount.Load())
		assert.Len(t, got, 1)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, _ url.Values) (*http.Response, error) {
			return nil, errors.New("network error")
		}

		s := issue.NewService(method)
		seq, err := s.All(ctx, 10)
		require.NoError(t, err)
		for iss, err := range seq {
			assert.Nil(t, iss)
			require.Error(t, err)
			break
		}
	})

	t.Run("error-api-error", func(t *testing.T) {
		t.Parallel()

		method := mock.NewMethod(t)
		method.Get = func(_ context.Context, _ string, _ url.Values) (*http.Response, error) {
			return nil, &core.APIResponseError{}
		}

		s := issue.NewService(method)
		seq, err := s.All(ctx, 10)
		require.NoError(t, err)
		for iss, err := range seq {
			assert.Nil(t, iss)
			require.Error(t, err)
			assert.IsType(t, &core.APIResponseError{}, err)
			break
		}
	})

	t.Run("error-invalid-count", func(t *testing.T) {
		t.Parallel()

		s := issue.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 0)
		require.Error(t, err)
		assert.IsType(t, &core.ValidationError{}, err)
	})

	t.Run("error-invalid-option", func(t *testing.T) {
		t.Parallel()

		s := issue.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, mock.NewInvalidTypeOption())
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})

	t.Run("error-offset-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := issue.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, o.WithOffset(5))
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})

	t.Run("error-count-passed-to-all", func(t *testing.T) {
		t.Parallel()

		o := &core.OptionService{}
		s := issue.NewService(mock.NewMethod(t))
		_, err := s.All(ctx, 10, o.WithCount(50))
		require.Error(t, err)
		assert.IsType(t, &core.InvalidOptionKeyError{}, err)
	})
}
