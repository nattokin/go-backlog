package backlog

import (
	"context"
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
// This method supports options returned by methods in "*Client.Project.Option",
// such as:
//   - WithQueryAll
//   - WithQueryArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(ctx context.Context, opts ...RequestOption) ([]*Project, error) {

	query := url.Values{}
	validTypes := []apiParamOptionType{paramAll, paramArchived}
	if err := applyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "projects", query)
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
func (s *ProjectService) One(ctx context.Context, projectIDOrKey string) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Get(ctx, spath, nil)
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
// This method supports options returned by methods in "*Client.Project.Option",
// such as:
//   - WithChartEnabled
//   - WithProjectLeaderCanEditProjectLeader
//   - WithSubtaskingEnabled
//   - WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(ctx context.Context, key, name string, opts ...RequestOption) (*Project, error) {

	form := url.Values{}
	validTypes := []apiParamOptionType{paramKey, paramName, paramChartEnabled, paramSubtaskingEnabled, paramProjectLeaderCanEditProjectLeader, paramTextFormattingRule}
	options := append([]RequestOption{s.Option.base.WithKey(key), s.Option.base.WithName(name)}, opts...)
	if err := applyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	resp, err := s.method.Post(ctx, "projects", form)
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
// This method supports options returned by methods in "*Client.Project.Option",
// such as:
//   - WithArchived
//   - WithChartEnabled
//   - WithKey
//   - WithName
//   - WithProjectLeaderCanEditProjectLeader
//   - WithSubtaskingEnabled
//   - WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(ctx context.Context, projectIDOrKey string, opts ...RequestOption) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []apiParamOptionType{
		paramKey, paramName, paramChartEnabled, paramSubtaskingEnabled,
		paramProjectLeaderCanEditProjectLeader, paramTextFormattingRule, paramArchived,
	}
	if err := applyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(ctx, spath, form)
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
func (s *ProjectService) Delete(ctx context.Context, projectIDOrKey string) (*Project, error) {
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := Project{}
	if err := decodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}
