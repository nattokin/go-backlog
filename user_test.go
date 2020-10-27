package backlog_test

import (
	"errors"
	"net/http"
	"os"
	"strconv"
	"testing"

	backlog "github.com/nattokin/go-backlog"
	"github.com/stretchr/testify/assert"
)

func TestUserService_One_getUser(t *testing.T) {
	userID := "admin"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "users/1", spath)
			assert.Nil(t, params)

			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
	user, err := s.One(1)
	assert.Nil(t, err)
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

	projectIDOrKey := "TEST"
	excludeGroupMembers := false
	bj, err := os.Open("testdata/json/user_list.json")
	if err != nil {
		t.Fatal(err)
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "projects/"+projectIDOrKey+"/users", spath)
			assert.Equal(t, strconv.FormatBool(excludeGroupMembers), params.Get("excludeGroupMembers"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectUserService(m)
	users, err := s.All(projectIDOrKey, excludeGroupMembers)
	assert.Nil(t, err)
	assert.Equal(t, userID, users[0].UserID)
	assert.Equal(t, name, users[0].Name)
	assert.Equal(t, mailAddress, users[0].MailAddress)
	assert.Equal(t, roleType, users[0].RoleType)
}

func TestUserService_Add_addUser(t *testing.T) {
	userID := "admin"
	password := "password"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	m := &backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "users", spath)
			assert.Equal(t, userID, params.Get("userId"))
			assert.Equal(t, password, params.Get("password"))
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, mailAddress, params.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(int(roleType)), params.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
	user, err := s.Add(userID, password, name, mailAddress, roleType)
	assert.Nil(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestProjectUserService_Delete_deleteUser(t *testing.T) {
	userID := "admin"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator

	projectIDOrKey := "TEST"
	id := 1
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "projects/"+projectIDOrKey+"/users", spath)
			assert.Equal(t, strconv.Itoa(id), params.Get("userId"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewProjectUserService(m)
	users, err := s.Delete(projectIDOrKey, id)
	assert.Nil(t, err)
	assert.Equal(t, userID, users.UserID)
	assert.Equal(t, name, users.Name)
	assert.Equal(t, mailAddress, users.MailAddress)
	assert.Equal(t, roleType, users.RoleType)
}

func TestUserService_Update_updateUser(t *testing.T) {
	id := 1
	userID := "admin"
	password := "password"
	name := "admin"
	mailAddress := "eguchi@nulab.example"
	roleType := backlog.RoleAdministrator
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	m := &backlog.ExportMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "users/"+strconv.Itoa(id), spath)
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, password, params.Get("password"))
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, mailAddress, params.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(int(roleType)), params.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
	o := s.Option
	user, err := s.Update(
		id, o.WithPassword(password), o.WithName(name), o.WithMailAddress(mailAddress), o.WithRoleType(roleType),
	)
	assert.Nil(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestUserService_All(t *testing.T) {
	want := struct {
		spath string
	}{
		spath: "users",
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewUserService(m)
	users, err := s.All()
	assert.Nil(t, users)
	assert.Error(t, err)
}

func TestUserService_All_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
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
			m := &backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(m)
			s.One(tc.id)
		})
	}
}

