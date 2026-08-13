package wiki

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

// HistorySevice handles wiki history-related Backlog API calls.
type HistorySevice struct {
	method *client.Method
}

// List returns the version history of a wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-history/
func (s *HistorySevice) List(ctx context.Context, wikiID int) ([]*model.WikiHistory, error) {
	var ves validation.Errors
	if ve := validate.ValidateWikiID(wikiID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "history")
	resp, err := s.method.Get(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := []*model.WikiHistory{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewHistoryService(method *client.Method) *HistorySevice {
	return &HistorySevice{method: method}
}
