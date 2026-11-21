package wiki

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// HistoryService handles wiki history-related Backlog API calls.
type HistoryService struct {
	method *client.Method
}

// List returns the version history of a wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-history/
func (s *HistoryService) List(ctx context.Context, wikiID int, opts ...*option.APIParamOption) ([]*model.WikiHistory, error) {
	query := url.Values{}
	validOptionKeys := []option.APIParamOptionType{option.ParamMinID, option.ParamMaxID, option.ParamCount, option.ParamOrder}

	var ves validation.Errors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(query, validOptionKeys, opts...)); err != nil {
		return nil, err
	}

	spath := client.PathJoin("wikis", strconv.Itoa(wikiID), "history")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.WikiHistory{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// NewHistoryService returns a *HistoryService that sends requests through method.
func NewHistoryService(method *client.Method) *HistoryService {
	return &HistoryService{method: method}
}
