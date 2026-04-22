package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/star"
)

// ──────────────────────────────────────────────────────────────
//  StarService
// ──────────────────────────────────────────────────────────────

// StarService handles communication with the star-related methods of the Backlog API.
type StarService struct {
	base   *star.Service
	option *core.OptionService

	Option *StarOptionService
}

// Add adds a star to a resource.
//
// Exactly one of the following options returned by methods in "*Client.Star.Option"
// must be provided:
//   - WithIssueID
//   - WithCommentID
//   - WithWikiPageID
//   - WithPullRequestID
//   - WithPullRequestCommentID
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *StarService) Add(ctx context.Context, opts ...RequestOption) error {
	return convertError(s.base.Add(ctx, toCoreOptions(opts)...))
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *StarService) Remove(ctx context.Context, starID int) error {
	return convertError(s.base.Remove(ctx, starID))
}

// ──────────────────────────────────────────────────────────────
//  StarOptionService
// ──────────────────────────────────────────────────────────────

// StarOptionService provides a domain-specific set of option builders
// for operations within the StarService.
type StarOptionService struct {
	base *core.OptionService
}

// WithCommentID sets the comment ID to add a star to.
func (s *StarOptionService) WithCommentID(id int) RequestOption {
	return s.base.WithCommentID(id)
}

// WithIssueID sets the issue ID to add a star to.
func (s *StarOptionService) WithIssueID(id int) RequestOption {
	return s.base.WithIssueID(id)
}

// WithPullRequestCommentID sets the pull request comment ID to add a star to.
func (s *StarOptionService) WithPullRequestCommentID(id int) RequestOption {
	return s.base.WithPullRequestCommentID(id)
}

// WithPullRequestID sets the pull request ID to add a star to.
func (s *StarOptionService) WithPullRequestID(id int) RequestOption {
	return s.base.WithPullRequestID(id)
}

// WithWikiPageID sets the wiki page ID to add a star to.
func (s *StarOptionService) WithWikiPageID(id int) RequestOption {
	return s.base.WithWikiPageID(id)
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

func newStarService(method *core.Method, option *core.OptionService) *StarService {
	return &StarService{
		base:   star.NewService(method),
		option: option,
		Option: &StarOptionService{base: option},
	}
}
