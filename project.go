package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/user"
)

// Project represents a Backlog project.
type Project struct {
	ID                       int
	ProjectKey               string
	Name                     string
	ChartEnabled             bool
	SubtaskingEnabled        bool
	ProjectLeaderCanEditProjectLeader bool
	TextFormattingRule       string
	Archived                 bool
	DisplayOrder             int
}

// DiskUsageProject represents disk usage of a project.
type DiskUsageProject struct {
	ProjectID int
	Issue     int
	Wiki      int
	File      int
	Subversion int
	Git        int
	GitLFS     int
}

// ──────────────────────────────────────────────────────────────
//  ProjectService
// ──────────────────────────────────────────────────────────────

// ProjectService handles communication with the project-related methods of the Backlog API.
type ProjectService struct {
	base *project.Service

	Category  *ProjectCategoryService
	IssueType *ProjectIssueTypeService
	Status    *ProjectStatusService
	User      *ProjectUserService

	Option *ProjectOptionService
}

// All returns a list of projects.
//
// This method supports options returned by methods in "*Client.Project.Option",
// such as:
//   - WithAll
//   - WithArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *ProjectService) All(ctx context.Context, opts ...RequestOption) ([]*Project, error) {
	v, err := s.base.All(ctx, toCoreOptions(opts)...)
	return projectsFromModel(v), convertError(err)
}

// One returns a single project by ID or key.
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
//   - WithSubtaskingEnabled
//   - WithProjectLeaderCanEditProjectLeader
//   - WithTextFormattingRule
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *ProjectService) Create(ctx context.Context, key, name string, opts ...RequestOption) (*Project, error) {
	v, err := s.base.Create(ctx, key, name, toCoreOptions(opts)...)
	return projectFromModel(v), convertError(err)
}

// Update updates an existing project.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.Option", such as:
//   - WithKey
//   - WithName
//   - WithChartEnabled
//   - WithSubtaskingEnabled
//   - WithProjectLeaderCanEditProjectLeader
//   - WithTextFormattingRule
//   - WithArchived
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

// DiskUsage returns disk usage of a project.
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
func (s *ProjectCategoryService) Create(ctx context.Context, projectIDOrKey, name string) (*Category, error) {
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
//  ProjectIssueTypeService
// ──────────────────────────────────────────────────────────────

// ProjectIssueTypeService handles communication with the project issue-type-related methods of the Backlog API.
type ProjectIssueTypeService struct {
	base *project.IssueTypeService
}

// All returns a list of issue types in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-type-list
func (s *ProjectIssueTypeService) All(ctx context.Context, projectIDOrKey string) ([]*IssueType, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return issueTypesFromModel(v), convertError(err)
}

// Create adds a new issue type to a project.
//
// This method supports options returned by methods in "*Client.Project.IssueType.Option",
// such as:
//   - WithTemplateSummary
//   - WithTemplateDescription
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue-type
func (s *ProjectIssueTypeService) Create(ctx context.Context, projectIDOrKey, name, color string, opts ...RequestOption) (*IssueType, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name, color, toCoreOptions(opts)...)
	return issueTypeFromModel(v), convertError(err)
}

// Update updates an issue type in a project.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.IssueType.Option", such as:
//   - WithName
//   - WithColor
//   - WithTemplateSummary
//   - WithTemplateDescription
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue-type
func (s *ProjectIssueTypeService) Update(ctx context.Context, projectIDOrKey string, issueTypeID int, option RequestOption, opts ...RequestOption) (*IssueType, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, issueTypeID, option, toCoreOptions(opts)...)
	return issueTypeFromModel(v), convertError(err)
}

// Delete deletes an issue type from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-type
func (s *ProjectIssueTypeService) Delete(ctx context.Context, projectIDOrKey string, issueTypeID, substituteIssueTypeID int) (*IssueType, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, issueTypeID, substituteIssueTypeID)
	return issueTypeFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectStatusService
// ──────────────────────────────────────────────────────────────

// ProjectStatusService handles communication with the project status-related methods of the Backlog API.
type ProjectStatusService struct {
	base *project.StatusService
}

// All returns a list of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project
func (s *ProjectStatusService) All(ctx context.Context, projectIDOrKey string) ([]*Status, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return statusesFromModel(v), convertError(err)
}

// Create adds a new status to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-status
func (s *ProjectStatusService) Create(ctx context.Context, projectIDOrKey, name, color string) (*Status, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name, color)
	return statusFromModel(v), convertError(err)
}

