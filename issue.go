package backlog

import (
	"context"
	"iter"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/issue"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  Issue models
// ──────────────────────────────────────────────────────────────

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
	StartDate      Date
	DueDate        Date
	EstimatedHours float64
	ActualHours    float64
	ParentIssueID  int
	CreatedUser    *User
	Created        Timestamp
	UpdatedUser    *User
	Updated        Timestamp
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
	Comment    *IssueCommentService
	SharedFile *IssueSharedFileService
	Star       *IssueStarService

	Option *IssueOptionService
}

// List returns a list of issues.
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
func (s *IssueService) List(ctx context.Context, opts ...RequestOption) ([]*Issue, error) {
	v, err := s.base.List(ctx, toCoreOptions(opts)...)
	return issuesFromModel(v), convertError(err)
}

// All returns an iterator that lazily fetches all issues with automatic
// pagination, along with any validation error encountered at call time.
//
// perPage controls how many issues are fetched per API call (1-100).
// Iteration stops automatically when all issues have been returned.
// The caller must not pass WithCount or WithOffset in opts; those are managed
// internally. If they are passed, an error is returned immediately.
//
// This method supports filter options returned by methods in "*Client.Issue.Option",
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
func (s *IssueService) All(ctx context.Context, perPage int, opts ...RequestOption) (iter.Seq2[*Issue, error], error) {
	seq, err := s.base.All(ctx, perPage, toCoreOptions(opts)...)
	if err != nil {
		return nil, convertError(err)
	}
	return func(yield func(*Issue, error) bool) {
		for v, err := range seq {
			if !yield(issueFromModel(v), convertError(err)) {
				return
			}
		}
	}, nil
}

// Count returns the total count of issues matching the given filters.
//
// This method supports the same filter options as All, except WithIssueSort,
// WithOrder, WithOffset, and WithCount.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-issue
func (s *IssueService) Count(ctx context.Context, opts ...RequestOption) (int, error) {
	count, err := s.base.Count(ctx, toCoreOptions(opts)...)
	return count, convertError(err)
}

// One returns a single issue by its ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue
func (s *IssueService) One(ctx context.Context, issueIDOrKey string) (*Issue, error) {
	v, err := s.base.One(ctx, issueIDOrKey)
	return issueFromModel(v), convertError(err)
}

// Create creates a new issue.
//
// This method supports options returned by methods in "*Client.Issue.Option",
// such as:
//   - WithDescription
//   - WithStartDate
//   - WithDueDate
//   - WithEstimatedHours
//   - WithActualHours
//   - WithCategoryIDs
//   - WithVersionIDs
//   - WithMilestoneIDs
//   - WithAssigneeID
//   - WithParentIssueID
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//   - WithCustomField
//   - WithCustomFieldItems
//   - WithCustomFieldOther
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue
func (s *IssueService) Create(ctx context.Context, projectID int, summary string, issueTypeID int, priorityID int, opts ...RequestOption) (*Issue, error) {
	v, err := s.base.Create(ctx, projectID, summary, issueTypeID, priorityID, toCoreOptions(opts)...)
	return issueFromModel(v), convertError(err)
}

// Update updates an existing issue.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.Issue.Option", such as:
//   - WithSummary
//   - WithDescription
//   - WithIssueTypeID
//   - WithCategoryIDs
//   - WithVersionIDs
//   - WithMilestoneIDs
//   - WithStartDate
//   - WithDueDate
//   - WithEstimatedHours
//   - WithActualHours
//   - WithAssigneeID
//   - WithParentIssueID
//   - WithPriorityID
//   - WithStatusID
//   - WithResolutionID
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//   - WithComment
//   - WithCustomField
//   - WithCustomFieldItems
//   - WithCustomFieldOther
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue
func (s *IssueService) Update(ctx context.Context, issueIDOrKey string, option RequestOption, opts ...RequestOption) (*Issue, error) {
	v, err := s.base.Update(ctx, issueIDOrKey, option, toCoreOptions(opts)...)
	return issueFromModel(v), convertError(err)
}

