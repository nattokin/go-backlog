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

// All returns a list of projects.
//
// By default, it returns only projects the authenticated user has joined.
// The "WithQueryAll" option is available only for administrators; when set to true,
// it returns all projects in the space. When false (default), it returns joined projects only.
//
// This method supports options returned by methods in "*Client.Project.Option",
// such as:
//
//   - WithQueryAll
//   - WithQueryArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(opts ...*QueryOption) ([]*Project, error) {
	validOptions := []queryType{queryAll, queryArchived}
	for _, option := range opts {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	o := s.Option.support.query
	query := NewQueryParams()
	err := o.applyOptions(query, opts...)
	if err != nil {
		return nil, err
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
//
//	WithFormChartEnabled
//	WithFormSubtaskingEnabled
//	WithFormProjectLeaderCanEditProjectLeader
//	WithFormTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, opts ...*FormOption) (*Project, error) {
	validOptions := []formType{formChartEnabled, formSubtaskingEnabled, formProjectLeaderCanEditProjectLeader, formTextFormattingRule}
	for _, option := range opts {
		if err := option.validate(validOptions); err != nil {
			return nil, err
		}
	}

	o := s.Option.support.form
	form := NewFormParams()
	err := o.applyOptions(form, append(
		[]*FormOption{o.WithKey(key), o.WithName(name)}, opts...,
	)...)
	if err != nil {
		return nil, err
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
//
//	WithFormKey
//	WithFormName
//	WithFormChartEnabled
//	WithFormSubtaskingEnabled
//	WithFormProjectLeaderCanEditProjectLeader
//	WithFormTextFormattingRule
//	WithFormArchived
//	WithFormTextFormattingRule
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

	o := s.Option.support.form
	form := NewFormParams()
	err := o.applyOptions(form, options...)
	if err != nil {
		return nil, err
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
