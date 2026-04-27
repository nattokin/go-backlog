package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/comment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/sharedfile"
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
	Comment    *IssueCommentService
	SharedFile *IssueSharedFileService
	Star       *IssueStarService

	Option *IssueOptionService
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
//   - WithStatusID
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
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
//  IssueCommentService
// ──────────────────────────────────────────────────────────────

// IssueCommentService handles communication with the issue comment-related methods of the Backlog API.
type IssueCommentService struct {
	base   *comment.IssueService
	Option *IssueCommentOptionService
}

// All returns a list of comments on an issue.
//
// This method supports options returned by methods in "*Client.Issue.Comment.Option",
// such as:
//   - WithMinID
//   - WithMaxID
//   - WithCount
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment-list
func (s *IssueCommentService) All(ctx context.Context, issueIDOrKey string, opts ...RequestOption) ([]*Comment, error) {
	v, err := s.base.All(ctx, issueIDOrKey, toCoreOptions(opts)...)
	return commentsFromModel(v), convertError(err)
}

// Add adds a comment to an issue.
//
// This method supports options returned by methods in "*Client.Issue.Comment.Option",
// such as:
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment
func (s *IssueCommentService) Add(ctx context.Context, issueIDOrKey string, content string, opts ...RequestOption) (*Comment, error) {
	v, err := s.base.Add(ctx, issueIDOrKey, content, toCoreOptions(opts)...)
	return commentFromModel(v), convertError(err)
}

// Count returns the number of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-comment
func (s *IssueCommentService) Count(ctx context.Context, issueIDOrKey string) (int, error) {
	count, err := s.base.Count(ctx, issueIDOrKey)
	return count, convertError(err)
}

// One returns information about a specific comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment
func (s *IssueCommentService) One(ctx context.Context, issueIDOrKey string, commentID int) (*Comment, error) {
	v, err := s.base.One(ctx, issueIDOrKey, commentID)
	return commentFromModel(v), convertError(err)
}

// Delete deletes a comment from an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-comment
func (s *IssueCommentService) Delete(ctx context.Context, issueIDOrKey string, commentID int) (*Comment, error) {
	v, err := s.base.Delete(ctx, issueIDOrKey, commentID)
	return commentFromModel(v), convertError(err)
}

// Update updates a comment on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-comment
func (s *IssueCommentService) Update(ctx context.Context, issueIDOrKey string, commentID int, content string) (*Comment, error) {
	v, err := s.base.Update(ctx, issueIDOrKey, commentID, content)
	return commentFromModel(v), convertError(err)
}

// Notifications returns a list of notifications on a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-comment-notifications
func (s *IssueCommentService) Notifications(ctx context.Context, issueIDOrKey string, commentID int) ([]*Notification, error) {
	v, err := s.base.Notifications(ctx, issueIDOrKey, commentID)
	return notificationsFromModel(v), convertError(err)
}

