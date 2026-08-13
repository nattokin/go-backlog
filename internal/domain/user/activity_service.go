package user

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/shared/activity"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// ActivityService handles user activity-related Backlog API calls.
type ActivityService struct {
	base   *activity.Service
	method *client.Method
}

// List returns a list of activities for the user.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *ActivityService) List(ctx context.Context, userID int, opts ...*option.APIParamOption) ([]*model.Activity, error) {
	query := url.Values{}

	var ves validation.Errors
	if ve := validate.ValidateUserID(userID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, s.base.ApplyOptions(query, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return s.base.Fetch(ctx, spath, query)
}

func NewActivityService(method *client.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
