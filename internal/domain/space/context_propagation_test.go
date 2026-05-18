package space_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/space"
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

	makeMockUploadFn := func(t *testing.T) func(context.Context, string, string, io.Reader) (*http.Response, error) {
		return func(got context.Context, _, _ string, _ io.Reader) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.Info", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := space.NewService(m)
			s.Info(ctx) //nolint:errcheck
		}},
		{"Service.DiskUsage", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := space.NewService(m)
			s.DiskUsage(ctx) //nolint:errcheck
		}},
		{"Service.Notification", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := space.NewService(m)
			s.Notification(ctx) //nolint:errcheck
		}},
		{"Service.UpdateNotification", func(t *testing.T, m *core.Method) {
			m.Put = makeMockFn(t)
			s := space.NewService(m)
			s.UpdateNotification(ctx, "content") //nolint:errcheck
		}},
		{"ActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := space.NewActivityService(m)
			s.List(ctx) //nolint:errcheck
		}},
		{"ActivityService.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := space.NewActivityService(m)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"AttachmentService.Upload", func(t *testing.T, m *core.Method) {
			m.Upload = makeMockUploadFn(t)
			s := space.NewAttachmentService(m)
			s.Upload(ctx, "file.txt", nil) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
