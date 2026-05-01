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

// ──────────────────────────────────────────────────────────────
//  Service
// ──────────────────────────────────────────────────────────────

// Service handles communication with the project-related methods of the Backlog API.
type Service struct {
	method *core.Method
}

// One returns a single project by its ID or key.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project
func (s *Service) One(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// All returns a list of projects.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *Service) All(ctx context.Context, opts ...core.RequestOption) ([]*model.Project, error) {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamArchived}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "projects", form)
	if err != nil {
		return nil, err
	}

	v := []*model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create creates a new project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *Service) Create(ctx context.Context, name, key string, opts ...core.RequestOption) (*model.Project, error) {
	opt := &core.OptionService{}
	nameOpt := opt.WithName(name)
	if err := nameOpt.Check(); err != nil {
		return nil, err
	}
	keyOpt := opt.WithKey(key)
	if err := keyOpt.Check(); err != nil {
		return nil, err
	}

	form := url.Values{}
	nameOpt.Set(form)
	keyOpt.Set(form)

	validTypes := []core.APIParamOptionType{
		core.ParamChartEnabled,
		core.ParamSubtaskingEnabled,
		core.ParamProjectLeaderCanEditProjectLeader,
		core.ParamTextFormattingRule,
	}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Post(ctx, "projects", form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *Service) Update(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamName,
		core.ParamKey,
		core.ParamChartEnabled,
		core.ParamSubtaskingEnabled,
		core.ParamProjectLeaderCanEditProjectLeader,
		core.ParamTextFormattingRule,
		core.ParamArchived,
	}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *Service) Delete(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// DiskUsage returns disk usage of a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-disk-usage
func (s *Service) DiskUsage(ctx context.Context, projectIDOrKey string) (*model.DiskUsageProject, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "diskUsage")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.DiskUsageProject{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  CategoryService
// ──────────────────────────────────────────────────────────────

// CategoryService handles communication with the category-related methods of the Backlog API.
type CategoryService struct {
	method *core.Method
}

// All returns a list of categories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-issue-category-list
func (s *CategoryService) All(ctx context.Context, projectIDOrKey string) ([]*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "categories")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Category{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a category to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-issue-category
func (s *CategoryService) Create(ctx context.Context, projectIDOrKey, name string) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	opt := &core.OptionService{}
	nameOpt := opt.WithName(name)
	if err := nameOpt.Check(); err != nil {
		return nil, err
	}

	form := url.Values{}
	nameOpt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "categories")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a category in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-issue-category
func (s *CategoryService) Update(ctx context.Context, projectIDOrKey string, categoryID int, name string) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if categoryID < 1 {
		return nil, core.NewValidationError("categoryId must not be less than 1")
	}

	opt := &core.OptionService{}
	nameOpt := opt.WithName(name)
	if err := nameOpt.Check(); err != nil {
		return nil, err
	}

	form := url.Values{}
	nameOpt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "categories", strconv.Itoa(categoryID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a category from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-category
func (s *CategoryService) Delete(ctx context.Context, projectIDOrKey string, categoryID int) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if categoryID < 1 {
		return nil, core.NewValidationError("categoryId must not be less than 1")
	}

	spath := path.Join("projects", projectIDOrKey, "categories", strconv.Itoa(categoryID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  StatusService
// ──────────────────────────────────────────────────────────────

// StatusService handles communication with the status-related methods of the Backlog API.
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
func (s *StatusService) Update(ctx context.Context, projectIDOrKey string, statusID int, opts ...core.RequestOption) (*model.Status, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if statusID < 1 {
		return nil, core.NewValidationError("statusId must not be less than 1")
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamColor}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
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

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{method: method}
}

func NewCategoryService(method *core.Method) *CategoryService {
	return &CategoryService{method: method}
}

func NewStatusService(method *core.Method) *StatusService {
	return &StatusService{method: method}
}
