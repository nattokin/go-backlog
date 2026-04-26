package user_test

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
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/user"
)

func TestProjectUserService_All(t *testing.T) {
	cases := map[string]struct {
		projectKey          string
		excludeGroupMembers bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantUsers   []*model.User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey:          "TEST",
			excludeGroupMembers: false,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},

			wantUsers: []*model.User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    model.RoleAdministrator,
				},
			},
		},
		"success-excludeGroupMembers-true": {
			projectKey:          "TEST2",
			excludeGroupMembers: true,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "true", query.Get("excludeGroupMembers"))
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},

			wantUsers: []*model.User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    model.RoleAdministrator,
				},
			},
		},
		"success-excludeGroupMembers-false": {
			projectKey:          "TEST3",
			excludeGroupMembers: false,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},

			wantUsers: []*model.User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    model.RoleAdministrator,
				},
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &core.ValidationError{},
		},
		"error-response-invalid-json": {
			projectKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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
			s := user.NewProjectService(method)

			users, err := s.All(context.Background(), tc.projectKey, tc.excludeGroupMembers)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, users)
				return
			}

			assert.NoError(t, err)
			require.Len(t, users, 4)
			require.NotNil(t, users[0])
			assert.Equal(t, tc.wantUsers[0].UserID, users[0].UserID)
			assert.Equal(t, tc.wantUsers[0].Name, users[0].Name)
			assert.Equal(t, tc.wantUsers[0].MailAddress, users[0].MailAddress)
			assert.Equal(t, tc.wantUsers[0].RoleType, users[0].RoleType)
		})
	}
}

func TestProjectUserService_Add(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantUser    *model.User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &core.ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &core.ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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
			s := user.NewProjectService(method)

			user, err := s.Add(context.Background(), tc.projectKey, tc.userID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestProjectUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantUser    *model.User
		wantErrType error
	}{
		"success-delete-user": {
			projectKey: "TEST",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"success-projectIDOrKey-number": {
			projectKey: "1234",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",
			userID:     1,

			wantErrType: &core.ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &core.ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := user.NewProjectService(method)

			user, err := s.Delete(context.Background(), tc.projectKey, tc.userID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestProjectUserService_AddAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantUser    *model.User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &core.ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &core.ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},

			wantUser: &model.User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    model.RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
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
			s := user.NewProjectService(method)

			user, err := s.AddAdmin(context.Background(), tc.projectKey, tc.userID)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestProjectUserService_AdminAll(t *testing.T) {
	cases := map[string]struct {
		projectKey string

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &core.ValidationError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewProjectService(method)

			users, err := s.AdminAll(context.Background(), tc.projectKey)

			assert.Error(t, err)
			assert.IsType(t, tc.wantErrType, err)
			assert.Nil(t, users)
		})
	}
}

func TestProjectUserService_DeleteAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey string
		userID     int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-validation-projectKey-empty": {
			projectKey: "",
			userID:     1,

			wantErrType: &core.ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &core.ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := user.NewProjectService(method)

			user, err := s.DeleteAdmin(context.Background(), tc.projectKey, tc.userID)

			assert.Error(t, err)
			assert.IsType(t, tc.wantErrType, err)
			assert.Nil(t, user)
		})
	}
}
