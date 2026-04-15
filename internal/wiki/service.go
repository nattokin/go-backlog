package wiki

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

type Service struct {
	method *core.Method
}

func (s *Service) All(ctx context.Context, projectIDOrKey string, opts ...core.RequestOption) ([]*model.Wiki, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamKeyword}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis", query)
	if err != nil {
		return nil, err
	}

	v := []*model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) Count(ctx context.Context, projectIDOrKey string) (int, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return 0, err
	}

	query := url.Values{}
	query.Set("projectIdOrKey", projectIDOrKey)

	resp, err := s.method.Get(ctx, "wikis/count", query)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

func (s *Service) One(ctx context.Context, wikiID int) (*model.Wiki, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Create(ctx context.Context, projectID int, name, content string, opts ...core.RequestOption) (*model.Wiki, error) {
	if err := validate.ValidateProjectID(projectID); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamContent, core.ParamMailNotify}
	options := append([]core.RequestOption{option.WithName(name), option.WithContent(content)}, opts...)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	form.Set("projectId", strconv.Itoa(projectID))

	resp, err := s.method.Post(ctx, "wikis", form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Update(ctx context.Context, wikiID int, option core.RequestOption, opts ...core.RequestOption) (*model.Wiki, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamName, core.ParamContent, core.ParamMailNotify}
	options := append([]core.RequestOption{option}, opts...)

	if !core.HasRequiredOption(options, []core.APIParamOptionType{core.ParamName, core.ParamContent}) {
		return nil, core.NewValidationError("requires an option to modify wiki content or name (WithName or WithContent)")
	}

	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Delete(ctx context.Context, wikiID int, opts ...core.RequestOption) (*model.Wiki, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{core.ParamMailNotify}
	if err := core.ApplyOptions(form, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID))
	resp, err := s.method.Delete(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Wiki{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method, option *core.OptionService) *Service {
	return &Service{
		method: method,
	}
}
