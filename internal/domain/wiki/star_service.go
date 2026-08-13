package wiki

import (
	"context"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// StarService handles wiki star-related Backlog API calls.
type StarService struct {
	method *client.Method
}

// List returns a list of stars on the wiki page.
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func NewStarService(method *client.Method) *StarService {
	return &StarService{method: method}
}
