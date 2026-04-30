package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/sharedfile"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/version"
	"github.com/nattokin/go-backlog/internal/webhook"
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

// ProjectStatus represents a status defined within a project.
type ProjectStatus struct {
	ID           int
	ProjectID    int
	Name         string
	Color        string
	DisplayOrder int
}

// ──────────────────────────────────────────────────────────────
//  ProjectService
// ──────────────────────────────────────────────────────────────

// ProjectService handles communication with the project-related methods of the Backlog API.
type ProjectService struct {
	base *project.Service

	Activity   *ProjectActivityService
	Category   *ProjectCategoryService
	Status     *ProjectStatusService
	User       *ProjectUserService
	SharedFile *ProjectSharedFileService
	Webhook    *ProjectWebhookService
	Version    *ProjectVersionService

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

// DiskUsage returns information about the disk usage of your project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-disk-usage
func (s *ProjectService) DiskUsage(ctx context.Context, projectIDOrKey string) (*DiskUsageProject, error) {
	v, err := s.base.DiskUsage(ctx, projectIDOrKey)
	return diskUsageProjectFromModel(v), convertError(err)
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
//  ProjectStatusService
// ──────────────────────────────────────────────────────────────

// ProjectStatusService handles communication with the project status-related methods of the Backlog API.
type ProjectStatusService struct {
	base   *project.StatusService
	Option *ProjectStatusOptionService
}

// All returns a list of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project
func (s *ProjectStatusService) All(ctx context.Context, projectIDOrKey string) ([]*ProjectStatus, error) {
	v, err := s.base.All(ctx, projectIDOrKey)
	return projectStatusesFromModel(v), convertError(err)
}

// Create adds a new status to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-status
func (s *ProjectStatusService) Create(ctx context.Context, projectIDOrKey, name, color string) (*ProjectStatus, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name, color)
	return projectStatusFromModel(v), convertError(err)
}

// Update updates a status in a project.
//
// This method supports options returned by methods in "*Client.Project.Status.Option",
// such as:
//   - WithColor
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-status
func (s *ProjectStatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, opts ...RequestOption) (*ProjectStatus, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, statusID, toCoreOptions(opts)...)
	return projectStatusFromModel(v), convertError(err)
}

// Delete deletes a status from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-status
func (s *ProjectStatusService) Delete(ctx context.Context, projectIDOrKey string, statusID, substituteStatusID int) (*ProjectStatus, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, statusID, substituteStatusID)
	return projectStatusFromModel(v), convertError(err)
}

