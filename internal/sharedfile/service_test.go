package sharedfile_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/sharedfile"
)

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"IssueService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"IssueService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.Link(ctx, "TEST-1", []int{1}) //nolint:errcheck
		}},
		{"IssueService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := sharedfile.NewIssueService(m)
			s.Unlink(ctx, "TEST-1", 1) //nolint:errcheck
		}},
		{"WikiService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"WikiService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.Link(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"WikiService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := sharedfile.NewWikiService(m)
			s.Unlink(ctx, 1, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
