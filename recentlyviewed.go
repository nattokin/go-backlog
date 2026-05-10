package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/domain/recentlyviewed"
)

// ──────────────────────────────────────────────────────────────
//  UserRecentlyViewedService
// ──────────────────────────────────────────────────────────────

// RecentlyViewedService handles communication with the recently-viewed methods of the Backlog API.
// All endpoints are scoped to the authenticated user (myself), so no userID argument is needed.
type RecentlyViewedService struct {
	base *recentlyviewed.Service

	Option *RecentlyViewedOptionService
}

// ListIssues returns a list of issues recently viewed by the authenticated user.
//
// This method supports options returned by methods in "*Client.User.RecentlyViewed.Option",
// such as:
//   - WithCount
//   - WithOffset
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-issues
func (s *RecentlyViewedService) ListIssues(ctx context.Context, opts ...RequestOption) ([]*Issue, error) {
	v, err := s.base.ListIssues(ctx, toCoreOptions(opts)...)
	return issuesFromModel(v), convertError(err)
}

// AddIssue adds an issue to the recently viewed list of the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-issue
func (s *RecentlyViewedService) AddIssue(ctx context.Context, issueID int) (*Issue, error) {
	v, err := s.base.AddIssue(ctx, issueID)
	return issueFromModel(v), convertError(err)
}

// ListProjects returns a list of projects recently viewed by the authenticated user.
//
// This method supports options returned by methods in "*Client.User.RecentlyViewed.Option",
// such as:
//   - WithCount
//   - WithOffset
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-projects
func (s *RecentlyViewedService) ListProjects(ctx context.Context, opts ...RequestOption) ([]*Project, error) {
	v, err := s.base.ListProjects(ctx, toCoreOptions(opts)...)
	return projectsFromModel(v), convertError(err)
}

// ListWikis returns a list of Wiki pages recently viewed by the authenticated user.
//
// This method supports options returned by methods in "*Client.User.RecentlyViewed.Option",
// such as:
//   - WithCount
//   - WithOffset
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-wikis
func (s *RecentlyViewedService) ListWikis(ctx context.Context, opts ...RequestOption) ([]*Wiki, error) {
	v, err := s.base.ListWikis(ctx, toCoreOptions(opts)...)
	return wikisFromModel(v), convertError(err)
}

// AddWiki adds a Wiki page to the recently viewed list of the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-wiki
func (s *RecentlyViewedService) AddWiki(ctx context.Context, wikiID int) (*Wiki, error) {
	v, err := s.base.AddWiki(ctx, wikiID)
	return wikiFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  UserRecentlyViewedOptionService
// ──────────────────────────────────────────────────────────────

// RecentlyViewedOptionService provides a domain-specific set of option builders
// for operations within the UserRecentlyViewedService.
type RecentlyViewedOptionService struct {
	base *core.OptionService
}

// WithCount sets the number of results to return (1-100).
func (s *RecentlyViewedOptionService) WithCount(count int) RequestOption {
	return s.base.WithCount(count)
}

// WithOffset sets the number of items to skip.
func (s *RecentlyViewedOptionService) WithOffset(offset int) RequestOption {
	return s.base.WithOffset(offset)
}

// WithOrder sets the sort order of results.
func (s *RecentlyViewedOptionService) WithOrder(order Order) RequestOption {
	return s.base.WithOrder(string(order))
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newRecentlyViewedService(method *core.Method, option *core.OptionService) *RecentlyViewedService {
	return &RecentlyViewedService{
		base:   recentlyviewed.NewService(method),
		Option: newRecentlyViewedOptionService(option),
	}
}

func newRecentlyViewedOptionService(option *core.OptionService) *RecentlyViewedOptionService {
	return &RecentlyViewedOptionService{
		base: option,
	}
}
