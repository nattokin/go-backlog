package backlog

import (
	"encoding/json"
	"errors"
	"fmt"
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
	if i <= 0 {
		return "", errors.New("id must be greater than 0")
	}
	return strconv.Itoa(int(i)), nil
}

func (k ProjectKey) getProjectIDOrKey() (string, error) {
	if k == "" {
		return "", errors.New("key must not be empty")
	}
	return string(k), nil
}

// ProjectOption is type of functional option for ProjectService.
type ProjectOption func(p *requestParams) error

// WithKey returns option. the option sets `key` for project.
func (*ProjectOptionService) WithKey(key string) ProjectOption {
	return func(p *requestParams) error {
		if key == "" {
			return errors.New("key must not be empty")
		}
		p.Set("key", key)
		return nil
	}
}

// WithName returns option. the option sets `name` for project.
func (*ProjectOptionService) WithName(name string) ProjectOption {
	return func(p *requestParams) error {
		if name == "" {
			return errors.New("name must not be empty")
		}
		p.Set("name", name)
		return nil
	}
}

// WithChartEnabled returns option. the option sets `chartEnabled` for project.
func (*ProjectOptionService) WithChartEnabled(enabeld bool) ProjectOption {
	return func(p *requestParams) error {
		p.Set("chartEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

// WithSubtaskingEnabled returns option. the option sets `subtaskingEnabled` for project.
func (*ProjectOptionService) WithSubtaskingEnabled(enabeld bool) ProjectOption {
	return func(p *requestParams) error {
		p.Set("subtaskingEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

// WithProjectLeaderCanEditProjectLeader returns option. the option sets `projectLeaderCanEditProjectLeader` for project.
func (*ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabeld bool) ProjectOption {
	return func(p *requestParams) error {
		p.Set("projectLeaderCanEditProjectLeader", strconv.FormatBool(enabeld))
		return nil
	}
}

// WithTextFormattingRule returns option. the option sets `textFormattingRule` for project.
func (*ProjectOptionService) WithTextFormattingRule(format string) ProjectOption {
	return func(p *requestParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			return fmt.Errorf("format must be only '%s' or '%s'", FormatBacklog, FormatMarkdown)
		}
		p.Set("textFormattingRule", format)
		return nil
	}
}

// WithArchived returns option. the option sets `archived` for project.
func (*ProjectOptionService) WithArchived(archived bool) ProjectOption {
	return func(p *requestParams) error {
		p.Set("archived", strconv.FormatBool(archived))
		return nil
	}
}

// Joined returns all of joining projects.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) Joined() ([]*Project, error) {
	params := newRequestParams()
	params.Set("all", "false")

	return s.getList(params)
}

// All returns all of projects. This is limited to admin.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All() ([]*Project, error) {
	params := newRequestParams()
	params.Set("all", "true")

	return s.getList(params)
}

// Archived returns all of joining projects archived.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) Archived() ([]*Project, error) {
	params := newRequestParams()
	params.Set("archived", "true")
	params.Set("all", "false")

	return s.getList(params)
}

// AllArchived returns all of projects archived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllArchived() ([]*Project, error) {
	params := newRequestParams()
	params.Set("archived", "true")
	params.Set("all", "true")

	return s.getList(params)
}

// Unarchived returns all of joining projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) Unarchived() ([]*Project, error) {
	params := newRequestParams()
	params.Set("archived", "false")
	params.Set("all", "false")

	return s.getList(params)
}

// AllUnarchived returns all of projects unarchived.
// If you are not an admin, only joining projects returned.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) AllUnarchived() ([]*Project, error) {
	params := newRequestParams()
	params.Set("archived", "false")
	params.Set("all", "true")

	return s.getList(params)
}

func (s *ProjectService) getList(params *requestParams) ([]*Project, error) {
	resp, err := s.clientMethod.Get("projects", params)
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
func (s *ProjectService) One(target ProjectIDOrKeyGetter) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	spath := "projects/" + projectIDOrKey
	resp, err := s.clientMethod.Get(spath, nil)
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
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(key, name string, options ...ProjectOption) (*Project, error) {
	if key == "" {
		return nil, errors.New("key must not be empty")
	}
	if name == "" {
		return nil, errors.New("name must not be empty")
	}

	params := newRequestParams()

	// Set default options.
	params.Set("chartEnabled", "false")
	params.Set("subtaskingEnabled", "false")
	params.Set("projectLeaderCanEditProjectLeader", "false")
	params.Set("textFormattingRule", FormatMarkdown)

	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	// Disable invalid options.
	params.Set("key", key)
	params.Set("name", name)
	params.Del("archived")

	resp, err := s.clientMethod.Post("projects", params)
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
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(target ProjectIDOrKeyGetter, options ...ProjectOption) (*Project, error) {
	projectIDOrKey, err := target.getProjectIDOrKey()
	if err != nil {
		return nil, err
	}
	params := newRequestParams()
	for _, option := range options {
		if err := option(params); err != nil {
			return nil, err
		}
	}

	spath := "projects/" + projectIDOrKey
	resp, err := s.clientMethod.Patch(spath, params)
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
	resp, err := s.clientMethod.Delete(spath, newRequestParams())
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
// 	resp, err := s.clientMethod.Get(spath, nil)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return resp.Body, nil
// }
