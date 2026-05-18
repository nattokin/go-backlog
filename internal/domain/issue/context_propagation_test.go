package issue_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
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

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewService(m)
			s.List(ctx) //nolint:errcheck
		}},
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewService(m)
			seq, err := s.All(ctx, 10)
			require.NoError(t, err)
			for range seq {
				break
			}
		}},
		{"Service.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewService(m)
			s.Count(ctx) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewService(m)
			s.One(ctx, "PRJ-1") //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := issue.NewService(m)
			s.Create(ctx, 10, "summary", 2, 3) //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := issue.NewService(m)
			s.Update(ctx, "PRJ-1", o.WithSummary("x")) //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := issue.NewService(m)
			s.Delete(ctx, "PRJ-1") //nolint:errcheck
		}},
		{"Service.Participants", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewService(m)
			s.Participants(ctx, "PRJ-1") //nolint:errcheck
		}},
		{"AttachmentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewAttachmentService(m)
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"AttachmentService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := issue.NewAttachmentService(m)
			s.Remove(ctx, "TEST-1", 1) //nolint:errcheck
		}},
		{"AttachmentService.Download", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := issue.NewAttachmentService(m)
			s.Download(ctx, "TEST-1", 1) //nolint:errcheck
		}},
		{"CommentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.List(ctx, "ISSUE-1") //nolint:errcheck
		}},
		{"CommentService.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Add(ctx, "ISSUE-1", "comment") //nolint:errcheck
		}},
		{"CommentService.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Count(ctx, "ISSUE-1") //nolint:errcheck
		}},
		{"CommentService.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.One(ctx, "ISSUE-1", 1) //nolint:errcheck
		}},
		{"CommentService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Delete(ctx, "ISSUE-1", 1) //nolint:errcheck
		}},
		{"CommentService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Update(ctx, "ISSUE-1", 1, "content") //nolint:errcheck
		}},
		{"CommentService.Notifications", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Notifications(ctx, "ISSUE-1", 1) //nolint:errcheck
		}},
		{"CommentService.Notify", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := issue.NewCommentService(m)
			s.Notify(ctx, "ISSUE-1", 1, []int{1, 2}) //nolint:errcheck
		}},
		{"SharedFileService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := issue.NewSharedFileService(m)
			s.List(ctx, "TEST-1") //nolint:errcheck
		}},
		{"SharedFileService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := issue.NewSharedFileService(m)
			s.Link(ctx, "TEST-1", []int{1}) //nolint:errcheck
		}},
		{"SharedFileService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := issue.NewSharedFileService(m)
			s.Unlink(ctx, "TEST-1", 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
