package activity

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

func GetActivityList(ctx context.Context, m *core.Method, spath string, opts ...core.RequestOption) ([]*model.Activity, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamActivityTypeIDs, core.ParamMinID, core.ParamMaxID, core.ParamCount, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := m.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Activity{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

type ProjectActivityService struct {
	method *core.Method
}

func (s *ProjectActivityService) List(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
	return GetActivityList(ctx, s.method, spath, opts...)
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
func (s *SpaceActivityService) List(ctx context.Context, opts ...core.RequestOption) ([]*model.Activity, error) {
	return GetActivityList(ctx, s.method, "space/activities", opts...)
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
func (s *UserActivityService) List(ctx context.Context, userID int, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return GetActivityList(ctx, s.method, spath, opts...)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewProjectActivityService(method *core.Method) *ProjectActivityService {
	return &ProjectActivityService{
		method: method,
	}
}

func NewSpaceActivityService(method *core.Method, option *core.OptionService) *SpaceActivityService {
	return &SpaceActivityService{
		method: method,
		Option: &ActivityOptionService{},
	}
}

func NewUserActivityService(method *core.Method, option *core.OptionService) *UserActivityService {
	return &UserActivityService{
		method: method,
		Option: &ActivityOptionService{},
	}
}
