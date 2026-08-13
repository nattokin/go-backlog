// Package wiki implements the Backlog Wiki API service.
package wiki

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
)

type Service struct {
	method *client.Method
}

// List returns a list of wiki pages in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list
func (s *Service) List(ctx context.Context, projectIDOrKey string, opts ...*option.APIParamOption) ([]*model.Wiki, error) {
	query := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamKeyword}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validTypes, opts...)); err != nil {
		return nil, err
	}

	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Wiki{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of wiki pages in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-wiki-page
func (s *Service) Count(ctx context.Context, projectIDOrKey string) (int, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return 0, ves
	}

	query := url.Values{}
	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis/count", query)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns a single wiki page by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page
func (s *Service) One(ctx context.Context, wikiID int) (*model.Wiki, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Create creates a new wiki page in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/create-wiki-page
func (s *Service) Create(ctx context.Context, projectID int, name, content string, opts ...*option.APIParamOption) (*model.Wiki, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamName, option.ParamContent, option.ParamMailNotify}
	options := append([]*option.APIParamOption{optSvc.WithName(name), optSvc.WithContent(content)}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectID(projectID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	form.Set("projectId", strconv.Itoa(projectID))

	resp, err := s.method.Post(ctx, "wikis", form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an existing wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-wiki-page
func (s *Service) Update(ctx context.Context, wikiID int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.Wiki, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamName, option.ParamContent, option.ParamMailNotify}
	options := append([]*option.APIParamOption{opt}, opts...)

	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.ApplyOptions(form, validTypes, options...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	// Only check name/content presence when there are no opt errors.
	if len(ves) == 0 && !form.Has("name") && !form.Has("content") {
		ves = append(ves, core.NewValidationError("", "requires an opt to modify wiki content or name (WithName or WithContent)"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a wiki page by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page
func (s *Service) Delete(ctx context.Context, wikiID int, opts ...*option.APIParamOption) (*model.Wiki, error) {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamMailNotify}

	var ves core.ValidationErrors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewService(method *client.Method) *Service {
	return &Service{
		method: method,
	}
}
