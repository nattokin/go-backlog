// Package activity implements shared HTTP logic for activity-related Backlog API endpoints.
package activity

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ValidOptionTypes are the option types accepted by activity list endpoints.
var ValidOptionTypes = []option.APIParamOptionType{
	option.ParamActivityTypeIDs,
	option.ParamMinID,
	option.ParamMaxID,
	option.ParamCount,
	option.ParamOrder,
}

// Service holds shared HTTP logic for activity-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *client.Method
}

// ApplyOptions validates and applies opts to query.
// Callers should call this before Fetch to separate option validation from HTTP.
func (s *Service) ApplyOptions(query url.Values, opts ...*option.APIParamOption) error {
	return option.ApplyOptions(query, ValidOptionTypes, opts...)
}

// Fetch executes the GET request with pre-built query values.
func (s *Service) Fetch(ctx context.Context, spath string, query url.Values) ([]*model.Activity, error) {
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Activity{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewService(method *client.Method) *Service {
	return &Service{
		method: method,
	}
}
