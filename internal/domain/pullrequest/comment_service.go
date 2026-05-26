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

// validateListArgs returns ValidationErrors for the three path arguments,
// or nil if all are valid.
func validateListArgs(projectIDOrKey, repoIDOrName string, prNumber int) core.ValidationErrors {
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
	return ves
}

// applyOptsForArgVes runs ApplyOptions with empty validTypes to trigger
// fail-fast before returning the collected argument errors.
func applyOptsForArgVes(argVes core.ValidationErrors, opts []*core.APIParamOption) error {
	var dummy [0]core.APIParamOptionType
	if err := core.ApplyOptions(nil, dummy[:], opts...); err != nil {
		var ves core.ValidationErrors
		if errors.As(err, &ves) {
			ves = append(ves, argVes...)
			return ves
		}
		// fail-fast (InvalidOptionError): takes priority
		return err
	}
	return argVes
}

// List returns a list of comments on a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-pull-request-comment
func (s *CommentService) List(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, opts ...*core.APIParamOption) ([]*model.Comment, error) {
	if argVes := validateListArgs(projectIDOrKey, repoIDOrName, prNumber); len(argVes) > 0 {
		return nil, applyOptsForArgVes(argVes, opts)
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.List(ctx, spath, opts...)
}

// Add adds a comment to a pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment
func (s *CommentService) Add(ctx context.Context, projectIDOrKey string, repoIDOrName string, prNumber int, content string, opts ...*core.APIParamOption) (*model.Comment, error) {
	if argVes := validateListArgs(projectIDOrKey, repoIDOrName, prNumber); len(argVes) > 0 {
		return nil, applyOptsForArgVes(argVes, opts)
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments")
	return s.base.Add(ctx, spath, content, opts...)
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
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repoIDOrName, "pullRequests", strconv.Itoa(prNumber), "comments", strconv.Itoa(commentID))
	result, err := s.base.Update(ctx, spath, content)
	if err != nil {
		var ve *core.ValidationError
		if errors.As(err, &ve) {
			return nil, core.ValidationErrors{ve}
		}
		return nil, err
	}

	return result, nil
}

func NewCommentService(method *core.Method) *CommentService {
	return &CommentService{
		base: comment.NewService(method),
	}
}
