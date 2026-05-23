// Package star implements the Backlog Star API service.
package star

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles star-related Backlog API calls.
type Service struct {
	method *core.Method
}

// Add adds a star to a resource (issue, comment, wiki page, pull request, or pull request comment).
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *Service) Add(ctx context.Context, option core.RequestOption) error {
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamIssueID,
		core.ParamCommentID,
		core.ParamWikiID,
		core.ParamPullRequestID,
		core.ParamPullRequestCommentID,
	}
	if err := core.ApplyOptions(form, validTypes, option); err != nil {
		return err
	}

	if _, err := s.method.Post(ctx, "stars", form); err != nil {
		return err
	}

	return nil
}

// Remove removes a star by its ID.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/remove-star
func (s *Service) Remove(ctx context.Context, id int) error {
	if err := validate.ValidateStarID(id); err != nil {
		return err
	}

	form := url.Values{}
	form.Set("id", strconv.Itoa(id))

	if _, err := s.method.Delete(ctx, "stars", form); err != nil {
		return err
	}

	return nil
}

func NewService(method *core.Method) *Service {
	return &Service{method: method}
}
