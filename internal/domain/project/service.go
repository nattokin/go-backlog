// Package project implements the Backlog Project API service.
package project

import (
	"context"
	"net/url"
	"path"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles project-related Backlog API calls.
type Service struct {
	method *core.Method
}

// List returns a list of projects in the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *Service) List(ctx context.Context, opts ...core.RequestOption) ([]*model.Project, error) {

	query := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamAll, core.ParamArchived}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "projects", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a single project.
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

// Create creates a new project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *Service) Create(ctx context.Context, key, name string, opts ...core.RequestOption) (*model.Project, error) {
	option := &core.OptionService{}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamKey, core.ParamName, core.ParamChartEnabled, core.ParamSubtaskingEnabled, core.ParamProjectLeaderCanEditProjectLeader, core.ParamTextFormattingRule}
	options := append([]core.RequestOption{option.WithKey(key), option.WithName(name)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
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
func (s *Service) Update(ctx context.Context, projectIDOrKey string, option core.RequestOption, opts ...core.RequestOption) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamKey, core.ParamName, core.ParamChartEnabled, core.ParamSubtaskingEnabled,
		core.ParamProjectLeaderCanEditProjectLeader, core.ParamTextFormattingRule, core.ParamArchived,
	}
	options := append([]core.RequestOption{option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
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
	resp, err := s.method.Delete(ctx, spath, url.Values{})
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

// Icon returns the icon image of a project.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-icon
func (s *Service) Icon(ctx context.Context, projectIDOrKey string) (*model.FileData, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "image")
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return core.DownloadResponse(resp)
}

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
