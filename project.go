package backlog

import (
	"encoding/json"
	"path"
)

func validateProjectID(projectID int) error {
	if projectID < 1 {
		return newValidationError("projectID must not be less than 1")
	}
	return nil
}

func validateProjectIDOrKey(projectIDOrKey string) error {
	if projectIDOrKey == "" {
		return newValidationError("projectIDOrKey must not be empty")
	}
	if projectIDOrKey == "0" {
		return newValidationError("projectIDOrKey must not be '0'")
	}
	return nil
}

// ProjectService handles communication with the project-related methods of the Backlog API.
type ProjectService struct {
	method *method

	Activity *ProjectActivityService
	User     *ProjectUserService
	Option   *ProjectOptionService
}

// All returns a list of all projects.
//
// This method supports options returned by methods in "*Client.Project.Option".
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

	query := NewQueryParams()
	for _, option := range options {
		if err := option.set(query); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Get("projects", query)
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
// This method supports options returned by methods in "*Client.Project.Option".
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
func (s *ProjectService) One(projectIDOrKey string) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
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
// This method supports options returned by methods in "*Client.Project.Option".
//
// Use the following methods:
//   WithFormChartEnabled
//   WithFormSubtaskingEnabled
//   WithFormProjectLeaderCanEditProjectLeader
//   WithFormTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, options ...*FormOption) (*Project, error) {
	form := NewFormParams()
	if err := withFormKey(key).set(form); err != nil {
		return nil, err
	}
	if err := withFormName(name).set(form); err != nil {
		return nil, err
	}

	validOptions := []formType{formChartEnabled, formSubtaskingEnabled, formProjectLeaderCanEditProjectLeader, formTextFormattingRule}
	for _, option := range options {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Post("projects", form)
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
// This method supports options returned by methods in "*Client.Project.Option".
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
func (s *ProjectService) Update(projectIDOrKey string, options ...*FormOption) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
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

	form := NewFormParams()
	for _, option := range options {
		if err := option.set(form); err != nil {
			return nil, err
		}
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(spath, form)
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
func (s *ProjectService) Delete(projectIDOrKey string) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Delete(spath, NewFormParams())
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

// Icon returns the icon image of the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-icon
// func (s *ProjectService) Icon(projectIDOrKey string) (io.ReadCloser, error) {
// 	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
// 		return nil, err
// 	}
//
// 	spath := path.Join("projects", projectIDOrKey, "image")
// 	resp, err := s.method.Get(spath, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Body, nil
// }
