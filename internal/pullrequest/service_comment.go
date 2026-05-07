package pullrequest

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/comment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// CommentService handles communication with the pull request comment-related methods of the Backlog API.
type CommentService struct {
	base   *comment.Service
	method *core.Method
}

// All returns a list of comments on a pull request.
//
// This method supports options:
//   - WithMinID
//   - WithMaxID
//   - WithCount
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-comment
func (s *CommentService) All(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, opts ...core.RequestOption) ([]*model.Comment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.All(ctx, spath, opts...)
}

// Add adds a comment to a pull request.
//
// This method supports options:
//   - WithNotifiedUserIDs
//   - WithAttachmentIDs
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment
func (s *CommentService) Add(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, content string, opts ...core.RequestOption) (*model.Comment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.Add(ctx, spath, content, opts...)
}

// Count returns the number of comments on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-request-comments
func (s *CommentService) Count(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int) (int, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return 0, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return 0, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return 0, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", "count")
	return s.base.Count(ctx, spath)
}

// TODO
// func (s *CommentService) One(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int) (*model.Comment, error) {
// 	return nil, nil
// }

// Update updates a comment on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-pull-request-comment-information
func (s *CommentService) Update(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, commentID int, content string) (*model.Comment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", strconv.Itoa(commentID))
	return s.base.Update(ctx, spath, content)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewCommentService creates and returns a new pullrequest CommentService.
func NewCommentService(method *core.Method) *CommentService {
	return &CommentService{
		base:   comment.NewService(method),
		method: method,
	}
}
