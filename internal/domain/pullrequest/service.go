// Package pullrequest implements the Backlog Pull Request API service.
package pullrequest

import (
	"context"
	"iter"
	"maps"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
)

var filterValidTypes = []option.APIParamOptionType{
	option.ParamStatusIDs,
	option.ParamAssigneeIDs,
	option.ParamIssueIDs,
	option.ParamCreatedUserIDs,
}

var listValidTypes = append(filterValidTypes,
	option.ParamOffset,
	option.ParamCount,
)

// Service handles pull request-related Backlog API calls.
type Service struct {
	method *core.Method
}

func (s *Service) list(ctx context.Context, projectIDOrKey string, repoIDOrName string, query url.Values) ([]*model.PullRequest, error) {
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

// List returns a list of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...*option.APIParamOption) ([]*model.PullRequest, error) {
	query := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, listValidTypes, opts...)); err != nil {
		return nil, err
	}

	return s.list(ctx, projectIDOrKey, repoIDOrName, query)
}

// All returns an iterator that lazily fetches all pull requests with automatic pagination.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) All(ctx context.Context, perPage int, projectIDOrKey string, repoIDOrName string, opts ...*option.APIParamOption) (iter.Seq2[*model.PullRequest, error], error) {
	o := &option.OptionService{}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	countOpt := o.WithCount(perPage)
	if ve := countOpt.Check(); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	baseQuery := url.Values{}
	countOpt.Set(baseQuery)
	if err := option.ApplyOptions(baseQuery, filterValidTypes, opts...); err != nil {
		return nil, err
	}

	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.PullRequest, error) {
		q := maps.Clone(baseQuery)
		q.Set(option.ParamOffset.Value(), strconv.Itoa(offset))
		return s.list(ctx, projectIDOrKey, repoIDOrName, q)
	}), nil
}

// Count returns the number of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests
func (s *Service) Count(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...*option.APIParamOption) (int, error) {
	query := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamStatusIDs,
		option.ParamAssigneeIDs,
		option.ParamIssueIDs,
		option.ParamCreatedUserIDs,
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validTypes, opts...)); err != nil {
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
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
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
func (s *Service) Create(ctx context.Context, projectIDOrKey string, repoIDOrName string, summary string, description string, base string, branch string, opts ...*option.APIParamOption) (*model.PullRequest, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamSummary,
		option.ParamDescription,
		option.ParamBase,
		option.ParamBranch,
		option.ParamIssueID,
		option.ParamAssigneeID,
		option.ParamNotifiedUserIDs,
		option.ParamAttachmentIDs,
	}
	options := append(
		[]*option.APIParamOption{
			optSvc.WithSummary(summary),
			optSvc.WithDescription(description),
			optSvc.WithBase(base),
			optSvc.WithBranch(branch),
		},
		opts...,
	)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
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
func (s *Service) Update(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.PullRequest, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamSummary,
		option.ParamDescription,
		option.ParamIssueID,
		option.ParamAssigneeID,
		option.ParamNotifiedUserIDs,
		option.ParamComment,
	}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
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