// UpdateOrder updates the display order of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-order-of-status
func (s *ProjectStatusService) UpdateOrder(ctx context.Context, projectIDOrKey string, statusIDs []int) ([]*ProjectStatus, error) {
	v, err := s.base.UpdateOrder(ctx, projectIDOrKey, statusIDs)
	return projectStatusesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectSharedFileService
// ──────────────────────────────────────────────────────────────

// ProjectSharedFileService handles communication with the project shared-file-related methods of the Backlog API.
type ProjectSharedFileService struct {
	base *sharedfile.ProjectService
}

// List returns a list of shared files in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files
func (s *ProjectSharedFileService) List(ctx context.Context, projectIDOrKey string) ([]*SharedFile, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return sharedFilesFromModel(v), convertError(err)
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	base *user.ProjectService
}

// All returns all users in the project.
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

// Delete deletes a user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-user
func (s *ProjectUserService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AddAdmin adds a admin user to the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project-administrator
func (s *ProjectUserService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.AddAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// AdminAll returns a list of all admin users in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators
func (s *ProjectUserService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*User, error) {
	v, err := s.base.AdminAll(ctx, projectIDOrKey)
	return usersFromModel(v), convertError(err)
}

// DeleteAdmin removes an admin user from the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator
func (s *ProjectUserService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*User, error) {
	v, err := s.base.DeleteAdmin(ctx, projectIDOrKey, userID)
	return userFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectWebhookService
// ──────────────────────────────────────────────────────────────

// ProjectWebhookService handles communication with the project webhook-related methods of the Backlog API.
type ProjectWebhookService struct {
	base   *webhook.Service
	Option *ProjectWebhookOptionService
}

// All returns a list of webhooks in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks
func (s *ProjectWebhookService) All(ctx context.Context, projectIDOrKey string) ([]*Webhook, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return webhooksFromModel(v), convertError(err)
}

// Create adds a webhook to the project.
//
// This method supports options returned by methods in "*Client.Project.Webhook.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithAllEvent
//   - WithDescription
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-webhook
func (s *ProjectWebhookService) Create(ctx context.Context, projectIDOrKey, name, hookURL string, opts ...RequestOption) (*Webhook, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, name, hookURL, toCoreOptions(opts)...)
	return webhookFromModel(v), convertError(err)
}

// One returns a webhook in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-webhook
func (s *ProjectWebhookService) One(ctx context.Context, projectIDOrKey string, webhookID int) (*Webhook, error) {
	v, err := s.base.Get(ctx, projectIDOrKey, webhookID)
	return webhookFromModel(v), convertError(err)
}

// Update updates a webhook.
//
// This method supports options returned by methods in "*Client.Project.Webhook.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithAllEvent
//   - WithDescription
//   - WithHookURL
//   - WithName
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-webhook
func (s *ProjectWebhookService) Update(ctx context.Context, projectIDOrKey string, webhookID int, opts ...RequestOption) (*Webhook, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, webhookID, toCoreOptions(opts)...)
	return webhookFromModel(v), convertError(err)
}

// Delete deletes a webhook.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-webhook
func (s *ProjectWebhookService) Delete(ctx context.Context, projectIDOrKey string, webhookID int) (*Webhook, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, webhookID)
	return webhookFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  ProjectVersionService
// ──────────────────────────────────────────────────────────────

// ProjectVersionService handles communication with the project version/milestone-related methods of the Backlog API.
type ProjectVersionService struct {
	base   *version.Service
	Option *ProjectVersionOptionService
}

// All returns a list of versions/milestones in the project.
//
// This method supports options returned by methods in "*Client.Project.Version.Option",
// such as:
//   - WithArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-version-milestone-list/
func (s *ProjectVersionService) All(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Version, error) {
	v, err := s.base.All(ctx, projectIDOrKey, toCoreOptions(opts)...)
	return versionsFromModel(v), convertError(err)
}

// Create adds a version/milestone to the project.
//
// This method supports options returned by methods in "*Client.Project.Version.Option",
// such as:
//   - WithDescription
//   - WithReleaseDueDate
//   - WithStartDate
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-version-milestone/
func (s *ProjectVersionService) Create(ctx context.Context, projectIDOrKey, name string, opts ...RequestOption) (*Version, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, name, toCoreOptions(opts)...)
	return versionFromModel(v), convertError(err)
}

// Update updates a version/milestone.
//
// This method supports options returned by methods in "*Client.Project.Version.Option",
// such as:
//   - WithArchived
//   - WithDescription
//   - WithName
//   - WithReleaseDueDate
//   - WithStartDate
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-version-milestone/
func (s *ProjectVersionService) Update(ctx context.Context, projectIDOrKey string, versionID int, opts ...RequestOption) (*Version, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, versionID, toCoreOptions(opts)...)
	return versionFromModel(v), convertError(err)
}

