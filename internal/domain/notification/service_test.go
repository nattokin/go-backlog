package notification_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/domain/notification"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/validation"
)

func newNoContentResponse() *http.Response {
	return &http.Response{
		StatusCode: http.StatusNoContent,
		Body:       http.NoBody,
	}
}

func TestService_List(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantValidationErrCount int
		wantInvalidOptionError bool
		wantLen                int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications", spath)
				assert.Empty(t, query)
				return mock.NewResponse(fixture.Notification.ListJSON), nil
			},
			wantLen: 1,
		},
		"success-with-options": {
			opts: []*option.APIParamOption{
				o.WithMinID(10),
				o.WithMaxID(20),
				o.WithCount(5),
				o.WithOrder("asc"),
				o.WithSenderID(2),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications", spath)
				assert.Equal(t, "10", query.Get("minId"))
				assert.Equal(t, "20", query.Get("maxId"))
				assert.Equal(t, "5", query.Get("count"))
				assert.Equal(t, "asc", query.Get("order"))
				assert.Equal(t, "2", query.Get("senderId"))
				return mock.NewResponse(fixture.Notification.ListJSON), nil
			},
			wantLen: 1,
		},
		"success-empty": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return mock.NewResponse(`[]`), nil
			},
			wantLen: 0,
		},

		"error-validation-opt-single": {
			opts:                   []*option.APIParamOption{o.WithCount(0)},
			wantValidationErrCount: 1,
		},
		"error-validation-opt-all": {
			opts:                   []*option.APIParamOption{o.WithCount(0), o.WithSenderID(0)},
			wantValidationErrCount: 2,
		},
		"error-nil-option": {
			opts:                   []*option.APIParamOption{o.WithCount(5), nil},
			wantInvalidOptionError: true,
		},
		"error-option-invalid-type": {
			opts:        []*option.APIParamOption{o.WithAlreadyRead(true)},
			wantErrType: &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
			s := notification.NewService(method)
			got, err := s.List(context.Background(), tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Nil(t, got)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantValidationErrCount > 0 {
				assert.Error(t, err)
				assert.Nil(t, got)
				var ves validation.Errors
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
			assert.Len(t, got, tc.wantLen)
		})
	}
}

func TestService_List_decode(t *testing.T) {
	method := mock.NewMethod(t)
	method.Get = func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
		return mock.NewResponse(fixture.Notification.ListJSON), nil
	}
	s := notification.NewService(method)

	got, err := s.List(context.Background())
	require.NoError(t, err)
	require.Len(t, got, 1)

	n := got[0]
	assert.Equal(t, 299, n.ID)
	assert.True(t, n.AlreadyRead)
	assert.Equal(t, 2, n.Reason)
	assert.True(t, n.ResourceAlreadyRead)
	require.NotNil(t, n.Project)
	assert.Equal(t, "TEST2", n.Project.ProjectKey)
	require.NotNil(t, n.Issue)
	assert.Equal(t, "TEST2-17", n.Issue.IssueKey)
	require.NotNil(t, n.Comment)
	assert.Equal(t, 7007, n.Comment.ID)
	assert.Nil(t, n.PullRequest)
	assert.Nil(t, n.PullRequestComment)
	require.NotNil(t, n.Sender)
	assert.Equal(t, "eguchi", n.Sender.UserID)
}

func TestService_Count(t *testing.T) {
	o := &option.OptionService{}

	cases := map[string]struct {
		opts []*option.APIParamOption

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType            error
		wantInvalidOptionError bool
		wantCount              int
	}{
		"success-no-options": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications/count", spath)
				assert.Empty(t, query)
				return mock.NewResponse(fixture.Notification.CountJSON), nil
			},
			wantCount: 138,
		},
		"success-with-options": {
			opts: []*option.APIParamOption{
				o.WithAlreadyRead(false),
				o.WithResourceAlreadyRead(true),
			},
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications/count", spath)
				assert.Equal(t, "false", query.Get("alreadyRead"))
				assert.Equal(t, "true", query.Get("resourceAlreadyRead"))
				return mock.NewResponse(fixture.Notification.CountJSON), nil
			},
			wantCount: 138,
		},

		"error-nil-option": {
			opts:                   []*option.APIParamOption{o.WithAlreadyRead(false), nil},
			wantInvalidOptionError: true,
		},
		"error-option-invalid-type": {
			opts:        []*option.APIParamOption{o.WithCount(5)},
			wantErrType: &option.InvalidOptionKeyError{},
		},

		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
			s := notification.NewService(method)
			got, err := s.Count(context.Background(), tc.opts...)

			if tc.wantInvalidOptionError {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				var target *option.InvalidOptionError
				assert.ErrorAs(t, err, &target)
				return
			}

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, got)
		})
	}
}

func TestService_ResetUnreadCount(t *testing.T) {
	cases := map[string]struct {
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
		wantCount   int
	}{
		"success": {
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications/markAsRead", spath)
				assert.Empty(t, form)
				return mock.NewResponse(fixture.Notification.ResetUnreadCountJSON), nil
			},
			wantCount: 4,
		},
		"error-client-network": {
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
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
			method.Post = tc.mockPostFn
			s := notification.NewService(method)
			got, err := s.ResetUnreadCount(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Equal(t, 0, got)
				assert.ErrorAs(t, err, &tc.wantErrType)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.wantCount, got)
		})
	}
}

func TestService_MarkAsRead(t *testing.T) {
	cases := map[string]struct {
		id         int
		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErr                bool
		wantValidationErrCount int
	}{
		"success": {
			id: 299,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "notifications/299/markAsRead", spath)
				assert.Empty(t, form)
				return newNoContentResponse(), nil
			},
		},
		"error-validation-id-zero": {
			id:                     0,
			wantValidationErrCount: 1,
		},
		"error-validation-id-negative": {
			id:                     -1,
			wantValidationErrCount: 1,
		},
		"error-client-network": {
			id: 1,
			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErr: true,
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := notification.NewService(method)
			err := s.MarkAsRead(context.Background(), tc.id)

			if tc.wantValidationErrCount > 0 {
				var ves validation.Errors
				if assert.ErrorAs(t, err, &ves) {
					assert.Len(t, ves, tc.wantValidationErrCount)
				}
				return
			}

			if tc.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
		})
	}
}