// Notify sends notifications for a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment-notification
func (s *IssueCommentService) Notify(ctx context.Context, issueIDOrKey string, commentID int, userIDs []int) (*Comment, error) {
	v, err := s.base.Notify(ctx, issueIDOrKey, commentID, userIDs)
	return commentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  IssueCommentOptionService
// ──────────────────────────────────────────────────────────────

// IssueCommentOptionService provides a domain-specific set of option builders
// for operations within the IssueCommentService.
type IssueCommentOptionService struct {
	base *core.OptionService
}

// WithCount sets the number of comments to retrieve (1-100).
func (s *IssueCommentOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithMaxID filters comments with ID at or below the given value.
func (s *IssueCommentOptionService) WithMaxID(id int) RequestOption {
	return s.base.WithMaxID(id)
}

// WithMinID filters comments with ID at or above the given value.
func (s *IssueCommentOptionService) WithMinID(id int) RequestOption {
	return s.base.WithMinID(id)
}

// WithOrder sets the sort order of results.
func (s *IssueCommentOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(model.Order(order))
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *IssueCommentOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *IssueCommentOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

// ──────────────────────────────────────────────────────────────
//  IssueSharedFileService
// ──────────────────────────────────────────────────────────────

// IssueSharedFileService handles communication with the issue shared-file-related methods of the Backlog API.
type IssueSharedFileService struct {
	base *sharedfile.IssueService
}

// List returns a list of shared files linked to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-linked-shared-files
func (s *IssueSharedFileService) List(ctx context.Context, issueIDOrKey string) ([]*SharedFile, error) {
	v, err := s.base.List(ctx, issueIDOrKey)
	return sharedFilesFromModel(v), convertError(err)
}

// Link links shared files to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue
func (s *IssueSharedFileService) Link(ctx context.Context, issueIDOrKey string, fileIDs []int) ([]*SharedFile, error) {
	v, err := s.base.Link(ctx, issueIDOrKey, fileIDs)
	return sharedFilesFromModel(v), convertError(err)
}

// Unlink removes a shared file link from the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue
func (s *IssueSharedFileService) Unlink(ctx context.Context, issueIDOrKey string, fileID int) (*SharedFile, error) {
	v, err := s.base.Unlink(ctx, issueIDOrKey, fileID)
	return sharedFileFromModel(v), convertError(err)
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
//  IssueOptionService
// ──────────────────────────────────────────────────────────────

// IssueOptionService provides a domain-specific set of option builders
// for operations within the IssueService.
type IssueOptionService struct {
	base *core.OptionService
}

// WithActualHours returns an option to set the `actualHours` parameter.
func (s *IssueOptionService) WithActualHours(hours int) RequestOption {
	return s.base.WithActualHours(hours)
}

// WithAssigneeID returns an option to set the `assigneeId` parameter.
func (s *IssueOptionService) WithAssigneeID(id int) RequestOption {
	return s.base.WithAssigneeID(id)
}

// WithAssigneeIDs filters issues by assignee user IDs.
func (s *IssueOptionService) WithAssigneeIDs(ids []int) RequestOption {
	return s.base.WithAssigneeIDs(ids)
}

// WithAttachment filters to include only issues with attachments.
func (s *IssueOptionService) WithAttachment(enabled bool) RequestOption {
	return s.base.WithAttachment(enabled)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *IssueOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

// WithCategoryIDs filters issues by category IDs.
func (s *IssueOptionService) WithCategoryIDs(ids []int) RequestOption {
	return s.base.WithCategoryIDs(ids)
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

// WithCreatedUserIDs filters issues by created user IDs.
func (s *IssueOptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return s.base.WithCreatedUserIDs(ids)
}

// WithDescription returns an option to set the `description` parameter.
func (s *IssueOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithDueDate returns an option to set the `dueDate` parameter.
func (s *IssueOptionService) WithDueDate(t time.Time) RequestOption {
	return s.base.WithDueDate(t)
}

// WithDueDateSince filters issues with a due date on or after the given date.
func (s *IssueOptionService) WithDueDateSince(t time.Time) RequestOption {
	return s.base.WithDueDateSince(t)
}

// WithDueDateUntil filters issues with a due date on or before the given date.
func (s *IssueOptionService) WithDueDateUntil(t time.Time) RequestOption {
	return s.base.WithDueDateUntil(t)
}

// WithEstimatedHours returns an option to set the `estimatedHours` parameter.
func (s *IssueOptionService) WithEstimatedHours(hours int) RequestOption {
	return s.base.WithEstimatedHours(hours)
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

// WithIssueSort sets the field to sort issue list results by.
func (s *IssueOptionService) WithIssueSort(sort IssueSort) RequestOption {
	return s.base.WithIssueSort(model.IssueSort(sort))
}

// WithIssueTypeID returns an option to set the `issueTypeId` parameter.
func (s *IssueOptionService) WithIssueTypeID(id int) RequestOption {
	return s.base.WithIssueTypeID(id)
}

// WithIssueTypeIDs filters issues by issue type IDs.
func (s *IssueOptionService) WithIssueTypeIDs(ids []int) RequestOption {
	return s.base.WithIssueTypeIDs(ids)
}

// WithKeyword filters issues by keyword.
func (s *IssueOptionService) WithKeyword(keyword string) RequestOption {
	return s.base.WithKeyword(keyword)
}

// WithMilestoneIDs filters issues by milestone IDs.
func (s *IssueOptionService) WithMilestoneIDs(ids []int) RequestOption {
	return s.base.WithMilestoneIDs(ids)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *IssueOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithOffset sets the number of issues to skip.
func (s *IssueOptionService) WithOffset(offset int) RequestOption {
	return s.base.WithOffset(offset)
}

// WithOrder sets the sort order of results.
func (s *IssueOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(model.Order(order))
}

// WithParentChild filters issues by subtask relationship.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *IssueOptionService) WithParentChild(parentChild int) RequestOption {
	return s.base.WithParentChild(parentChild)
}

// WithParentIssueID returns an option to set the `parentIssueId` parameter.
func (s *IssueOptionService) WithParentIssueID(id int) RequestOption {
	return s.base.WithParentIssueID(id)
}

// WithParentIssueIDs filters issues by parent issue IDs.
func (s *IssueOptionService) WithParentIssueIDs(ids []int) RequestOption {
	return s.base.WithParentIssueIDs(ids)
}

// WithPriorityID returns an option to set the `priorityId` parameter.
func (s *IssueOptionService) WithPriorityID(id int) RequestOption {
	return s.base.WithPriorityID(id)
}

// WithPriorityIDs filters issues by priority IDs.
func (s *IssueOptionService) WithPriorityIDs(ids []int) RequestOption {
	return s.base.WithPriorityIDs(ids)
}

// WithProjectIDs filters issues by project IDs.
func (s *IssueOptionService) WithProjectIDs(ids []int) RequestOption {
	return s.base.WithProjectIDs(ids)
}

// WithResolutionID returns an option to set the `resolutionId` parameter.
func (s *IssueOptionService) WithResolutionID(id int) RequestOption {
	return s.base.WithResolutionID(id)
}

// WithResolutionIDs filters issues by resolution IDs.
func (s *IssueOptionService) WithResolutionIDs(ids []int) RequestOption {
	return s.base.WithResolutionIDs(ids)
}

// WithSharedFile filters to include only issues with shared files.
func (s *IssueOptionService) WithSharedFile(enabled bool) RequestOption {
	return s.base.WithSharedFile(enabled)
}

// WithStartDate returns an option to set the `startDate` parameter.
func (s *IssueOptionService) WithStartDate(t time.Time) RequestOption {
	return s.base.WithStartDate(t)
}

// WithStartDateSince filters issues with a start date on or after the given date.
func (s *IssueOptionService) WithStartDateSince(t time.Time) RequestOption {
	return s.base.WithStartDateSince(t)
}

// WithStartDateUntil filters issues with a start date on or before the given date.
func (s *IssueOptionService) WithStartDateUntil(t time.Time) RequestOption {
	return s.base.WithStartDateUntil(t)
}

// WithStatusID returns an option to set the `statusId` parameter.
func (s *IssueOptionService) WithStatusID(id int) RequestOption {
	return s.base.WithStatusID(id)
}

// WithStatusIDs filters issues by status IDs.
func (s *IssueOptionService) WithStatusIDs(ids []int) RequestOption {
	return s.base.WithStatusIDs(ids)
}

// WithSummary returns an option to set the `summary` parameter.
func (s *IssueOptionService) WithSummary(summary string) RequestOption {
	return s.base.WithSummary(summary)
}

// WithUpdatedSince filters issues updated on or after the given date.
func (s *IssueOptionService) WithUpdatedSince(t time.Time) RequestOption {
	return s.base.WithUpdatedSince(t)
}

// WithUpdatedUntil filters issues updated on or before the given date.
func (s *IssueOptionService) WithUpdatedUntil(t time.Time) RequestOption {
	return s.base.WithUpdatedUntil(t)
}

// WithVersionIDs filters issues by version IDs.
func (s *IssueOptionService) WithVersionIDs(ids []int) RequestOption {
	return s.base.WithVersionIDs(ids)
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

func newIssueAttachmentService(method *core.Method) *IssueAttachmentService {
	return &IssueAttachmentService{
		base: attachment.NewIssueService(method),
	}
}

func newIssueCommentService(method *core.Method, option *core.OptionService) *IssueCommentService {
	return &IssueCommentService{
		base:   comment.NewIssueService(method),
		Option: &IssueCommentOptionService{base: option},
	}
}

func newIssueSharedFileService(method *core.Method) *IssueSharedFileService {
	return &IssueSharedFileService{
		base: sharedfile.NewIssueService(method),
	}
}

func newIssueStarService(method *core.Method, option *core.OptionService) *IssueStarService {
	return &IssueStarService{star: newStarService(method, option)}
}

func newIssueOptionService(option *core.OptionService) *IssueOptionService {
	return &IssueOptionService{base: option}
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
		Resolutions:    resolutionsFromModel(m.Resolutions),
		Priority:       priority,
		Status:         statusFromModel(m.Status),
		Assignee:       userFromModel(m.Assignee),
		Category:       categoriesFromModel(m.Category),
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
