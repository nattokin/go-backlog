package backlog

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
	"github.com/nattokin/go-backlog/internal/user"
)

func TestUserService_One(t *testing.T) {
	cases := map[string]struct {
		id int

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-id-1": {
			id: 1,

			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"success-id-100": {
			id: 100,

			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
				}, nil
			},

			wantUser: &User{},
		},
		"error-validation-id-zero": {
			id: 0,

			wantErrType: &ValidationError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewUserService(method, nil)

			user, err := s.One(context.Background(), tc.id)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestUserService_Add(t *testing.T) {
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    Role

		mockPostFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-add-user": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    RoleAdministrator,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Equal(t, "admin", form.Get("userId"))
				assert.Equal(t, "password", form.Get("password"))
				assert.Equal(t, "admin", form.Get("name"))
				assert.Equal(t, "eguchi@nulab.example", form.Get("mailAddress"))
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-client-network": {
			userID:      "errorUser",
			password:    "password",
			name:        "error",
			mailAddress: "error@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				return nil, errors.New("network failure")
			},

			wantErrType: errors.New(""),
		},
		"error-validation-userID-empty": {
			userID:      "",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: mock.NewUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"error-validation-password-empty": {
			userID:      "admin",
			password:    "",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: mock.NewUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"error-validation-name-empty": {
			userID:      "admin",
			password:    "password",
			name:        "",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: mock.NewUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"error-validation-mailAddress-empty": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "",
			roleType:    RoleAdministrator,

			mockPostFn: mock.NewUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"error-validation-multiple-empty": {
			userID:      "test",
			password:    "",
			name:        "",
			mailAddress: "",
			roleType:    RoleAdministrator,

			mockPostFn: mock.NewUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"error-response-invalid-json": {
			userID:      "userID",
			password:    "password",
			name:        "name",
			mailAddress: "mailAdress",
			roleType:    1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Post: mock.NewUnexpectedPostFn(t),
			}
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := user.NewUserService(method, nil)

			user, err := s.Add(context.Background(), tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestUserService_All(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantLen     int
		wantFirst   *User
		wantErrType error
	}{
		"success-get-users": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},

			wantLen: 4,
			wantFirst: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewUserService(method, nil)

			users, err := s.All(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, users)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.Len(t, users, tc.wantLen)

			require.NotNil(t, users[0])
			assert.Equal(t, tc.wantFirst.UserID, users[0].UserID)
			assert.Equal(t, tc.wantFirst.Name, users[0].Name)
			assert.Equal(t, tc.wantFirst.MailAddress, users[0].MailAddress)
			assert.Equal(t, tc.wantFirst.RoleType, users[0].RoleType)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	o := user.NewUserOptionService(&core.OptionService{})

	cases := map[string]struct {
		id   int
		opts []RequestOption

		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-update-user": {
			id: 1,
			opts: []RequestOption{
				o.WithPassword("password"),
				o.WithName("admin"),
				o.WithMailAddress("eguchi@nulab.example"),
				o.WithRoleType(RoleAdministrator),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "password", form.Get("password"))
				assert.Equal(t, "admin", form.Get("name"))
				assert.Equal(t, "eguchi@nulab.example", form.Get("mailAddress"))
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-validation-id-zero": {
			id: 0,

			wantErrType: &ValidationError{},
		},
		"error-response-invalid-json": {
			id: 1234,

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
		"success-option-withName": {
			id: 1,
			opts: []RequestOption{
				o.WithName("testname"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testname", form.Get("name"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"success-option-withPassword": {
			id: 1,
			opts: []RequestOption{
				o.WithPassword("testpassword"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword", form.Get("password"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"success-option-withMailAddress": {
			id: 1,
			opts: []RequestOption{
				o.WithMailAddress("test@test.com"),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "test@test.com", form.Get("mailAddress"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"success-option-withRoleType": {
			id: 1,
			opts: []RequestOption{
				o.WithRoleType(RoleAdministrator),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"success-option-multiple": {
			id: 1,
			opts: []RequestOption{
				o.WithPassword("testpassword1"),
				o.WithName("testname1"),
				o.WithMailAddress("test1@test.com"),
				o.WithRoleType(RoleAdministrator),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword1", form.Get("password"))
				assert.Equal(t, "testname1", form.Get("name"))
				assert.Equal(t, "test1@test.com", form.Get("mailAddress"))
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-option-invalid-value": {
			id: 1,
			opts: []RequestOption{
				o.WithName(""),
			},

			wantErrType: &ValidationError{},
		},
		"error-option-invalid-type": {
			id:   1,
			opts: []RequestOption{mock.NewInvalidTypeOption()},

			wantErrType: &InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			id:          1,
			opts:        []RequestOption{mock.NewFailingSetOption(core.ParamName)},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Patch: mock.NewUnexpectedPatchFn(t),
			}
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}
			s := user.NewUserService(method, nil)

			user, err := s.Update(context.Background(), tc.id, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, tc.wantUser.UserID, user.UserID)
			assert.Equal(t, tc.wantUser.Name, user.Name)
			assert.Equal(t, tc.wantUser.MailAddress, user.MailAddress)
			assert.Equal(t, tc.wantUser.RoleType, user.RoleType)
		})
	}
}

func TestUserService_Own(t *testing.T) {
	cases := map[string]struct {
		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-get-own-user": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-client-network": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedPatchFn(t),
			}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewUserService(method, nil)

			user, err := s.Own(context.Background())

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

func TestUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		id int

		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)

		wantErrType error
	}{
		"success-id-1": {
			id: 1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"success-id-100": {
			id: 100,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"error-validation-id-zero": {
			id: 0,

			wantErrType: &ValidationError{},
		},
		"error-response-invalid-json": {
			id: 1234,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Delete: mock.NewUnexpectedPatchFn(t),
			}
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := user.NewUserService(method, nil)

			user, err := s.Delete(context.Background(), tc.id)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
		})
	}
}

func TestProjectUserService_All(t *testing.T) {
	cases := map[string]struct {
		projectKey          string
		excludeGroupMembers bool

		mockGetFn func(ctx context.Context, spath string, query url.Values) (*http.Response, error)

		wantUsers   []*User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey:          "TEST",
			excludeGroupMembers: false,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},

			wantUsers: []*User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    RoleAdministrator,
				},
			},
		},
		"success-excludeGroupMembers-true": {
			projectKey:          "TEST2",
			excludeGroupMembers: true,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "true", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},

			wantUsers: []*User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    RoleAdministrator,
				},
			},
		},
		"success-excludeGroupMembers-false": {
			projectKey:          "TEST3",
			excludeGroupMembers: false,

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},

			wantUsers: []*User{
				{
					UserID:      "admin",
					Name:        "admin",
					MailAddress: "eguchi@nulab.example",
					RoleType:    RoleAdministrator,
				},
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"error-response-invalid-json": {
			projectKey: "TEST",

			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewProjectUserService(method, nil)

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

		wantUser    *User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Post: mock.NewUnexpectedPostFn(t),
			}
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := user.NewProjectUserService(method, nil)

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

		wantUser    *User
		wantErrType error
	}{
		"success-delete-user": {
			projectKey: "TEST",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"success-projectIDOrKey-number": {
			projectKey: "1234",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/1234/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",
			userID:     1,

			wantErrType: &ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Delete: mock.NewUnexpectedDeleteFn(t),
			}
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := user.NewProjectUserService(method, nil)

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

		wantUser    *User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-validation-projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"error-response-invalid-json": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Post: mock.NewUnexpectedPostFn(t),
			}
			if tc.mockPostFn != nil {
				method.Post = tc.mockPostFn
			}
			s := user.NewProjectUserService(method, nil)

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

			wantErrType: &ValidationError{},
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			// default: unexpected API call
			method := &core.Method{
				Get: mock.NewUnexpectedGetFn(t),
			}
			if tc.mockGetFn != nil {
				method.Get = tc.mockGetFn
			}
			s := user.NewProjectUserService(method, nil)

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

			wantErrType: &ValidationError{},
		},
		"error-validation-userID-zero": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
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

			// default: unexpected API call
			method := &core.Method{
				Delete: mock.NewUnexpectedDeleteFn(t),
			}
			if tc.mockDeleteFn != nil {
				method.Delete = tc.mockDeleteFn
			}
			s := user.NewProjectUserService(method, nil)

			user, err := s.DeleteAdmin(context.Background(), tc.projectKey, tc.userID)

			assert.Error(t, err)
			assert.IsType(t, tc.wantErrType, err)
			assert.Nil(t, user)
		})
	}
}

