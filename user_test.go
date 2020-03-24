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

func TestBaseUserService_Get(t *testing.T) {
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
	us := backlog.ExportNewUserService(cm)
	user, err := us.One(1)
	assert.Nil(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestBaseUserService_GetList(t *testing.T) {
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
	pus := backlog.ExportNewProjectUserService(cm)
	users, err := pus.All(projectIDOrKey, excludeGroupMembers)
	assert.Nil(t, err)
	assert.Equal(t, userID, users[0].UserID)
	assert.Equal(t, name, users[0].Name)
	assert.Equal(t, mailAddress, users[0].MailAddress)
	assert.Equal(t, roleType, users[0].RoleType)
}

func TestBaseUserService_Post(t *testing.T) {
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
	us := backlog.ExportNewUserService(cm)
	user, err := us.Add(userID, password, name, mailAddress, roleType)
	assert.Nil(t, err)
	assert.Equal(t, userID, user.UserID)
	assert.Equal(t, name, user.Name)
	assert.Equal(t, mailAddress, user.MailAddress)
	assert.Equal(t, roleType, user.RoleType)
}

func TestBaseUserService_Delete(t *testing.T) {
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
	pus := backlog.ExportNewProjectUserService(cm)
	users, err := pus.Delete(projectIDOrKey, id)
	assert.Nil(t, err)
	assert.Equal(t, userID, users.UserID)
	assert.Equal(t, name, users.Name)
	assert.Equal(t, mailAddress, users.MailAddress)
	assert.Equal(t, roleType, users.RoleType)
}

func TestBaseUserService_Patch(t *testing.T) {
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
	us := backlog.ExportNewUserService(cm)
	user, err := us.Update(
		id, us.Option.WithPassword(password), us.Option.WithName(name), us.Option.WithMailAddress(mailAddress), us.Option.WithRoleType(roleType),
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
	us := backlog.ExportNewUserService(cm)
	us.All()
}

func TestUserService_One(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		id       int
		hasError bool
		want     want
	}{
		"id_1": {
			id:       1,
			hasError: false,
			want: want{
				spath: "users/1",
			},
		},
		"id_100": {
			id:       100,
			hasError: false,
			want: want{
				spath: "users/100",
			},
		},

		"id_0": {
			id:       0,
			hasError: true,
			want:     want{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.hasError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			us := backlog.ExportNewUserService(cm)
			us.One(tc.id)
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
	us := backlog.ExportNewUserService(cm)
	us.Own()
}

func TestUserService_Add(t *testing.T) {
	wantSpath := "users"
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    int
		hasError    bool
	}{
		"no_error": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    2,
			hasError:    false,
		},
		"userID_empty": {
			userID:      "",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			hasError:    true,
		},
		"password_empty": {
			userID:      "testid",
			password:    "",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			hasError:    true,
		},
		"name_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "",
			mailAddress: "test@test.com",
			roleType:    1,
			hasError:    true,
		},
		"mailAddress_empty": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "",
			roleType:    1,
			hasError:    true,
		},
		"roleType_0": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    0,
			hasError:    true,
		},
		"roleType_1": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    1,
			hasError:    false,
		},
		"roleType_6": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    6,
			hasError:    false,
		},
		"roleType_7": {
			userID:      "testid",
			password:    "testpass",
			name:        "testname",
			mailAddress: "test@test.com",
			roleType:    7,
			hasError:    true,
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			cm := &backlog.ExportClientMethod{
				Post: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.hasError {
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
			us := backlog.ExportNewUserService(cm)

			_, err := us.Add(tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	uos := &backlog.UserOptionService{}
	type options struct {
		password    string
		name        string
		mailAddress string
		roleType    string
	}
	cases := map[string]struct {
		id       int
		options  []backlog.UserOption
		hasError bool
		want     options
	}{
		"no-option": {
			id:       1,
			options:  []backlog.UserOption{},
			hasError: false,
			want:     options{},
		},
		"id_0": {
			id:       0,
			options:  []backlog.UserOption{},
			hasError: true,
			want:     options{},
		},
		"option-password": {
			id: 2,
			options: []backlog.UserOption{
				uos.WithPassword("testpasword"),
			},
			hasError: false,
			want: options{
				password: "testpasword",
			},
		},
		"option-password_empty": {
			id: 3,
			options: []backlog.UserOption{
				uos.WithPassword(""),
			},
			hasError: true,
			want:     options{},
		},
		"option-name": {
			id: 4,
			options: []backlog.UserOption{
				uos.WithName("testname"),
			},
			hasError: false,
			want: options{
				name: "testname",
			},
		},
		"option-name_empty": {
			id: 5,
			options: []backlog.UserOption{
				uos.WithName(""),
			},
			hasError: true,
			want:     options{},
		},
		"option-mailAddress": {
			id: 6,
			options: []backlog.UserOption{
				uos.WithMailAddress("test@test.com"),
			},
			hasError: false,
			want: options{
				mailAddress: "test@test.com",
			},
		},
		"option-mailAddress_empty": {
			id: 7,
			options: []backlog.UserOption{
				uos.WithMailAddress(""),
			},
			hasError: true,
			want:     options{},
		},
		"option-roleType_0": {
			id: 8,
			options: []backlog.UserOption{
				uos.WithRoleType(0),
			},
			hasError: true,
			want:     options{},
		},
		"option-roleType_1": {
			id: 9,
			options: []backlog.UserOption{
				uos.WithRoleType(1),
			},
			hasError: false,
			want: options{
				roleType: "1",
			},
		},
		"option-roleType_6": {
			id: 10,
			options: []backlog.UserOption{
				uos.WithRoleType(6),
			},
			hasError: false,
			want: options{
				roleType: "6",
			},
		},
		"option-roleType_7": {
			id: 11,
			options: []backlog.UserOption{
				uos.WithRoleType(7),
			},
			hasError: true,
			want:     options{},
		},
		"multi-option": {
			id: 1,
			options: []backlog.UserOption{
				uos.WithPassword("testpasword1"),
				uos.WithName("testname1"),
				uos.WithMailAddress("test1@test.com"),
				uos.WithRoleType(1),
			},
			hasError: false,
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
					if tc.hasError {
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
			us := backlog.ExportNewUserService(cm)

			_, err := us.Update(tc.id, tc.options...)
			assert.Error(t, err)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		id       int
		hasError bool
		want     want
	}{
		"id_1": {
			id:       1,
			hasError: false,
			want: want{
				spath: "users/1",
			},
		},
		"id_100": {
			id:       100,
			hasError: false,
			want: want{
				spath: "users/100",
			},
		},

		"id_0": {
			id:       0,
			hasError: true,
			want:     want{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			cm := &backlog.ExportClientMethod{
				Delete: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.hasError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			us := backlog.ExportNewUserService(cm)
			us.Delete(tc.id)
		})
	}
}

func TestProjectUserService_All(t *testing.T) {
	type want struct {
		spath               string
		excludeGroupMembers string
	}
	cases := map[string]struct {
		projectIDOrKey      string
		excludeGroupMembers bool
		hasError            bool
		want                want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey:      "TEST",
			excludeGroupMembers: false,
			hasError:            false,
			want: want{
				spath:               "projects/TEST/users",
				excludeGroupMembers: "false",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey:      "1234",
			excludeGroupMembers: false,
			hasError:            false,
			want: want{
				spath:               "projects/1234/users",
				excludeGroupMembers: "false",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey:      "",
			excludeGroupMembers: false,
			hasError:            true,
			want:                want{},
		},
		"excludeGroupMembers_true": {
			projectIDOrKey:      "TEST2",
			excludeGroupMembers: true,
			hasError:            false,
			want: want{
				spath:               "projects/TEST2/users",
				excludeGroupMembers: "true",
			},
		},
		"excludeGroupMembers_false": {
			projectIDOrKey:      "TEST3",
			excludeGroupMembers: false,
			hasError:            false,
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
					if tc.hasError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.excludeGroupMembers, params.Get("excludeGroupMembers"))
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.All(tc.projectIDOrKey, tc.excludeGroupMembers)
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
		hasError       bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/1234/users",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			hasError:       true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			hasError:       true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			hasError:       false,
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
					if tc.hasError {
						t.Error("clientMethod.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.Add(tc.projectIDOrKey, tc.userID)
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
		hasError       bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/TEST/users",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/1234/users",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			hasError:       true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			hasError:       true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			hasError:       false,
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
					if tc.hasError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.Delete(tc.projectIDOrKey, tc.userID)
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
		hasError       bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/1234/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			hasError:       true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			hasError:       true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			hasError:       false,
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
					if tc.hasError {
						t.Error("clientMethod.Post must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.AddAdmin(tc.projectIDOrKey, tc.userID)
		})
	}
}
func TestProjectUserService_AdminAll(t *testing.T) {
	type want struct {
		spath string
	}
	cases := map[string]struct {
		projectIDOrKey string
		hasError       bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			hasError:       false,
			want: want{
				spath: "projects/TEST/administrators",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			hasError:       false,
			want: want{
				spath: "projects/1234/administrators",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			hasError:       true,
			want:           want{},
		},
	}
	for n, tc := range cases {
		tc := tc
		t.Run(n, func(t *testing.T) {
			cm := &backlog.ExportClientMethod{
				Get: func(spath string, params *backlog.ExportRequestParams) (*backlog.ExportResponse, error) {
					if tc.hasError {
						t.Error("clientMethod.Get must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Nil(t, params)
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.AdminAll(tc.projectIDOrKey)
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
		hasError       bool
		want           want
	}{
		"projectIDOrKey_string": {
			projectIDOrKey: "TEST",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/TEST/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_number": {
			projectIDOrKey: "1234",
			userID:         1,
			hasError:       false,
			want: want{
				spath:  "projects/1234/administrators",
				userID: "1",
			},
		},
		"projectIDOrKey_empty": {
			projectIDOrKey: "",
			userID:         1,
			hasError:       true,
			want:           want{},
		},
		"userID_0": {
			projectIDOrKey: "TEST1",
			userID:         0,
			hasError:       true,
			want:           want{},
		},
		"userID_1": {
			projectIDOrKey: "TEST2",
			userID:         1,
			hasError:       false,
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
					if tc.hasError {
						t.Error("clientMethod.Delete must never be called")
					} else {
						assert.Equal(t, tc.want.spath, spath)
						assert.Equal(t, tc.want.userID, params.Get("userId"))
					}
					return nil, errors.New("error")
				},
			}
			pus := backlog.ExportNewProjectUserService(cm)
			pus.DeleteAdmin(tc.projectIDOrKey, tc.userID)
		})
	}
}
