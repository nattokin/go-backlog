package user

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/activity"
	"github.com/nattokin/go-backlog/internal/validate"
)

// ActivityService handles user activity-related Backlog API calls.
// It delegates HTTP operations to the shared activity.Service.
type ActivityService struct {
	base   *activity.Service
	method *core.Method
}

// List returns a list of activities for the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *ActivityService) List(ctx context.Context, userID int, opts ...*core.APIParamOption) ([]*model.Activity, error) {
	query := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := s.base.ApplyOptions(query, opts...); err != nil {
		if !errors.As(err, &ves) {
			return nil, err
		}
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return s.base.Fetch(ctx, spath, query)
}

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
