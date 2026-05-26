package pullrequest

import (
	"context"
	"errors"
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
	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	result, err := s.base.List(ctx, spath, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidatePRNumber(prNumber); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

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
		return nil, ves
	}

	return result, nil
}

// Add adds a comment to a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment
func (s *CommentService) Add(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, content string, opts ...*core.APIParamOption) (*model.Comment, error) {
	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	result, err := s.base.Add(ctx, spath, content, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateProjectIDOrKey(projectIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateRepositoryIDOrName(repoIDOrName); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidatePRNumber(prNumber); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

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
		return nil, ves
	}

	return result, nil
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

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", strconv.Itoa(commentID))
	result, err := s.base.Update(ctx, spath, content)
	if err != nil {
		var updateVes core.ValidationErrors
		if !errors.As(err, &updateVes) {
			if len(ves) > 0 {
				return nil, ves
			}
			return nil, err
		}
		ves = append(ves, updateVes...)
		return nil, ves
	}

	if len(ves) > 0 {
		return nil, ves
	}

	return result, nil
}

func NewCommentService(method *core.Method) *CommentService {
	return &CommentService{
		base: comment.NewService(method),
	}
}
