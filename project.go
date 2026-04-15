package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/project"
)

// ProjectService handles communication with the project-related methods of the Backlog API.
type ProjectService struct {
	base *project.Service

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
func (s *ProjectService) All(ctx context.Context, opts ...core.RequestOption) ([]*model.Project, error) {
	return s.base.All(ctx, opts...)
}

// One returns one of the projects searched by ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project
func (s *ProjectService) One(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	return s.base.One(ctx, projectIDOrKey)
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
func (s *ProjectService) Create(ctx context.Context, key, name string, opts ...core.RequestOption) (*model.Project, error) {
	return s.base.Create(ctx, key, name, opts...)
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
func (s *ProjectService) Update(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) (*model.Project, error) {
	return s.base.Update(ctx, projectIDOrKey, opts...)
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *ProjectService) Delete(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	return s.base.Delete(ctx, projectIDOrKey)
}

// ProjectActivityService handles communication with the project activities-related methods of the Backlog API.
type ProjectActivityService struct {
	base *activity.ProjectService
}

// List returns a list of activities in the project.
//
// This method supports options returned by methods in "*Client.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Activity, error) {
	return s.base.List(ctx, projectIDOrKey, opts...)
}

// ProjectOptionService provides a domain-specific set of option builders
// for operations within the ProjectService.
type ProjectOptionService struct {
	base *core.OptionService
}

// WithAll sets whether to include all projects.
func (s *ProjectOptionService) WithAll(enabled bool) core.RequestOption {
	return s.base.WithAll(enabled)
}

// WithArchived sets whether to include archived projects.
func (s *ProjectOptionService) WithArchived(enabled bool) core.RequestOption {
	return s.base.WithArchived(enabled)
}

// WithChartEnabled sets whether the project uses a chart.
func (s *ProjectOptionService) WithChartEnabled(enabled bool) core.RequestOption {
	return s.base.WithChartEnabled(enabled)
}

// WithKey sets the project key.
func (s *ProjectOptionService) WithKey(key string) core.RequestOption {
	return s.base.WithKey(key)
}

// WithName sets the project name.
func (s *ProjectOptionService) WithName(name string) core.RequestOption {
	return s.base.WithName(name)
}

// WithProjectLeaderCanEditProjectLeader sets whether a project leader can edit other project leaders.
func (s *ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabled bool) core.RequestOption {
	return s.base.WithProjectLeaderCanEditProjectLeader(enabled)
}

// WithSubtaskingEnabled sets whether subtasking is enabled.
func (s *ProjectOptionService) WithSubtaskingEnabled(enabled bool) core.RequestOption {
	return s.base.WithSubtaskingEnabled(enabled)
}

// WithTextFormattingRule sets the text formatting rule.
func (s *ProjectOptionService) WithTextFormattingRule(format model.Format) core.RequestOption {
	return s.base.WithTextFormattingRule(format)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectService(method *core.Method, option *core.OptionService) *ProjectService {
	return &ProjectService{
		base:     project.NewService(method),
		Activity: newProjectActivityService(method),
		User:     newProjectUserService(method, option),
		Option:   newProjectOptionService(option),
	}
}

func newProjectActivityService(method *core.Method) *ProjectActivityService {
	return &ProjectActivityService{
		base: activity.NewProjectService(method),
	}
}

func newProjectOptionService(option *core.OptionService) *ProjectOptionService {
	return &ProjectOptionService{
		base: option,
	}
}
