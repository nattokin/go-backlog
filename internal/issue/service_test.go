package issue_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_All(t *testing.T) {
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-projectIDs": {
			opts: []core.RequestOption{o.WithProjectIDs([]int{10, 20})},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, []string{"10", "20"}, query["projectId[]"])
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-keyword": {
			opts: []core.RequestOption{o.WithKeyword("bug")},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "bug", query.Get("keyword"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-sort-and-order": {
			opts: []core.RequestOption{
				o.WithIssueSort(model.IssueSortCreated),
				o.WithOrder(model.OrderAsc),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "created", query.Get("sort"))
				assert.Equal(t, "asc", query.Get("order"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
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
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-date-filters": {
			opts: []core.RequestOption{
				o.WithCreatedSince(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
				o.WithCreatedUntil(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "2024-01-01", query.Get("createdSince"))
				assert.Equal(t, "2024-12-31", query.Get("createdUntil"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-parentChild": {
			opts: []core.RequestOption{o.WithParentChild(1)},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "1", query.Get("parentChild"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.Issue.ListJSON))),
				}, nil
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
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(fixture.InvalidJSON))),
				}, nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := &core.Method{Get: mock.NewUnexpectedGetFn(t)}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}

			s := issue.NewService(method)

			issues, err := s.All(context.Background(), tc.opts...)

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

func TestIssueService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"All", func(t *testing.T, m *core.Method) {
			m.Get = func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
				assert.Same(t, sentinel, got.Value(ctxKey{}))
				return nil, errors.New("stop")
			}
			s := issue.NewService(m)
			s.All(ctx) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
