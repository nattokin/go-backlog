// Package issue implements the Backlog Issue API service.
package issue

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

var countValidTypes = []option.APIParamOptionType{
	option.ParamProjectIDs,
	option.ParamIssueTypeIDs,
	option.ParamCategoryIDs,
	option.ParamVersionIDs,
	option.ParamMilestoneIDs,
	option.ParamStatusIDs,
	option.ParamPriorityIDs,
	option.ParamAssigneeIDs,
	option.ParamCreatedUserIDs,
	option.ParamResolutionIDs,
	option.ParamParentChild,
	option.ParamAttachment,
	option.ParamSharedFile,
	option.ParamCreatedSince,
	option.ParamCreatedUntil,
	option.ParamUpdatedSince,
	option.ParamUpdatedUntil,
	option.ParamStartDateSince,
	option.ParamStartDateUntil,
	option.ParamDueDateSince,
	option.ParamDueDateUntil,
	option.ParamHasDueDate,
	option.ParamIDs,
	option.ParamParentIssueIDs,
	option.ParamKeyword,
}

var filterValidTypes = append(countValidTypes,
	option.ParamSort,
	option.ParamOrder,
)

var listValidTypes = append(filterValidTypes,
	option.ParamOffset,
	option.ParamCount,
)

var createValidTypes = []option.APIParamOptionType{
	option.ParamSummary,
	option.ParamIssueTypeID,
	option.ParamPriorityID,
	option.ParamDescription,
	option.ParamStartDate,
	option.ParamDueDate,
	option.ParamEstimatedHours,
	option.ParamActualHours,
	option.ParamCategoryIDs,
	option.ParamVersionIDs,
	option.ParamMilestoneIDs,
	option.ParamAssigneeID,
	option.ParamParentIssueID,
	option.ParamNotifiedUserIDs,
	option.ParamAttachmentIDs,
	option.ParamCustomField,
}

var updateValidTypes = append(createValidTypes,
	option.ParamStatusID,
	option.ParamResolutionID,
	option.ParamComment,
)

type Service struct {
	method *core.Method
}

func (s *Service) list(ctx context.Context, query url.Values) ([]*model.Issue, error) {
	resp, err := s.method.Get(ctx, "issues", query)
	if err != nil {
		return nil, err
	}
	v := []*model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}
	return v, nil
}

// List returns a list of issues.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-list
func (s *Service) List(ctx context.Context, opts ...*option.APIParamOption) ([]*model.Issue, error) {
	query := url.Values{}
	if err := option.ApplyOptions(query, listValidTypes, opts...); err != nil {
		return nil, err
	}
	return s.list(ctx, query)
}

// All returns an iterator that lazily fetches all issues with automatic pagination.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-list
func (s *Service) All(ctx context.Context, perPage int, opts ...*option.APIParamOption) (iter.Seq2[*model.Issue, error], error) {
	o := &option.OptionService{}

	countOpt := o.WithCount(perPage)
	if ve := countOpt.Check(); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	baseQuery := url.Values{}
	countOpt.Set(baseQuery)
	if err := option.ApplyOptions(baseQuery, filterValidTypes, opts...); err != nil {
		return nil, err
	}

	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.Issue, error) {
		q := maps.Clone(baseQuery)
		q.Set(option.ParamOffset.Value(), strconv.Itoa(offset))
		return s.list(ctx, q)
	}), nil
}

// Count returns the total count of issues matching the given filters.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-issue
func (s *Service) Count(ctx context.Context, opts ...*option.APIParamOption) (int, error) {
	query := url.Values{}
	if err := option.ApplyOptions(query, countValidTypes, opts...); err != nil {
		return 0, err
	}

	resp, err := s.method.Get(ctx, "issues/count", query)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns a single issue by its ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue
func (s *Service) One(ctx context.Context, issueIDOrKey string) (*model.Issue, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey)
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue
func (s *Service) Create(ctx context.Context, projectID int, summary string, issueTypeID int, priorityID int, opts ...*option.APIParamOption) (*model.Issue, error) {
	o := &option.OptionService{}
	form := url.Values{}
	options := append(
		[]*option.APIParamOption{
			o.WithSummary(summary),
			o.WithIssueTypeID(issueTypeID),
			o.WithPriorityID(priorityID),
		},
		opts...,
	)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectID(projectID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, createValidTypes, options...)); err != nil {
		return nil, err
	}

	form.Set("projectId", strconv.Itoa(projectID))

	resp, err := s.method.Post(ctx, "issues", form)
	if err != nil {
		return nil, err
	}

	v := model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an existing issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue
func (s *Service) Update(ctx context.Context, issueIDOrKey string, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Issue, error) {
	form := url.Values{}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, updateValidTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey)
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes an issue by its ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue
func (s *Service) Delete(ctx context.Context, issueIDOrKey string) (*model.Issue, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey)
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Issue{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Participants returns a list of participants on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-participant-list
func (s *Service) Participants(ctx context.Context, issueIDOrKey string) ([]*model.User, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "participants")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
