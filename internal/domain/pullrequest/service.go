// Package pullrequest implements the Backlog Pull Request API service.
package pullrequest

import (
	"context"
	"iter"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles pull request-related Backlog API calls.
type Service struct {
	method *core.Method
}

// List returns a list of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) ([]*model.PullRequest, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}

	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamStatusIDs,
		core.ParamAssigneeIDs,
		core.ParamIssueIDs,
		core.ParamCreatedUserIDs,
		core.ParamOffset,
		core.ParamCount,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.PullRequest{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// All returns an iterator that lazily fetches all pull requests with automatic pagination.
//
// perPage controls how many pull requests are fetched per API call (1-100).
// Iteration stops automatically when all pull requests have been returned.
// The caller must not pass WithCount or WithOffset in opts; those are managed internally.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) All(ctx context.Context, perPage int, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) iter.Seq2[*model.PullRequest, error] {
	o := &core.OptionService{}
	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.PullRequest, error) {
		return s.List(ctx, projectIDOrKey, repoIDOrName, append(opts,
			o.WithCount(perPage),
			o.WithOffset(offset),
		)...)
	})
}

// Count returns the number of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests
func (s *Service) Count(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) (int, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return 0, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return 0, err
	}

	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamStatusIDs,
		core.ParamAssigneeIDs,
		core.ParamIssueIDs,
		core.ParamCreatedUserIDs,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return 0, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", "count")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns a single pull request by its number.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request
func (s *Service) One(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int) (*model.PullRequest, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.PullRequest{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request
func (s *Service) Create(ctx context.Context, projectIDOrKey string, repoIDOrName string, summary string, description string, base string, branch string, opts ...core.RequestOption) (*model.PullRequest, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamSummary,
		core.ParamDescription,
		core.ParamBase,
		core.ParamBranch,
		core.ParamIssueID,
		core.ParamAssigneeID,
		core.ParamNotifiedUserIDs,
		core.ParamAttachmentIDs,
	}
	options := append(
		[]core.RequestOption{
			option.WithSummary(summary),
			option.WithDescription(description),
			option.WithBase(base),
			option.WithBranch(branch),
		},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.PullRequest{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an existing pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-pull-request
func (s *Service) Update(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, option core.RequestOption, opts ...core.RequestOption) (*model.PullRequest, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repoIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamSummary,
		core.ParamDescription,
		core.ParamIssueID,
		core.ParamAssigneeID,
		core.ParamNotifiedUserIDs,
		core.ParamComment,
	}
	options := append([]core.RequestOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.PullRequest{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
