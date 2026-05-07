package wiki

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// StarService handles communication with the wiki star-related methods of the Backlog API.
type StarService struct {
	method *core.Method
}

// List returns a list of stars on the wiki page with the given ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-star
func (s *StarService) List(ctx context.Context, wikiID int) ([]*model.Star, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "stars")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Star{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewStarService creates and returns a new wiki StarService.
func NewStarService(method *core.Method) *StarService {
	return &StarService{method: method}
}
