package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

// UserOptionService has methods to make functional option for UserService.
type UserOptionService struct {
}

// UserOption is type of functional option for UserService.
type UserOption func(p *requestParams) error

// WithPassword returns option. the option sets `password` for user.
func (*UserOptionService) WithPassword(password string) UserOption {
	return func(p *requestParams) error {
		if password == "" {
			return errors.New("password must not be empty")
		}
		p.Set("password", password)
		return nil
	}
}

// WithName returns option. the option sets `password` for user.
func (*UserOptionService) WithName(name string) UserOption {
	return func(p *requestParams) error {
		if name == "" {
			return errors.New("name must not be empty")
		}
		p.Set("name", name)
		return nil
	}
}

// WithMailAddress returns option. the option sets `mailAddress` for user.
func (*UserOptionService) WithMailAddress(mailAddress string) UserOption {
	// ToDo: validate mailAddress
	return func(p *requestParams) error {
		if mailAddress == "" {
			return errors.New("mailAddress must not be empty")
		}
		p.Set("mailAddress", mailAddress)
		return nil
	}
}

// WithRoleType returns option. the option sets `roleType` for user.
func (*UserOptionService) WithRoleType(roleType int) UserOption {
	return func(p *requestParams) error {
		if roleType < 1 || 6 < roleType {
			return errors.New("roleType must be between 1 and 7")
		}
		p.Add("roleType", strconv.Itoa(roleType))
		return nil
	}
}

type baseUserService struct {
	clientMethod *clientMethod
}

func newBaseUserService(cm *clientMethod) *baseUserService {
	return &baseUserService{
		clientMethod: cm,
	}
}

func (s *baseUserService) get(spath string) (*User, error) {
	resp, err := s.clientMethod.Get(spath, nil)
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

func (s *baseUserService) getList(spath string, params *requestParams) ([]*User, error) {
	resp, err := s.clientMethod.Get(spath, params)
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

func (s *baseUserService) post(spath string, params *requestParams) (*User, error) {
	resp, err := s.clientMethod.Post(spath, params)
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

func (s *baseUserService) patch(spath string, params *requestParams) (*User, error) {
	resp, err := s.clientMethod.Patch(spath, params)
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

func (s *baseUserService) delete(spath string, params *requestParams) (*User, error) {
	resp, err := s.clientMethod.Delete(spath, params)
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
	*baseUserService

	Activity *UserActivityService
	Option   *UserOptionService
}

func newUserService(cm *clientMethod) *UserService {
	return &UserService{
		baseUserService: newBaseUserService(cm),
		Activity:        newUserActivityService(cm),
		Option:          &UserOptionService{},
	}
}

// All returns all users in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *UserService) All() ([]*User, error) {
	spath := "users"
	return s.getList(spath, nil)
}

// One returns a user in your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *UserService) One(id int) (*User, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id)
	return s.get(spath)
}

// Own returns your own user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *UserService) Own() (*User, error) {
	spath := "users/myself"
	return s.get(spath)
}

// ToDo: func (s *UserService) Icon()
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon

// Add adds a user to your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *UserService) Add(userID, password, name, mailAddress string, roleType int) (*User, error) {
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
	if roleType < 1 || 6 < roleType {
		return nil, errors.New("roleType must be between 1 and 7")
	}

	params := newRequestParams()
	params.Add("userId", userID)
	params.Add("password", password)
	params.Add("name", name)
	params.Add("mailAddress", mailAddress)
	params.Add("roleType", strconv.Itoa(roleType))

	spath := "users"
	return s.post(spath, params)
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

	return s.patch(spath, params)
}

// Delete deletes a user from your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *UserService) Delete(id int) (*User, error) {
	if id < 1 {
		return nil, errors.New("id must be greater than 1")
	}

	spath := "users/" + strconv.Itoa(id)
	return s.delete(spath, nil)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	*baseUserService
}

func newProjectUserService(cm *clientMethod) *ProjectUserService {
	return &ProjectUserService{
		baseUserService: newBaseUserService(cm),
	}
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
	return s.getList(spath, params)
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
	return s.post(spath, params)
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
	return s.delete(spath, params)
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
	return s.post(spath, params)
}

// AdminAll returns all of admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(projectIDOrKey string) ([]*User, error) {
	if projectIDOrKey == "" {
		return nil, errors.New("projectIDOrKey must not be empty")
	}

	spath := "projects/" + projectIDOrKey + "/administrators"
	return s.getList(spath, nil)
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
	return s.delete(spath, params)
}
