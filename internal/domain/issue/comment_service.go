package issue

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

// CommentService handles issue comment-related Backlog API calls.
// It delegates HTTP operations to the shared comment.Service and is
// responsible only for validation and spath construction.
type CommentService struct {
	base   *comment.Service
	method *core.Method
}

// List returns a list of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment-list
func (s *CommentService) List(ctx context.Context, issueIDOrKey string, opts ...*core.APIParamOption) ([]*model.Comment, error) {
	spath := path.Join("issues", issueIDOrKey, "comments")
	result, err := s.base.List(ctx, spath, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
				return nil, core.ValidationErrors{ve}
			}
			return nil, err
		}
		if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	return result, nil
}

// Add adds a comment to an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment
func (s *CommentService) Add(ctx context.Context, issueIDOrKey string, content string, opts ...*core.APIParamOption) (*model.Comment, error) {
	spath := path.Join("issues", issueIDOrKey, "comments")
	result, err := s.base.Add(ctx, spath, content, opts...)
	if err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
				return nil, core.ValidationErrors{ve}
			}
			return nil, err
		}
		if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		return nil, core.ValidationErrors{ve}
	}

	return result, nil
}

// Count returns the number of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-comment
func (s *CommentService) Count(ctx context.Context, issueIDOrKey string) (int, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return 0, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", "count")
	return s.base.Count(ctx, spath)
}

// One returns a single comment on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment
func (s *CommentService) One(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	return s.base.One(ctx, spath)
}

// Delete deletes a comment from an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-comment
func (s *CommentService) Delete(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	resp, err := s.method.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Update updates a comment on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-comment
func (s *CommentService) Update(ctx context.Context, issueIDOrKey string, commentID int, content string) (*model.Comment, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	result, err := s.base.Update(ctx, spath, content)
	if err != nil {
		var ve *core.ValidationError
		if errors.As(err, &ve) {
			ves = append(ves, ve)
			return nil, ves
		}
		if len(ves) > 0 {
			return nil, ves
		}
		return nil, err
	}

	if len(ves) > 0 {
		return nil, ves
	}

	return result, nil
}

// Notifications returns a list of notifications on a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-comment-notifications
func (s *CommentService) Notifications(ctx context.Context, issueIDOrKey string, commentID int) ([]*model.Notification, error) {
	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID), "notifications")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Notification{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Notify sends notifications for a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment-notification
func (s *CommentService) Notify(ctx context.Context, issueIDOrKey string, commentID int, userIDs []int) (*model.Comment, error) {
	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamNotifiedUserIDs,
	}
	if err := core.ApplyOptions(form, validTypes, option.WithNotifiedUserIDs(userIDs)); err != nil {
		var ves core.ValidationErrors
		if !errors.As(err, &ves) {
			return nil, err
		}
		if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
			ves = append(ves, ve)
		}
		if ve := validate.ValidateCommentID(commentID); ve != nil {
			ves = append(ves, ve)
		}
		return nil, ves
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID), "notifications")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewCommentService(method *core.Method) *CommentService {
	return &CommentService{
		base:   comment.NewService(method),
		method: method,
	}
}
