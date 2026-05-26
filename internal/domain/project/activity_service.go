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
	argVe := validate.ValidateProjectIDOrKey(projectIDOrKey)

	spath := path.Join("projects", projectIDOrKey, "activities")
	result, err := s.base.List(ctx, spath, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			if argVe != nil {
				return nil, core.ValidationErrors{argVe}
			}
			return nil, err
		}
		if argVe != nil {
			ves = append(ves, argVe)
		}
		return nil, ves
	}

	if argVe != nil {
		return nil, core.ValidationErrors{argVe}
	}

	return result, nil
}

func NewActivityService(method *core.Method) *ActivityService {
	return &ActivityService{
		base:   activity.NewService(method),
		method: method,
	}
}
