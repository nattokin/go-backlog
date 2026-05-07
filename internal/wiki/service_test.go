package wiki_test

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/wiki"
)

func TestWikiService_All(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		expectError    bool
		mockGetFn      func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "TEST", query.Get("projectIdOrKey"))
				return mock.NewJSONResponse(fixture.Wiki.ListJSON), nil
			},
		},
		"error-validation-empty": {
			projectIDOrKey: "",
			expectError:    true,
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
		"error-client": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewService(method)

			wikis, err := s.All(context.Background(), tc.projectIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, wikis)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, wikis)
			assert.Len(t, wikis, len(fixture.Wiki.List))

			for i, w := range fixture.Wiki.List {
				assert.Equal(t, w.ID, wikis[i].ID)
				assert.Equal(t, w.Name, wikis[i].Name)
			}
		})
	}
}

func TestWikiService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		expectError    bool
		wantCount      int
		mockGetFn      func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			projectIDOrKey: "TEST",
			wantCount:      fixture.Wiki.Count,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/count", spath)
				return mock.NewJSONResponse(fixture.Wiki.CountJSON), nil
			},
		},
		"error-validation-empty": {
			projectIDOrKey: "",
			expectError:    true,
			mockGetFn:      mock.NewUnexpectedGetFn(t),
		},
		"error-client": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectIDOrKey: "TEST",
			expectError:    true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewService(method)

			count, err := s.Count(context.Background(), tc.projectIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestWikiService_One(t *testing.T) {
	cases := map[string]struct {
		wikiID      int
		expectError bool
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return mock.NewJSONResponse(fixture.Wiki.SingleJSON), nil
			},
		},
		"error-validation-zero": {
			wikiID:      0,
			expectError: true,
			mockGetFn:   mock.NewUnexpectedGetFn(t),
		},
		"error-client": {
			wikiID:      1,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			wikiID:      1,
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Get = tc.mockGetFn
			s := wiki.NewService(method)

			v, err := s.One(context.Background(), tc.wikiID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Wiki.Single.ID, v.ID)
			assert.Equal(t, fixture.Wiki.Single.Name, v.Name)
		})
	}
}

func TestWikiService_Create(t *testing.T) {
	o := &core.OptionService{}
	cases := map[string]struct {
		projectID   int
		name        string
		content     string
		expectError bool
		mockPostFn  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			projectID: 1,
			name:      "test wiki",
			content:   "content",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "test wiki", form.Get("name"))
				assert.Equal(t, "content", form.Get("content"))
				return mock.NewJSONResponse(fixture.Wiki.SingleJSON), nil
			},
		},
		"error-validation-projectID-zero": {
			projectID:   0,
			name:        "test",
			content:     "c",
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-validation-name-empty": {
			projectID:   1,
			name:        "",
			content:     "c",
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-client": {
			projectID:   1,
			name:        "test",
			content:     "c",
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn
			s := wiki.NewService(method)

			v, err := s.Create(context.Background(), tc.projectID, tc.name, tc.content, o.WithMailNotify(false))

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Wiki.Single.ID, v.ID)
		})
	}
}

func TestWikiService_Update(t *testing.T) {
	o := &core.OptionService{}
	cases := map[string]struct {
		wikiID      int
		expectError bool
		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return mock.NewJSONResponse(fixture.Wiki.SingleJSON), nil
			},
		},
		"error-validation-zero": {
			wikiID:      0,
			expectError: true,
			mockPatchFn: mock.NewUnexpectedPatchFn(t),
		},
		"error-client": {
			wikiID:      1,
			expectError: true,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Patch = tc.mockPatchFn
			s := wiki.NewService(method)

			v, err := s.Update(context.Background(), tc.wikiID, o.WithName("updated"))

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Wiki.Single.ID, v.ID)
		})
	}
}

func TestWikiService_Delete(t *testing.T) {
	o := &core.OptionService{}
	cases := map[string]struct {
		wikiID       int
		expectError  bool
		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			wikiID: 1,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "wikis/1", spath)
				return mock.NewJSONResponse(fixture.Wiki.SingleJSON), nil
			},
		},
		"error-validation-zero": {
			wikiID:       0,
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},
		"error-client": {
			wikiID:       1,
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := wiki.NewService(method)

			v, err := s.Delete(context.Background(), tc.wikiID, o.WithMailNotify(false))

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Wiki.Single.ID, v.ID)
		})
	}
}

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
			s := wiki.NewService(m)
			s.All(ctx, "TEST") //nolint:errcheck
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
			s.Create(ctx, 10, "name", "content") //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := wiki.NewService(m)
			s.Update(ctx, 1, o.WithName("x")) //nolint:errcheck
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
			s.List(ctx, 1) //nolint:errcheck
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
