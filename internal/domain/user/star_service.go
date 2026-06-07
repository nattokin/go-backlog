package user

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// StarService handles user star-related Backlog API calls.
type StarService struct {
	method *core.Method
}

// List returns a list of stars received by the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-received-star-list
func (s *StarService) List(ctx context.Context, userID int, opts ...*core.APIParamOption) ([]*model.Star, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamMinID, core.ParamMaxID, core.ParamCount, core.ParamOrder}

	var ves core.ValidationErrors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		if !errors.As(err, &ves) {
			return nil, err
		}
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("users", strconv.Itoa(userID), "stars")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Star{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Count returns the number of stars received by the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-user-received-stars
func (s *StarService) Count(ctx context.Context, userID int) (int, error) {
	var ves core.ValidationErrors
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
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v.Count, nil
}

func NewStarService(method *core.Method) *StarService {
	return &StarService{method: method}
}
