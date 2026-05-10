package issue

import (
	"context"
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

// All returns a list of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment-list
func (s *CommentService) All(ctx context.Context, issueIDOrKey string, opts ...core.RequestOption) ([]*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
	return s.base.All(ctx, spath, opts...)
}

// Add adds a comment to an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment
func (s *CommentService) Add(ctx context.Context, issueIDOrKey string, content string, opts ...core.RequestOption) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
	return s.base.Add(ctx, spath, content, opts...)
}

// Count returns the number of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-comment
func (s *CommentService) Count(ctx context.Context, issueIDOrKey string) (int, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return 0, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", "count")
	return s.base.Count(ctx, spath)
}

// One returns a single comment on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment
func (s *CommentService) One(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	return s.base.One(ctx, spath)
}

// Delete deletes a comment from an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-comment
func (s *CommentService) Delete(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
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
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	return s.base.Update(ctx, spath, content)
}

// Notifications returns a list of notifications on a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-comment-notifications
func (s *CommentService) Notifications(ctx context.Context, issueIDOrKey string, commentID int) ([]*model.Notification, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
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
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamNotifiedUserIDs,
	}
	if err := core.ApplyOptions(form, validTypes, option.WithNotifiedUserIDs(userIDs)); err != nil {
		return nil, err
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
