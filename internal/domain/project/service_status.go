package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// StatusService handles status-related Backlog API calls for a project.
type StatusService struct {
	method *core.Method
}

// All returns a list of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project
func (s *StatusService) All(ctx context.Context, projectIDOrKey string) ([]*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "statuses")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new status to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-status
func (s *StatusService) Create(ctx context.Context, projectIDOrKey, name, color string) (*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	opt := &core.OptionService{}
	nameOpt := opt.WithName(name)
	if err := nameOpt.Check(); err != nil {
		return nil, err
	}
	colorOpt := opt.WithColor(color)
	if err := colorOpt.Check(); err != nil {
		return nil, err
	}

	form := url.Values{}
	nameOpt.Set(form)
	colorOpt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "statuses")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a status in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-status
func (s *StatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, option core.RequestOption, opts ...core.RequestOption) (*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if statusID < 1 {
		return nil, core.NewValidationError("statusId must not be less than 1")
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor}
	options := append([]core.RequestOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "statuses", strconv.Itoa(statusID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a status from a project.
// substituteStatusID specifies the status to migrate existing issues to.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-status
func (s *StatusService) Delete(ctx context.Context, projectIDOrKey string, statusID, substituteStatusID int) (*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if statusID < 1 {
		return nil, core.NewValidationError("statusId must not be less than 1")
	}
	if substituteStatusID < 1 {
		return nil, core.NewValidationError("substituteStatusId must not be less than 1")
	}

	form := url.Values{}
	form.Set("substituteStatusId", strconv.Itoa(substituteStatusID))

	spath := path.Join("projects", projectIDOrKey, "statuses", strconv.Itoa(statusID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// UpdateOrder updates the display order of statuses in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-order-of-status
func (s *StatusService) UpdateOrder(ctx context.Context, projectIDOrKey string, statusIDs []int) ([]*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if len(statusIDs) == 0 {
		return nil, core.NewValidationError("statusIDs must not be empty")
	}
	for _, id := range statusIDs {
		if id < 1 {
			return nil, core.NewValidationError("each statusId must not be less than 1")
		}
	}

	form := url.Values{}
	for _, id := range statusIDs {
		form.Add("statusId[]", strconv.Itoa(id))
	}

	spath := path.Join("projects", projectIDOrKey, "statuses", "updateDisplayOrder")
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.Status{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewStatusService(method *core.Method) *StatusService {
	return &StatusService{method: method}
}
