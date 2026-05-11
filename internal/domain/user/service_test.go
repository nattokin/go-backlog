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
	"github.com/nattokin/go-backlog/internal/domain/user"
	"github.com/nattokin/go-backlog/internal/testutil/fixture"
	"github.com/nattokin/go-backlog/internal/testutil/mock"
)

func TestUserService_One(t *testing.T) {
	cases := map[string]struct {
		id          int
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErrType error
	}{
		"success-id-1": {
			id: 1,

			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
		},
		"success-id-100": {
			id: 100,

			mockGetFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
		},
		"error-validation-id-zero": {
			id:          0,
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
			s := user.NewService(method)

			user, err := s.One(context.Background(), tc.id)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, 1, user.RoleType)
		})
	}
}

func TestUserService_Add(t *testing.T) {
	cases := map[string]struct {
		userID      string
		password    string
		name        string
		mailAddress string
		roleType    int
		mockPostFn  func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErrType error
	}{
		"success-roleType-1": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Equal(t, "admin", form.Get("userId"))
				assert.Equal(t, "password", form.Get("password"))
				assert.Equal(t, "admin", form.Get("name"))
				assert.Equal(t, "eguchi@nulab.example", form.Get("mailAddress"))
				assert.Equal(t, "1", form.Get("roleType"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
		},
		"success-roleType-6": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "eguchi@nulab.example",
			roleType:    6,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "6", form.Get("roleType"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
		},
		"error-validation-userID-empty": {
			userID:      "",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    1,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-password-empty": {
			userID:      "admin",
			password:    "",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    1,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-name-empty": {
			userID:      "admin",
			password:    "password",
			name:        "",
			mailAddress: "admin@example.com",
			roleType:    1,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-mailAddress-empty": {
			userID:      "admin",
			password:    "password",
			name:        "admin",
			mailAddress: "",
			roleType:    1,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-multiple-empty": {
			userID:      "test",
			password:    "",
			name:        "",
			mailAddress: "",
			roleType:    1,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-roleType-0": {
			userID:      "test",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    0,
			wantErrType: &core.ValidationError{},
		},
		"error-validation-roleType-7": {
			userID:      "test",
			password:    "password",
			name:        "admin",
			mailAddress: "admin@example.com",
			roleType:    7,
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			userID:      "errorUser",
			password:    "password",
			name:        "error",
			mailAddress: "error@example.com",
			roleType:    1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				return nil, errors.New("network failure")
			},
			wantErrType: errors.New(""),
		},
		"error-response-invalid-json": {
			userID:      "userID",
			password:    "password",
			name:        "name",
			mailAddress: "mailAdress",
			roleType:    1,

			mockPostFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
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
			s := user.NewService(method)

			user, err := s.Add(context.Background(), tc.userID, tc.password, tc.name, tc.mailAddress, tc.roleType)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, 1, user.RoleType)
		})
	}
}

func TestUserService_All(t *testing.T) {
	cases := map[string]struct {
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantLen     int
		wantErrType error
	}{
		"success-get-users": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.User.ListJSON), nil
			},
			wantLen: 4,
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
			s := user.NewService(method)

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
			assert.Equal(t, "admin", users[0].UserID)
			assert.Equal(t, "admin", users[0].Name)
			assert.Equal(t, "eguchi@nulab.example", users[0].MailAddress)
			assert.Equal(t, 1, users[0].RoleType)
		})
	}
}

func TestUserService_Update(t *testing.T) {
	o := &core.OptionService{}

	cases := map[string]struct {
		id          int
		option      core.RequestOption
		opts        []core.RequestOption
		mockPatchFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErrType error
	}{
		"success-update-user": {
			id:     1,
			option: o.WithPassword("password"),
			opts: []core.RequestOption{
				o.WithName("admin"),
				o.WithMailAddress("eguchi@nulab.example"),
				o.WithRoleType(1),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "password", form.Get("password"))
				assert.Equal(t, "admin", form.Get("name"))
				assert.Equal(t, "eguchi@nulab.example", form.Get("mailAddress"))
				assert.Equal(t, "1", form.Get("roleType"))
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
		},
		"error-response-invalid-json": {
			id:     1234,
			option: o.WithName("test"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)
				return mock.NewJSONResponse(fixture.InvalidJSON), nil
			},
			wantErrType: &json.SyntaxError{},
		},
		"success-option-withName": {
			id:     1,
			option: o.WithName("testname"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testname", form.Get("name"))
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"success-option-withPassword": {
			id:     1,
			option: o.WithPassword("testpassword"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword", form.Get("password"))
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"success-option-withMailAddress": {
			id:     1,
			option: o.WithMailAddress("test@test.com"),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "test@test.com", form.Get("mailAddress"))
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"success-option-withRoleType": {
			id:     1,
			option: o.WithRoleType(1),

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "1", form.Get("roleType"))
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"success-option-multiple": {
			id:     1,
			option: o.WithPassword("testpassword1"),
			opts: []core.RequestOption{
				o.WithName("testname1"),
				o.WithMailAddress("test1@test.com"),
				o.WithRoleType(1),
			},

			mockPatchFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Equal(t, "testpassword1", form.Get("password"))
				assert.Equal(t, "testname1", form.Get("name"))
				assert.Equal(t, "test1@test.com", form.Get("mailAddress"))
				assert.Equal(t, "1", form.Get("roleType"))
				return nil, errors.New("error")
			},
			wantErrType: errors.New(""),
		},
		"error-option-invalid-value": {
			id:          1,
			option:      o.WithName(""),
			wantErrType: &core.ValidationError{},
		},
		"error-option-invalid-type": {
			id:          1,
			option:      mock.NewInvalidTypeOption(),
			wantErrType: &core.InvalidOptionKeyError{},
		},
		"error-option-set-faild": {
			id:          1,
			option:      mock.NewFailingSetOption(core.ParamName),
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockPatchFn != nil {
				method.Patch = tc.mockPatchFn
			}
			s := user.NewService(method)

			user, err := s.Update(context.Background(), tc.id, tc.option, tc.opts...)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, user)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, user)

			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, 1, user.RoleType)
		})
	}
}

