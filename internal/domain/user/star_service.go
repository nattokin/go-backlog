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

// StarService handles user star-related Backlog API calls.
type StarService struct {
	method *core.Method
}

// List returns a list of stars received by the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-received-star-list
func (s *StarService) List(ctx context.Context, userID int, opts ...core.RequestOption) ([]*model.Star, error) {
	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamMinID, core.ParamMaxID, core.ParamCount, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
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
	if err := validate.ValidateUserID(userID); err != nil {
		return 0, err
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
