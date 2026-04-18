package issue

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
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

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
