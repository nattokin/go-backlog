package user_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/user"
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
			s := user.NewService(m)
			s.List(ctx) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewService(m)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"Service.Me", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewService(m)
			s.Me(ctx) //nolint:errcheck
		}},
		{"Service.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := user.NewService(m)
			s.Add(ctx, "u", "p", "n", "m@m.com", 1) //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := user.NewService(m)
			s.Update(ctx, 1, o.WithName("n")) //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := user.NewService(m)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
		{"Service.Icon", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := user.NewService(m)
			s.Icon(ctx, 1) //nolint:errcheck
		}},
		{"ActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewActivityService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"StarService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewStarService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"StarService.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewStarService(m)
			s.Count(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
