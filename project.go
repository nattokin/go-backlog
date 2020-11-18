package backlog

import (
	"encoding/json"
	"errors"
	"strconv"
)

// ProjectIDOrKeyGetter has method to get ProjectIDOrKey and validation errror.
type ProjectIDOrKeyGetter interface {
	getProjectIDOrKey() (string, error)
}

// ProjectID implements ProjectIDOrKeyGetter interface.
type ProjectID int

// ProjectKey implements ProjectIDOrKeyGetter interface.
type ProjectKey string

func (i ProjectID) getProjectIDOrKey() (string, error) {
	if i < 1 {
		return "", errors.New("id must not be less than 1")
	}
	return strconv.Itoa(int(i)), nil
}

func (k ProjectKey) getProjectIDOrKey() (string, error) {
	if k == "" {
		return "", errors.New("key must not be empty")
	}
	return string(k), nil
}

// ProjectService has methods for Project.
type ProjectService struct {
	method *method

	Activity *ProjectActivityService
	User     *ProjectUserService
	Option   *ProjectOptionService
}

// All returns all of projects.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(options ...*ProjectOption) ([]*Project, error) {
	validOptions := []optionType{optionAll, optionArchived}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Get("projects", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := []*Project{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return v, nil
}

// AdminAll returns all of projects. This is limited to admin.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAll(options ...*ProjectOption) ([]*Project, error) {
	validOptions := []optionType{optionArchived}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	return s.All(append(options, s.Option.WithAll(true))...)
}

// AllUnarchived returns all of joining projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllUnarchived() ([]*Project, error) {
	params := newRequestParams()
	params.Set("archived", "false")

	return s.All(s.Option.WithArchived(false))
}

// AdminAllUnarchived returns all of projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAllUnarchived() ([]*Project, error) {
	return s.All(s.Option.WithAll(true), s.Option.WithArchived(false))
}

// AllArchived returns all of joining projects archived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllArchived() ([]*Project, error) {
	return s.All(s.Option.WithArchived(true))
}

// AdminAllArchived returns all of projects archived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAllArchived() ([]*Project, error) {
	return s.All(s.Option.WithAll(true), s.Option.WithArchived(true))
}

// One returns one of the projects searched by ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project
func (s *ProjectService) One(target ProjectIDOrKeyGetter) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	spath := "projects/" + projectIDOrKey
	resp, err := s.method.Get(spath, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Project{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new project.
//
// This method can specify the options returned by methods in "*Client.Project.Option".
//
// Use the following methods:
//   WithChartEnabled
//   WithSubtaskingEnabled
//   WithProjectLeaderCanEditProjectLeader
//   WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, options ...*ProjectOption) (*Project, error) {
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	validOptions := []optionType{optionChartEnabled, optionSubtaskingEnabled, optionProjectLeaderCanEditProjectLeader, optionTextFormattingRule}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
			return nil, err
		}
	}
	params.Set("key", key)
	params.Set("name", name)

	resp, err := s.method.Post("projects", params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Project{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a project.
//
// This method can specify the options returned by methods in "*Client.Project.Option".
//
// Use the following methods:
//   WithKey
//   WithName
//   WithChartEnabled
//   WithSubtaskingEnabled
//   WithProjectLeaderCanEditProjectLeader
//   WithTextFormattingRule
//   WithArchived
//   WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(target ProjectIDOrKeyGetter, options ...*ProjectOption) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	validOptions := []optionType{
		optionKey, optionName, optionChartEnabled, optionSubtaskingEnabled, optionProjectLeaderCanEditProjectLeader, optionTextFormattingRule, optionArchived,
	}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := newRequestParams()
	for _, option := range options {
		if err := option.set(params); err != nil {
			return nil, err
		}
	}

	spath := "projects/" + projectIDOrKey
	resp, err := s.method.Patch(spath, params)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Project{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *ProjectService) Delete(target ProjectIDOrKeyGetter) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	spath := "projects/" + projectIDOrKey
	resp, err := s.method.Delete(spath, newRequestParams())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	v := Project{}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}

	return &v, nil
}

// TODO: Icon returns icon image of the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-icon
// func (s *ProjectService) Icon(target ProjectIDOrKeyGetter) (io.ReadCloser, error) {
// 	projectIDOrKey, err := target.getProjectIDOrKey()
// 	if err != nil {
// 		return nil, err
// 	}
// 	spath := "projects/" + projectIDOrKey + "/image"
// 	resp, err := s.method.Get(spath, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Body, nil
// }
