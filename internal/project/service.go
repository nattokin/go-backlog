package project

import (
	"context"
	"net/url"
	"path"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

type Service struct {
	method *core.Method
}

func (s *Service) All(ctx context.Context, opts ...core.RequestOption) ([]*model.Project, error) {

	query := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamAll, core.ParamArchived}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, "projects", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) One(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Create(ctx context.Context, key, name string, opts ...core.RequestOption) (*model.Project, error) {
	option := &core.OptionService{}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamKey, core.ParamName, core.ParamChartEnabled, core.ParamSubtaskingEnabled, core.ParamProjectLeaderCanEditProjectLeader, core.ParamTextFormattingRule}
	options := append([]core.RequestOption{option.WithKey(key), option.WithName(name)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	resp, err := s.method.Post(ctx, "projects", form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Update(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamKey, core.ParamName, core.ParamChartEnabled, core.ParamSubtaskingEnabled,
		core.ParamProjectLeaderCanEditProjectLeader, core.ParamTextFormattingRule, core.ParamArchived,
	}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Delete(ctx context.Context, projectIDOrKey string) (*model.Project, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey)
	resp, err := s.method.Delete(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := model.Project{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