func TestUserService_Own(t *testing.T) {
	cases := map[string]struct {
		mockGetFn   func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErrType error
	}{
		"success-get-own-user": {
			mockGetFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/myself", spath)
				assert.Nil(t, query)
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
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
			s := user.NewService(method)

			user, err := s.Own(context.Background())

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.IsType(t, tc.wantErrType, err)
				assert.Nil(t, user)
				return
			}

			assert.NoError(t, err)
			require.NotNil(t, user)
			assert.Equal(t, "admin", user.UserID)
			assert.Equal(t, "admin", user.Name)
			assert.Equal(t, "eguchi@nulab.example", user.MailAddress)
			assert.Equal(t, 1, user.RoleType)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	cases := map[string]struct {
		id           int
		mockDeleteFn func(ctx context.Context, spath string, form url.Values) (*http.Response, error)
		wantErrType  error
	}{
		"success-id-1": {
			id: 1,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1", spath)
				assert.Nil(t, form)
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"success-id-100": {
			id: 100,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100", spath)
				assert.Nil(t, form)
				return mock.NewJSONResponse(fixture.User.SingleJSON), nil
			},
			wantErrType: nil,
		},
		"error-validation-id-zero": {
			id:          0,
			wantErrType: &core.ValidationError{},
		},
		"error-response-invalid-json": {
			id: 1234,

			mockDeleteFn: func(ctx context.Context, spath string, form url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1234", spath)
				assert.Nil(t, form)
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
			s := user.NewService(method)

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

func TestUserService_Icon(t *testing.T) {
	cases := map[string]struct {
		id              int
		mockDownloadFn  func(ctx context.Context, spath string, query url.Values) (*http.Response, error)
		wantErrType     error
		wantFilename    string
		wantContentType string
	}{
		"success-id-1": {
			id: 1,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/1/icon", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("avatar.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "avatar.png",
			wantContentType: "image/png",
		},
		"success-id-100": {
			id: 100,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				assert.Equal(t, "users/100/icon", spath)
				assert.Nil(t, query)
				return mock.NewBinaryResponse("avatar.png", "image/png", []byte("PNG")), nil
			},
			wantFilename:    "avatar.png",
			wantContentType: "image/png",
		},
		"error-validation-id-zero": {
			id:          0,
			wantErrType: &core.ValidationError{},
		},
		"error-client-network": {
			id: 1,
			mockDownloadFn: func(ctx context.Context, spath string, query url.Values) (*http.Response, error) {
				return nil, errors.New("network error")
			},
			wantErrType: errors.New(""),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			method := mock.NewMethod(t)
			if tc.mockDownloadFn != nil {
				method.Download = tc.mockDownloadFn
			}
			s := user.NewService(method)

			got, err := s.Icon(context.Background(), tc.id)

			if tc.wantErrType != nil {
				assert.Error(t, err)
				assert.Nil(t, got)
				assert.IsType(t, tc.wantErrType, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tc.wantFilename, got.Filename)
			assert.Equal(t, tc.wantContentType, got.ContentType)
			require.NotNil(t, got.Body)
			got.Body.Close()
		})
	}
}

func Test_contextPropagation(t *testing.T) {
	type ctxKey struct{}
	sentinel := &struct{}{}
	ctx := context.WithValue(context.Background(), ctxKey{}, sentinel)

	makeMockFn := func(t *testing.T) func(context.Context, string, url.Values) (*http.Response, error) {
		return func(got context.Context, _ string, _ url.Values) (*http.Response, error) {
			assert.Same(t, sentinel, got.Value(ctxKey{}))
			return nil, errors.New("stop")
		}
	}

	o := &core.OptionService{}

	cases := []struct {
		name string
		call func(t *testing.T, m *core.Method)
	}{
		{"Service.All", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewService(m)
			s.All(ctx) //nolint:errcheck
		}},
		{"Service.One", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewService(m)
			s.One(ctx, 1) //nolint:errcheck
		}},
		{"Service.Own", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewService(m)
			s.Own(ctx) //nolint:errcheck
		}},
		{"Service.Add", func(t *testing.T, m *core.Method) {
			m.Post = makeMockFn(t)
			s := user.NewService(m)
			s.Add(ctx, "u", "p", "n", "m@m.com", 1) //nolint:errcheck
		}},
		{"Service.Update", func(t *testing.T, m *core.Method) {
			m.Patch = makeMockFn(t)
			s := user.NewService(m)
			s.Update(ctx, 1, o.WithName("n")) //nolint:errcheck
		}},
		{"Service.Delete", func(t *testing.T, m *core.Method) {
			m.Delete = makeMockFn(t)
			s := user.NewService(m)
			s.Delete(ctx, 1) //nolint:errcheck
		}},
		{"Service.Icon", func(t *testing.T, m *core.Method) {
			m.Download = makeMockFn(t)
			s := user.NewService(m)
			s.Icon(ctx, 1) //nolint:errcheck
		}},
		{"ActivityService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewActivityService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"StarService.List", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewStarService(m)
			s.List(ctx, 1) //nolint:errcheck
		}},
		{"StarService.Count", func(t *testing.T, m *core.Method) {
			m.Get = makeMockFn(t)
			s := user.NewStarService(m)
			s.Count(ctx, 1) //nolint:errcheck
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			tc.call(t, &core.Method{})
		})
	}
}
