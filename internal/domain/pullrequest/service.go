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

// validateListArgs validates the path arguments shared by List and All.
func (s *Service) validateListArgs(projectIDOrKey string, repoIDOrName string) error {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return err
	}
	return validate.ValidateRepositoryIDOrName(repoIDOrName)
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
func (s *Service) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) ([]*model.PullRequest, error) {
	if err := s.validateListArgs(projectIDOrKey, repoIDOrName); err != nil {
		return nil, err
	}

	query := url.Values{}
	if err := core.ApplyOptions(query, listValidTypes, opts...); err != nil {
		return nil, err
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
func (s *Service) All(ctx context.Context, perPage int, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) (iter.Seq2[*model.PullRequest, error], error) {
	o := &core.OptionService{}
	countOpt := o.WithCount(perPage)
	if err := countOpt.Check(); err != nil {
		return nil, err
	}

	if err := s.validateListArgs(projectIDOrKey, repoIDOrName); err != nil {
		return nil, err
	}

	baseQuery := url.Values{}
	if err := core.ApplyOptions(baseQuery, filterValidTypes, opts...); err != nil {
		return nil, err
	}
	if err := countOpt.Set(baseQuery); err != nil {
		return nil, err
	}

	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.PullRequest, error) {
		q := cloneQuery(baseQuery)
		q.Set(core.ParamOffset.Value(), strconv.Itoa(offset))
		return s.list(ctx, projectIDOrKey, repoIDOrName, q)
	}), nil
}

// cloneQuery returns a shallow copy of url.Values.
func cloneQuery(src url.Values) url.Values {
	dst := make(url.Values, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
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
