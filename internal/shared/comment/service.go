// Package comment implements shared HTTP logic for comment-related Backlog API endpoints.
package comment

import (
	"context"
	"net/url"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
)

// ListValidTypes are the option types accepted by comment list endpoints.
var ListValidTypes = []option.APIParamOptionType{
	option.ParamMinID,
	option.ParamMaxID,
	option.ParamCount,
	option.ParamOrder,
}

// AddValidTypes are the option types accepted by comment add endpoints.
var AddValidTypes = []option.APIParamOptionType{
	option.ParamContent,
	option.ParamNotifiedUserIDs,
	option.ParamAttachmentIDs,
}

// Service holds shared HTTP logic for comment-related Backlog API endpoints.
// It is spath-agnostic: callers supply the full sub-path and are responsible
// for validation and path construction.
type Service struct {
	method *client.Method
}

// ApplyListOptions validates and applies opts to query for List.
func (s *Service) ApplyListOptions(query url.Values, opts ...*option.APIParamOption) error {
	return option.ApplyOptions(query, ListValidTypes, opts...)
}

// FetchList executes the GET request for comment listing.
func (s *Service) FetchList(ctx context.Context, spath string, query url.Values) ([]*model.Comment, error) {
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Comment{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// ApplyAddOptions validates and applies content + opts to form for Add.
func (s *Service) ApplyAddOptions(form url.Values, content string, opts ...*option.APIParamOption) error {
	optSvc := &option.OptionService{}
	options := append(
		[]*option.APIParamOption{optSvc.WithContent(content)},
		opts...,
	)
	return option.ApplyOptions(form, AddValidTypes, options...)
}

// FetchAdd executes the POST request for adding a comment.
func (s *Service) FetchAdd(ctx context.Context, spath string, form url.Values) (*model.Comment, error) {
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := client.DecodeResponse(resp, &v); err != nil {
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
	if err := client.DecodeResponse(resp, &v); err != nil {
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ApplyUpdateOptions validates content and writes it into form.
// Returns a *ValidationError if content is invalid, nil otherwise.
func (s *Service) ApplyUpdateOptions(form url.Values, content string) *core.ValidationError {
	opt := (&option.OptionService{}).WithContent(content)
	if ve := opt.Check(); ve != nil {
		return ve
	}
	opt.Set(form)
	return nil
}

// FetchUpdate executes the PATCH request for updating a comment.
func (s *Service) FetchUpdate(ctx context.Context, spath string, form url.Values) (*model.Comment, error) {
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewService(method *client.Method) *Service {
	return &Service{method: method}
}
