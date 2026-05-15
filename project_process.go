package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
)

// ──────────────────────────────────────────────────────────────
//  ProjectStatusService
// ──────────────────────────────────────────────────────────────

// ProjectStatusService handles communication with the project status-related methods of the Backlog API.
type ProjectStatusService struct {
	base   *project.StatusService
	Option *ProjectStatusOptionService
}

// List returns a list of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project
func (s *ProjectStatusService) List(ctx context.Context, projectIDOrKey string) ([]*Status, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
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
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.Status.Option", such as:
//   - WithColor
//   - WithName
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
//  ProjectVersionService
// ──────────────────────────────────────────────────────────────

// ProjectVersionService handles communication with the project version/milestone-related methods of the Backlog API.
type ProjectVersionService struct {
	base   *project.VersionService
	Option *ProjectVersionOptionService
}

// List returns a list of versions/milestones in the project.
//
// This method supports options returned by methods in "*Client.Project.Version.Option",
// such as:
//   - WithArchived
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-version-milestone-list/
func (s *ProjectVersionService) List(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Version, error) {
	v, err := s.base.List(ctx, projectIDOrKey, toCoreOptions(opts)...)
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
// At least one option is required. This method supports options returned by
// methods in "*Client.Project.Version.Option", such as:
//   - WithArchived
//   - WithDescription
//   - WithName
//   - WithReleaseDueDate
//   - WithStartDate
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-version-milestone/
func (s *ProjectVersionService) Update(ctx context.Context, projectIDOrKey string, versionID int, option RequestOption, opts ...RequestOption) (*Version, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, versionID, option, toCoreOptions(opts)...)
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
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *ProjectVersionOptionService) WithReleaseDueDate(date string) RequestOption {
	return s.base.WithReleaseDueDate(date)
}

// WithStartDate sets the version start date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *ProjectVersionOptionService) WithStartDate(date string) RequestOption {
	return s.base.WithStartDate(date)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectStatusService(method *core.Method, option *core.OptionService) *ProjectStatusService {
	return &ProjectStatusService{
		base:   project.NewStatusService(method),
		Option: newProjectStatusOptionService(option),
	}
}

func newProjectVersionService(method *core.Method, option *core.OptionService) *ProjectVersionService {
	return &ProjectVersionService{
		base:   project.NewVersionService(method),
		Option: newVersionOptionService(option),
	}
}

func newProjectStatusOptionService(option *core.OptionService) *ProjectStatusOptionService {
	return &ProjectStatusOptionService{base: option}
}

func newVersionOptionService(option *core.OptionService) *ProjectVersionOptionService {
	return &ProjectVersionOptionService{base: option}
}
