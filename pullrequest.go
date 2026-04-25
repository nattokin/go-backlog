package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/pullrequest"
)

// PullRequest represents pull request of Backlog git.
type PullRequest struct {
	ID           int
	ProjectID    int
	RepositoryID int
	Number       int
	Summary      string
	Description  string
	Base         string
	Branch       string
	Status       *Status
	Assignee     *User
	Issue        *Issue
	BaseCommit   string
	BranchCommit string
	CloseAt      time.Time
	MergeAt      time.Time
	CreatedUser  *User
	Created      time.Time
	UpdatedUser  *User
	Updated      time.Time
	Attachments  []*Attachment
	Stars        []*Star
}

// ──────────────────────────────────────────────────────────────
//  PullRequestService
// ──────────────────────────────────────────────────────────────

// PullRequestService handles communication with the pull request-related methods of the Backlog API.
type PullRequestService struct {
	base *pullrequest.Service

	Attachment *PullRequestAttachmentService
	Option     *PullRequestOptionService
	Star       *PullRequestStarService
}

// All returns a list of pull requests.
//
// This method supports options returned by methods in "*Client.PullRequest.Option",
// such as:
//   - WithStatusIDs
//   - WithAssigneeIDs
//   - WithIssueIDs
//   - WithCreatedUserIDs
//   - WithOffset
//   - WithCount
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *PullRequestService) All(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, opts ...RequestOption) ([]*PullRequest, error) {
	v, err := s.base.All(ctx, projectIDOrKey, repositoryIDOrName, toCoreOptions(opts)...)
	return pullRequestsFromModel(v), convertError(err)
}

// Count returns the number of pull requests.
//
// This method supports options returned by methods in "*Client.PullRequest.Option",
// such as:
//   - WithStatusIDs
//   - WithAssigneeIDs
//   - WithIssueIDs
//   - WithCreatedUserIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests
func (s *PullRequestService) Count(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, opts ...RequestOption) (int, error) {
	count, err := s.base.Count(ctx, projectIDOrKey, repositoryIDOrName, toCoreOptions(opts)...)
	return count, convertError(err)
}

// One returns a single pull request by its number.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request
func (s *PullRequestService) One(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) (*PullRequest, error) {
	v, err := s.base.One(ctx, projectIDOrKey, repositoryIDOrName, prNumber)
	return pullRequestFromModel(v), convertError(err)
}

// Create creates a new pull request.
//
// This method supports options returned by methods in "*Client.PullRequest.Option",
// such as:
//   - WithIssueID
//   - WithAssigneeID
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request
func (s *PullRequestService) Create(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, summary string, description string, base string, branch string, opts ...RequestOption) (*PullRequest, error) {
	v, err := s.base.Create(ctx, projectIDOrKey, repositoryIDOrName, summary, description, base, branch, toCoreOptions(opts)...)
	return pullRequestFromModel(v), convertError(err)
}

// Update updates an existing pull request.
//
// At least one option is required. This method supports options returned by
// methods in "*Client.PullRequest.Option", such as:
//   - WithSummary
//   - WithDescription
//   - WithIssueID
//   - WithAssigneeID
//   - WithNotifiedUserIDs
//   - WithComment
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-pull-request
func (s *PullRequestService) Update(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, option RequestOption, opts ...RequestOption) (*PullRequest, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, repositoryIDOrName, prNumber, option, toCoreOptions(opts)...)
	return pullRequestFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  PullRequestAttachmentService
// ──────────────────────────────────────────────────────────────

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	base *attachment.PullRequestService
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*Attachment, error) {
	v, err := s.base.List(ctx, projectIDOrKey, repositoryIDOrName, prNumber)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, projectIDOrKey, repositoryIDOrName, prNumber, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  PullRequestStarService
// ──────────────────────────────────────────────────────────────

// PullRequestStarService handles communication with the pull request star-related methods of the Backlog API.
type PullRequestStarService struct {
	star *StarService
}

// Add adds a star to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *PullRequestStarService) Add(ctx context.Context, pullRequestID int) error {
	return s.star.Add(ctx, s.star.Option.WithPullRequestID(pullRequestID))
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *PullRequestStarService) Remove(ctx context.Context, starID int) error {
	return s.star.Remove(ctx, starID)
}

// ──────────────────────────────────────────────────────────────
//  PullRequestOptionService
// ──────────────────────────────────────────────────────────────

