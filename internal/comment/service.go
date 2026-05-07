package comment

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service holds shared HTTP logic for comment-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

func (s *Service) All(ctx context.Context, spath string, opts ...core.RequestOption) ([]*model.Comment, error) {
	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamMinID,
		core.ParamMaxID,
		core.ParamCount,
		core.ParamOrder,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *Service) Add(ctx context.Context, spath, content string, opts ...core.RequestOption) (*model.Comment, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamContent,
		core.ParamNotifiedUserIDs,
		core.ParamAttachmentIDs,
	}
	options := append(
		[]core.RequestOption{option.WithContent(content)},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Count(ctx context.Context, spath string) (int, error) {
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

func (s *Service) One(ctx context.Context, spath string) (*model.Comment, error) {
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func (s *Service) Update(ctx context.Context, spath, content string) (*model.Comment, error) {
	option := (&core.OptionService{}).WithContent(content)
	if err := option.Check(); err != nil {
		return nil, err
	}
	form := url.Values{}
	option.Set(form)

	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewService creates and returns a new comment Service.
func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
