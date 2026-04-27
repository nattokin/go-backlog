package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
)

// ──────────────────────────────────────────────────────────────
//  Project models
// ──────────────────────────────────────────────────────────────

// Project represents a project of Backlog.
type Project struct {
	ID                                int
	ProjectKey                        string
	Name                              string
	ChartEnabled                      bool
	SubtaskingEnabled                 bool
	ProjectLeaderCanEditProjectLeader bool
	TextFormattingRule                Format
	Archived                          bool
}

// ──────────────────────────────────────────────────────────────
//  ProjectService
// ──────────────────────────────────────────────────────────────

// ProjectService handles communication with the project-related methods of the Backlog API.
type ProjectService struct {
	base *project.Service

	Activity *ProjectActivityService
	Category *ProjectCategoryService
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
	v, err := s.base.All(ctx, toCoreOptions(opts)...)
	return projectsFromModel(v), convertError(err)
}

// One returns one of the projects searched by ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project
func (s *ProjectService) One(ctx context.Context, projectIDOrKey string) (*Project, error) {
	v, err := s.base.One(ctx, projectIDOrKey)
	return projectFromModel(v), convertError(err)
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
	v, err := s.base.Create(ctx, key, name, toCoreOptions(opts)...)
	return projectFromModel(v), convertError(err)
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
	v, err := s.base.Update(ctx, projectIDOrKey, toCoreOptions(opts)...)
	return projectFromModel(v), convertError(err)
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *ProjectService) Delete(ctx context.Context, projectIDOrKey string) (*Project, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey)
	return projectFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectActivityService
// ──────────────────────────────────────────────────────────────

// ProjectActivityService handles communication with the project activities-related methods of the Backlog API.
type ProjectActivityService struct {
	base *activity.ProjectService

	Option *ActivityOptionService
}

// List returns a list of activities in the project.
//
// This method supports options returned by methods in "*Client.Project.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Activity, error) {
	v, err := s.base.List(ctx, projectIDOrKey, toCoreOptions(opts)...)
	return activitiesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectCategoryService
// ──────────────────────────────────────────────────────────────

// ProjectCategoryService handles communication with the project category-related methods of the Backlog API.
type ProjectCategoryService struct {
	base *project.CategoryService
}

// All returns a list of categories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-category-list
func (s *ProjectCategoryService) All(ctx context.Context, projectIDOrKey string) ([]*Category, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return categoriesFromModel(v), convertError(err)
}

// Create adds a new category to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-category
func (s *ProjectCategoryService) Create(ctx context.Context, projectIDOrKey string, name string) (*Category, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name)
	return categoryFromModel(v), convertError(err)
}

// Update updates a category in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-category
func (s *ProjectCategoryService) Update(ctx context.Context, projectIDOrKey string, categoryID int, name string) (*Category, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, categoryID, name)
	return categoryFromModel(v), convertError(err)
}

// Delete deletes a category from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-category
func (s *ProjectCategoryService) Delete(ctx context.Context, projectIDOrKey string, categoryID int) (*Category, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, categoryID)
	return categoryFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectOptionService
// ──────────────────────────────────────────────────────────────

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	base *core.OptionService
}

// WithAll sets whether to include all projects.
func (s *ProjectOptionService) WithAll(enabled bool) RequestOption {
	return s.base.WithAll(enabled)
}

// WithArchived sets whether to include archived projects.
func (s *ProjectOptionService) WithArchived(enabled bool) RequestOption {
	return s.base.WithArchived(enabled)
}

// WithChartEnabled sets whether the project uses a chart.
func (s *ProjectOptionService) WithChartEnabled(enabled bool) RequestOption {
	return s.base.WithChartEnabled(enabled)
}

// WithKey sets the project key.
func (s *ProjectOptionService) WithKey(key string) RequestOption {
	return s.base.WithKey(key)
}

// WithName sets the project name.
func (s *ProjectOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// WithProjectLeaderCanEditProjectLeader sets whether a project leader can edit other project leaders.
func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return s.base.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithSubtaskingEnabled sets whether subtasking is enabled.
func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return s.base.WithSubtaskingEnabled(enabled)
}

// WithTextFormattingRule sets the text formatting rule.
func (s *ProjectOptionService) WithTextFormattingRule(format model.Format) RequestOption {
	return s.base.WithTextFormattingRule(format)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectService(method *core.Method, option *core.OptionService) *ProjectService {
	return &ProjectService{
		base:     project.NewService(method),
		Activity: newProjectActivityService(method, option),
		Category: newProjectCategoryService(method),
		User:     newProjectUserService(method, option),
		Option:   newProjectOptionService(option),
	}
}

func newProjectActivityService(method *core.Method, option *core.OptionService) *ProjectActivityService {
	return &ProjectActivityService{
		base:   activity.NewProjectService(method),
		Option: newActivityOptionService(option),
	}
}

func newProjectCategoryService(method *core.Method) *ProjectCategoryService {
	return &ProjectCategoryService{
		base: project.NewCategoryService(method),
	}
}

func newProjectOptionService(option *core.OptionService) *ProjectOptionService {
	return &ProjectOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func projectFromModel(m *model.Project) *Project {
	if m == nil {
		return nil
	}
	return &Project{
		ID:                                m.ID,
		ProjectKey:                        m.ProjectKey,
		Name:                              m.Name,
		ChartEnabled:                      m.ChartEnabled,
		SubtaskingEnabled:                 m.SubtaskingEnabled,
		ProjectLeaderCanEditProjectLeader: m.ProjectLeaderCanEditProjectLeader,
		TextFormattingRule:                Format(m.TextFormattingRule),
		Archived:                          m.Archived,
	}
}

func projectsFromModel(ms []*model.Project) []*Project {
	result := make([]*Project, len(ms))
	for i, v := range ms {
		result[i] = projectFromModel(v)
	}
	return result
}