// TestUserService_contextPropagation verifies that the context passed to each
// UserService and ProjectUserService method is correctly relayed to the
// underlying method call.
// A sentinel value is embedded in the context and its pointer identity is
// asserted inside the mock to catch any ctx substitution (e.g. context.Background()).
func TestUserService_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	o := user.NewUserOptionService(&core.OptionService{})

	cases := []struct {
		name string
		call func(t *testing.T)
	}{
		{"UserService.All", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.All(ctx) //nolint:errcheck
		}},
		{"UserService.One", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"UserService.Own", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Own(ctx) //nolint:errcheck
		}},
		{"UserService.Add", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Post: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Add(ctx, "u", "p", "n", "m@m.com", RoleAdministrator) //nolint:errcheck
		}},
		{"UserService.Update", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Patch: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Update(ctx, 1, o.WithName("n")) //nolint:errcheck
		}},
		{"UserService.Delete", func(t *testing.T) {
			s := user.NewUserService(&core.Method{
				Delete: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
		{"ProjectUserService.All", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.All(ctx, "TEST", false) //nolint:errcheck
		}},
		{"ProjectUserService.Add", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Post: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Add(ctx, "TEST", 1) //nolint:errcheck
		}},
		{"ProjectUserService.Delete", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Delete: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.Delete(ctx, "TEST", 1) //nolint:errcheck
		}},
		{"ProjectUserService.AddAdmin", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Post: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.AddAdmin(ctx, "TEST", 1) //nolint:errcheck
		}},
		{"ProjectUserService.AdminAll", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Get: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.AdminAll(ctx, "TEST") //nolint:errcheck
		}},
		{"ProjectUserService.DeleteAdmin", func(t *testing.T) {
			s := user.NewProjectUserService(&core.Method{
				Delete: func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
					assert.Same(t, sentinel, got.Value(ctxKey{}))
					return nil, errors.New("stop")
				},
			}, nil)
			s.DeleteAdmin(ctx, "TEST", 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t)
		})
	}
}
