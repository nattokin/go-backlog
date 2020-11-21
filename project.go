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
// This method can specify the options returned by methods in "*Client.Project.Option".
//
// Use the following methods:
//    WithQueryAll
//    WithQueryArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(options ...*QueryOption) ([]*Project, error) {
	validOptions := []queryType{queryAll, queryArchived}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	params := NewQueryParams()
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
// This method can specify the options returned by methods in "*Client.Project.Option".
//
// Use the following methods:
//    WithQueryArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAll(options ...*QueryOption) ([]*Project, error) {
	validOptions := []queryType{queryArchived}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	return s.All(append(options, s.Option.WithQueryAll(true))...)
}

// AllUnarchived returns all of joining projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllUnarchived() ([]*Project, error) {
	return s.All(s.Option.WithQueryArchived(false))
}

// AdminAllUnarchived returns all of projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAllUnarchived() ([]*Project, error) {
	return s.All(s.Option.WithQueryAll(true), s.Option.WithQueryArchived(false))
}

// AllArchived returns all of joining projects archived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllArchived() ([]*Project, error) {
	return s.All(s.Option.WithQueryArchived(true))
}

// AdminAllArchived returns all of projects archived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AdminAllArchived() ([]*Project, error) {
	return s.All(s.Option.WithQueryAll(true), s.Option.WithQueryArchived(true))
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
//   WithFormChartEnabled
//   WithFormSubtaskingEnabled
//   WithFormProjectLeaderCanEditProjectLeader
//   WithFormTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, options ...*FormOption) (*Project, error) {
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	validOptions := []formType{formChartEnabled, formSubtaskingEnabled, formProjectLeaderCanEditProjectLeader, formTextFormattingRule}
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
//   WithFormKey
//   WithFormName
//   WithFormChartEnabled
//   WithFormSubtaskingEnabled
//   WithFormProjectLeaderCanEditProjectLeader
//   WithFormTextFormattingRule
//   WithFormArchived
//   WithFormTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(target ProjectIDOrKeyGetter, options ...*FormOption) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}

	validOptions := []formType{
		formKey, formName, formChartEnabled, formSubtaskingEnabled, formProjectLeaderCanEditProjectLeader, formTextFormattingRule, formArchived,
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
