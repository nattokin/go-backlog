package backlog

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserService_One(t *testing.T) {
	cases := map[string]struct {
		id int

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-id-1": {
			id: 1,

			mockGetFn: newMockGetFn(t, "users/1", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}),

			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
		},
		"success-id-100": {
			id: 100,

			mockGetFn: newMockGetFn(t, "users/100", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
			}),

			wantUser: &User{},
		},
		"invalid-id-0": {
			id: 0,

			wantErrType: &ValidationError{},
		},
	}

	for name, tc := range cases {
t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			user, err := s.One(tc.id)

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

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-add-user": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    RoleAdministrator,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"post-error": {
			userID:      "errorUser",
			password:    "password",
			name:        "error",
			mailAddress: "error@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				return nil, errors.New("network failure")
			},

			wantErrType: errors.New(""),
		},
		"invalid-empty-userID": {
			userID:      "",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: newUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"invalid-empty-password": {
			userID:      "admin",
			password:    "",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: newUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"invalid-empty-name": {
			userID:      "admin",
			password:    "password",
			name:        "",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,

			mockPostFn: newUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"invalid-empty-mailAddress": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "",
			roleType:    RoleAdministrator,

			mockPostFn: newUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"invalid-option-validation-error": {
			userID:      "test",
			password:    "",
			name:        "",
			mailAddress: "",
			roleType:    RoleAdministrator,

			mockPostFn: newUnexpectedPostFn(t),

			wantErrType: &ValidationError{},
		},
		"invalid-json-response": {
			userID:      "userID",
			password:    "password",
			name:        "name",
			mailAddress: "mailAdress",
			roleType:    1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newUserService()

			// default: unexpected API call
			s.method.Post = newUnexpectedPostFn(t)

			if tc.mockPostFn != nil {
				s.method.Post = tc.mockPostFn
			}

			user, err := s.Add(tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)

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
		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantLen     int
		wantFirst   *User
		wantErrType error
	}{
		"success-get-users": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		"request-error": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"invalid-json-response": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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

			s := newUserService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			users, err := s.All()

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
	o := newUserOptionService()

	cases := map[string]struct {
		id      int
		options []*FormOption

		mockPatchFn func(spath string, form *FormParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-update-user": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("password"),
				o.WithFormName("admin"),
				o.WithFormMailAddress("eguchi@nulab.example"),
				o.WithFormRoleType(RoleAdministrator),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"invalid-id": {
			id: 0,

			wantErrType: &ValidationError{},
		},
		"invalid-json-response": {
			id: 1234,

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)

				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},

			wantErrType: &json.SyntaxError{},
		},
		"option-WithName": {
			id: 1,
			options: []*FormOption{
				o.WithFormName("testname"),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testname", form.Get("name"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"option-WithPassword": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("testpassword"),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword", form.Get("password"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"option-WithMailAddress": {
			id: 1,
			options: []*FormOption{
				o.WithFormMailAddress("test@test.com"),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "test@test.com", form.Get("mailAddress"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"option-WithRoleType": {
			id: 1,
			options: []*FormOption{
				o.WithFormRoleType(RoleAdministrator),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"option-multi-WithFormName-WithFormPassword": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("testpassword1"),
				o.WithFormName("testname1"),
				o.WithFormMailAddress("test1@test.com"),
				o.WithFormRoleType(RoleAdministrator),
			},

			mockPatchFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword1", form.Get("password"))
				assert.Equal(t, "testname1", form.Get("name"))
				assert.Equal(t, "test1@test.com", form.Get("mailAddress"))
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"option-error": {
			id: 1,
			options: []*FormOption{
				o.WithFormName(""),
			},

			wantErrType: &ValidationError{},
		},
		"option-invalid": {
			id:      1,
			options: []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},

			wantErrType: &InvalidOptionError[formType]{},
		},
	}

	for name, tc := range cases {
t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()

			// default: unexpected API call
			s.method.Patch = newUnexpectedPatchFn(t)

			if tc.mockPatchFn != nil {
				s.method.Patch = tc.mockPatchFn
			}

			user, err := s.Update(tc.id, tc.options...)

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
		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-get-own-user": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		"request-error": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"invalid-json-response": {
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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

			s := newUserService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			user, err := s.Own()

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

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
	}{
		"success-id-1": {
			id: 1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},

			wantErrType: nil,
		},
		"invalid-id-0": {
			id: 0,

			wantErrType: &ValidationError{},
		},
		"invalid-json-response": {
			id: 1234,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newUserService()

			// default: unexpected API call
			s.method.Delete = newUnexpectedDeleteFn(t)

			if tc.mockDeleteFn != nil {
				s.method.Delete = tc.mockDeleteFn
			}

			user, err := s.Delete(tc.id)

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

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantUsers   []*User
		wantErrType error
	}{
		"success-projectKey-valid": {
			projectKey:          "TEST",
			excludeGroupMembers: false,

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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
		"invalid-projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"invalid-json-response": {
			projectKey: "TEST",

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
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

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			users, err := s.All(tc.projectKey, tc.excludeGroupMembers)

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

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"userID-0": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"invalid-json-response": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Post = newUnexpectedPostFn(t)

			if tc.mockPostFn != nil {
				s.method.Post = tc.mockPostFn
			}

			user, err := s.Add(tc.projectKey, tc.userID)

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

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"success-delete-user": {
			projectKey: "TEST",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"success-projectKey-number": {
			projectKey: "1234",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"invalid-projectKey-empty": {
			projectKey: "",
			userID:     1,

			wantErrType: &ValidationError{},
		},
		"invalid-userID-0": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"success-userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"invalid-json-response": {
			projectKey: "TEST3",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Delete = newUnexpectedDeleteFn(t)

			if tc.mockDeleteFn != nil {
				s.method.Delete = tc.mockDeleteFn
			}

			user, err := s.Delete(tc.projectKey, tc.userID)

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

		mockPostFn func(spath string, form *FormParams) (*http.Response, error)

		wantUser    *User
		wantErrType error
	}{
		"projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
		"userID-0": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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
		"invalid-json-response": {
			projectKey: "TEST3",
			userID:     1,

			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Post = newUnexpectedPostFn(t)

			if tc.mockPostFn != nil {
				s.method.Post = tc.mockPostFn
			}

			user, err := s.AddAdmin(tc.projectKey, tc.userID)

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

		mockGetFn func(spath string, query *QueryParams) (*http.Response, error)

		wantErrType error
	}{
		"projectKey-valid": {
			projectKey: "TEST",

			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"projectKey-empty": {
			projectKey: "",

			wantErrType: &ValidationError{},
		},
	}

	for name, tc := range cases {
t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Get = newUnexpectedGetFn(t)

			if tc.mockGetFn != nil {
				s.method.Get = tc.mockGetFn
			}

			users, err := s.AdminAll(tc.projectKey)

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

		mockDeleteFn func(spath string, form *FormParams) (*http.Response, error)

		wantErrType error
	}{
		"projectKey-valid": {
			projectKey: "TEST",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/administrators", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return nil, errors.New("error")
			},

			wantErrType: errors.New(""),
		},
		"projectKey-empty": {
			projectKey: "",
			userID:     1,

			wantErrType: &ValidationError{},
		},
		"userID-0": {
			projectKey: "TEST1",
			userID:     0,

			wantErrType: &ValidationError{},
		},
		"userID-1": {
			projectKey: "TEST2",
			userID:     1,

			mockDeleteFn: func(spath string, form *FormParams) (*http.Response, error) {
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

			s := newProjectUserService()

			// default: unexpected API call
			s.method.Delete = newUnexpectedDeleteFn(t)

			if tc.mockDeleteFn != nil {
				s.method.Delete = tc.mockDeleteFn
			}

			user, err := s.DeleteAdmin(tc.projectKey, tc.userID)

			assert.Error(t, err)
			assert.IsType(t, tc.wantErrType, err)
			assert.Nil(t, user)
		})
	}
}
