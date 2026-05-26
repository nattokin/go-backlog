// Package pullrequest implements the Backlog Pull Request API service.
package pullrequest

import (
	"context"
	"errors"
	"iter"
	"maps"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// filterValidTypes are the options accepted by both List and All (filter params only).
var filterValidTypes = []core.APIParamOptionType{
	core.ParamStatusIDs,
	core.ParamAssigneeIDs,
	core.ParamIssueIDs,
	core.ParamCreatedUserIDs,
}

// listValidTypes are the options accepted by List (filter params + pagination).
var listValidTypes = append(filterValidTypes,
	core.ParamOffset,
	core.ParamCount,
)

// Service handles pull request-related Backlog API calls.
type Service struct {
	method *core.Method
}

// list fetches a page of pull requests using the given pre-built query.
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
func (s *Service) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...*core.APIParamOption) ([]*model.PullRequest, error) {
	query := url.Values{}
	if err := core.ApplyOptions(query, listValidTypes, opts...); err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

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

	return s.list(ctx, projectIDOrKey, repoIDOrName, query)
}

// All returns an iterator that lazily fetches all pull requests with automatic
// pagination, along with any validation error encountered at call time.
//
// perPage controls how many pull requests are fetched per API call (1-100).
// Iteration stops automatically when all pull requests have been returned.
// Passing WithCount or WithOffset in opts returns an error immediately.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) All(ctx context.Context, perPage int, projectIDOrKey string, repoIDOrName string, opts ...*core.APIParamOption) (iter.Seq2[*model.PullRequest, error], error) {
	o := &core.OptionService{}

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
	if err := core.ApplyOptions(baseQuery, filterValidTypes, opts...); err != nil {
		return nil, err
	}

	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.PullRequest, error) {
		q := maps.Clone(baseQuery)
		q.Set(core.ParamOffset.Value(), strconv.Itoa(offset))
		return s.list(ctx, projectIDOrKey, repoIDOrName, q)
	}), nil
}

// Count returns the number of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests
func (s *Service) Count(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...*core.APIParamOption) (int, error) {
	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamStatusIDs,
		core.ParamAssigneeIDs,
		core.ParamIssueIDs,
		core.ParamCreatedUserIDs,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return 0, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		return 0, ves
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return 0, ves
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
func (s *Service) Create(ctx context.Context, projectIDOrKey string, repoIDOrName string, summary string, description string, base string, branch string, opts ...*core.APIParamOption) (*model.PullRequest, error) {
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
		[]*core.APIParamOption{
			option.WithSummary(summary),
			option.WithDescription(description),
			option.WithBase(base),
			option.WithBranch(branch),
		},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

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
func (s *Service) Update(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.PullRequest, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamSummary,
		core.ParamDescription,
		core.ParamIssueID,
		core.ParamAssigneeID,
		core.ParamNotifiedUserIDs,
		core.ParamComment,
	}
	options := append([]*core.APIParamOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidatePRNumber(prNumber); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

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
