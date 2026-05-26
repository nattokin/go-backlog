package project

import (
	"context"
	"errors"
	"path"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/activity"
	"github.com/nattokin/go-backlog/internal/validate"
)

// ActivityService handles project activity-related Backlog API calls.
// It delegates HTTP operations to the shared activity.Service.
type ActivityService struct {
	base   *activity.Service
	method *core.Method
}

// List returns a list of activities in the project.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ActivityService) List(ctx context.Context, projectIDOrKey string, opts ...*core.APIParamOption) ([]*model.Activity, error) {
	// Validate opts via ApplyOptions first (fail-fast on nil/invalidKey).
	// If it returns ValidationErrors, append the argument error and return.
	spath := path.Join("projects", projectIDOrKey, "activities")
	result, err := s.base.List(ctx, spath, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			// fail-fast or HTTP error: validate arg and prepend if invalid
			if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
				return nil, core.ValidationErrors{ve}
			}
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	return result, nil
}

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
