package backlog

import (
	"context"
	"iter"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  PullRequest models
// ──────────────────────────────────────────────────────────────

// PullRequest represents a pull request.
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
	MergeCommit  string
	CloseAt      Timestamp
	MergeAt      Timestamp
	CreatedUser  *User
	Created      Timestamp
	UpdatedUser  *User
	Updated      Timestamp
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
	Comment    *PullRequestCommentService
	Star       *PullRequestStarService

	Option *PullRequestOptionService
}

// List returns a list of pull requests.
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
func (s *PullRequestService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, opts ...RequestOption) ([]*PullRequest, error) {
	v, err := s.base.List(ctx, projectIDOrKey, repositoryIDOrName, toCoreOptions(opts)...)
	return pullRequestsFromModel(v), convertError(err)
}

// All returns an iterator that lazily fetches all pull requests with automatic
// pagination, along with any validation error encountered at call time.
//
// perPage controls how many pull requests are fetched per API call (1-100).
// Iteration stops automatically when all pull requests have been returned.
// The caller must not pass WithCount or WithOffset in opts; those are managed
// internally. If they are passed, an error is returned immediately.
//
// This method supports filter options returned by methods in "*Client.PullRequest.Option",
// such as:
//   - WithStatusIDs
//   - WithAssigneeIDs
//   - WithIssueIDs
//   - WithCreatedUserIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *PullRequestService) All(ctx context.Context, perPage int, projectIDOrKey string, repositoryIDOrName string, opts ...RequestOption) (iter.Seq2[*PullRequest, error], error) {
	seq, err := s.base.All(ctx, perPage, projectIDOrKey, repositoryIDOrName, toCoreOptions(opts)...)
	if err != nil {
		return nil, convertError(err)
	}
	return func(yield func(*PullRequest, error) bool) {
		for v, err := range seq {
			if !yield(pullRequestFromModel(v), convertError(err)) {
				return
			}
		}
	}, nil
}

// Count returns the number of pull requests.
//
// This method supports the same filter options as All, except WithOffset and WithCount.
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
//  Constructors
// ──────────────────────────────────────────────────────────────

func newPullRequestService(method *core.Method, option *core.OptionService) *PullRequestService {
	return &PullRequestService{
		base: pullrequest.NewService(method),

		Attachment: newPullRequestAttachmentService(method),
		Comment:    newPullRequestCommentService(method, option),
		Star:       newPullRequestStarService(method, option),

		Option: newPullRequestOptionService(option),
	}
}

func newPullRequestStarService(method *core.Method, option *core.OptionService) *PullRequestStarService {
	return &PullRequestStarService{star: newStarService(method, option)}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func pullRequestFromModel(m *model.PullRequest) *PullRequest {
	if m == nil {
		return nil
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
		MergeCommit:  m.MergeCommit,
		CloseAt:      Timestamp{m.CloseAt},
		MergeAt:      Timestamp{m.MergeAt},
		CreatedUser:  userFromModel(m.CreatedUser),
		Created:      Timestamp{m.Created},
		UpdatedUser:  userFromModel(m.UpdatedUser),
		Updated:      Timestamp{m.Updated},
		Attachments:  attachmentsFromModel(m.Attachments),
		Stars:        starsFromModel(m.Stars),
	}
}

func pullRequestsFromModel(m []*model.PullRequest) []*PullRequest {
	if m == nil {
		return nil
	}
	result := make([]*PullRequest, len(m))
	for i, v := range m {
		result[i] = pullRequestFromModel(v)
	}
	return result
}
