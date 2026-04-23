package star

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
)

// Service handles communication with the star-related methods of the Backlog API.
type Service struct {
	method *core.Method
}

// Add adds a star to a resource (issue, comment, wiki page, pull request, or pull request comment).
//
// Exactly one of the following options must be provided:
//   - WithIssueID
//   - WithCommentID
//   - WithWikiPageID
//   - WithPullRequestID
//   - WithPullRequestCommentID
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *Service) Add(ctx context.Context, opts ...core.RequestOption) error {
	validOptions := []core.APIParamOptionType{
		core.ParamIssueID,
		core.ParamCommentID,
		core.ParamWikiID,
		core.ParamPullRequestID,
		core.ParamPullRequestCommentID,
	}
	requiredOptions := []core.APIParamOptionType{
		core.ParamIssueID,
		core.ParamCommentID,
		core.ParamWikiID,
		core.ParamPullRequestID,
		core.ParamPullRequestCommentID,
	}

	if !core.HasRequiredOption(opts, requiredOptions) {
		return core.NewValidationError("one of issueId, commentId, wikiPageId, pullRequestId, or pullRequestCommentId is required")
	}

	form := url.Values{}
	if err := core.ApplyOptions(form, validOptions, opts...); err != nil {
		return err
	}

	_, err := s.method.Post(ctx, "stars", form)
	if err != nil {
		return err
	}

	return nil
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *Service) Remove(ctx context.Context, starID int) error {
	option := &core.OptionService{}
	form := url.Values{}
	withStarID := option.WithStarID(starID)
	if err := withStarID.Check(); err != nil {
		return err
	}
	if err := withStarID.Set(form); err != nil {
		return err
	}

	_, err := s.method.Delete(ctx, "stars", form)
	if err != nil {
		return err
	}

	return nil
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewService creates and returns a new star Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
