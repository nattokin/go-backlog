package pullrequest

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

type Service struct {
	method *core.Method
}

// All returns a list of pull requests.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list
func (s *Service) All(ctx context.Context, projectIDOrKey string, repoIDOrName string, opts ...core.RequestOption) ([]*model.PullRequest, error) {
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

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
