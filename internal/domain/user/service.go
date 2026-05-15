// Package user implements the Backlog User API service.
package user

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// getUser is a shared helper that fetches a single user from the given spath.
func getUser(ctx context.Context, m *core.Method, spath string) (*model.User, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

type Service struct {
	method *core.Method
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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// One returns a single user by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user
func (s *Service) One(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return getUser(ctx, s.method, spath)
}

// Own returns the currently authenticated user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-own-user
func (s *Service) Own(ctx context.Context) (*model.User, error) {
	return getUser(ctx, s.method, "users/myself")
}

// Add adds a new user to the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-user
func (s *Service) Add(ctx context.Context, userID, password, name, mailAddress string, roleType int) (*model.User, error) {
	if userID == "" {
		return nil, core.NewValidationError("userID must not be empty")
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamPassword, core.ParamName, core.ParamMailAddress, core.ParamRoleType}
	options := []core.RequestOption{
		option.WithPassword(password),
		option.WithName(name),
		option.WithMailAddress(mailAddress),
		option.WithRoleType(roleType),
	}
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	form.Set("userId", userID)

	resp, err := s.method.Post(ctx, "users", form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates an existing user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-user
func (s *Service) Update(ctx context.Context, id int, option core.RequestOption, opts ...core.RequestOption) (*model.User, error) {
	baseOpt := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamUserID, core.ParamName, core.ParamPassword, core.ParamMailAddress, core.ParamRoleType}
	options := append([]core.RequestOption{baseOpt.WithUserID(id), option}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))

	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a user by ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-user
func (s *Service) Delete(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Icon downloads the icon image of the user.
// The caller is responsible for closing FileData.Body after use.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-icon
func (s *Service) Icon(ctx context.Context, id int) (*model.FileData, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id), "icon")
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
