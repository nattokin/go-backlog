package backlog

import (
	"net/url"
	"path"
	"strconv"
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

func getUser(m *method, spath string) (*User, error) {
	resp, err := m.Get(spath, nil)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func getUserList(m *method, spath string, query url.Values) ([]*User, error) {
	resp, err := m.Get(spath, query)
	if err != nil {
		return nil, err
	}

	v := []*User{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func addUser(m *method, spath string, form url.Values) (*User, error) {
	resp, err := m.Post(spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func updateUser(m *method, spath string, form url.Values) (*User, error) {
	resp, err := m.Patch(spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func deleteUser(m *method, spath string, form url.Values) (*User, error) {
	resp, err := m.Delete(spath, form)
	if err != nil {
		return nil, err
	}

	v := User{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UserService has methods for user
type UserService struct {
	method *method

	Activity *UserActivityService
	Option   *UserOptionService
}

// All returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) All() ([]*User, error) {
	return getUserList(s.method, "users", nil)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(id int) (*User, error) {
	uID := UserID(id)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("users", uID.String())
	return getUser(s.method, spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own() (*User, error) {
	return getUser(s.method, "users/myself")
}

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(userID, password, name, mailAddress string, roleType Role) (*User, error) {
	if userID == "" {
		return nil, newValidationError("userID must not be empty")
	}

	o := s.Option.registry.option
	form := url.Values{}
	for _, opt := range []RequestOption{
		o.WithPassword(password),
		o.WithName(name),
		o.WithMailAddress(mailAddress),
		o.WithRoleType(roleType),
	} {
		if err := opt.Check(); err != nil {
			return nil, err
		}
		if err := opt.Set(form); err != nil {
			return nil, err
		}
	}

	form.Set("userId", userID)

	return addUser(s.method, "users", form)
}

// Update updates a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(id int, opts ...RequestOption) (*User, error) {
	o := s.Option.registry.option
	form := url.Values{}
	userIDOpt := o.WithUserID(id)
	if err := userIDOpt.Check(); err != nil {
		return nil, err
	}
	if err := userIDOpt.Set(form); err != nil {
		return nil, err
	}

	validTypes := []formType{formName, formPassword, formMailAddress, formRoleType}
	if err := o.applyFormOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return updateUser(s.method, spath, form)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(id int) (*User, error) {
	uID := UserID(id)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("users", uID.String())
	return deleteUser(s.method, spath, nil)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	method *method
}

// All returns all users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *ProjectUserService) All(projectIDOrKey string, excludeGroupMembers bool) ([]*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := path.Join("projects", projectIDOrKey, "users")
	return getUserList(s.method, spath, query)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *ProjectUserService) Add(projectIDOrKey string, userID int) (*User, error) {
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
	return addUser(s.method, spath, form)
}

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(projectIDOrKey string, userID int) (*User, error) {
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
	return deleteUser(s.method, spath, form)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(projectIDOrKey string, userID int) (*User, error) {
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
	return addUser(s.method, spath, form)
}

// AdminAll returns a list of all admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(projectIDOrKey string) ([]*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getUserList(s.method, spath, nil)
}

// DeleteAdmin removes an admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(projectIDOrKey string, userID int) (*User, error) {
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
	return deleteUser(s.method, spath, form)
}
