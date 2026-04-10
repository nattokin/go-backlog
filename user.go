package backlog

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
)

// UserID is the unique identifier for a user.
type UserID int

func (id UserID) validate() error {
	if id < 1 {
		return newValidationError("userID must not be less than 1")
	}
	return nil
}

func (id UserID) String() string {
	return strconv.Itoa(int(id))
}

func getUser(ctx context.Context, m *core.Method, spath string) (*User, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func getUserList(ctx context.Context, m *core.Method, spath string, query url.Values) ([]*User, error) {
	resp, err := m.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func addUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*User, error) {
	resp, err := m.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func updateUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*User, error) {
	resp, err := m.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func deleteUser(ctx context.Context, m *core.Method, spath string, form url.Values) (*User, error) {
	resp, err := m.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UserService has methods for user
type UserService struct {
	method *core.Method

	Activity *UserActivityService
	Option   *UserOptionService
}

// All returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) All(ctx context.Context) ([]*User, error) {
	return getUserList(ctx, s.method, "users", nil)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(ctx context.Context, id int) (*User, error) {
	uID := UserID(id)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("users", uID.String())
	return getUser(ctx, s.method, spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own(ctx context.Context) (*User, error) {
	return getUser(ctx, s.method, "users/myself")
}

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(ctx context.Context, userID, password, name, mailAddress string, roleType Role) (*User, error) {
	if userID == "" {
		return nil, newValidationError("userID must not be empty")
	}

	form := url.Values{}
	validTypes := []apiParamOptionType{paramPassword, paramName, paramMailAddress, paramRoleType}
	options := []RequestOption{
		s.Option.base.WithPassword(password),
		s.Option.base.WithName(name),
		s.Option.base.WithMailAddress(mailAddress),
		s.Option.base.WithRoleType(roleType),
	}
	if err := applyOptions(form, validTypes, options...); err != nil {
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
//   - WithRoleType
//   - WithUserID
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(ctx context.Context, id int, opts ...RequestOption) (*User, error) {

	form := url.Values{}
	validTypes := []apiParamOptionType{paramUserID, paramName, paramPassword, paramMailAddress, paramRoleType}
	options := append([]RequestOption{s.Option.base.WithUserID(id)}, opts...)
	if err := applyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return updateUser(ctx, s.method, spath, form)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(ctx context.Context, id int) (*User, error) {
	uID := UserID(id)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("users", uID.String())
	return deleteUser(ctx, s.method, spath, nil)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	method *core.Method
}

// All returns all users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *ProjectUserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
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
func (s *ProjectUserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	uID := UserID(userID)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "users")
	return addUser(ctx, s.method, spath, form)
}

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	uID := UserID(userID)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "users")
	return deleteUser(ctx, s.method, spath, form)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	uID := UserID(userID)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return addUser(ctx, s.method, spath, form)
}

// AdminAll returns a list of all admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getUserList(ctx, s.method, spath, nil)
}

// DeleteAdmin removes an admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	uID := UserID(userID)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return deleteUser(ctx, s.method, spath, form)
}
