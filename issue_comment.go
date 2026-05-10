package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
)

// IssueCommentService handles communication with the issue comment-related methods of the Backlog API.
type IssueCommentService struct {
	base   *issue.CommentService
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
	return s.base.WithOrder(string(order))
}

// WithNotifiedUserIDs returns an option to set multiple `notifiedUserId[]` parameters.
func (s *IssueCommentOptionService) WithNotifiedUserIDs(ids []int) RequestOption {
	return s.base.WithNotifiedUserIDs(ids)
}

// WithAttachmentIDs returns an option to set multiple `attachmentId[]` parameters.
func (s *IssueCommentOptionService) WithAttachmentIDs(ids []int) RequestOption {
	return s.base.WithAttachmentIDs(ids)
}

func newIssueCommentService(method *core.Method, option *core.OptionService) *IssueCommentService {
	return &IssueCommentService{
		base:   issue.NewCommentService(method),
		Option: &IssueCommentOptionService{base: option},
	}
}
