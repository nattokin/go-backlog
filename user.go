package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

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

func getUserList(get clientGet, spath string, params *requestParams) ([]*User, error) {
	resp, err := get(spath, params)
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

func addUser(post clientPost, spath string, params *requestParams) (*User, error) {
	resp, err := post(spath, params)
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

func updateUser(patch clientPatch, spath string, params *requestParams) (*User, error) {
	resp, err := patch(spath, params)
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

func deleteUser(delete clientDelete, spath string, params *requestParams) (*User, error) {
	resp, err := delete(spath, params)
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
	spath := "users"
	return getUserList(s.method.Get, spath, nil)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(id int) (*User, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id)
	return getUser(s.method.Get, spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own() (*User, error) {
	spath := "users/myself"
	return getUser(s.method.Get, spath)
}

// ToDo: func (s *UserService) Icon()
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(userID, password, name, mailAddress string, roleType role) (*User, error) {
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	if password == "" {
		return nil, errors.New("password must not be empty")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}
	if mailAddress == "" {
		return nil, errors.New("mailAddress must not be empty")
	}

	params := newRequestParams()
	params.Add("userId", userID)
	params.Add("password", password)
	params.Add("name", name)
	params.Add("mailAddress", mailAddress)
	params.Add("roleType", strconv.Itoa(int(roleType)))

	spath := "users"
	return addUser(s.method.Post, spath, params)
}

// Update updates a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *UserService) Update(id int, options ...UserOption) (*User, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id)

	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	return updateUser(s.method.Patch, spath, params)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(id int) (*User, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id)
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
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}

	params := newRequestParams()
	params.Add("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := "projects/" + projectIDOrKey + "/users"
	return getUserList(s.method.Get, spath, params)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *ProjectUserService) Add(projectIDOrKey string, userID int) (*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}
	if userID < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	params := newRequestParams()
	params.Add("userId", strconv.Itoa(userID))

	spath := "projects/" + projectIDOrKey + "/users"
	return addUser(s.method.Post, spath, params)
}

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(projectIDOrKey string, userID int) (*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}
	if userID < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	params := newRequestParams()
	params.Add("userId", strconv.Itoa(userID))

	spath := "projects/" + projectIDOrKey + "/users"
	return deleteUser(s.method.Delete, spath, params)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(projectIDOrKey string, userID int) (*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}
	if userID < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	params := newRequestParams()
	params.Add("userId", strconv.Itoa(userID))

	spath := "projects/" + projectIDOrKey + "/administrators"
	return addUser(s.method.Post, spath, params)
}

// AdminAll returns all of admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(projectIDOrKey string) ([]*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}

	spath := "projects/" + projectIDOrKey + "/administrators"
	return getUserList(s.method.Get, spath, nil)
}

// DeleteAdmin deletes a admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(projectIDOrKey string, userID int) (*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}
	if userID < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	params := newRequestParams()
	params.Add("userId", strconv.Itoa(userID))

	spath := "projects/" + projectIDOrKey + "/administrators"
	return deleteUser(s.method.Delete, spath, params)
}
