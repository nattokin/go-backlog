package project

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
)

// CategoryService handles category-related Backlog API calls for a project.
type CategoryService struct {
	method *client.Method
}

// List returns a list of categories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-category-list
func (s *CategoryService) List(ctx context.Context, projectIDOrKey string) ([]*model.Category, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "categories")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Category{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a new category to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-category
func (s *CategoryService) Create(ctx context.Context, projectIDOrKey string, name string) (*model.Category, error) {
	opt := (&option.OptionService{}).WithName(name)
	if ve := opt.Check(); ve != nil {
		var ves core.ValidationErrors
		ves = append(ves, ve)
		if ve2 := validate.ValidateProjectIDOrKey(projectIDOrKey); ve2 != nil {
			ves = append(ves, ve2)
		}
		return nil, ves
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	opt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "categories")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a category in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-category
func (s *CategoryService) Update(ctx context.Context, projectIDOrKey string, categoryID int, name string) (*model.Category, error) {
	opt := (&option.OptionService{}).WithName(name)
	var ves core.ValidationErrors
	if ve := opt.Check(); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if categoryID < 1 {
		ves = append(ves, core.NewValidationError("categoryId", "categoryId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	opt.Set(form)

	spath := path.Join("projects", projectIDOrKey, "categories", strconv.Itoa(categoryID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a category from a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-category
func (s *CategoryService) Delete(ctx context.Context, projectIDOrKey string, categoryID int) (*model.Category, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if categoryID < 1 {
		ves = append(ves, core.NewValidationError("categoryId", "categoryId must not be less than 1"))
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "categories", strconv.Itoa(categoryID))
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewCategoryService(method *client.Method) *CategoryService {
	return &CategoryService{method: method}
}
