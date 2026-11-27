package watching_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/watching"
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
		call func(t *testing.T, m *client.Method)
	}{
		{"Service.List", func(t *testing.T, m *client.Method) {
			m.Get = makeMockFn(t)
			s := watching.NewService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"Service.Count", func(t *testing.T, m *client.Method) {
			m.Get = makeMockFn(t)
			s := watching.NewService(m)
			s.Count(ctx, 1) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *client.Method) {
			m.Get = makeMockFn(t)
			s := watching.NewService(m)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *client.Method) {
			m.Post = makeMockFn(t)
			s := watching.NewService(m)
			s.Create(ctx, "PRJ-1") //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *client.Method) {
			m.Patch = makeMockFn(t)
			s := watching.NewService(m)
			s.Update(ctx, 1, "note") //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *client.Method) {
			m.Delete = makeMockFn(t)
			s := watching.NewService(m)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
		{"Service.MarkAsRead", func(t *testing.T, m *client.Method) {
			m.Post = makeMockFn(t)
			s := watching.NewService(m)
			s.MarkAsRead(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &client.Method{})
		})
	}
}