// PullRequestOptionService provides a domain-specific set of option builders
// for operations within the PullRequestService.
type PullRequestOptionService struct {
	base *core.OptionService
}

// WithAssigneeID returns an option to set the `assigneeId` parameter.
func (s *PullRequestOptionService) WithAssigneeID(id int) RequestOption {
	return s.base.WithAssigneeID(id)
}

// WithAssigneeIDs filters pull requests by assignee user IDs.
func (s *PullRequestOptionService) WithAssigneeIDs(ids []int) RequestOption {
	return s.base.WithAssigneeIDs(ids)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *PullRequestOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

// WithComment returns an option to set the `comment` parameter.
func (s *PullRequestOptionService) WithComment(comment string) RequestOption {
	return s.base.WithComment(comment)
}

// WithCount sets the number of pull requests to retrieve.
func (s *PullRequestOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithCreatedUserIDs filters pull requests by created user IDs.
func (s *PullRequestOptionService) WithCreatedUserIDs(ids []int) RequestOption {
	return s.base.WithCreatedUserIDs(ids)
}

// WithDescription returns an option to set the `description` parameter.
func (s *PullRequestOptionService) WithDescription(description string) RequestOption {
	return s.base.WithDescription(description)
}

// WithIssueID returns an option to set the `issueId` parameter.
func (s *PullRequestOptionService) WithIssueID(id int) RequestOption {
	return s.base.WithIssueID(id)
}

// WithIssueIDs filters pull requests by issue IDs.
func (s *PullRequestOptionService) WithIssueIDs(ids []int) RequestOption {
	return s.base.WithIssueIDs(ids)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *PullRequestOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithOffset sets the number of pull requests to skip.
func (s *PullRequestOptionService) WithOffset(offset int) RequestOption {
	return s.base.WithOffset(offset)
}

// WithStatusIDs filters pull requests by status IDs.
func (s *PullRequestOptionService) WithStatusIDs(ids []int) RequestOption {
	return s.base.WithStatusIDs(ids)
}

// WithSummary returns an option to set the `summary` parameter.
func (s *PullRequestOptionService) WithSummary(summary string) RequestOption {
	return s.base.WithSummary(summary)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newPullRequestService(method *core.Method, option *core.OptionService) *PullRequestService {
	starSvc := newStarService(method, option)
	return &PullRequestService{
		base:       pullrequest.NewService(method),
		Attachment: newPullRequestAttachmentService(method),
		Option:     newPullRequestOptionService(option),
		Star:       newPullRequestStarService(starSvc),
	}
}

func newPullRequestAttachmentService(method *core.Method) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		base: attachment.NewPullRequestService(method),
	}
}

func newPullRequestStarService(starSvc *StarService) *PullRequestStarService {
	return &PullRequestStarService{star: starSvc}
}

func newPullRequestOptionService(option *core.OptionService) *PullRequestOptionService {
	return &PullRequestOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func pullRequestFromModel(m *model.PullRequest) *PullRequest {
	if m == nil {
		return nil
	}
	attachments := make([]*Attachment, len(m.Attachments))
	for i, v := range m.Attachments {
		attachments[i] = attachmentFromModel(v)
	}
	stars := make([]*Star, len(m.Stars))
	for i, v := range m.Stars {
		stars[i] = starFromModel(v)
	}
	return &PullRequest{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		RepositoryID: m.RepositoryID,
		Number:       m.Number,
		Summary:      m.Summary,
		Description:  m.Description,
		Base:         m.Base,
		Branch:       m.Branch,
		Status:       statusFromModel(m.Status),
		Assignee:     userFromModel(m.Assignee),
		Issue:        issueFromModel(m.Issue),
		BaseCommit:   m.BaseCommit,
		BranchCommit: m.BranchCommit,
		CloseAt:      m.CloseAt,
		MergeAt:      m.MergeAt,
		CreatedUser:  userFromModel(m.CreatedUser),
		Created:      m.Created,
		UpdatedUser:  userFromModel(m.UpdatedUser),
		Updated:      m.Updated,
		Attachments:  attachments,
		Stars:        stars,
	}
}

func pullRequestsFromModel(ms []*model.PullRequest) []*PullRequest {
	result := make([]*PullRequest, len(ms))
	for i, v := range ms {
		result[i] = pullRequestFromModel(v)
	}
	return result
}