// Delete deletes a version/milestone.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-version/
func (s *ProjectVersionService) Delete(ctx context.Context, projectIDOrKey string, versionID int) (*Version, error) {
	v, err := s.base.Delete(ctx, projectIDOrKey, versionID)
	return versionFromModel(v), convertError(err)
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
//  ProjectStatusOptionService
// ──────────────────────────────────────────────────────────────

// ProjectStatusOptionService provides a domain-specific set of option builders
// for operations within the ProjectStatusService.
type ProjectStatusOptionService struct {
	base *core.OptionService
}

// WithColor sets the status color.
func (s *ProjectStatusOptionService) WithColor(color string) RequestOption {
	return s.base.WithColor(color)
}

// WithName sets the status name.
func (s *ProjectStatusOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// ──────────────────────────────────────────────────────────────
//  ProjectWebhookOptionService
// ──────────────────────────────────────────────────────────────

// ProjectWebhookOptionService provides a domain-specific set of option builders
// for operations within the ProjectWebhookService.
type ProjectWebhookOptionService struct {
	base *core.OptionService
}

// WithActivityTypeIDs sets activity type IDs for webhook events.
func (s *ProjectWebhookOptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return s.base.WithActivityTypeIDs(typeIDs)
}

// WithAllEvent sets whether the webhook receives all events.
func (s *ProjectWebhookOptionService) WithAllEvent(enabled bool) RequestOption {
	return s.base.WithAllEvent(enabled)
}

// WithDescription sets the webhook description.
func (s *ProjectWebhookOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithHookURL sets the webhook URL.
func (s *ProjectWebhookOptionService) WithHookURL(hookURL string) RequestOption {
	return s.base.WithHookURL(hookURL)
}

// WithName sets the webhook name.
func (s *ProjectWebhookOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// ──────────────────────────────────────────────────────────────
//  ProjectVersionOptionService
// ──────────────────────────────────────────────────────────────

// ProjectVersionOptionService provides a domain-specific set of option builders
// for operations within the ProjectVersionService.
type ProjectVersionOptionService struct {
	base *core.OptionService
}

// WithArchived sets whether to include archived versions.
func (s *ProjectVersionOptionService) WithArchived(enabled bool) RequestOption {
	return s.base.WithArchived(enabled)
}

// WithDescription sets the version description.
func (s *ProjectVersionOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithName sets the version name.
func (s *ProjectVersionOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// WithReleaseDueDate sets the release due date.
func (s *ProjectVersionOptionService) WithReleaseDueDate(t time.Time) RequestOption {
	return s.base.WithReleaseDueDate(t)
}

// WithStartDate sets the version start date.
func (s *ProjectVersionOptionService) WithStartDate(t time.Time) RequestOption {
	return s.base.WithStartDate(t)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectService(method *core.Method, option *core.OptionService) *ProjectService {
	return &ProjectService{
		base:       project.NewService(method),
		Activity:   newProjectActivityService(method, option),
		Category:   newProjectCategoryService(method),
		Status:     newProjectStatusService(method, option),
		User:       newProjectUserService(method, option),
		SharedFile: newProjectSharedFileService(method),
		Webhook:    newProjectWebhookService(method, option),
		Version:    newProjectVersionService(method, option),
		Option:     newProjectOptionService(option),
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

func newProjectStatusService(method *core.Method, option *core.OptionService) *ProjectStatusService {
	return &ProjectStatusService{
		base:   project.NewStatusService(method),
		Option: newProjectStatusOptionService(option),
	}
}

func newProjectSharedFileService(method *core.Method) *ProjectSharedFileService {
	return &ProjectSharedFileService{
		base: sharedfile.NewProjectService(method),
	}
}

func newProjectWebhookService(method *core.Method, option *core.OptionService) *ProjectWebhookService {
	return &ProjectWebhookService{
		base:   webhook.NewService(method),
		Option: newWebhookOptionService(option),
	}
}

func newProjectVersionService(method *core.Method, option *core.OptionService) *ProjectVersionService {
	return &ProjectVersionService{
		base:   version.NewService(method),
		Option: newVersionOptionService(option),
	}
}

func newProjectOptionService(option *core.OptionService) *ProjectOptionService {
	return &ProjectOptionService{
		base: option,
	}
}

func newProjectStatusOptionService(option *core.OptionService) *ProjectStatusOptionService {
	return &ProjectStatusOptionService{base: option}
}

func newWebhookOptionService(option *core.OptionService) *ProjectWebhookOptionService {
	return &ProjectWebhookOptionService{base: option}
}

func newVersionOptionService(option *core.OptionService) *ProjectVersionOptionService {
	return &ProjectVersionOptionService{base: option}
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

func projectStatusFromModel(m *model.ProjectStatus) *ProjectStatus {
	if m == nil {
		return nil
	}
	return &ProjectStatus{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		Name:         m.Name,
		Color:        m.Color,
		DisplayOrder: m.DisplayOrder,
	}
}

func projectStatusesFromModel(ms []*model.ProjectStatus) []*ProjectStatus {
	if ms == nil {
		return nil
	}
	result := make([]*ProjectStatus, len(ms))
	for i, v := range ms {
		result[i] = projectStatusFromModel(v)
	}
	return result
}

func webhookFromModel(m *model.Webhook) *Webhook {
	if m == nil {
		return nil
	}

	return &Webhook{
		ID:              m.ID,
		Name:            m.Name,
		Description:     m.Description,
		HookURL:         m.HookURL,
		AllEvent:        m.AllEvent,
		ActivityTypeIDs: m.ActivityTypeIDs,
		CreatedUser:     userFromModel(m.CreatedUser),
		Created:         m.Created,
		UpdatedUser:     userFromModel(m.UpdatedUser),
		Updated:         m.Updated,
	}
}

func webhooksFromModel(ms []*model.Webhook) []*Webhook {
	if ms == nil {
		return nil
	}

	result := make([]*Webhook, len(ms))
	for i, v := range ms {
		result[i] = webhookFromModel(v)
	}
	return result
}
