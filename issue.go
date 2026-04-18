package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  Issue models
// ──────────────────────────────────────────────────────────────

// Category represents an issue category.
type Category struct {
	ID           int
	Name         string
	DisplayOrder int
}

// Issue represents a issue of Backlog.
type Issue struct {
	ID             int
	ProjectID      int
	IssueKey       string
	KeyID          int
	IssueType      *IssueType
	Summary        string
	Description    string
	Resolutions    []*Resolution
	Priority       *Priority
	Status         *Status
	Assignee       *User
	Category       []*Category
	Versions       []*Version
	Milestone      []*Version
	StartDate      time.Time
	DueDate        time.Time
	EstimatedHours int
	ActualHours    int
	ParentIssueID  int
	CreatedUser    *User
	Created        time.Time
	UpdatedUser    *User
	Updated        time.Time
	CustomFields   []*CustomField
	Attachments    []*Attachment
	SharedFiles    []*SharedFile
	Stars          []*Star
}

// IssueType represents type of Issue.
type IssueType struct {
	ID           int
	ProjectID    int
	Name         string
	Color        string
	DisplayOrder int
}

// Priority represents a priority.
type Priority struct {
	ID   int
	Name string
}

// Resolution represents a resolution.
type Resolution struct {
	ID   int
	Name string
}

// ──────────────────────────────────────────────────────────────
//  IssueService
// ──────────────────────────────────────────────────────────────

// IssueService handles communication with the issue-related methods of the Backlog API.
type IssueService struct {
	base *issue.Service

	Attachment *IssueAttachmentService
	Option     *IssueOptionService
}

// All returns a list of issues.
//
// This method supports options returned by methods in "*Client.Issue.Option",
// such as:
//   - WithProjectIDs
//   - WithIssueTypeIDs
//   - WithCategoryIDs
//   - WithVersionIDs
//   - WithMilestoneIDs
//   - WithStatusIDs
//   - WithPriorityIDs
//   - WithAssigneeIDs
//   - WithCreatedUserIDs
//   - WithResolutionIDs
//   - WithParentChild
//   - WithAttachment
//   - WithSharedFile
//   - WithIssueSort
//   - WithOrder
//   - WithOffset
//   - WithCount
//   - WithCreatedSince
//   - WithCreatedUntil
//   - WithUpdatedSince
//   - WithUpdatedUntil
//   - WithStartDateSince
//   - WithStartDateUntil
//   - WithDueDateSince
//   - WithDueDateUntil
//   - WithHasDueDate
//   - WithIDs
//   - WithParentIssueIDs
//   - WithKeyword
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-list
func (s *IssueService) All(ctx context.Context, opts ...RequestOption) ([]*Issue, error) {
	v, err := s.base.All(ctx, toCoreOptions(opts)...)
	return issuesFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  IssueAttachmentService
// ──────────────────────────────────────────────────────────────

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	base *attachment.IssueService
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*Attachment, error) {
	v, err := s.base.List(ctx, issueIDOrKey)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, issueIDOrKey, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  IssueOptionService
// ──────────────────────────────────────────────────────────────

// IssueOptionService provides a domain-specific set of option builders
// for operations within the IssueService.
type IssueOptionService struct {
	base *core.OptionService
}

// WithProjectIDs filters issues by project IDs.
func (s *IssueOptionService) WithProjectIDs(ids []int) RequestOption {
	return s.base.WithProjectIDs(ids)
}

// WithIssueTypeIDs filters issues by issue type IDs.
func (s *IssueOptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return s.base.WithIssueTypeIDs(ids)
}

// WithCategoryIDs filters issues by category IDs.
func (s *IssueOptionService) WithCategoryIDs(ids []int) RequestOption {
	return s.base.WithCategoryIDs(ids)
}

// WithVersionIDs filters issues by version IDs.
func (s *IssueOptionService) WithVersionIDs(ids []int) RequestOption {
	return s.base.WithVersionIDs(ids)
}

// WithMilestoneIDs filters issues by milestone IDs.
func (s *IssueOptionService) WithMilestoneIDs(ids []int) RequestOption {
	return s.base.WithMilestoneIDs(ids)
}

// WithStatusIDs filters issues by status IDs.
func (s *IssueOptionService) WithStatusIDs(ids []int) RequestOption {
	return s.base.WithStatusIDs(ids)
}

// WithPriorityIDs filters issues by priority IDs.
func (s *IssueOptionService) WithPriorityIDs(ids []int) RequestOption {
	return s.base.WithPriorityIDs(ids)
}

// WithAssigneeIDs filters issues by assignee user IDs.
func (s *IssueOptionService) WithAssigneeIDs(ids []int) RequestOption {
	return s.base.WithAssigneeIDs(ids)
}

// WithCreatedUserIDs filters issues by created user IDs.
func (s *IssueOptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return s.base.WithCreatedUserIDs(ids)
}

// WithResolutionIDs filters issues by resolution IDs.
func (s *IssueOptionService) WithResolutionIDs(ids []int) RequestOption {
	return s.base.WithResolutionIDs(ids)
}

// WithParentChild filters issues by subtask relationship.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *IssueOptionService) WithParentChild(parentChild int) RequestOption {
	return s.base.WithParentChild(parentChild)
}

