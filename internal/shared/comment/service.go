// Package comment implements shared HTTP logic for comment-related Backlog API endpoints.
package comment

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
)

// ListValidTypes are the option types accepted by comment list endpoints.
var ListValidTypes = []core.APIParamOptionType{
	core.ParamMinID,
	core.ParamMaxID,
	core.ParamCount,
	core.ParamOrder,
}

// AddValidTypes are the option types accepted by comment add endpoints.
var AddValidTypes = []core.APIParamOptionType{
	core.ParamContent,
	core.ParamNotifiedUserIDs,
	core.ParamAttachmentIDs,
}

// Service holds shared HTTP logic for comment-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *core.Method
}

// ApplyListOptions validates and applies opts to query for List.
func (s *Service) ApplyListOptions(query url.Values, opts ...*core.APIParamOption) error {
	return core.ApplyOptions(query, ListValidTypes, opts...)
}

// FetchList executes the GET request for comment listing.
func (s *Service) FetchList(ctx context.Context, spath string, query url.Values) ([]*model.Comment, error) {
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

// ApplyAddOptions validates and applies content + opts to form for Add.
func (s *Service) ApplyAddOptions(form url.Values, content string, opts ...*core.APIParamOption) error {
	option := &core.OptionService{}
	options := append(
		[]*core.APIParamOption{option.WithContent(content)},
		opts...,
	)
	return core.ApplyOptions(form, AddValidTypes, options...)
}

// FetchAdd executes the POST request for adding a comment.
func (s *Service) FetchAdd(ctx context.Context, spath string, form url.Values) (*model.Comment, error) {
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

// ApplyUpdateOptions validates content and writes it into form.
// Returns a *ValidationError if content is invalid, nil otherwise.
func (s *Service) ApplyUpdateOptions(form url.Values, content string) *core.ValidationError {
	option := (&core.OptionService{}).WithContent(content)
	if ve := option.Check(); ve != nil {
		return ve
	}
	option.Set(form)
	return nil
}

// FetchUpdate executes the PATCH request for updating a comment.
func (s *Service) FetchUpdate(ctx context.Context, spath string, form url.Values) (*model.Comment, error) {
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

func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
