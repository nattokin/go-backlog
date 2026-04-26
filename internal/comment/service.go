package comment

import (
	"context"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

// IssueService handles communication with the issue comment-related methods of the Backlog API.
type IssueService struct {
	method *core.Method
}

// All returns a list of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment-list
func (s *IssueService) All(ctx context.Context, issueIDOrKey string, opts ...core.RequestOption) ([]*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	query := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamMinID,
		core.ParamMaxID,
		core.ParamCount,
		core.ParamOrder,
	}
	if err := core.ApplyOptions(query, validTypes, opts...); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
	resp, err := s.method.Get(ctx, spath, query)
	if err != nil {
		return nil, err
	}

	v := []*model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

// Add adds a comment to an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/add-comment
func (s *IssueService) Add(ctx context.Context, issueIDOrKey string, content string, opts ...core.RequestOption) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	option := &core.OptionService{}
	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamContent,
		core.ParamNotifiedUserIDs,
		core.ParamAttachmentIDs,
	}
	options := append(
		[]core.RequestOption{option.WithContent(content)},
		opts...,
	)
	if err := core.ApplyOptions(form, validTypes, options...); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments")
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

// Count returns the number of comments on an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/count-comment
func (s *IssueService) Count(ctx context.Context, issueIDOrKey string) (int, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return 0, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", "count")
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return 0, err
	}

	v := map[string]int{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return 0, err
	}

	return v["count"], nil
}

// One returns information about a specific comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-comment
func (s *IssueService) One(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	resp, err := s.method.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Delete deletes a comment from an issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-comment
func (s *IssueService) Delete(ctx context.Context, issueIDOrKey string, commentID int) (*model.Comment, error) {
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
func (s *IssueService) Update(ctx context.Context, issueIDOrKey string, commentID int, option core.RequestOption) (*model.Comment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateCommentID(commentID); err != nil {
		return nil, err
	}

	form := url.Values{}
	validTypes := []core.APIParamOptionType{
		core.ParamContent,
	}
	if err := core.ApplyOptions(form, validTypes, option); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "comments", strconv.Itoa(commentID))
	resp, err := s.method.Patch(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := model.Comment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// Notifications returns a list of notifications on a comment.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-comment-notifications
func (s *IssueService) Notifications(ctx context.Context, issueIDOrKey string, commentID int) ([]*model.Notification, error) {
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
func (s *IssueService) Notify(ctx context.Context, issueIDOrKey string, commentID int, userIDs []int) (*model.Comment, error) {
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

// ──────────────────────────────────────────────────────────────
//  Constructor
// ──────────────────────────────────────────────────────────────

// NewIssueService creates and returns a new comment IssueService.
func NewIssueService(method *core.Method) *IssueService {
	return &IssueService{method: method}
}
