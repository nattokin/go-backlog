package user

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

func get(ctx context.Context, m *core.Method, spath string) (*model.User, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func getList(ctx context.Context, m *core.Method, spath string, query url.Values) ([]*model.User, error) {
	resp, err := m.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func add(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func update(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func delete(ctx context.Context, m *core.Method, spath string, form url.Values) (*model.User, error) {
	resp, err := m.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.User{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

type Service struct {
	method *core.Method
}

func (s *Service) All(ctx context.Context) ([]*model.User, error) {
	return getList(ctx, s.method, "users", nil)
}

func (s *Service) One(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return get(ctx, s.method, spath)
}

func (s *Service) Own(ctx context.Context) (*model.User, error) {
	return get(ctx, s.method, "users/myself")
}

func (s *Service) Add(ctx context.Context, userID, password, name, mailAddress string, roleType model.Role) (*model.User, error) {
	if userID == "" {
		return nil, core.NewValidationError("userID must not be empty")
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamPassword, core.ParamName, core.ParamMailAddress, core.ParamRoleType}
	options := []core.RequestOption{
		option.WithPassword(password),
		option.WithName(name),
		option.WithMailAddress(mailAddress),
		option.WithRoleType(roleType),
	}
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	form.Set("userId", userID)

	return add(ctx, s.method, "users", form)
}

func (s *Service) Update(ctx context.Context, id int, opts ...core.RequestOption) (*model.User, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamUserID, core.ParamName, core.ParamPassword, core.ParamMailAddress, core.ParamRoleType}
	options := append([]core.RequestOption{option.WithUserID(id)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return update(ctx, s.method, spath, form)
}

func (s *Service) Delete(ctx context.Context, id int) (*model.User, error) {
	if err := validate.ValidateUserID(id); err != nil {
		return nil, err
	}

	spath := path.Join("users", strconv.Itoa(id))
	return delete(ctx, s.method, spath, nil)
}

type ProjectService struct {
	method *core.Method
}

func (s *ProjectService) All(ctx context.Context, projectIDOrKey string, excludeGroupMembers bool) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	query.Set("excludeGroupMembers", strconv.FormatBool(excludeGroupMembers))

	spath := path.Join("projects", projectIDOrKey, "users")
	return getList(ctx, s.method, spath, query)
}

func (s *ProjectService) Add(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "users")
	return add(ctx, s.method, spath, form)
}

func (s *ProjectService) Delete(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "users")
	return delete(ctx, s.method, spath, form)
}

func (s *ProjectService) AddAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return add(ctx, s.method, spath, form)
}

func (s *ProjectService) AdminAll(ctx context.Context, projectIDOrKey string) ([]*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return getList(ctx, s.method, spath, nil)
}

func (s *ProjectService) DeleteAdmin(ctx context.Context, projectIDOrKey string, userID int) (*model.User, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	if err := validate.ValidateUserID(userID); err != nil {
		return nil, err
	}

	form := url.Values{}
	form.Set("userId", strconv.Itoa(userID))

	spath := path.Join("projects", projectIDOrKey, "administrators")
	return delete(ctx, s.method, spath, form)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}

func NewProjectService(method *core.Method) *ProjectService {
	return &ProjectService{
		method: method,
	}
}
