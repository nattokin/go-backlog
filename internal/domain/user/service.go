// Package user implements the Backlog User API service.
package user

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

// getUser is a shared helper that fetches a single user from the given spath.
func getUser(ctx context.Context, m *client.Method, spath string) (*model.User, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

type Service struct {
	method *client.Method
}

// List returns a list of all users in the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-list
func (s *Service) List(ctx context.Context) ([]*model.User, error) {
	resp, err := s.method.Get(ctx, "users", nil)
	if err != nil {
		return nil, err
	}

	v := []*model.User{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a single user by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *Service) One(ctx context.Context, id int) (*model.User, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateUserID(id); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("users", strconv.Itoa(id))
	return getUser(ctx, s.method, spath)
}

// Me returns the currently authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *Service) Me(ctx context.Context) (*model.User, error) {
	return getUser(ctx, s.method, "users/myself")
}

// Add adds a new user to the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *Service) Add(ctx context.Context, userID, password, name, mailAddress string, roleType int) (*model.User, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamPassword, option.ParamName, option.ParamMailAddress, option.ParamRoleType}
	options := []*option.APIParamOption{
		optSvc.WithPassword(password),
		optSvc.WithName(name),
		optSvc.WithMailAddress(mailAddress),
		optSvc.WithRoleType(roleType),
	}
	var ves core.ValidationErrors
	if userID == "" {
		ves = append(ves, core.NewValidationError("userID", "userID must not be empty"))
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, options...)); err != nil {
		return nil, err
	}

	form.Set("userId", userID)

	resp, err := s.method.Post(ctx, "users", form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an existing user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *Service) Update(ctx context.Context, id int, opt *option.APIParamOption, opts ...*option.APIParamOption) (*model.User, error) {
	baseOpt := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{option.ParamUserID, option.ParamName, option.ParamPassword, option.ParamMailAddress, option.ParamRoleType}
	options := append([]*option.APIParamOption{baseOpt.WithUserID(id), opt}, opts...)
	if err := option.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))

	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a user by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *Service) Delete(ctx context.Context, id int) (*model.User, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateUserID(id); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("users", strconv.Itoa(id))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Icon downloads the icon image of the user.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon
func (s *Service) Icon(ctx context.Context, id int) (*model.FileData, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateUserID(id); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("users", strconv.Itoa(id), "icon")
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
