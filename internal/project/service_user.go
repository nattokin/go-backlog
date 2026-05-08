package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

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

func addUser(ctx context.Context, m *core.Method, spath string, userID int) (*model.User, error) {
	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

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

func deleteUser(ctx context.Context, m *core.Method, spath string, userID int) (*model.User, error) {
	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

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

// UserService handles communication with the project user-related methods of the Backlog API.
type UserService struct {
	method *core.Method
}

// All returns a list of users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *UserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := path.Join("projects", projectIDOrKey, "users")
	return getUserList(ctx, s.method, spath, query)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *UserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "users")
	return addUser(ctx, s.method, spath, userID)
}

// Delete removes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *UserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "users")
	return deleteUser(ctx, s.method, spath, userID)
}

// AddAdmin adds a user as an administrator of the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *UserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return addUser(ctx, s.method, spath, userID)
}

// AdminAll returns a list of project administrators.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *UserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getUserList(ctx, s.method, spath, nil)
}

// DeleteAdmin removes a user from the project administrators.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *UserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return deleteUser(ctx, s.method, spath, userID)
}

// NewUserService creates and returns a new project UserService.
func NewUserService(method *core.Method) *UserService {
	return &UserService{method: method}
}