// Update updates a status in a project.
//
// At least one option is required.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-status
func (s *ProjectStatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, option RequestOption, opts ...RequestOption) (*Status, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, statusID, option, toCoreOptions(opts)...)
	return statusFromModel(v), convertError(err)
}

// Delete deletes a status from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-status
func (s *ProjectStatusService) Delete(ctx context.Context, projectIDOrKey string, statusID, substituteStatusID int) (*Status, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, statusID, substituteStatusID)
	return statusFromModel(v), convertError(err)
}

// UpdateOrder updates the display order of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-order-of-status
func (s *ProjectStatusService) UpdateOrder(ctx context.Context, projectIDOrKey string, statusIDs []int) ([]*Status, error) {
	v, err := s.base.UpdateOrder(ctx, projectIDOrKey, statusIDs)
	return statusesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectUserService
// ──────────────────────────────────────────────────────────────

// ProjectUserService handles communication with the project user-related methods of the Backlog API.
type ProjectUserService struct {
	base *user.ProjectService
}

// All returns a list of users in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-user-list
func (s *ProjectUserService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*User, error) {
	v, err := s.base.All(ctx, projectIDOrKey, excludeGroupMembers)
	return usersFromModel(v), convertError(err)
}

// Add adds a user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-user
func (s *ProjectUserService) Add(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// Delete removes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AddAdmin adds a user as administrator to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.AddAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AdminAll returns a list of users who are administrators of the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*User, error) {
	v, err := s.base.AdminAll(ctx, projectIDOrKey)
	return usersFromModel(v), convertError(err)
}

// DeleteAdmin removes a user from the administrators of the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.DeleteAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
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

// WithArchived filters projects by archived state.
func (s *ProjectOptionService) WithArchived(archived bool) RequestOption {
	return s.base.WithArchived(archived)
}

// WithChartEnabled sets whether the chart feature is enabled.
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

// WithProjectLeaderCanEditProjectLeader sets whether the project leader can edit the project leader.
func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) RequestOption {
	return s.base.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithSubtaskingEnabled sets whether subtasking is enabled.
func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) RequestOption {
	return s.base.WithSubtaskingEnabled(enabled)
}

// WithTextFormattingRule sets the text formatting rule.
func (s *ProjectOptionService) WithTextFormattingRule(rule TextFormattingRule) RequestOption {
	return s.base.WithTextFormattingRule(model.TextFormattingRule(rule))
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectService(method *core.Method, option *core.OptionService) *ProjectService {
	return &ProjectService{
		base:      project.NewService(method),
		Category:  newProjectCategoryService(method),
		IssueType: newProjectIssueTypeService(method),
		Status:    newProjectStatusService(method),
		User:      newProjectUserService(method, option),
		Option:    newProjectOptionService(option),
	}
}

func newProjectCategoryService(method *core.Method) *ProjectCategoryService {
	return &ProjectCategoryService{base: project.NewCategoryService(method)}
}

func newProjectIssueTypeService(method *core.Method) *ProjectIssueTypeService {
	return &ProjectIssueTypeService{base: project.NewIssueTypeService(method)}
}

func newProjectStatusService(method *core.Method) *ProjectStatusService {
	return &ProjectStatusService{base: project.NewStatusService(method)}
}

func newProjectOptionService(option *core.OptionService) *ProjectOptionService {
	return &ProjectOptionService{base: option}
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
		TextFormattingRule:                m.TextFormattingRule,
		Archived:                          m.Archived,
		DisplayOrder:                      m.DisplayOrder,
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

func diskUsageProjectFromModel(m *model.DiskUsageProject) *DiskUsageProject {
	if m == nil {
		return nil
	}
	return &DiskUsageProject{
		ProjectID:  m.ProjectID,
		Issue:      m.Issue,
		Wiki:       m.Wiki,
		File:       m.File,
		Subversion: m.Subversion,
		Git:        m.Git,
		GitLFS:     m.GitLFS,
	}
}

func issueTypeFromModel(m *model.IssueType) *IssueType {
	if m == nil {
		return nil
	}
	return &IssueType{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		Name:         m.Name,
		Color:        m.Color,
		DisplayOrder: m.DisplayOrder,
	}
}

func issueTypesFromModel(ms []*model.IssueType) []*IssueType {
	if ms == nil {
		return nil
	}
	result := make([]*IssueType, len(ms))
	for i, v := range ms {
		result[i] = issueTypeFromModel(v)
	}
	return result
}
