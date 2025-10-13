package backlog_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strconv"
	"testing"

	"github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestUserService_One_getUser(t *testing.T) {
	t.Parallel()

	userID := "admin"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, "users/1", spath)
			assert.Nil(t, query)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.One(1)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestProjectUserService_All_getUserList(t *testing.T) {
	userID := "admin"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	projectKey := "TEST"
	excludeGroupMembers := false

	s := backlog.ExportNewProjectUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, "projects/"+projectKey+"/users", spath)
			assert.Equal(t, strconv.FormatBool(excludeGroupMembers), query.Get("excludeGroupMembers"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserListJSON))),
			}
			return resp, nil
		},
	})

	users, err := s.All(projectKey, excludeGroupMembers)
	assert.NoError(t, err)
	assert.Equal(t, userID, users[0].UserID)
	assert.Equal(t, name, users[0].Name)
	assert.Equal(t, mailAddress, users[0].MailAddress)
	assert.Equal(t, roleType, users[0].RoleType)
}

func TestUserService_Add_addUser(t *testing.T) {
	t.Parallel()

	userID := "admin"
	password := "password"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			assert.Equal(t, "users", spath)
			assert.Equal(t, userID, form.Get("userId"))
			assert.Equal(t, password, form.Get("password"))
			assert.Equal(t, name, form.Get("name"))
			assert.Equal(t, mailAddress, form.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(int(roleType)), form.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.Add(userID, password, name, mailAddress, roleType)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestProjectUserService_Delete_deleteUser(t *testing.T) {
	t.Parallel()

	userID := "admin"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	projectKey := "TEST"
	id := 1

	s := backlog.ExportNewProjectUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			assert.Equal(t, "projects/"+projectKey+"/users", spath)
			assert.Equal(t, strconv.Itoa(id), form.Get("userId"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}
			return resp, nil
		},
	})

	users, err := s.Delete(projectKey, id)
	assert.NoError(t, err)
	assert.Equal(t, userID, users.UserID)
	assert.Equal(t, name, users.Name)
	assert.Equal(t, mailAddress, users.MailAddress)
	assert.Equal(t, roleType, users.RoleType)
}

func TestUserService_Update_updateUser(t *testing.T) {
	t.Parallel()

	id := 1
	userID := "admin"
	password := "password"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			assert.Equal(t, "users/"+strconv.Itoa(id), spath)
			assert.Equal(t, name, form.Get("name"))
			assert.Equal(t, password, form.Get("password"))
			assert.Equal(t, name, form.Get("name"))
			assert.Equal(t, mailAddress, form.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(int(roleType)), form.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataUserJSON))),
			}
			return resp, nil
		},
	})

	option := s.Option
	user, err := s.Update(
		id, option.WithFormPassword(password), option.WithFormName(name), option.WithFormMailAddress(mailAddress), option.WithFormRoleType(roleType),
	)
	assert.NoError(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestUserService_All(t *testing.T) {
	t.Parallel()

	want := struct {
		spath string
	}{
		spath: "users",
	}

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)
			return nil, errors.New("error")
		},
	})

	users, err := s.All()
	assert.Nil(t, users)
	assert.Error(t, err)
}

func TestUserService_All_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	users, err := s.All()
	assert.Nil(t, users)
	assert.Error(t, err)
}

func TestUserService_One(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		id        int
		wantError bool
		want      want
	}{
		"id_1": {
			id:        1,
			wantError: false,
			want: want{
				spath: "users/1",
			},
		},
		"id_100": {
			id:        100,
			wantError: false,
			want: want{
				spath: "users/100",
			},
		},

		"id_0": {
			id:        0,
			wantError: true,
			want:      want{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, query)
					}
					return nil, errors.New("error")
				},
			})

			s.One(tc.id)
		})

	}
}

func TestUserService_Own(t *testing.T) {
	t.Parallel()

	want := struct {
		spath string
	}{
		spath: "users/myself",
	}

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, query)
			return nil, errors.New("error")
		},
	})

	user, err := s.Own()
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Own_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.Own()
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Add(t *testing.T) {
	wantSpath := "users"
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    backlog.Role
		wantError   bool
	}{
		"no_error": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    backlog.RoleAdministrator,
			wantError:   false,
		},
		"userID_empty": {
			userID:      "",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    backlog.RoleAdministrator,
			wantError:   true,
		},
		"password_empty": {
			userID:      "testid",
			password:    "",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    backlog.RoleAdministrator,
			wantError:   true,
		},
		"name_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "",
			mailAddress: "test@test.com",
			roleType:    backlog.RoleAdministrator,
			wantError:   true,
		},
		"mailAddress_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "",
			roleType:    backlog.RoleAdministrator,
			wantError:   true,
		},
		"roleType_invalid": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    0,
			wantError:   true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, wantSpath, spath)
						assert.Equal(t, tc.userID, form.Get("userId"))
						assert.Equal(t, tc.password, form.Get("password"))
						assert.Equal(t, tc.name, form.Get("name"))
						assert.Equal(t, tc.mailAddress, form.Get("mailAddress"))
						assert.Equal(t, strconv.Itoa(int(tc.roleType)), form.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			})

			user, err := s.Add(tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)
			assert.Nil(t, user)
			assert.Error(t, err)
		})

	}
}

