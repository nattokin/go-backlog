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
	// Run option validation first (fail-fast on nil/invalidKey via ApplyOptions
	// inside base.List). If that returns a non-ValidationErrors error, we still
	// need to check the argument and surface it as ValidationErrors.
	// However, to avoid firing HTTP when the argument is invalid, we validate
	// the argument up front and return early before calling the shared service.
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		// Still run ApplyOptions to respect fail-fast (nil option takes priority).
		var dummy [0]core.APIParamOptionType
		if err := core.ApplyOptions(nil, dummy[:], opts...); err != nil {
			var ves core.ValidationErrors
			if errors.As(err, &ves) {
				ves = append(ves, ve)
				return nil, ves
			}
			// fail-fast (InvalidOptionError or InvalidOptionKeyError): return as-is
			return nil, err
		}
		return nil, core.ValidationErrors{ve}
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
	return s.base.List(ctx, spath, opts...)
}

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
