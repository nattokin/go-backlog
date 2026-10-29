package backlog

import (
	"context"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/domain/recentlyviewed"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ──────────────────────────────────────────────────────────────
//  RecentlyViewed models
// ──────────────────────────────────────────────────────────────

// RecentlyViewedIssue represents an entry in the recently viewed issues list.
// Updated is the time the issue was last viewed.
type RecentlyViewedIssue struct {
	Issue   *Issue
	Updated Timestamp
}

// RecentlyViewedProject represents an entry in the recently viewed projects list.
// Updated is the time the project was last viewed.
type RecentlyViewedProject struct {
	Project *Project
	Updated Timestamp
}

// RecentlyViewedWiki represents an entry in the recently viewed wikis list.
// Updated is the time the Wiki page was last viewed.
type RecentlyViewedWiki struct {
	Page    *Wiki
	Updated Timestamp
}

// ──────────────────────────────────────────────────────────────
//  UserRecentlyViewedService
// ──────────────────────────────────────────────────────────────

// RecentlyViewedService handles communication with the recently-viewed methods of the Backlog API.
// All endpoints are scoped to the authenticated user (myself), so no userID argument is needed.
type RecentlyViewedService struct {
	base *recentlyviewed.Service

	// Option provides the options accepted by [RecentlyViewedService] methods.
	Option *RecentlyViewedOptionService
}

// ListIssues returns a list of issues recently viewed by the authenticated user.
//
// This method supports options returned by methods in [RecentlyViewedOptionService],
// such as:
//   - [RecentlyViewedOptionService.WithCount]
//   - [RecentlyViewedOptionService.WithOffset]
//   - [RecentlyViewedOptionService.WithOrder]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-issues
func (s *RecentlyViewedService) ListIssues(ctx context.Context, opts ...RequestOption) ([]*RecentlyViewedIssue, error) {
	v, err := s.base.ListIssues(ctx, toInnerOptions(opts)...)
	return recentlyViewedIssuesFromModel(v), convertError(err)
}

// AddIssue adds an issue to the recently viewed list of the authenticated user.
// issueIDOrKey can be an issue ID or an issue key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-issue
func (s *RecentlyViewedService) AddIssue(ctx context.Context, issueIDOrKey string) (*Issue, error) {
	v, err := s.base.AddIssue(ctx, issueIDOrKey)
	return issueFromModel(v), convertError(err)
}

// ListProjects returns a list of projects recently viewed by the authenticated user.
//
// This method supports options returned by methods in [RecentlyViewedOptionService],
// such as:
//   - [RecentlyViewedOptionService.WithCount]
//   - [RecentlyViewedOptionService.WithOffset]
//   - [RecentlyViewedOptionService.WithOrder]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-projects
func (s *RecentlyViewedService) ListProjects(ctx context.Context, opts ...RequestOption) ([]*RecentlyViewedProject, error) {
	v, err := s.base.ListProjects(ctx, toInnerOptions(opts)...)
	return recentlyViewedProjectsFromModel(v), convertError(err)
}

// ListWikis returns a list of Wiki pages recently viewed by the authenticated user.
//
// This method supports options returned by methods in [RecentlyViewedOptionService],
// such as:
//   - [RecentlyViewedOptionService.WithCount]
//   - [RecentlyViewedOptionService.WithOffset]
//   - [RecentlyViewedOptionService.WithOrder]
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-wikis
func (s *RecentlyViewedService) ListWikis(ctx context.Context, opts ...RequestOption) ([]*RecentlyViewedWiki, error) {
	v, err := s.base.ListWikis(ctx, toInnerOptions(opts)...)
	return recentlyViewedWikisFromModel(v), convertError(err)
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
	base *option.OptionService
}

// WithCount sets the number of results to return (1-100).
func (s *RecentlyViewedOptionService) WithCount(count int) RequestOption {
	return &requestOption{opt: s.base.WithCount(count)}
}

// WithOffset sets the number of items to skip.
func (s *RecentlyViewedOptionService) WithOffset(offset int) RequestOption {
	return &requestOption{opt: s.base.WithOffset(offset)}
}

// WithOrder sets the sort order of results.
func (s *RecentlyViewedOptionService) WithOrder(order Order) RequestOption {
	return &requestOption{opt: s.base.WithOrder(string(order))}
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newRecentlyViewedService(method *client.Method, option *option.OptionService) *RecentlyViewedService {
	return &RecentlyViewedService{
		base:   recentlyviewed.NewService(method),
		Option: newRecentlyViewedOptionService(option),
	}
}

func newRecentlyViewedOptionService(option *option.OptionService) *RecentlyViewedOptionService {
	return &RecentlyViewedOptionService{
		base: option,
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func recentlyViewedIssueFromModel(m *model.RecentlyViewedIssue) *RecentlyViewedIssue {
	if m == nil {
		return nil
	}
	return &RecentlyViewedIssue{
		Issue:   issueFromModel(m.Issue),
		Updated: Timestamp{m.Updated},
	}
}

func recentlyViewedIssuesFromModel(ms []*model.RecentlyViewedIssue) []*RecentlyViewedIssue {
	if ms == nil {
		return nil
	}
	result := make([]*RecentlyViewedIssue, len(ms))
	for i, v := range ms {
		result[i] = recentlyViewedIssueFromModel(v)
	}
	return result
}

func recentlyViewedProjectFromModel(m *model.RecentlyViewedProject) *RecentlyViewedProject {
	if m == nil {
		return nil
	}
	return &RecentlyViewedProject{
		Project: projectFromModel(m.Project),
		Updated: Timestamp{m.Updated},
	}
}

func recentlyViewedProjectsFromModel(ms []*model.RecentlyViewedProject) []*RecentlyViewedProject {
	if ms == nil {
		return nil
	}
	result := make([]*RecentlyViewedProject, len(ms))
	for i, v := range ms {
		result[i] = recentlyViewedProjectFromModel(v)
	}
	return result
}

func recentlyViewedWikiFromModel(m *model.RecentlyViewedWiki) *RecentlyViewedWiki {
	if m == nil {
		return nil
	}
	return &RecentlyViewedWiki{
		Page:    wikiFromModel(m.Page),
		Updated: Timestamp{m.Updated},
	}
}

func recentlyViewedWikisFromModel(ms []*model.RecentlyViewedWiki) []*RecentlyViewedWiki {
	if ms == nil {
		return nil
	}
	result := make([]*RecentlyViewedWiki, len(ms))
	for i, v := range ms {
		result[i] = recentlyViewedWikiFromModel(v)
	}
	return result
}