func TestUserService_Add_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.Add("userID", "password", "name", "mailAdress", 1)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Update(t *testing.T) {
	cases := map[string]struct {
		id        int
		wantError bool
	}{
		"valid": {
			id:        1,
			wantError: false,
		},
		"invalid": {
			id:        0,
			wantError: true,
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Patch must never be called")
					} else {
						assert.Equal(t, "users/"+strconv.Itoa(tc.id), spath)
					}
					return nil, errors.New("error")
				},
			})

			user, err := s.Update(tc.id)
			assert.Nil(t, user)
			assert.Error(t, err)
		})

	}
}

func TestUserService_Update_option(t *testing.T) {
	o := backlog.ExportNewUserOptionService()

	type options struct {
		password    string
		name        string
		mailAddress string
		roleType    string
	}
	cases := map[string]struct {
		options   []*backlog.FormOption
		wantError bool
		want      options
	}{
		"WithoutOption": {
			options:   []*backlog.FormOption{},
			wantError: false,
			want:      options{},
		},
		"WithName": {
			options: []*backlog.FormOption{
				o.WithFormName("testname"),
			},
			wantError: false,
			want: options{
				name: "testname",
			},
		},
		"WithPassword": {
			options: []*backlog.FormOption{
				o.WithFormPassword("testpasword"),
			},
			wantError: false,
			want: options{
				password: "testpasword",
			},
		},
		"WithMailAddress": {
			options: []*backlog.FormOption{
				o.WithFormMailAddress("test@test.com"),
			},
			wantError: false,
			want: options{
				mailAddress: "test@test.com",
			},
		},
		"WithRoleType": {
			options: []*backlog.FormOption{
				o.WithFormRoleType(backlog.RoleAdministrator),
			},
			wantError: false,
			want: options{
				roleType: "1",
			},
		},
		"MultiOptions": {
			options: []*backlog.FormOption{
				o.WithFormPassword("testpasword1"),
				o.WithFormName("testname1"),
				o.WithFormMailAddress("test1@test.com"),
				o.WithFormRoleType(backlog.RoleAdministrator),
			},
			wantError: false,
			want: options{
				password:    "testpasword1",
				name:        "testname1",
				mailAddress: "test1@test.com",
				roleType:    "1",
			},
		},
		"OptionError": {
			options: []*backlog.FormOption{
				o.WithFormName(""),
			},
			wantError: true,
			want:      options{},
		},
		"InvalidOption": {
			options: []*backlog.FormOption{
				backlog.ExportNewFormOption(0, nil, func(p *backlog.FormParams) error {
					return nil
				}),
			},
			wantError: true,
			want:      options{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			id := 1

			s := backlog.ExportNewUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Patch: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Patch must never be called")
					} else {
						assert.Equal(t, "users/"+strconv.Itoa(id), spath)
						assert.Equal(t, tc.want.password, form.Get("password"))
						assert.Equal(t, tc.want.name, form.Get("name"))
						assert.Equal(t, tc.want.mailAddress, form.Get("mailAddress"))
						assert.Equal(t, tc.want.roleType, form.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			})

			user, err := s.Update(id, tc.options...)
			assert.Nil(t, user)
			assert.Error(t, err)
		})

	}
}

func TestUserService_Update_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Patch: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.Update(1234)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Delete(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		id        int
		wantError bool
		want      want
	}{
		"id_1": {
			id:        1,
			wantError: false,
			want: want{
				spath: "users/1",
			},
		},
		"id_100": {
			id:        100,
			wantError: false,
			want: want{
				spath: "users/100",
			},
		},

		"id_0": {
			id:        0,
			wantError: true,
			want:      want{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, form)
					}
					return nil, errors.New("error")
				},
			})

			user, err := s.Delete(tc.id)
			assert.Nil(t, user)
			assert.Error(t, err)
		})

	}
}

