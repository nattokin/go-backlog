package backlog

import (
	"context"

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

	Activity    *ProjectActivityService
	Category    *ProjectCategoryService
	CustomField *ProjectCustomFieldService
	IssueType   *ProjectIssueTypeService
	Status      *ProjectStatusService
	User        *ProjectUserService
	SharedFile  *ProjectSharedFileService
	Webhook     *ProjectWebhookService
	Version     *ProjectVersionService

	Option *ProjectOptionService
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
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.Option", such as:
//   - WithArchived
//   - WithChartEnabled
//   - WithKey
//   - WithName
//   - WithProjectLeaderCanEditProjectLeader
//   - WithSubtaskingEnabled
//   - WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *ProjectService) Update(ctx context.Context, projectIDOrKey string, option RequestOption, opts ...RequestOption) (*Project, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, option, toCoreOptions(opts)...)
	return projectFromModel(v), convertError(err)
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *ProjectService) Delete(ctx context.Context, projectIDOrKey string) (*Project, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey)
	return projectFromModel(v), convertError(err)
}

// DiskUsage returns information about the disk usage of your project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-disk-usage
func (s *ProjectService) DiskUsage(ctx context.Context, projectIDOrKey string) (*DiskUsageProject, error) {
	v, err := s.base.DiskUsage(ctx, projectIDOrKey)
	return diskUsageProjectFromModel(v), convertError(err)
}

// Icon returns the icon image of a project.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-icon
func (s *ProjectService) Icon(ctx context.Context, projectIDOrKey string) (*FileData, error) {
	v, err := s.base.Icon(ctx, projectIDOrKey)
	return fileDataFromCore(v), convertError(err)
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
		base:        project.NewService(method),
		Activity:    newProjectActivityService(method, option),
		Category:    newProjectCategoryService(method),
		CustomField: newProjectCustomFieldService(method, option),
		IssueType:   newProjectIssueTypeService(method, option),
		Status:      newProjectStatusService(method, option),
		User:        newProjectUserService(method, option),
		SharedFile:  newProjectSharedFileService(method),
		Webhook:     newProjectWebhookService(method, option),
		Version:     newProjectVersionService(method, option),
		Option:      newProjectOptionService(option),
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
	if ms == nil {
		return nil
	}

	result := make([]*Project, len(ms))
	for i, v := range ms {
		result[i] = projectFromModel(v)
	}
	return result
}
