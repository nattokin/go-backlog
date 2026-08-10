package pullrequest_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestCommentService_List(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		opts           []*core.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantIDs                []int
	}{
		"success-no-options": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-count-and-order": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts: []*core.APIParamOption{
				o.WithCount(20),
				o.WithOrder("asc"),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "20", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},
		"success-with-minID-maxID": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts: []*core.APIParamOption{
				o.WithMinID(10),
				o.WithMaxID(100),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "100", query.Get("maxId"))
				return mock.NewResponse(fixture.Comment.ListJSON), nil
			},
			wantIDs: []int{1, 2},
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			repoIDOrName:           "repo",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "0",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               0,
			wantValidationErrCount: 1,
		},
		"error-validation-opt": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamCount)},
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			opts:                   []*core.APIParamOption{mock.NewFailingCheckOption(core.ParamCount)},
			wantValidationErrCount: 4,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := pullrequest.NewCommentService(method)
			got, err := s.List(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Len(t, got, len(tc.wantIDs))
			for i := range got {
				assert.Equal(t, tc.wantIDs[i], got[i].ID)
			}
		})
	}
}

func TestCommentService_Add(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		content        string
		opts           []*core.APIParamOption

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantID                 int
	}{
		"success-required-only": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "This is a comment.",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments", spath)
				assert.Equal(t, "This is a comment.", form.Get("content"))
				return mock.NewCreatedResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-notifiedUserIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "Notifying users.",
			opts:           []*core.APIParamOption{o.WithNotifiedUserIDs([]int{5, 6})},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, []string{"5", "6"}, form["notifiedUserId[]"])
				return mock.NewCreatedResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},
		"success-with-attachmentIDs": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "Attaching files.",
			opts:           []*core.APIParamOption{o.WithAttachmentIDs([]int{10, 11})},
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "Attaching files.", form.Get("content"))
				assert.Equal(t, []string{"10", "11"}, form["attachmentId[]"])
				return mock.NewCreatedResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo",
			prNumber:               1,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               0,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-content-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			content:                "",
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			content:                "",
			wantValidationErrCount: 4,
		},

		// --- fail-fast: nil option ---
		"error-nil-option-with-valid-values": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			content:                "x",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},
		"error-nil-option-with-invalid-values": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			content:                "",
			opts:                   []*core.APIParamOption{nil},
			wantInvalidOptionError: true,
		},

		// --- fail-fast: invalid option key ---
		"error-option-invalid-type": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			opts:           []*core.APIParamOption{mock.NewInvalidTypeOption()},
			wantErrType:    &core.InvalidOptionKeyError{},
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			content:        "x",
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := pullrequest.NewCommentService(method)
			got, err := s.Add(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.content, tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *core.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}

func TestCommentService_Count(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantCount              int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments/count", spath)
				return mock.NewResponse(`{"count":7}`), nil
			},
			wantCount: 7,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-projectIDOrKey-zero": {
			projectIDOrKey:         "0",
			repoIDOrName:           "repo",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "0",
			prNumber:               1,
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               0,
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			wantValidationErrCount: 3,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := pullrequest.NewCommentService(method)
			count, err := s.Count(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Zero(t, count)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Zero(t, count)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, count)
		})
	}
}

func TestCommentService_Update(t *testing.T) {
	cases := map[string]struct {
		projectIDOrKey string
		repoIDOrName   string
		prNumber       int
		commentID      int
		content        string

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantID                 int
	}{
		"success": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "Updated content.",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/PRJ/git/repositories/repo/pullRequests/1/comments/42", spath)
				assert.Equal(t, "Updated content.", form.Get("content"))
				return mock.NewResponse(fixture.Comment.SingleJSON), nil
			},
			wantID: 1,
		},

		// --- validation errors ---
		"error-validation-projectIDOrKey-empty": {
			projectIDOrKey:         "",
			repoIDOrName:           "repo",
			prNumber:               1,
			commentID:              1,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-repoIDOrName-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "",
			prNumber:               1,
			commentID:              1,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-prNumber-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               0,
			commentID:              1,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-commentID-zero": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			commentID:              0,
			content:                "x",
			wantValidationErrCount: 1,
		},
		"error-validation-content-empty": {
			projectIDOrKey:         "PRJ",
			repoIDOrName:           "repo",
			prNumber:               1,
			commentID:              1,
			content:                "",
			wantValidationErrCount: 1,
		},
		"error-validation-all": {
			projectIDOrKey:         "",
			repoIDOrName:           "",
			prNumber:               0,
			commentID:              0,
			content:                "",
			wantValidationErrCount: 5,
		},

		// --- other errors ---
		"error-client-network": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			projectIDOrKey: "PRJ",
			repoIDOrName:   "repo",
			prNumber:       1,
			commentID:      42,
			content:        "x",
			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return mock.NewResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}
			s := pullrequest.NewCommentService(method)
			got, err := s.Update(context.Background(), tc.projectIDOrKey, tc.repoIDOrName, tc.prNumber, tc.commentID, tc.content)

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves core.ValidationErrors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantID, got.ID)
		})
	}
}
