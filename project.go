package backlog

import (
	"net/url"
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
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(opts ...RequestOption) ([]*Project, error) {
	o := s.Option.registry.option
	validTypes := []queryType{queryAll, queryArchived}
	query := url.Values{}
	if err := o.applyQueryOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get("projects", query)
	if err != nil {
		return nil, err
	}

	v := []*Project{}
	if err := decodeResponse(resp, &v); err != nil {
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

	v := Project{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, opts ...RequestOption) (*Project, error) {
	o := s.Option.registry.option
	validTypes := []formType{formChartEnabled, formSubtaskingEnabled, formProjectLeaderCanEditProjectLeader, formTextFormattingRule}
	form := url.Values{}
	if err := o.applyFormOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}
	// apply mandatory key/name
	for _, opt := range []RequestOption{o.WithKey(key), o.WithName(name)} {
		if err := opt.Check(); err != nil {
			return nil, err
		}
		if err := opt.Set(form); err != nil {
			return nil, err
		}
	}

	resp, err := s.method.Post("projects", form)
	if err != nil {
		return nil, err
	}

	v := Project{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(projectIDOrKey string, options ...RequestOption) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	o := s.Option.registry.option
	validTypes := []formType{
		formKey, formName, formChartEnabled, formSubtaskingEnabled,
		formProjectLeaderCanEditProjectLeader, formTextFormattingRule, formArchived,
	}
	form := url.Values{}
	if err := o.applyFormOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(spath, form)
	if err != nil {
		return nil, err
	}

	v := Project{}
	if err := decodeResponse(resp, &v); err != nil {
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
	resp, err := s.method.Delete(spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := Project{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
