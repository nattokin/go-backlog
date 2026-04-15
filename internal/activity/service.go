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

func getList(ctx context.Context, m *core.Method, spath string, opts ...core.RequestOption) ([]*model.Activity, error) {
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

type ProjectService struct {
	method *core.Method
}

func (s *ProjectService) List(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "activities")
	return getList(ctx, s.method, spath, opts...)
}

type SpaceService struct {
	method *core.Method
}

func (s *SpaceService) List(ctx context.Context, opts ...core.RequestOption) ([]*model.Activity, error) {
	return getList(ctx, s.method, "space/activities", opts...)
}

type UserService struct {
	method *core.Method
}

func (s *UserService) List(ctx context.Context, userID int, opts ...core.RequestOption) ([]*model.Activity, error) {
	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(userID), "activities")
	return getList(ctx, s.method, spath, opts...)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewProjectService(method *core.Method) *ProjectService {
	return &ProjectService{
		method: method,
	}
}

func NewSpaceService(method *core.Method) *SpaceService {
	return &SpaceService{
		method: method,
	}
}

func NewUserService(method *core.Method) *UserService {
	return &UserService{
		method: method,
	}
}
