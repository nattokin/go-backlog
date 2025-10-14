package backlog

import (
	"bytes"
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
		id          int
		expectError bool
		wantUser    *User
		mockGetFn   func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success-id-1": {
			id:          1,
			expectError: false,
			wantUser: &User{
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
				RoleType:    RoleAdministrator,
			},
			mockGetFn: newMockGetFn(t, "users/1", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}),
		},
		"success-id-100": {
			id:          100,
			expectError: false,
			wantUser:    &User{},
			mockGetFn: newMockGetFn(t, "users/100", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(`{}`))),
			}),
		},
		"invalid-id-0": {
			id:          0,
			expectError: true,
			mockGetFn:   newUnexpectedGetFn(t, "id is invalid"),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Get = tc.mockGetFn

			user, err := s.One(tc.id)
			if tc.expectError {
				assert.Error(t, err)
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

func TestUserService_Add(t *testing.T) {
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    Role
		expectError bool
		mockPostFn  func(spath string, form *FormParams) (*http.Response, error)
		wantUser    *User
	}{
		"success-add-user": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    RoleAdministrator,
			expectError: false,
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
			expectError: true,
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				return nil, errors.New("network failure")
			},
		},
		"invalid-empty-userID": {
			userID:      "",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t, "userID is empty"),
		},
		"invalid-empty-password": {
			userID:      "admin",
			password:    "",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t, "password is empty"),
		},
		"invalid-empty-name": {
			userID:      "admin",
			password:    "password",
			name:        "",
			mailAddress: "admin@example.com",
			roleType:    RoleAdministrator,
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t, "name is empty"),
		},
		"invalid-empty-mailAddress": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "",
			roleType:    RoleAdministrator,
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t, "mailAddress is empty"),
		},
		"invalid-option-validation-error": {
			userID:      "test",
			password:    "",
			name:        "",
			mailAddress: "",
			roleType:    RoleAdministrator,
			expectError: true,
			mockPostFn:  newUnexpectedPostFn(t, "option validation fails"),
		},
		"invalid-json-response": {
			userID:      "userID",
			password:    "password",
			name:        "name",
			mailAddress: "mailAdress",
			roleType:    1,
			expectError: true,
			mockPostFn: func(spath string, form *FormParams) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Post = tc.mockPostFn

			user, err := s.Add(tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)
			if tc.expectError {
				assert.Error(t, err)
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

func TestUserService_All(t *testing.T) {
	cases := map[string]struct {
		expectError bool
		mockGetFn   func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success-get-users": {
			expectError: false,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},
		},
		"request-error": {
			expectError: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},
		},
		"invalid-json-response": {
			expectError: true,
			mockGetFn: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Get = tc.mockGetFn

			users, err := s.All()
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
				return
			}

			assert.NoError(t, err)
			require.Len(t, users, 4)
			require.NotNil(t, users[0])
			assert.Equal(t, "admin", users[0].UserID)
			assert.Equal(t, "admin", users[0].Name)
			assert.Equal(t, "eguchi@nulab.example", users[0].MailAddress)
			assert.Equal(t, RoleAdministrator, users[0].RoleType)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	o := newUserOptionService()

	cases := map[string]struct {
		id          int
		options     []*FormOption
		expectError bool
		mockPatch   func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success-update-user": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("password"),
				o.WithFormName("admin"),
				o.WithFormMailAddress("eguchi@nulab.example"),
				o.WithFormRoleType(RoleAdministrator),
			},
			expectError: false,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
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
		},
		"invalid-id": {
			id:          0,
			expectError: true,
			mockPatch:   newUnexpectedPatchFn(t, "id is invalid"),
		},
		"invalid-json-response": {
			id:          1234,
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
		"option-WithName": {
			id: 1,
			options: []*FormOption{
				o.WithFormName("testname"),
			},
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testname", form.Get("name"))
				return nil, errors.New("error")
			},
		},
		"option-WithPassword": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("testpassword"),
			},
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword", form.Get("password"))
				return nil, errors.New("error")
			},
		},
		"option-WithMailAddress": {
			id: 1,
			options: []*FormOption{
				o.WithFormMailAddress("test@test.com"),
			},
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "test@test.com", form.Get("mailAddress"))
				return nil, errors.New("error")
			},
		},
		"option-WithRoleType": {
			id: 1,
			options: []*FormOption{
				o.WithFormRoleType(RoleAdministrator),
			},
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},
		},
		"option-multi-WithFormName-WithFormPassword": {
			id: 1,
			options: []*FormOption{
				o.WithFormPassword("testpassword1"),
				o.WithFormName("testname1"),
				o.WithFormMailAddress("test1@test.com"),
				o.WithFormRoleType(RoleAdministrator),
			},
			expectError: true,
			mockPatch: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword1", form.Get("password"))
				assert.Equal(t, "testname1", form.Get("name"))
				assert.Equal(t, "test1@test.com", form.Get("mailAddress"))
				assert.Equal(t, strconv.Itoa(int(RoleAdministrator)), form.Get("roleType"))
				return nil, errors.New("error")
			},
		},
		"option-error": {
			id: 1,
			options: []*FormOption{
				o.WithFormName(""),
			},
			expectError: true,
			mockPatch:   newUnexpectedPatchFn(t, "option validation fails"),
		},
		"option-invalid": {
			id:          1,
			options:     []*FormOption{{0, nil, func(p *FormParams) error { return nil }}},
			expectError: true,
			mockPatch:   newUnexpectedPatchFn(t, "invalid option is passed"),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Patch = tc.mockPatch

			user, err := s.Update(tc.id, tc.options...)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, RoleAdministrator, user.RoleType)
		})
	}
}

