package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/project"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  ProjectIssueTypeService
// ──────────────────────────────────────────────────────────────

// ProjectIssueTypeService handles communication with the project issue-type-related methods of the Backlog API.
type ProjectIssueTypeService struct {
	base   *project.IssueTypeService
	Option *ProjectIssueTypeOptionService
}

// List returns a list of issue types in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-type-list
func (s *ProjectIssueTypeService) List(ctx context.Context, projectIDOrKey string) ([]*IssueType, error) {
	v, err := s.base.List(ctx, projectIDOrKey)
	return issueTypesFromModel(v), convertError(err)
}

// Create adds a new issue type to a project.
//
// This method supports options returned by methods in "*Client.Project.IssueType.Option",
// such as:
//   - WithTemplateDescription
//   - WithTemplateSummary
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue-type
func (s *ProjectIssueTypeService) Create(ctx context.Context, projectIDOrKey, name, color string, opts ...RequestOption) (*IssueType, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, name, color, toCoreOptions(opts)...)
	return issueTypeFromModel(v), convertError(err)
}

// Update updates an issue type in a project.
//
// This method supports options returned by methods in "*Client.Project.IssueType.Option",
// such as:
//   - WithColor
//   - WithName
//   - WithTemplateDescription
//   - WithTemplateSummary
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
//  ProjectIssueTypeOptionService
// ──────────────────────────────────────────────────────────────

// ProjectIssueTypeOptionService provides a domain-specific set of option builders
// for operations within the ProjectIssueTypeService.
type ProjectIssueTypeOptionService struct {
	base *core.OptionService
}

// WithColor sets the issue type color.
func (s *ProjectIssueTypeOptionService) WithColor(color string) RequestOption {
	return s.base.WithColor(color)
}

// WithName sets the issue type name.
func (s *ProjectIssueTypeOptionService) WithName(name string) RequestOption {
	return s.base.WithName(name)
}

// WithTemplateDescription sets the default description template for new issues of this type.
func (s *ProjectIssueTypeOptionService) WithTemplateDescription(description string) RequestOption {
	return s.base.WithTemplateDescription(description)
}

// WithTemplateSummary sets the default summary template for new issues of this type.
func (s *ProjectIssueTypeOptionService) WithTemplateSummary(summary string) RequestOption {
	return s.base.WithTemplateSummary(summary)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newProjectIssueTypeService(method *core.Method, option *core.OptionService) *ProjectIssueTypeService {
	return &ProjectIssueTypeService{
		base:   project.NewIssueTypeService(method),
		Option: &ProjectIssueTypeOptionService{base: option},
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

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
