package activity

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// Service holds shared HTTP logic for activity-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

func (s *Service) List(ctx context.Context, spath string, opts ...core.RequestOption) ([]*model.Activity, error) {
	query := url.Values{}
	validOptionKeys := []core.APIParamOptionType{core.ParamActivityTypeIDs, core.ParamMinID, core.ParamMaxID, core.ParamCount, core.ParamOrder}
	if err := core.ApplyOptions(query, validOptionKeys, opts...); err != nil {
		return nil, err
	}

	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Activity{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewService(method *core.Method) *Service {
	return &Service{
		method: method,
	}
}