func TestUserService_Own(t *testing.T) {
	cases := map[string]struct {
		expectError bool
		mockGet     func(spath string, query *QueryParams) (*http.Response, error)
		wantUser    *User
	}{
		"success-get-own-user": {
			expectError: false,
			mockGet: newMockGetFn(t, "users/myself", &http.Response{
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
		"request-error": {
			expectError: true,
			mockGet: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			},
		},
		"invalid-json-response": {
			expectError: true,
			mockGet: newMockGetFn(t, "users/myself", &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Get = tc.mockGet

			user, err := s.Own()
			if tc.expectError {
				assert.Error(t, err)
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
		id          int
		expectError bool
		mockDelete  func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success-id-1": {
			id:          1,
			expectError: false,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},
		},
		"success-id-100": {
			id:          100,
			expectError: false,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},
		},
		"invalid-id-0": {
			id:          0,
			expectError: true,
			mockDelete:  newUnexpectedDeleteFn(t, "id is invalid"),
		},
		"invalid-json-response": {
			id:          1234,
			expectError: true,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)
				assert.Nil(t, form)
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newUserService()
			s.method.Delete = tc.mockDelete

			user, err := s.Delete(tc.id)
			if tc.expectError {
				assert.Error(t, err)
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
		expectError         bool
		mockGet             func(spath string, query *QueryParams) (*http.Response, error)
	}{
		"success-projectKey-valid": {
			projectKey:          "TEST",
			excludeGroupMembers: false,
			expectError:         false,
			mockGet: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},
		},
		"success-excludeGroupMembers-true": {
			projectKey:          "TEST2",
			excludeGroupMembers: true,
			expectError:         false,
			mockGet: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "true", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},
		},
		"success-excludeGroupMembers-false": {
			projectKey:          "TEST3",
			excludeGroupMembers: false,
			expectError:         false,
			mockGet: func(spath string, query *QueryParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST3/users", spath)
				assert.Equal(t, "false", query.Get("excludeGroupMembers"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
				}, nil
			},
		},
		"invalid-projectKey-empty": {
			projectKey:  "",
			expectError: true,
			mockGet:     newUnexpectedGetFn(t, "projectKey is empty"),
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Get = tc.mockGet

			users, err := s.All(tc.projectKey, tc.excludeGroupMembers)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, users)
				return
			}

			assert.NoError(t, err)
			require.Len(t, users, 4)
			require.NotNil(t, users[0])
			assert.Equal(t, "admin", users[0].UserID)
			assert.Equal(t, "admin", users[0].Name)
			assert.Equal(t, "eguchi@nulab.example", users[0].MailAddress)
			assert.Equal(t, RoleAdministrator, users[0].RoleType)
		})
	}
}

func TestProjectUserService_Add(t *testing.T) {
	cases := map[string]struct {
		projectKey  string
		userID      int
		expectError bool
		wantSpath   string
		wantUserID  string
	}{
		"projectKey-valid": {
			projectKey:  "TEST",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST/users",
			wantUserID:  "1",
		},
		"projectKey-empty": {
			projectKey:  "",
			userID:      1,
			expectError: true,
		},
		"userID-0": {
			projectKey:  "TEST1",
			userID:      0,
			expectError: true,
		},
		"userID-1": {
			projectKey:  "TEST2",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST2/users",
			wantUserID:  "1",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
				if tc.expectError {
					t.Error("s.method.Post must never be called")
					return nil, errors.New("unexpected call")
				}
				assert.Equal(t, tc.wantSpath, spath)
				assert.Equal(t, tc.wantUserID, form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			}

			user, err := s.Add(tc.projectKey, tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, RoleAdministrator, user.RoleType)
		})
	}
}

func TestProjectUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		projectKey  string
		userID      int
		expectError bool
		mockDelete  func(spath string, form *FormParams) (*http.Response, error)
	}{
		"success-delete-user": {
			projectKey:  "TEST",
			userID:      1,
			expectError: false,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},
		},
		"success-projectKey-number": {
			projectKey:  "1234",
			userID:      1,
			expectError: false,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/1234/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},
		},
		"invalid-projectKey-empty": {
			projectKey:  "",
			userID:      1,
			expectError: true,
			mockDelete:  newUnexpectedDeleteFn(t, "projectKey is empty"),
		},
		"invalid-userID-0": {
			projectKey:  "TEST1",
			userID:      0,
			expectError: true,
			mockDelete:  newUnexpectedDeleteFn(t, "userID is 0"),
		},
		"success-userID-1": {
			projectKey:  "TEST2",
			userID:      1,
			expectError: false,
			mockDelete: func(spath string, form *FormParams) (*http.Response, error) {
				assert.Equal(t, "projects/TEST2/users", spath)
				assert.Equal(t, "1", form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			},
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Delete = tc.mockDelete

			user, err := s.Delete(tc.projectKey, tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, RoleAdministrator, user.RoleType)
		})
	}
}

func TestProjectUserService_AddAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey  string
		userID      int
		expectError bool
		wantSpath   string
		wantUserID  string
	}{
		"projectKey-valid": {
			projectKey:  "TEST",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST/administrators",
			wantUserID:  "1",
		},
		"projectKey-empty": {
			projectKey:  "",
			userID:      1,
			expectError: true,
		},
		"userID-0": {
			projectKey:  "TEST1",
			userID:      0,
			expectError: true,
		},
		"userID-1": {
			projectKey:  "TEST2",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST2/administrators",
			wantUserID:  "1",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Post = func(spath string, form *FormParams) (*http.Response, error) {
				if tc.expectError {
					t.Error("s.method.Post must never be called")
					return nil, errors.New("unexpected call")
				}
				assert.Equal(t, tc.wantSpath, spath)
				assert.Equal(t, tc.wantUserID, form.Get("userId"))
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
				}, nil
			}

			user, err := s.AddAdmin(tc.projectKey, tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, RoleAdministrator, user.RoleType)
		})
	}
}

func TestProjectUserService_AdminAll(t *testing.T) {
	cases := map[string]struct {
		projectKey  string
		expectError bool
		wantSpath   string
	}{
		"projectKey-valid": {
			projectKey:  "TEST",
			expectError: false,
			wantSpath:   "projects/TEST/administrators",
		},
		"projectKey-empty": {
			projectKey:  "",
			expectError: true,
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Get = func(spath string, query *QueryParams) (*http.Response, error) {
				if tc.expectError {
					t.Error("s.method.Get must never be called")
					return nil, errors.New("unexpected call")
				}
				assert.Equal(t, tc.wantSpath, spath)
				assert.Nil(t, query)
				return nil, errors.New("error")
			}

			s.AdminAll(tc.projectKey)
		})
	}
}

func TestProjectUserService_DeleteAdmin(t *testing.T) {
	cases := map[string]struct {
		projectKey  string
		userID      int
		expectError bool
		wantSpath   string
		wantUserID  string
	}{
		"projectKey-valid": {
			projectKey:  "TEST",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST/administrators",
			wantUserID:  "1",
		},
		"projectKey-empty": {
			projectKey:  "",
			userID:      1,
			expectError: true,
		},
		"userID-0": {
			projectKey:  "TEST1",
			userID:      0,
			expectError: true,
		},
		"userID-1": {
			projectKey:  "TEST2",
			userID:      1,
			expectError: false,
			wantSpath:   "projects/TEST2/administrators",
			wantUserID:  "1",
		},
	}

	for name, tc := range cases {
		tc := tc
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			s := newProjectUserService()
			s.method.Delete = func(spath string, form *FormParams) (*http.Response, error) {
				if tc.expectError {
					t.Error("s.method.Delete must never be called")
					return nil, errors.New("unexpected call")
				}
				assert.Equal(t, tc.wantSpath, spath)
				assert.Equal(t, tc.wantUserID, form.Get("userId"))
				return nil, errors.New("error")
			}

			s.DeleteAdmin(tc.projectKey, tc.userID)
		})
	}
}
