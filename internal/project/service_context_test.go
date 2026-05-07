package project_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
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
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.All(ctx) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.One(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewService(m)
			s.Create(ctx, "KEY", "name", o.WithChartEnabled(true)) //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewService(m)
			s.Update(ctx, "TEST", o.WithName("test")) //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewService(m)
			s.Delete(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.DiskUsage", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewService(m)
			s.DiskUsage(ctx, "TEST") //nolint:errcheck
		}},
		{"Service.Icon", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := project.NewService(m)
			s.Icon(ctx, "TEST") //nolint:errcheck
		}},
		{"CategoryService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"CategoryService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Create(ctx, "TEST", "Bug") //nolint:errcheck
		}},
		{"CategoryService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Update(ctx, "TEST", 12, "Bug Fixed") //nolint:errcheck
		}},
		{"CategoryService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewCategoryService(m)
			s.Delete(ctx, "TEST", 12) //nolint:errcheck
		}},
		{"CustomFieldService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"CustomFieldService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.Create(ctx, "TEST", model.CustomFieldTypeText, "Sprint") //nolint:errcheck
		}},
		{"CustomFieldService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.Update(ctx, "TEST", 1, o.WithName("Sprint Updated")) //nolint:errcheck
		}},
		{"CustomFieldService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.Delete(ctx, "TEST", 1) //nolint:errcheck
		}},
		{"CustomFieldService.AddListItem", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.AddListItem(ctx, "TEST", 1, "Item1") //nolint:errcheck
		}},
		{"CustomFieldService.UpdateListItem", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.UpdateListItem(ctx, "TEST", 1, 10, "Item1 Updated") //nolint:errcheck
		}},
		{"CustomFieldService.DeleteListItem", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewCustomFieldService(m)
			s.DeleteListItem(ctx, "TEST", 1, 10) //nolint:errcheck
		}},
		{"IssueTypeService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"IssueTypeService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Create(ctx, "TEST", "Bug", "#e30000") //nolint:errcheck
		}},
		{"IssueTypeService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Update(ctx, "TEST", 1, o.WithName("Bug Updated")) //nolint:errcheck
		}},
		{"IssueTypeService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewIssueTypeService(m)
			s.Delete(ctx, "TEST", 1, 2) //nolint:errcheck
		}},
		{"StatusService.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewStatusService(m)
			s.All(ctx, "TEST") //nolint:errcheck
		}},
		{"StatusService.Create", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Create(ctx, "TEST", "Open", "#ed8077") //nolint:errcheck
		}},
		{"StatusService.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Update(ctx, "TEST", 1, o.WithName("Open")) //nolint:errcheck
		}},
		{"StatusService.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := project.NewStatusService(m)
			s.Delete(ctx, "TEST", 1, 2) //nolint:errcheck
		}},
		{"StatusService.UpdateOrder", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := project.NewStatusService(m)
			s.UpdateOrder(ctx, "TEST", []int{1, 2}) //nolint:errcheck
		}},
		{"SharedFileService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := project.NewSharedFileService(m)
			s.List(ctx, "TEST") //nolint:errcheck
		}},
		{"SharedFileService.GetFile", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := project.NewSharedFileService(m)
			s.GetFile(ctx, "TEST", 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
