package issue

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/client"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/option"
	"github.com/nattokin/go-backlog/internal/shared/comment"
	"github.com/nattokin/go-backlog/internal/validate"
)

type CommentService struct {
	base   *comment.Service
	method *client.Method
}

// List returns a list of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment-list
func (s *CommentService) List(ctx context.Context, issueIDOrKey string, opts ...*option.APIParamOption) ([]*model.Comment, error) {
	query := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, s.base.ApplyListOptions(query, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
	return s.base.FetchList(ctx, spath, query)
}

// Add adds a comment to an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment
func (s *CommentService) Add(ctx context.Context, issueIDOrKey string, content string, opts ...*option.APIParamOption) (*model.Comment, error) {
	form := url.Values{}

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, s.base.ApplyAddOptions(form, content, opts...)); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
	return s.base.FetchAdd(ctx, spath, form)
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
	if err := client.DecodeResponse(resp, &v); err != nil {
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

	form := url.Values{}
	if ve := s.base.ApplyUpdateOptions(form, content); ve != nil {
		ves = append(ves, ve)
	}
	if len(ves) > 0 {
		return nil, ves
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	return s.base.FetchUpdate(ctx, spath, form)
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
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Notify sends notifications for a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment-notification
func (s *CommentService) Notify(ctx context.Context, issueIDOrKey string, commentID int, userIDs []int) (*model.Comment, error) {
	optSvc := &option.OptionService{}
	form := url.Values{}
	validTypes := []option.APIParamOptionType{
		option.ParamNotifiedUserIDs,
	}

	var ves core.ValidationErrors
	if ve := validate.ValidateIssueIDOrKey(issueIDOrKey); ve != nil {
		ves = append(ves, ve)
	}
	if ve := validate.ValidateCommentID(commentID); ve != nil {
		ves = append(ves, ve)
	}
	if err := option.MergeValidationErrors(ves, option.ApplyOptions(form, validTypes, optSvc.WithNotifiedUserIDs(userIDs))); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID), "notifications")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := client.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

func NewCommentService(method *client.Method) *CommentService {
	return &CommentService{
		base:   comment.NewService(method),
		method: method,
	}
}
