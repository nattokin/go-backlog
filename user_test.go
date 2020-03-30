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
	roleType := 1
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	cm := &backlog.ExportClientMethod{
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
	s := backlog.ExportNewUserService(cm)
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
	roleType := 1

	projectIDOrKey := "TEST"
	excludeGroupMembers := false
	bj, err := os.Open("testdata/json/user_list.json")
	if err != nil {
		t.Fatal(err)
	}
	cm := &backlog.ExportClientMethod{
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
	s := backlog.ExportNewProjectUserService(cm)
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
	roleType := 1
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "users", spath)
			assert.Equal(t, userID, params.Get("userId"))
			assert.Equal(t, password, params.Get("password"))
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, mailAddress, params.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(roleType), params.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
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
	roleType := 1

	projectIDOrKey := "TEST"
	id := 1
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	cm := &backlog.ExportClientMethod{
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
	s := backlog.ExportNewProjectUserService(cm)
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
	roleType := 1
	bj, err := os.Open("testdata/json/user.json")
	if err != nil {
		t.Fatal(err)
	}
	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, "users/"+strconv.Itoa(id), spath)
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, password, params.Get("password"))
			assert.Equal(t, name, params.Get("name"))
			assert.Equal(t, mailAddress, params.Get("mailAddress"))
			assert.Equal(t, strconv.Itoa(roleType), params.Get("roleType"))
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
	user, err := s.Update(
		id, s.Option.WithPassword(password), s.Option.WithName(name), s.Option.WithMailAddress(mailAddress), s.Option.WithRoleType(roleType),
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
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewUserService(cm)
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

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(cm)
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
	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			assert.Equal(t, want.spath, spath)
			assert.Nil(t, params)
			return nil, errors.New("error")
		},
	}
	s := backlog.ExportNewUserService(cm)
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

	cm := &backlog.ExportClientMethod{
		Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
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
		roleType    int
		wantError   bool
	}{
		"no_error": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    2,
			wantError:   false,
		},
		"userID_empty": {
			userID:      "",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			wantError:   true,
		},
		"password_empty": {
			userID:      "testid",
			password:    "",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			wantError:   true,
		},
		"name_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "",
			mailAddress: "test@test.com",
			roleType:    1,
			wantError:   true,
		},
		"mailAddress_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "",
			roleType:    1,
			wantError:   true,
		},
		"roleType_0": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    0,
			wantError:   true,
		},
		"roleType_1": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			wantError:   false,
		},
		"roleType_6": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    6,
			wantError:   false,
		},
		"roleType_7": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    7,
			wantError:   true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			cm := &backlog.ExportClientMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Post must never be called")
					} else {
						assert.Equal(t, wantSpath, spath)
						assert.Equal(t, tc.userID, params.Get("userId"))
						assert.Equal(t, tc.password, params.Get("password"))
						assert.Equal(t, tc.name, params.Get("name"))
						assert.Equal(t, tc.mailAddress, params.Get("mailAddress"))
						assert.Equal(t, strconv.Itoa(tc.roleType), params.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(cm)
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

	cm := &backlog.ExportClientMethod{
		Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
	user, err := s.Add("userid", "password", "name", "mailAdress", 1)
	assert.Nil(t, user)
	assert.Error(t, err)
}

func TestUserService_Update(t *testing.T) {
	ops := &backlog.UserOptionService{}
	type options struct {
		password    string
		name        string
		mailAddress string
		roleType    string
	}
	cases := map[string]struct {
		id        int
		options   []backlog.UserOption
		wantError bool
		want      options
	}{
		"no-option": {
			id:        1,
			options:   []backlog.UserOption{},
			wantError: false,
			want:      options{},
		},
		"id_0": {
			id:        0,
			options:   []backlog.UserOption{},
			wantError: true,
			want:      options{},
		},
		"option-password": {
			id: 2,
			options: []backlog.UserOption{
				ops.WithPassword("testpasword"),
			},
			wantError: false,
			want: options{
				password: "testpasword",
			},
		},
		"option-password_empty": {
			id: 3,
			options: []backlog.UserOption{
				ops.WithPassword(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-name": {
			id: 4,
			options: []backlog.UserOption{
				ops.WithName("testname"),
			},
			wantError: false,
			want: options{
				name: "testname",
			},
		},
		"option-name_empty": {
			id: 5,
			options: []backlog.UserOption{
				ops.WithName(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-mailAddress": {
			id: 6,
			options: []backlog.UserOption{
				ops.WithMailAddress("test@test.com"),
			},
			wantError: false,
			want: options{
				mailAddress: "test@test.com",
			},
		},
		"option-mailAddress_empty": {
			id: 7,
			options: []backlog.UserOption{
				ops.WithMailAddress(""),
			},
			wantError: true,
			want:      options{},
		},
		"option-roleType_0": {
			id: 8,
			options: []backlog.UserOption{
				ops.WithRoleType(0),
			},
			wantError: true,
			want:      options{},
		},
		"option-roleType_1": {
			id: 9,
			options: []backlog.UserOption{
				ops.WithRoleType(1),
			},
			wantError: false,
			want: options{
				roleType: "1",
			},
		},
		"option-roleType_6": {
			id: 10,
			options: []backlog.UserOption{
				ops.WithRoleType(6),
			},
			wantError: false,
			want: options{
				roleType: "6",
			},
		},
		"option-roleType_7": {
			id: 11,
			options: []backlog.UserOption{
				ops.WithRoleType(7),
			},
			wantError: true,
			want:      options{},
		},
		"multi-option": {
			id: 1,
			options: []backlog.UserOption{
				ops.WithPassword("testpasword1"),
				ops.WithName("testname1"),
				ops.WithMailAddress("test1@test.com"),
				ops.WithRoleType(1),
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
			cm := &backlog.ExportClientMethod{
				Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Patch must never be called")
					} else {
						assert.Equal(t, "users/"+strconv.Itoa(tc.id), spath)
						assert.Equal(t, tc.want.password, params.Get("password"))
						assert.Equal(t, tc.want.name, params.Get("name"))
						assert.Equal(t, tc.want.mailAddress, params.Get("mailAddress"))
						assert.Equal(t, tc.want.roleType, params.Get("roleType"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(cm)

			user, err := s.Update(tc.id, tc.options...)
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

	cm := &backlog.ExportClientMethod{
		Patch: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewUserService(cm)
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

	cm := &backlog.ExportClientMethod{
		Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
			resp := &http.Response{
				StatusCode: http.StatusOK,
				Body:       bj,
			}
			return backlog.ExportNewResponse(resp), nil
		},
	}
	s := backlog.ExportNewUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.excludeGroupMembers, params.Get("excludeGroupMembers"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
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
			cm := &backlog.ExportClientMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.wantError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			s := backlog.ExportNewProjectUserService(cm)
			s.DeleteAdmin(tc.projectIDOrKey, tc.userID)
		})
	}
}
