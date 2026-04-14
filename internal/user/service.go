package user

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

func getUser(ctx context.Context, m *core.Method, spath string) (*model.User, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func getUserList(ctx context.Context, m *core.Method, spath string, query url.Values) ([]*model.User, error) {
	resp, err := m.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func addUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func updateUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func deleteUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UserService has methods for user
type UserService struct {
	method *core.Method

	Activity *activity.UserActivityService
	Option   *UserOptionService
}

// All returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) All(ctx context.Context) ([]*model.User, error) {
	return getUserList(ctx, s.method, "users", nil)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return getUser(ctx, s.method, spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own(ctx context.Context) (*model.User, error) {
	return getUser(ctx, s.method, "users/myself")
}

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(ctx context.Context, userID, password, name, mailAddress string, roleType model.Role) (*model.User, error) {
	if userID == "" {
		return nil, core.NewValidationError("userID must not be empty")
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamPassword, core.ParamName, core.ParamMailAddress, core.ParamRoleType}
	options := []core.RequestOption{
		s.Option.base.WithPassword(password),
		s.Option.base.WithName(name),
		s.Option.base.WithMailAddress(mailAddress),
		s.Option.base.WithRoleType(roleType),
	}
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	form.Set("userId", userID)

	return addUser(ctx, s.method, "users", form)
}

// Update updates a user in your space.
//
// This method supports options returned by methods in "*Client.User.Option",
// such as:
//   - WithMailAddress
//   - WithName
//   - WithPassword
//   - Withmodel.RoleType
//   - WithUserID
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(ctx context.Context, id int, opts ...core.RequestOption) (*model.User, error) {

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamUserID, core.ParamName, core.ParamPassword, core.ParamMailAddress, core.ParamRoleType}
	options := append([]core.RequestOption{s.Option.base.WithUserID(id)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return updateUser(ctx, s.method, spath, form)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return deleteUser(ctx, s.method, spath, nil)
}

type ProjectUserService struct {
	method *core.Method
}

func (s *ProjectUserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := path.Join("projects", projectIDOrKey, "users")
	return getUserList(ctx, s.method, spath, query)
}

func (s *ProjectUserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "users")
	return addUser(ctx, s.method, spath, form)
}

func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "users")
	return deleteUser(ctx, s.method, spath, form)
}

func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return addUser(ctx, s.method, spath, form)
}

func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getUserList(ctx, s.method, spath, nil)
}

func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return deleteUser(ctx, s.method, spath, form)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewUserService(method *core.Method, option *core.OptionService) *UserService {
	return &UserService{
		method:   method,
		Activity: activity.NewUserActivityService(method, option),
		Option:   NewUserOptionService(option),
	}
}

func NewProjectUserService(method *core.Method, option *core.OptionService) *ProjectUserService {
	return &ProjectUserService{
		method: method,
	}
}
