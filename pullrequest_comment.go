package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/pullrequest"
)

// PullRequestCommentService handles communication with the pull request comment-related methods of the Backlog API.
type PullRequestCommentService struct {
	base   *pullrequest.CommentService
	Option *PullRequestCommentOptionService
}

// List returns a list of comments on a pull request.
//
// This method supports options returned by methods in "*Client.PullRequest.Comment.Option",
// such as:
//   - WithMinID
//   - WithMaxID
//   - WithCount
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-comment
func (s *PullRequestCommentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, opts ...RequestOption) ([]*Comment, error) {
	v, err := s.base.List(ctx, projectIDOrKey, repositoryIDOrName, prNumber, toCoreOptions(opts)...)
	return commentsFromModel(v), convertError(err)
}

// Add adds a comment to a pull request.
//
// This method supports options returned by methods in "*Client.PullRequest.Comment.Option",
// such as:
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment
func (s *PullRequestCommentService) Add(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, content string, opts ...RequestOption) (*Comment, error) {
	v, err := s.base.Add(ctx, projectIDOrKey, repositoryIDOrName, prNumber, content, toCoreOptions(opts)...)
	return commentFromModel(v), convertError(err)
}

// Count returns the number of comments on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-request-comments
func (s *PullRequestCommentService) Count(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) (int, error) {
	count, err := s.base.Count(ctx, projectIDOrKey, repositoryIDOrName, prNumber)
	return count, convertError(err)
}

// Update updates a comment on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-pull-request-comment-information
func (s *PullRequestCommentService) Update(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, commentID int, content string) (*Comment, error) {
	v, err := s.base.Update(ctx, projectIDOrKey, repositoryIDOrName, prNumber, commentID, content)
	return commentFromModel(v), convertError(err)
}

// PullRequestCommentOptionService provides a domain-specific set of option builders
// for operations within the PullRequestCommentService.
type PullRequestCommentOptionService struct {
	base *core.OptionService
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *PullRequestCommentOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

// WithCount sets the number of comments to retrieve (1-100).
func (s *PullRequestCommentOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithMaxID filters comments with ID at or below the given value.
func (s *PullRequestCommentOptionService) WithMaxID(id int) RequestOption {
	return s.base.WithMaxID(id)
}

// WithMinID filters comments with ID at or above the given value.
func (s *PullRequestCommentOptionService) WithMinID(id int) RequestOption {
	return s.base.WithMinID(id)
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *PullRequestCommentOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithOrder sets the sort order of results.
func (s *PullRequestCommentOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(string(order))
}

func newPullRequestCommentService(method *core.Method, option *core.OptionService) *PullRequestCommentService {
	return &PullRequestCommentService{
		base:   pullrequest.NewCommentService(method),
		Option: &PullRequestCommentOptionService{base: option},
	}
}
