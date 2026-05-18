package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/wiki"
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
			s := wiki.NewService(m)
			s.List(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewService(m)
			s.Count(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewService(m)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := wiki.NewService(m)
			s.Create(ctx, 1, "name", "content") //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := wiki.NewService(m)
			s.Update(ctx, 1, o.WithName("n")) //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := wiki.NewService(m)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
		{"AttachmentService.Attach", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := wiki.NewAttachmentService(m)
			s.Attach(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"AttachmentService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewAttachmentService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"AttachmentService.Remove", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := wiki.NewAttachmentService(m)
			s.Remove(ctx, 1, 1) //nolint:errcheck
		}},
		{"AttachmentService.Download", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := wiki.NewAttachmentService(m)
			s.Download(ctx, 1, 1) //nolint:errcheck
		}},
		{"HistoryService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewHistoryService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"StarService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewStarService(m)
			s.List(ctx, 34) //nolint:errcheck
		}},
		{"SharedFileService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := wiki.NewSharedFileService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"SharedFileService.Link", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := wiki.NewSharedFileService(m)
			s.Link(ctx, 1, []int{1}) //nolint:errcheck
		}},
		{"SharedFileService.Unlink", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := wiki.NewSharedFileService(m)
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