func TestUserService_Own(t *testing.T) {
	want := struct {
		spath string
	}{
		spath: "users/myself",
	}
	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewUserService(m)
	user, err := s.Own()
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Own_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
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
		roleType    backlog.ExportRole
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
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, wantSpath, spath)
						assert.Equal(t, tc.userID, params.Get("userId"))
						assert.Equal(t, tc.password, params.Get("password"))
						assert.Equal(t, tc.name, params.Get("name"))
						assert.Equal(t, tc.mailAddress, params.Get("mailAddress"))
						assert.Equal(t, strconv.Itoa(int(tc.roleType)), params.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(m)
			user, err := s.Add(tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Add_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
	user, err := s.Add("userid", "password", "name", "mailAdress", 1)
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
			m := &backlog.ExportMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Patch must never be called")
					} else {
						assert.Equal(t, "users/"+strconv.Itoa(tc.id), spath)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(m)

			user, err := s.Update(tc.id)
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Update_option(t *testing.T) {
	o := &backlog.UserOptionService{}
	id := 1

	type options struct {
		password    string
		name        string
		mailAddress string
		roleType    string
	}
	cases := map[string]struct {
		options   []backlog.UserOption
		wantError bool
		want      options
	}{
		"no-option": {
			options:   []backlog.UserOption{},
			wantError: false,
			want:      options{},
		},
		"option-password": {
			options: []backlog.UserOption{
				o.WithPassword("testpasword"),
			},
			wantError: false,
			want: options{
				password: "testpasword",
			},
		},
		"option-password_empty": {
			options: []backlog.UserOption{
				o.WithPassword(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-name": {
			options: []backlog.UserOption{
				o.WithName("testname"),
			},
			wantError: false,
			want: options{
				name: "testname",
			},
		},
		"option-name_empty": {
			options: []backlog.UserOption{
				o.WithName(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-mailAddress": {
			options: []backlog.UserOption{
				o.WithMailAddress("test@test.com"),
			},
			wantError: false,
			want: options{
				mailAddress: "test@test.com",
			},
		},
		"option-mailAddress_empty": {
			options: []backlog.UserOption{
				o.WithMailAddress(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-roleType_Administrator": {
			options: []backlog.UserOption{
				o.WithRoleType(backlog.RoleAdministrator),
			},
			wantError: false,
			want: options{
				roleType: "1",
			},
		},
		"multi-option": {
			options: []backlog.UserOption{
				o.WithPassword("testpasword1"),
				o.WithName("testname1"),
				o.WithMailAddress("test1@test.com"),
				o.WithRoleType(backlog.RoleAdministrator),
			},
			wantError: false,
			want: options{
				password:    "testpasword1",
				name:        "testname1",
				mailAddress: "test1@test.com",
				roleType:    "1",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Patch must never be called")
					} else {
						assert.Equal(t, "users/"+strconv.Itoa(id), spath)
						assert.Equal(t, tc.want.password, params.Get("password"))
						assert.Equal(t, tc.want.name, params.Get("name"))
						assert.Equal(t, tc.want.mailAddress, params.Get("mailAddress"))
						assert.Equal(t, tc.want.roleType, params.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(m)

			user, err := s.Update(id, tc.options...)
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Update_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
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
			m := &backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(m)
			user, err := s.Delete(tc.id)
			assert.Nil(t, user)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Delete_invaliedJson(t *testing.T) {
	bj, err := os.Open("testdata/json/invalied.json")
	if err != nil {
		t.Fatal(err)
	}
	defer bj.Close()

	m := &backlog.ExportMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(m)
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
		projectIDOrKey      string
		excludeGroupMembers bool
		wantError           bool
		want                want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey:      "TEST",
			excludeGroupMembers: false,
			wantError:           false,
			want: want{
				spath:               "projects/TEST/users",
				excludeGroupMembers: "false",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey:      "1234",
			excludeGroupMembers: false,
			wantError:           false,
			want: want{
				spath:               "projects/1234/users",
				excludeGroupMembers: "false",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey:      "",
			excludeGroupMembers: false,
			wantError:           true,
			want:                want{},
		},
		"excludeGroupMembers_true": {
			projectIDOrKey:      "TEST2",
			excludeGroupMembers: true,
			wantError:           false,
			want: want{
				spath:               "projects/TEST2/users",
				excludeGroupMembers: "true",
			},
		},
		"excludeGroupMembers_false": {
			projectIDOrKey:      "TEST3",
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
			m := &backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.excludeGroupMembers, params.Get("excludeGroupMembers"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.All(tc.projectIDOrKey, tc.excludeGroupMembers)
		})
	}
}

func TestProjectUserService_Add(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectIDOrKey string
		userID         int
		wantError      bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/1234/users",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			wantError:      true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			wantError:      true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST2/users",
				userID: "1",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.Add(tc.projectIDOrKey, tc.userID)
		})
	}
}
func TestProjectUserService_Delete(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectIDOrKey string
		userID         int
		wantError      bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/1234/users",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			wantError:      true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			wantError:      true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST2/users",
				userID: "1",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.Delete(tc.projectIDOrKey, tc.userID)
		})
	}
}
func TestProjectUserService_AddAdmin(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectIDOrKey string
		userID         int
		wantError      bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/1234/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			wantError:      true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			wantError:      true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST2/administrators",
				userID: "1",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.AddAdmin(tc.projectIDOrKey, tc.userID)
		})
	}
}
func TestProjectUserService_AdminAll(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		projectIDOrKey string
		wantError      bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			wantError:      false,
			want: want{
				spath: "projects/TEST/administrators",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			wantError:      false,
			want: want{
				spath: "projects/1234/administrators",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			wantError:      true,
			want:           want{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.AdminAll(tc.projectIDOrKey)
		})
	}
}
func TestProjectUserService_DeleteAdmin(t *testing.T) {
	type want struct {
		spath  string
		userID string
	}
	cases := map[string]struct {
		projectIDOrKey string
		userID         int
		wantError      bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/1234/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			wantError:      true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			wantError:      true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			wantError:      false,
			want: want{
				spath:  "projects/TEST2/administrators",
				userID: "1",
			},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			m := &backlog.ExportMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("s.method.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(m)
			s.DeleteAdmin(tc.projectIDOrKey, tc.userID)
		})
	}
}
