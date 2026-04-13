package backlog

import (
	"context"
	"path"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/validate"
)

// ProjectActivityService handles communication with the project activities-related methods of the Backlog API.
type ProjectActivityService struct {
	method *core.Method

	Option *ActivityOptionService
}

// List returns a list of activities in the project.
//
// This method supports options returned by methods in "*Client.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates
func (s *ProjectActivityService) List(ctx context.Context, projectIDOrKey string, opts ...RequestOption) ([]*Activity, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
	return activity.GetActivityList(ctx, s.method, spath, opts...)
}
