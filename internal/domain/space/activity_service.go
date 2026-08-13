package space

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
)

// ActivityService handles space activity-related Backlog API calls.
// It delegates list operations to the shared activity.Service.
type ActivityService struct {
	base   *activity.Service
	method *client.Method
}

// List returns a list of activities in the space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space-activities
func (s *ActivityService) List(ctx context.Context, opts ...*option.APIParamOption) ([]*model.Activity, error) {
	query := url.Values{}
	if err := s.base.ApplyOptions(query, opts...); err != nil {
		return nil, err
	}
	return s.base.Fetch(ctx, "space/activities", query)
}

// One returns a single activity by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-activity
func (s *ActivityService) One(ctx context.Context, activityID int) (*model.Activity, error) {
	if err := validate.ValidateActivityID(activityID); err != nil {
		return nil, err
	}

	spath := path.Join("activities", strconv.Itoa(activityID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Activity{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewActivityService(method *client.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
