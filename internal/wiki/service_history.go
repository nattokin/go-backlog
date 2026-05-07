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

// HistorySevice handles communication with the wiki history-related methods of the Backlog API.
type HistorySevice struct {
	method *core.Method
}

// List returns the version history of a wiki page.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-history/
func (s *HistorySevice) List(ctx context.Context, wikiID int) ([]*model.WikiHistory, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "history")
	resp, err := s.method.Get(ctx, spath, url.Values{})
	if err != nil {
		return nil, err
	}

	v := []*model.WikiHistory{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewHistoryService creates and returns a new history WikiService.
func NewHistoryService(method *core.Method) *HistorySevice {
	return &HistorySevice{method: method}
}