// Delete deletes an issue by its ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue
func (s *IssueService) Delete(ctx context.Context, issueIDOrKey string) (*Issue, error) {
	v, err := s.base.Delete(ctx, issueIDOrKey)
	return issueFromModel(v), convertError(err)
}

// Participants returns a list of participants on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-participant-list
func (s *IssueService) Participants(ctx context.Context, issueIDOrKey string) ([]*User, error) {
	v, err := s.base.Participants(ctx, issueIDOrKey)
	return usersFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  IssueStarService
// ──────────────────────────────────────────────────────────────

// IssueStarService handles communication with the issue star-related methods of the Backlog API.
type IssueStarService struct {
	star *StarService
}

// Add adds a star to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *IssueStarService) Add(ctx context.Context, issueID int) error {
	return s.star.Add(ctx, s.star.Option.WithIssueID(issueID))
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *IssueStarService) Remove(ctx context.Context, starID int) error {
	return s.star.Remove(ctx, starID)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newIssueService(method *core.Method, option *core.OptionService) *IssueService {
	return &IssueService{
		base: issue.NewService(method),

		Attachment: newIssueAttachmentService(method),
		Comment:    newIssueCommentService(method, option),
		SharedFile: newIssueSharedFileService(method),
		Star:       newIssueStarService(method, option),

		Option: newIssueOptionService(option),
	}
}

func newIssueStarService(method *core.Method, option *core.OptionService) *IssueStarService {
	return &IssueStarService{star: newStarService(method, option)}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func resolutionsFromModel(m []*model.Resolution) []*Resolution {
	if m == nil {
		return nil
	}
	result := make([]*Resolution, len(m))
	for i, v := range m {
		if v == nil {
			result[i] = nil
		} else {
			result[i] = &Resolution{ID: v.ID, Name: v.Name}
		}
	}
	return result
}

func versionFromModel(m *model.Version) *Version {
	if m == nil {
		return nil
	}
	return &Version{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		Name:           m.Name,
		Description:    m.Description,
		StartDate:      Date{value: m.StartDate},
		ReleaseDueDate: Date{value: m.ReleaseDueDate},
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

func commentsFromModel(m []*model.Comment) []*Comment {
	if m == nil {
		return nil
	}
	result := make([]*Comment, len(m))
	for i, v := range m {
		result[i] = commentFromModel(v)
	}
	return result
}

func notificationsFromModel(m []*model.Notification) []*Notification {
	if m == nil {
		return nil
	}
	result := make([]*Notification, len(m))
	for i, v := range m {
		result[i] = notificationFromModel(v)
	}
	return result
}

func issueFromModel(m *model.Issue) *Issue {
	if m == nil {
		return nil
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
		IssueType:      issueTypeFromModel(m.IssueType),
		Summary:        m.Summary,
		Description:    m.Description,
		Resolutions:    resolutionsFromModel(m.Resolutions),
		Priority:       priority,
		Status:         statusFromModel(m.Status),
		Assignee:       userFromModel(m.Assignee),
		Category:       categoriesFromModel(m.Category),
		Versions:       versionsFromModel(m.Versions),
		Milestone:      versionsFromModel(m.Milestone),
		StartDate:      Date{value: m.StartDate},
		DueDate:        Date{value: m.DueDate},
		EstimatedHours: m.EstimatedHours,
		ActualHours:    m.ActualHours,
		ParentIssueID:  m.ParentIssueID,
		CreatedUser:    userFromModel(m.CreatedUser),
		Created:        Timestamp{m.Created},
		UpdatedUser:    userFromModel(m.UpdatedUser),
		Updated:        Timestamp{m.Updated},
		CustomFields:   customFieldsFromModel(m.CustomFields),
		Attachments:    attachmentsFromModel(m.Attachments),
		SharedFiles:    sharedFilesFromModel(m.SharedFiles),
		Stars:          starsFromModel(m.Stars),
	}
}

func issuesFromModel(m []*model.Issue) []*Issue {
	if m == nil {
		return nil
	}
	result := make([]*Issue, len(m))
	for i, v := range m {
		result[i] = issueFromModel(v)
	}
	return result
}
