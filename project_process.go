package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  Status models
// ──────────────────────────────────────────────────────────────

// Status represents a project status that can be assigned to issues.
type Status struct {
	ID           int
	ProjectID    int
	Name         string
	Color        string
	DisplayOrder int
}

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
	v, err := s.base.Update(ctx, projectIDOrKey, statusID, toCoreOption(option), toInterOptions(opts)...)
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
	v, err := s.base.List(ctx, projectIDOrKey, toInterOptions(opts)...)
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
	v, err := s.base.Add(ctx, projectIDOrKey, name, toInterOptions(opts)...)
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
	v, err := s.base.Update(ctx, projectIDOrKey, versionID, toCoreOption(option), toInterOptions(opts)...)
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
	base *option.OptionService
}

// WithColor sets the status color.
func (s *ProjectStatusOptionService) WithColor(color string) RequestOption {
	return &requestOption{opt: s.base.WithColor(color)}
}

// WithName sets the status name.
func (s *ProjectStatusOptionService) WithName(name string) RequestOption {
	return &requestOption{opt: s.base.WithName(name)}
}

// ──────────────────────────────────────────────────────────────
//  ProjectVersionOptionService
// ──────────────────────────────────────────────────────────────

// ProjectVersionOptionService provides a domain-specific set of option builders
// for operations within the ProjectVersionService.
type ProjectVersionOptionService struct {
	base *option.OptionService
}

// WithArchived sets whether to include archived versions.
func (s *ProjectVersionOptionService) WithArchived(enabled bool) RequestOption {
	return &requestOption{opt: s.base.WithArchived(enabled)}
}

// WithDescription sets the version description.
func (s *ProjectVersionOptionService) WithDescription(description string) RequestOption {
	return &requestOption{opt: s.base.WithDescription(description)}
}

// WithName sets the version name.
func (s *ProjectVersionOptionService) WithName(name string) RequestOption {
	return &requestOption{opt: s.base.WithName(name)}
}

// WithReleaseDueDate sets the release due date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *ProjectVersionOptionService) WithReleaseDueDate(date string) RequestOption {
	return &requestOption{opt: s.base.WithReleaseDueDate(date)}
}

// WithStartDate sets the version start date.
// The date must be formatted as "yyyy-MM-dd" (e.g. "2024-01-20").
func (s *ProjectVersionOptionService) WithStartDate(date string) RequestOption {
	return &requestOption{opt: s.base.WithStartDate(date)}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectStatusService(method *client.Method, option *option.OptionService) *ProjectStatusService {
	return &ProjectStatusService{
		base:   project.NewStatusService(method),
		Option: newProjectStatusOptionService(option),
	}
}

func newProjectVersionService(method *client.Method, option *option.OptionService) *ProjectVersionService {
	return &ProjectVersionService{
		base:   project.NewVersionService(method),
		Option: newVersionOptionService(option),
	}
}

func newProjectStatusOptionService(option *option.OptionService) *ProjectStatusOptionService {
	return &ProjectStatusOptionService{base: option}
}

func newVersionOptionService(option *option.OptionService) *ProjectVersionOptionService {
	return &ProjectVersionOptionService{base: option}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func statusFromModel(m *model.Status) *Status {
	if m == nil {
		return nil
	}
	return &Status{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		Name:         m.Name,
		Color:        m.Color,
		DisplayOrder: m.DisplayOrder,
	}
}

func statusesFromModel(ms []*model.Status) []*Status {
	if ms == nil {
		return nil
	}
	result := make([]*Status, len(ms))
	for i, v := range ms {
		result[i] = statusFromModel(v)
	}
	return result
}
