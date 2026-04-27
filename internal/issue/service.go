package issue

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

// All returns a list of issues.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-list
func (s *Service) All(ctx context.Context, opts ...core.RequestOption) ([]*model.Issue, error) {
	query := url.Values{}
	validTypes := []core.APIParamOptionType{
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
		core.ParamSort,
		core.ParamOrder,
		core.ParamOffset,
		core.ParamCount,
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
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

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

// Count returns the total count of issues matching the given filters.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-issue
func (s *Service) Count(ctx context.Context, opts ...core.RequestOption) (int, error) {
	query := url.Values{}
	validTypes := []core.APIParamOptionType{
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
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
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
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
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
func (s *Service) Create(ctx context.Context, projectID int, summary string, issueTypeID int, priorityID int, opts ...core.RequestOption) (*model.Issue, error) {
	if err := validate.ValidateProjectID(projectID); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
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
		core.ParamStatusID,
		core.ParamNotifiedUserIDs,
		core.ParamAttachmentIDs,
	}
	options := append(
		[]core.RequestOption{
			option.WithSummary(summary),
			option.WithIssueTypeID(issueTypeID),
			option.WithPriorityID(priorityID),
		},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
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
func (s *Service) Update(ctx context.Context, issueIDOrKey string, option core.RequestOption, opts ...core.RequestOption) (*model.Issue, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamSummary,
		core.ParamDescription,
		core.ParamIssueTypeID,
		core.ParamCategoryIDs,
		core.ParamVersionIDs,
		core.ParamMilestoneIDs,
		core.ParamStartDate,
		core.ParamDueDate,
		core.ParamEstimatedHours,
		core.ParamActualHours,
		core.ParamAssigneeID,
		core.ParamParentIssueID,
		core.ParamPriorityID,
		core.ParamStatusID,
		core.ParamResolutionID,
		core.ParamNotifiedUserIDs,
		core.ParamAttachmentIDs,
	}
	options := append([]core.RequestOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
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
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
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
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
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

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
