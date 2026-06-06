package pullrequest

import (
	"context"
	"errors"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/shared/comment"
	"github.com/nattokin/go-backlog/internal/validate"
)

// CommentService handles pull request comment-related Backlog API calls.
// It delegates all HTTP operations to the shared comment.Service and is
// responsible only for validation and spath construction.
type CommentService struct {
	base *comment.Service
}

// List returns a list of comments on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-comment
func (s *CommentService) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, opts ...*core.APIParamOption) ([]*model.Comment, error) {
	query := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if err := s.base.ApplyListOptions(query, opts...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.FetchList(ctx, spath, query)
}

// Add adds a comment to a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment
func (s *CommentService) Add(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, content string, opts ...*core.APIParamOption) (*model.Comment, error) {
	form := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if err := s.base.ApplyAddOptions(form, content, opts...); err != nil {
		var optVes core.ValidationErrors
		if !errors.As(err, &optVes) {
			return nil, err
		}
		ves = append(ves, optVes...)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.FetchAdd(ctx, spath, form)
}

// Count returns the number of comments on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-request-comments
func (s *CommentService) Count(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int) (int, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return 0, ves
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", "count")
	return s.base.Count(ctx, spath)
}

// Update updates a comment on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-pull-request-comment-information
func (s *CommentService) Update(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, commentID int, content string) (*model.Comment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidatePRNumber(prNumber); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}

	form := url.Values{}
	if ve := s.base.ApplyUpdateOptions(form, content); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", strconv.Itoa(commentID))
	return s.base.FetchUpdate(ctx, spath, form)
}

func NewCommentService(method *core.Method) *CommentService {
	return &CommentService{
		base: comment.NewService(method),
	}
}
