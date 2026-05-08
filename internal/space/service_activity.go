package space

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

type ActivityService struct {
	base   *activity.Service
	method *core.Method
}

func (s *ActivityService) List(ctx context.Context, opts ...core.RequestOption) ([]*model.Activity, error) {
	return s.base.List(ctx, "space/activities", opts...)
}

// Get returns a single activity by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-activity
func (s *ActivityService) Get(ctx context.Context, activityID int) (*model.Activity, error) {
	if err := validate.ValidateActivityID(activityID); err != nil {
		return nil, err
	}

	spath := path.Join("activities", strconv.Itoa(activityID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Activity{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
