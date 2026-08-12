// Package star implements the Backlog Star API service.
package star

import (
	"context"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/validate"
)

// Service handles star-related Backlog API calls.
type Service struct {
	method *core.Method
}

// Add adds a star to a resource (issue, comment, wiki page, pull request, or pull request comment).
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-star
func (s *Service) Add(ctx context.Context, opt *option.APIParamOption) error {
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamIssueID,
		option.ParamCommentID,
		option.ParamWikiID,
		option.ParamPullRequestID,
		option.ParamPullRequestCommentID,
	}
	if err := option.ApplyOptions(form, validTypes, opt); err != nil {
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
