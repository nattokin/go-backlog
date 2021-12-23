package backlog

import (
	"encoding/json"
	"path"
	"strconv"
)

// UserID is ID of user.
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

func getUser(get clientGet, spath string) (*User, error) {
	resp, err := get(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := User{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func getUserList(get clientGet, spath string, query *QueryParams) ([]*User, error) {
	resp, err := get(spath, query)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*User{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

func addUser(post clientPost, spath string, form *FormParams) (*User, error) {
	resp, err := post(spath, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := User{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func updateUser(patch clientPatch, spath string, form *FormParams) (*User, error) {
	resp, err := patch(spath, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := User{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

func deleteUser(delete clientDelete, spath string, form *FormParams) (*User, error) {
	resp, err := delete(spath, form)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := User{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
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
	return getUserList(s.method.Get, "users", nil)
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
	return getUser(s.method.Get, spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own() (*User, error) {
	return getUser(s.method.Get, "users/myself")
}

// ToDo: func (s *UserService) Icon()
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(userID, password, name, mailAddress string, roleType Role) (*User, error) {
	if userID == "" {
		return nil, newValidationError("userID must not be empty")
	}

	form := NewFormParams()
	if err := withFormPassword(password).set(form); err != nil {
		return nil, err
	}
	if err := withFormName(name).set(form); err != nil {
		return nil, err
	}
	if err := withFormMailAddress(mailAddress).set(form); err != nil {
		return nil, err
	}
	if err := withFormRoleType(roleType).set(form); err != nil {
		return nil, err
	}

	form.Set("userId", userID)

	return addUser(s.method.Post, "users", form)
}

// Update updates a user in your space.
//
// This method can specify the options returned by methods in "*Client.User.Option".
//
// Use the following methods:
//   WithFormName
//   WithFormPassword
//   WithFormMailAddress
//   WithFormRoleType
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(id int, options ...*FormOption) (*User, error) {
	uID := UserID(id)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	validOptions := []formType{formName, formPassword, formMailAddress, formRoleType}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	form := NewFormParams()
	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}

	spath := path.Join("users", uID.String())

	return updateUser(s.method.Patch, spath, form)
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
	return deleteUser(s.method.Delete, spath, nil)
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

	query := NewQueryParams()
	query.Set("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := path.Join("projects", projectIDOrKey, "users")
	return getUserList(s.method.Get, spath, query)
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

	form := NewFormParams()
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "users")
	return addUser(s.method.Post, spath, form)
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

	form := NewFormParams()
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "users")
	return deleteUser(s.method.Delete, spath, form)
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

	form := NewFormParams()
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return addUser(s.method.Post, spath, form)
}

// AdminAll returns all of admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(projectIDOrKey string) ([]*User, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getUserList(s.method.Get, spath, nil)
}

// DeleteAdmin deletes a admin user from the project.
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

	form := NewFormParams()
	form.Set("userId", uID.String())

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return deleteUser(s.method.Delete, spath, form)
}