func TestUserService_Delete_invalidJson(t *testing.T) {
	t.Parallel()

	s := backlog.ExportNewUserService()
	s.ExportSetMethod(&backlog.ExportMethod{
		Delete: func(spath string, form *backlog.FormParams) (*http.Response, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewReader([]byte(testdataInvalidJSON))),
			}
			return resp, nil
		},
	})

	user, err := s.Delete(1234)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestProjectUserService_All(t *testing.T) {
	type want struct {
		spath               string
		excludeGroupMembers string
	}
	cases := map[string]struct {
		projectKey          string
		excludeGroupMembers bool
		wantError           bool
		want                want
	}{
		"projectKey_valid": {
			projectKey:          "TEST",
			excludeGroupMembers: false,
			wantError:           false,
			want: want{
				spath:               "projects/TEST/users",
				excludeGroupMembers: "false",
			},
		},
		"projectKey_empty": {
			projectKey:          "",
			excludeGroupMembers: false,
			wantError:           true,
			want:                want{},
		},
		"excludeGroupMembers_true": {
			projectKey:          "TEST2",
			excludeGroupMembers: true,
			wantError:           false,
			want: want{
				spath:               "projects/TEST2/users",
				excludeGroupMembers: "true",
			},
		},
		"excludeGroupMembers_false": {
			projectKey:          "TEST3",
			excludeGroupMembers: false,
			wantError:           false,
			want: want{
				spath:               "projects/TEST3/users",
				excludeGroupMembers: "false",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.excludeGroupMembers, query.Get("excludeGroupMembers"))
					}
					return nil, errors.New("error")
				},
			})

			s.All(tc.projectKey, tc.excludeGroupMembers)
		})

	}
}

func TestProjectUserService_Add(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectKey string
		userID     int
		wantError  bool
		want       want
	}{
		"projectKey_valid": {
			projectKey: "TEST",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectKey_empty": {
			projectKey: "",
			userID:     1,
			wantError:  true,
			want:       want{},
		},
		"userID_0": {
			projectKey: "TEST1",
			userID:     0,
			wantError:  true,
			want:       want{},
		},
		"userID_1": {
			projectKey: "TEST2",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST2/users",
				userID: "1",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, form.Get("userId"))
					}
					return nil, errors.New("error")
				},
			})

			s.Add(tc.projectKey, tc.userID)
		})

	}
}

func TestProjectUserService_Delete(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectKey string
		userID     int
		wantError  bool
		want       want
	}{
		"projectKey_string": {
			projectKey: "TEST",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectKey_number": {
			projectKey: "1234",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/1234/users",
				userID: "1",
			},
		},
		"projectKey_empty": {
			projectKey: "",
			userID:     1,
			wantError:  true,
			want:       want{},
		},
		"userID_0": {
			projectKey: "TEST1",
			userID:     0,
			wantError:  true,
			want:       want{},
		},
		"userID_1": {
			projectKey: "TEST2",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST2/users",
				userID: "1",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, form.Get("userId"))
					}
					return nil, errors.New("error")
				},
			})

			s.Delete(tc.projectKey, tc.userID)
		})

	}
}

func TestProjectUserService_AddAdmin(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectKey string
		userID     int
		wantError  bool
		want       want
	}{
		"projectKey_valid": {
			projectKey: "TEST",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectKey_empty": {
			projectKey: "",
			userID:     1,
			wantError:  true,
			want:       want{},
		},
		"userID_0": {
			projectKey: "TEST1",
			userID:     0,
			wantError:  true,
			want:       want{},
		},
		"userID_1": {
			projectKey: "TEST2",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST2/administrators",
				userID: "1",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Post: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, form.Get("userId"))
					}
					return nil, errors.New("error")
				},
			})

			s.AddAdmin(tc.projectKey, tc.userID)
		})

	}
}

func TestProjectUserService_AdminAll(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		projectKey string
		wantError  bool
		want       want
	}{
		"projectKey_valid": {
			projectKey: "TEST",
			wantError:  false,
			want: want{
				spath: "projects/TEST/administrators",
			},
		},
		"projectKey_empty": {
			projectKey: "",
			wantError:  true,
			want:       want{},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Get: func(spath string, query *backlog.QueryParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, query)
					}
					return nil, errors.New("error")
				},
			})

			s.AdminAll(tc.projectKey)
		})

	}
}

func TestProjectUserService_DeleteAdmin(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectKey string
		userID     int
		wantError  bool
		want       want
	}{
		"projectKey_valid": {
			projectKey: "TEST",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectKey_empty": {
			projectKey: "",
			userID:     1,
			wantError:  true,
			want:       want{},
		},
		"userID_0": {
			projectKey: "TEST1",
			userID:     0,
			wantError:  true,
			want:       want{},
		},
		"userID_1": {
			projectKey: "TEST2",
			userID:     1,
			wantError:  false,
			want: want{
				spath:  "projects/TEST2/administrators",
				userID: "1",
			},
		},
	}

	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			s := backlog.ExportNewProjectUserService()
			s.ExportSetMethod(&backlog.ExportMethod{
				Delete: func(spath string, form *backlog.FormParams) (*http.Response, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, form.Get("userId"))
					}
					return nil, errors.New("error")
				},
			})

			s.DeleteAdmin(tc.projectKey, tc.userID)
		})

	}
}
