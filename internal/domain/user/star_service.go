package user

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// StarService handles user star-related Backlog API calls.
type StarService struct {
	method *client.Method
}

// List returns a list of stars received by the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-received-star-list
func (s *StarService) List(ctx context.Context, userID int, opts ...*option.APIParamOption) ([]*model.Star, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{option.ParamMinID, option.ParamMaxID, option.ParamCount, option.ParamOrder}

	var ves validation.Errors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validOptionKeys, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "stars")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Star{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of stars received by the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-user-received-stars
func (s *StarService) Count(ctx context.Context, userID int) (int, error) {
	var ves validation.Errors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return 0, ves
	}

	spath := path.Join("users", strconv.Itoa(userID), "stars", "count")
	resp, err := s.method.Get(ctx, spath, nil)
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

func NewStarService(method *client.Method) *StarService {
	return &StarService{method: method}
}
