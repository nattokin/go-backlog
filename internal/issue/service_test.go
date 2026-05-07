package issue_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestIssueService_All(t *testing.T) {
	cases := map[string]struct {
		expectError bool
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				return mock.NewJSONResponse(fixture.Issue.ListJSON), nil
			},
		},
		"error-client": {
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
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
			s := issue.NewService(method)

			issues, err := s.All(context.Background())

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, issues)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, issues)
			assert.Len(t, issues, len(fixture.Issue.List))

			for i, w := range fixture.Issue.List {
				assert.Equal(t, w.ID, issues[i].ID)
				assert.Equal(t, w.IssueKey, issues[i].IssueKey)
				assert.Equal(t, w.Summary, issues[i].Summary)
			}
		})
	}
}

func TestIssueService_Count(t *testing.T) {
	cases := map[string]struct {
		expectError bool
		wantCount   int
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			wantCount: fixture.Issue.Count,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/count", spath)
				return mock.NewJSONResponse(fixture.Issue.CountJSON), nil
			},
		},
		"error-client": {
			expectError: true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
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
			s := issue.NewService(method)

			count, err := s.Count(context.Background())

			if tc.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestIssueService_One(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		expectError  bool
		mockGetFn    func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"error-validation-empty": {
			issueIDOrKey: "",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},
		"error-validation-zero": {
			issueIDOrKey: "0",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},
		"error-client": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
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
			s := issue.NewService(method)

			v, err := s.One(context.Background(), tc.issueIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Issue.Single.ID, v.ID)
			assert.Equal(t, fixture.Issue.Single.IssueKey, v.IssueKey)
			assert.Equal(t, fixture.Issue.Single.Summary, v.Summary)
		})
	}
}

func TestIssueService_Create(t *testing.T) {
	cases := map[string]struct {
		projectID   int
		summary     string
		issueTypeID int
		priorityID  int
		expectError bool
		mockPostFn  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			projectID:   1,
			summary:     "test issue",
			issueTypeID: 2,
			priorityID:  3,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues", spath)
				assert.Equal(t, "1", form.Get("projectId"))
				assert.Equal(t, "test issue", form.Get("summary"))
				assert.Equal(t, "2", form.Get("issueTypeId"))
				assert.Equal(t, "3", form.Get("priorityId"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"error-validation-projectID-zero": {
			projectID:   0,
			summary:     "test",
			issueTypeID: 2,
			priorityID:  3,
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-validation-summary-empty": {
			projectID:   1,
			summary:     "",
			issueTypeID: 2,
			priorityID:  3,
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-validation-issueTypeID-zero": {
			projectID:   1,
			summary:     "test",
			issueTypeID: 0,
			priorityID:  3,
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-validation-priorityID-zero": {
			projectID:   1,
			summary:     "test",
			issueTypeID: 2,
			priorityID:  0,
			expectError: true,
			mockPostFn:  mock.NewUnexpectedPostFn(t),
		},
		"error-client": {
			projectID:   1,
			summary:     "test",
			issueTypeID: 2,
			priorityID:  3,
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			projectID:   1,
			summary:     "test",
			issueTypeID: 2,
			priorityID:  3,
			expectError: true,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Post = tc.mockPostFn
			s := issue.NewService(method)

			v, err := s.Create(context.Background(), tc.projectID, tc.summary, tc.issueTypeID, tc.priorityID)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Issue.Single.ID, v.ID)
		})
	}
}

func TestIssueService_Update(t *testing.T) {
	o := &core.OptionService{}
	cases := map[string]struct {
		issueIDOrKey string
		expectError  bool
		mockPatchFn  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1", spath)
				assert.Equal(t, "updated summary", form.Get("summary"))
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"error-validation-empty": {
			issueIDOrKey: "",
			expectError:  true,
			mockPatchFn:  mock.NewUnexpectedPatchFn(t),
		},
		"error-client": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Patch = tc.mockPatchFn
			s := issue.NewService(method)

			v, err := s.Update(context.Background(), tc.issueIDOrKey, o.WithSummary("updated summary"))

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Issue.Single.ID, v.ID)
		})
	}
}

func TestIssueService_Delete(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		expectError  bool
		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1", spath)
				return mock.NewJSONResponse(fixture.Issue.SingleJSON), nil
			},
		},
		"error-validation-empty": {
			issueIDOrKey: "",
			expectError:  true,
			mockDeleteFn: mock.NewUnexpectedDeleteFn(t),
		},
		"error-client": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			method.Delete = tc.mockDeleteFn
			s := issue.NewService(method)

			v, err := s.Delete(context.Background(), tc.issueIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Equal(t, fixture.Issue.Single.ID, v.ID)
		})
	}
}

func TestIssueService_Participants(t *testing.T) {
	cases := map[string]struct {
		issueIDOrKey string
		expectError  bool
		mockGetFn    func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
	}{
		"success": {
			issueIDOrKey: "TEST-1",
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "issues/TEST-1/participants", spath)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
		},
		"error-validation-empty": {
			issueIDOrKey: "",
			expectError:  true,
			mockGetFn:    mock.NewUnexpectedGetFn(t),
		},
		"error-client": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("error")
			},
		},
		"error-invalid-json": {
			issueIDOrKey: "TEST-1",
			expectError:  true,
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
			s := issue.NewService(method)

			v, err := s.Participants(context.Background(), tc.issueIDOrKey)

			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, v)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, v)
			assert.Len(t, v, len(fixture.User.List))

			for i, w := range fixture.User.List {
				assert.Equal(t, w.ID, v[i].ID)
			}
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
			s := issue.NewService(m)
			s.All(ctx) //nolint:errcheck
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

func newTestIssue() *fixture.IssueData {
	return fixture.Issue.Single
}

func newTestTime() time.Time {
	return time.Date(2014, time.October, 14, 8, 16, 27, 0, time.UTC)
}

func init() {
	_ = newTestIssue
	_ = newTestTime
}
