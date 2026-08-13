// Package project implements the Backlog Project API service.
package project

import (
	"context"
	"net/url"
	"path"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

type Service struct {
	method *client.Method
}

// List returns a list of projects in the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-list
func (s *Service) List(ctx context.Context, opts ...*option.APIParamOption) ([]*model.Project, error) {
	query := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamAll, option.ParamArchived}
	if err := option.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "projects", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Project{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a single project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project
func (s *Service) One(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-project
func (s *Service) Create(ctx context.Context, key, name string, opts ...*option.APIParamOption) (*model.Project, error) {
	optSvc := &option.OptionService{}

	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamKey, option.ParamName, option.ParamChartEnabled, option.ParamSubtaskingEnabled, option.ParamProjectLeaderCanEditProjectLeader, option.ParamTextFormattingRule}
	options := append([]*option.APIParamOption{optSvc.WithKey(key), optSvc.WithName(name)}, opts...)
	if err := option.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	resp, err := s.method.Post(ctx, "projects", form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-project
func (s *Service) Update(ctx context.Context, projectIDOrKey string, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Project, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamKey, option.ParamName, option.ParamChartEnabled, option.ParamSubtaskingEnabled,
		option.ParamProjectLeaderCanEditProjectLeader, option.ParamTextFormattingRule, option.ParamArchived,
	}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-project
func (s *Service) Delete(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// DiskUsage returns disk usage of a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-disk-usage
func (s *Service) DiskUsage(ctx context.Context, projectIDOrKey string) (*model.DiskUsageProject, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "diskUsage")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.DiskUsageProject{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Icon returns the icon image of a project.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-icon
func (s *Service) Icon(ctx context.Context, projectIDOrKey string) (*model.FileData, error) {
	var ves validation.Errors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "image")
	resp, err := s.method.Download(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	return client.DownloadResponse(resp)
}

func NewService(method *client.Method) *Service {
	return &Service{
		method: method,
	}
}
