package recentlyviewed

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles communication with the recently-viewed methods of the Backlog API.
// All endpoints are scoped to the authenticated user (myself), so no userID argument is needed.
type Service struct {
	method *core.Method
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
func (s *Service) ListIssues(ctx context.Context, opts ...core.RequestOption) ([]*model.Issue, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamCount, core.ParamOffset, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "users/myself/recentlyViewedIssues", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// AddIssue adds an issue to the recently viewed list of the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-issue
func (s *Service) AddIssue(ctx context.Context, issueID int) (*model.Issue, error) {
	if issueID < 1 {
		return nil, core.NewValidationError("issueID must not be less than 1")
	}

	spath := path.Join("issues", strconv.Itoa(issueID), "recentlyViewedIssues")
	resp, err := s.method.Post(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Issue{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
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
func (s *Service) ListProjects(ctx context.Context, opts ...core.RequestOption) ([]*model.Project, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamCount, core.ParamOffset, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "users/myself/recentlyViewedProjects", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
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
func (s *Service) ListWikis(ctx context.Context, opts ...core.RequestOption) ([]*model.Wiki, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamCount, core.ParamOffset, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "users/myself/recentlyViewedWikis", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// AddWiki adds a Wiki page to the recently viewed list of the authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-wiki
func (s *Service) AddWiki(ctx context.Context, wikiID int) (*model.Wiki, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "recentlyViewedWikis")
	resp, err := s.method.Post(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Wiki{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// NewService creates and returns a new recently viewed Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
