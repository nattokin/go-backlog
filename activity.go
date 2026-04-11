package backlog

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
)

func getActivityList(ctx context.Context, m *core.Method, spath string, opts ...RequestOption) ([]*Activity, error) {
	query := url.Values{}
	validOptionKeys := []apiParamOptionType{core.ParamActivityTypeIDs, core.ParamMinID, core.ParamMaxID, core.ParamCount, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := m.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*Activity{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

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
	if err := validateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
	return getActivityList(ctx, s.method, spath, opts...)
}

// SpaceActivityService handles communication with the space activities-related methods of the Backlog API.
type SpaceActivityService struct {
	method *core.Method

	Option *ActivityOptionService
}

// List returns a list of activities in your space.
//
// This method supports options returned by methods in "*Client.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(ctx context.Context, opts ...RequestOption) ([]*Activity, error) {
	return getActivityList(ctx, s.method, "space/activities", opts...)
}

// UserActivityService handles communication with the user activities-related methods of the Backlog API.
type UserActivityService struct {
	method *core.Method

	Option *ActivityOptionService
}

// List returns a list of user activities.
//
// This method supports options returned by methods in "*Client.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates
func (s *UserActivityService) List(ctx context.Context, userID int, opts ...RequestOption) ([]*Activity, error) {
	uID := UserID(userID)
	if err := uID.validate(); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return getActivityList(ctx, s.method, spath, opts...)
}