// WithAttachment filters to include only issues with attachments.
func (s *IssueOptionService) WithAttachment(enabled bool) RequestOption {
	return s.base.WithAttachment(enabled)
}

// WithSharedFile filters to include only issues with shared files.
func (s *IssueOptionService) WithSharedFile(enabled bool) RequestOption {
	return s.base.WithSharedFile(enabled)
}

// WithIssueSort sets the field to sort issue list results by.
func (s *IssueOptionService) WithIssueSort(sort IssueSort) RequestOption {
	return s.base.WithIssueSort(model.IssueSort(sort))
}

// WithOrder sets the sort order of results.
func (s *IssueOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(model.Order(order))
}

// WithOffset sets the number of issues to skip.
func (s *IssueOptionService) WithOffset(offset int) RequestOption {
	return s.base.WithOffset(offset)
}

// WithCount sets the number of issues to retrieve (1-100).
func (s *IssueOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithCreatedSince filters issues created on or after the given date.
func (s *IssueOptionService) WithCreatedSince(t time.Time) RequestOption {
	return s.base.WithCreatedSince(t)
}

// WithCreatedUntil filters issues created on or before the given date.
func (s *IssueOptionService) WithCreatedUntil(t time.Time) RequestOption {
	return s.base.WithCreatedUntil(t)
}

// WithUpdatedSince filters issues updated on or after the given date.
func (s *IssueOptionService) WithUpdatedSince(t time.Time) RequestOption {
	return s.base.WithUpdatedSince(t)
}

// WithUpdatedUntil filters issues updated on or before the given date.
func (s *IssueOptionService) WithUpdatedUntil(t time.Time) RequestOption {
	return s.base.WithUpdatedUntil(t)
}

// WithStartDateSince filters issues with a start date on or after the given date.
func (s *IssueOptionService) WithStartDateSince(t time.Time) RequestOption {
	return s.base.WithStartDateSince(t)
}

// WithStartDateUntil filters issues with a start date on or before the given date.
func (s *IssueOptionService) WithStartDateUntil(t time.Time) RequestOption {
	return s.base.WithStartDateUntil(t)
}

// WithDueDateSince filters issues with a due date on or after the given date.
func (s *IssueOptionService) WithDueDateSince(t time.Time) RequestOption {
	return s.base.WithDueDateSince(t)
}

// WithDueDateUntil filters issues with a due date on or before the given date.
func (s *IssueOptionService) WithDueDateUntil(t time.Time) RequestOption {
	return s.base.WithDueDateUntil(t)
}

// WithHasDueDate filters to exclude issues without a due date.
// Note: Setting this to true is not supported by the Backlog API and will result in an error.
func (s *IssueOptionService) WithHasDueDate(enabled bool) RequestOption {
	return s.base.WithHasDueDate(enabled)
}

// WithIDs filters issues by issue IDs.
func (s *IssueOptionService) WithIDs(ids []int) RequestOption {
	return s.base.WithIDs(ids)
}

// WithParentIssueIDs filters issues by parent issue IDs.
func (s *IssueOptionService) WithParentIssueIDs(ids []int) RequestOption {
	return s.base.WithParentIssueIDs(ids)
}

// WithKeyword filters issues by keyword.
func (s *IssueOptionService) WithKeyword(keyword string) RequestOption {
	return s.base.WithKeyword(keyword)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newIssueService(method *core.Method, option *core.OptionService) *IssueService {
	return &IssueService{
		base:       issue.NewService(method),
		Attachment: newIssueAttachmentService(method),
		Option:     newIssueOptionService(option),
	}
}

func newIssueAttachmentService(method *core.Method) *IssueAttachmentService {
	return &IssueAttachmentService{
		base: attachment.NewIssueService(method),
	}
}

func newIssueOptionService(option *core.OptionService) *IssueOptionService {
	return &IssueOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func versionFromModel(m *model.Version) *Version {
	if m == nil {
		return nil
	}
	return &Version{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		Name:           m.Name,
		Description:    m.Description,
		StartDate:      m.StartDate,
		ReleaseDueDate: m.ReleaseDueDate,
		Archived:       m.Archived,
		DisplayOrder:   m.DisplayOrder,
	}
}

func versionsFromModel(m []*model.Version) []*Version {
	if m == nil {
		return nil
	}
	result := make([]*Version, len(m))
	for i, v := range m {
		result[i] = versionFromModel(v)
	}
	return result
}

func issueFromModel(m *model.Issue) *Issue {
	if m == nil {
		return nil
	}
	resolutions := make([]*Resolution, len(m.Resolutions))
	for i, v := range m.Resolutions {
		if v == nil {
			continue
		}
		resolutions[i] = &Resolution{ID: v.ID, Name: v.Name}
	}
	categories := make([]*Category, len(m.Category))
	for i, v := range m.Category {
		if v == nil {
			continue
		}
		categories[i] = &Category{ID: v.ID, Name: v.Name, DisplayOrder: v.DisplayOrder}
	}
	customFields := make([]*CustomField, len(m.CustomFields))
	for i, v := range m.CustomFields {
		customFields[i] = customFieldFromModel(v)
	}
	sharedFiles := make([]*SharedFile, len(m.SharedFiles))
	for i, v := range m.SharedFiles {
		sharedFiles[i] = sharedFileFromModel(v)
	}
	stars := make([]*Star, len(m.Stars))
	for i, v := range m.Stars {
		stars[i] = starFromModel(v)
	}
	var issueType *IssueType
	if m.IssueType != nil {
		issueType = &IssueType{
			ID:           m.IssueType.ID,
			ProjectID:    m.IssueType.ProjectID,
			Name:         m.IssueType.Name,
			Color:        m.IssueType.Color,
			DisplayOrder: m.IssueType.DisplayOrder,
		}
	}
	var priority *Priority
	if m.Priority != nil {
		priority = &Priority{ID: m.Priority.ID, Name: m.Priority.Name}
	}
	return &Issue{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		IssueKey:       m.IssueKey,
		KeyID:          m.KeyID,
		IssueType:      issueType,
		Summary:        m.Summary,
		Description:    m.Description,
		Resolutions:    resolutions,
		Priority:       priority,
		Status:         statusFromModel(m.Status),
		Assignee:       userFromModel(m.Assignee),
		Category:       categories,
		Versions:       versionsFromModel(m.Versions),
		Milestone:      versionsFromModel(m.Milestone),
		StartDate:      m.StartDate,
		DueDate:        m.DueDate,
		EstimatedHours: m.EstimatedHours,
		ActualHours:    m.ActualHours,
		ParentIssueID:  m.ParentIssueID,
		CreatedUser:    userFromModel(m.CreatedUser),
		Created:        m.Created,
		UpdatedUser:    userFromModel(m.UpdatedUser),
		Updated:        m.Updated,
		CustomFields:   customFields,
		Attachments:    attachmentsFromModel(m.Attachments),
		SharedFiles:    sharedFiles,
		Stars:          stars,
	}
}

func issuesFromModel(ms []*model.Issue) []*Issue {
	result := make([]*Issue, len(ms))
	for i, v := range ms {
		result[i] = issueFromModel(v)
	}
	return result
}
