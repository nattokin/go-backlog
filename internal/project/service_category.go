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

// CategoryService handles category-related Backlog API calls for a project.
type CategoryService struct {
	method *core.Method
}

// All returns a list of categories in a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-category-list
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

// Create adds a new category to a project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-category
func (s *CategoryService) Create(ctx context.Context, projectIDOrKey string, name string) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	option := (&core.OptionService{}).WithName(name)
	if err := option.Check(); err != nil {
		return nil, err
	}
	form := url.Values{}
	option.Set(form)

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
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-category
func (s *CategoryService) Update(ctx context.Context, projectIDOrKey string, categoryID int, name string) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if categoryID < 1 {
		return nil, core.NewValidationError("categoryId must not be less than 1")
	}

	option := (&core.OptionService{}).WithName(name)
	if err := option.Check(); err != nil {
		return nil, err
	}
	form := url.Values{}
	option.Set(form)

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
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-category
func (s *CategoryService) Delete(ctx context.Context, projectIDOrKey string, categoryID int) (*model.Category, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if categoryID < 1 {
		return nil, core.NewValidationError("categoryId must not be less than 1")
	}

	spath := path.Join("projects", projectIDOrKey, "categories", strconv.Itoa(categoryID))
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.Category{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewCategoryService(method *core.Method) *CategoryService {
	return &CategoryService{method: method}
}
