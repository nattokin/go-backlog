package notification_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/notification"
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
			s := notification.NewService(m)
			s.List(ctx) //nolint:errcheck
		}},
		{"Service.Count", func(t *testing.T, m *client.Method) {
			m.Get = makeMockFn(t)
			s := notification.NewService(m)
			s.Count(ctx) //nolint:errcheck
		}},
		{"Service.ResetUnreadCount", func(t *testing.T, m *client.Method) {
			m.Post = makeMockFn(t)
			s := notification.NewService(m)
			s.ResetUnreadCount(ctx) //nolint:errcheck
		}},
		{"Service.MarkAsRead", func(t *testing.T, m *client.Method) {
			m.Post = makeMockFn(t)
			s := notification.NewService(m)
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
