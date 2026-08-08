// Package issue implements the Backlog Issue API service.
package issue

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

var countValidTypes = []core.APIParamOptionType{
	core.ParamProjectIDs,
	core.ParamIssueTypeIDs,
	core.ParamCategoryIDs,
	core.ParamVersionIDs,
	core.ParamMilestoneIDs,
	core.ParamStatusIDs,
	core.ParamPriorityIDs,
	core.ParamAssigneeIDs,
	core.ParamCreatedUserIDs,
	core.ParamResolutionIDs,
	core.ParamParentChild,
	core.ParamAttachment,
	core.ParamSharedFile,
	core.ParamCreatedSince,
	core.ParamCreatedUntil,
	core.ParamUpdatedSince,
	core.ParamUpdatedUntil,
	core.ParamStartDateSince,
	core.ParamStartDateUntil,
	core.ParamDueDateSince,
	core.ParamDueDateUntil,
	core.ParamHasDueDate,
	core.ParamIDs,
	core.ParamParentIssueIDs,
	core.ParamKeyword,
}

var filterValidTypes = append(countValidTypes,
	core.ParamSort,
	core.ParamOrder,
)

var listValidTypes = append(filterValidTypes,
	core.ParamOffset,
	core.ParamCount,
)

var createValidTypes = []core.APIParamOptionType{
	core.ParamSummary,
	core.ParamIssueTypeID,
	core.ParamPriorityID,
	core.ParamDescription,
	core.ParamStartDate,
	core.ParamDueDate,
	core.ParamEstimatedHours,
	core.ParamActualHours,
	core.ParamCategoryIDs,
	core.ParamVersionIDs,
	core.ParamMilestoneIDs,
	core.ParamAssigneeID,
	core.ParamParentIssueID,
	core.ParamNotifiedUserIDs,
	core.ParamAttachmentIDs,
	core.ParamCustomField,
}

var updateValidTypes = append(createValidTypes,
	core.ParamStatusID,
	core.ParamResolutionID,
	core.ParamComment,
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
func (s *Service) List(ctx context.Context, opts ...*core.APIParamOption) ([]*model.Issue, error) {
	query := url.Values{}
	if err := core.ApplyOptions(query, listValidTypes, opts...); err != nil {
		return nil, err
	}
	return s.list(ctx, query)
}

// All returns an iterator that lazily fetches all issues with automatic pagination.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-list
func (s *Service) All(ctx context.Context, perPage int, opts ...*core.APIParamOption) (iter.Seq2[*model.Issue, error], error) {
	o := &core.OptionService{}

	countOpt := o.WithCount(perPage)
	if ve := countOpt.Check(); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	baseQuery := url.Values{}
	countOpt.Set(baseQuery)
	if err := core.ApplyOptions(baseQuery, filterValidTypes, opts...); err != nil {
		return nil, err
	}

	return core.AllSeq(ctx, perPage, func(ctx context.Context, offset int) ([]*model.Issue, error) {
		q := maps.Clone(baseQuery)
		q.Set(core.ParamOffset.Value(), strconv.Itoa(offset))
		return s.list(ctx, q)
	}), nil
}

// Count returns the total count of issues matching the given filters.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-issue
func (s *Service) Count(ctx context.Context, opts ...*core.APIParamOption) (int, error) {
	query := url.Values{}
	if err := core.ApplyOptions(query, countValidTypes, opts...); err != nil {
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
func (s *Service) Create(ctx context.Context, projectID int, summary string, issueTypeID int, priorityID int, opts ...*core.APIParamOption) (*model.Issue, error) {
	o := &core.OptionService{}
	form := url.Values{}
	options := append(
		[]*core.APIParamOption{
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
	if err := core.ApplyOptions(form, createValidTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if len(ves) > 0 {
		return nil, ves
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
func (s *Service) Update(ctx context.Context, issueIDOrKey string, option *core.APIParamOption, opts ...*core.APIParamOption) (*model.Issue, error) {
	form := url.Values{}
	options := append([]*core.APIParamOption{option}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.ApplyOptions(form, updateValidTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if len(ves) > 0 {
		return nil, ves
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
