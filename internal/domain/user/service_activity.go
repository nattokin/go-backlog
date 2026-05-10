package user

import (
	"context"
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
func (s *ActivityService) List(ctx context.Context, userID int, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return s.base.List(ctx, spath, opts...)
}

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
