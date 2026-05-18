package pullrequest_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
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
			s := pullrequest.NewService(m)
			s.List(ctx, testProject, testRepo) //nolint:errcheck
		}},
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			seq, err := s.All(ctx, 10, testProject, testRepo)
			require.NoError(t, err)
			for range seq {
				break
			}
		}},
		{"Service.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Count(ctx, testProject, testRepo) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.One(ctx, testProject, testRepo, 1) //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Create(ctx, testProject, testRepo, "summary", "desc", "main", "feature/foo") //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := pullrequest.NewService(m)
			s.Update(ctx, testProject, testRepo, 1, o.WithSummary("x")) //nolint:errcheck
		}},
		{"AttachmentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.List(ctx, "TEST", "repo", 1) //nolint:errcheck
		}},
		{"AttachmentService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.Remove(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
		{"AttachmentService.Download", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := pullrequest.NewAttachmentService(m)
			s.Download(ctx, "TEST", "repo", 1, 1) //nolint:errcheck
		}},
		{"CommentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.List(ctx, "PRJ-1", "REPO-1", 1) //nolint:errcheck
		}},
		{"CommentService.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Add(ctx, "PRJ-1", "REPO-1", 1, "comment") //nolint:errcheck
		}},
		{"CommentService.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Count(ctx, "PRJ-1", "REPO-1", 1) //nolint:errcheck
		}},
		{"CommentService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := pullrequest.NewCommentService(m)
			s.Update(ctx, "PRJ-1", "REPO-1", 1, 1, "content") //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
