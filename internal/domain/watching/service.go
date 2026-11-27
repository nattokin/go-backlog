// Package watching implements the Backlog Watching API service.
package watching

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// Service handles watching-related Backlog API calls.
type Service struct {
	method *client.Method
}

// List returns a list of watchings of the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-watching-list
func (s *Service) List(ctx context.Context, userID int, opts ...*option.APIParamOption) ([]*model.Watching, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{
		option.ParamOrder,
		option.ParamSort,
		option.ParamCount,
		option.ParamOffset,
		option.ParamResourceAlreadyRead,
		option.ParamIssueIDs,
	}

	var ves validation.Errors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validOptionKeys, opts...)); err != nil {
		return nil, err
	}

	spath := client.PathJoin("users", strconv.Itoa(userID), "watchings")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Watching{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of watchings of the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-watching
func (s *Service) Count(ctx context.Context, userID int, opts ...*option.APIParamOption) (int, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{option.ParamResourceAlreadyRead, option.ParamAlreadyRead}

	var ves validation.Errors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validOptionKeys, opts...)); err != nil {
		return 0, err
	}

	spath := client.PathJoin("users", strconv.Itoa(userID), "watchings", "count")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return 0, err
	}

	var v struct {
		Count int `json:"count"`
	}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v.Count, nil
}

// One returns a single watching.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-watching
func (s *Service) One(ctx context.Context, watchingID int) (*model.Watching, error) {
	var ves validation.Errors
	if ve := validate.ValidateWatchingID(watchingID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := client.PathJoin("watchings", strconv.Itoa(watchingID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Watching{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Create adds a watching for an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-watching
func (s *Service) Create(ctx context.Context, issueIDOrKey string, opts ...*option.APIParamOption) (*model.Watching, error) {
	form := url.Values{}
	validOptionKeys := []option.APIParamOptionType{option.ParamNote}

	var ves validation.Errors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validOptionKeys, opts...)); err != nil {
		return nil, err
	}

	form.Set("issueIdOrKey", issueIDOrKey)
	resp, err := s.method.Post(ctx, "watchings", form)
	if err != nil {
		return nil, err
	}

	v := &model.Watching{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Update updates the note of a watching.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-watching
func (s *Service) Update(ctx context.Context, watchingID int, note string) (*model.Watching, error) {
	var ves validation.Errors
	if ve := validate.ValidateWatchingID(watchingID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	form := url.Values{}
	form.Set("note", note)
	spath := client.PathJoin("watchings", strconv.Itoa(watchingID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := &model.Watching{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// Delete deletes a watching.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-watching
func (s *Service) Delete(ctx context.Context, watchingID int) (*model.Watching, error) {
	var ves validation.Errors
	if ve := validate.ValidateWatchingID(watchingID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := client.PathJoin("watchings", strconv.Itoa(watchingID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Watching{}
	if err := client.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// MarkAsRead marks a watching as read.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/mark-watching-as-read
func (s *Service) MarkAsRead(ctx context.Context, watchingID int) error {
	var ves validation.Errors
	if ve := validate.ValidateWatchingID(watchingID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return ves
	}

	spath := client.PathJoin("watchings", strconv.Itoa(watchingID), "markAsRead")
	if _, err := s.method.Post(ctx, spath, nil); err != nil {
		return err
	}

	return nil
}

// NewService returns a *Service that sends requests through method.
func NewService(method *client.Method) *Service {
	return &Service{method: method}
}
