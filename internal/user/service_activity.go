package user

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

func (s *ActivityService) List(ctx context.Context, userID int, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return s.base.List(ctx, spath, opts...)
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
