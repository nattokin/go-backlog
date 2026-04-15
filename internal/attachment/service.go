package attachment

import (
	"context"
	"errors"
	"io"
	"net/url"
	"path"
	"strconv"

	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/validate"
)

func ListAttachments(ctx context.Context, m *core.Method, spath string) ([]*model.Attachment, error) {
	resp, err := m.Get(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := []*model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func RemoveAttachment(ctx context.Context, m *core.Method, spath string) (*model.Attachment, error) {
	resp, err := m.Delete(ctx, spath, nil)
	if err != nil {
		return nil, err
	}

	v := &model.Attachment{}
	if err := core.DecodeResponse(resp, v); err != nil {
		return nil, err
	}

	return v, nil
}

// ──────────────────────────────────────────────────────────────
//  IssueAttachmentService
// ──────────────────────────────────────────────────────────────

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	method *core.Method
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*model.Attachment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return ListAttachments(ctx, s.method, spath)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.Attachment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments", strconv.Itoa(attachmentID))
	return RemoveAttachment(ctx, s.method, spath)
}

// ──────────────────────────────────────────────────────────────
//  PullRequestAttachmentService
// ──────────────────────────────────────────────────────────────

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	method *core.Method
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*model.Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments")
	return ListAttachments(ctx, s.method, spath)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*model.Attachment, error) {
	if err := validate.ValidateProjectIDOrKey(projectIDOrKey); err != nil {
		return nil, err
	}
	if err := validate.ValidateRepositoryIDOrName(repositoryIDOrName); err != nil {
		return nil, err
	}
	if err := validate.ValidatePRNumber(prNumber); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("projects", projectIDOrKey, "git", "repositories", repositoryIDOrName, "pullRequests", strconv.Itoa(prNumber), "attachments", strconv.Itoa(attachmentID))
	return RemoveAttachment(ctx, s.method, spath)
}

// ──────────────────────────────────────────────────────────────
//  SpaceAttachmentService
// ──────────────────────────────────────────────────────────────

// SpaceAttachmentService handles communication with the space attachment-related methods of the Backlog API.
type SpaceAttachmentService struct {
	method *core.Method
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *SpaceAttachmentService) Upload(ctx context.Context, fileName string, r io.Reader) (*model.Attachment, error) {
	resp, err := s.method.Upload(ctx, "space/attachment", fileName, r)
	if err != nil {
		return nil, err
	}

	v := model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return &v, nil
}

// ──────────────────────────────────────────────────────────────
//  WikiAttachmentService
// ──────────────────────────────────────────────────────────────

type WikiAttachmentService struct {
	method *core.Method
}

func (s *WikiAttachmentService) Attach(ctx context.Context, wikiID int, attachmentIDs []int) ([]*model.Attachment, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}
	if len(attachmentIDs) == 0 {
		return nil, errors.New("attachmentIDs must not be empty")
	}

	form := url.Values{}
	for _, id := range attachmentIDs {
		if err := validate.ValidateAttachmentID(id); err != nil {
			return nil, err
		}
		form.Add("attachmentId[]", strconv.Itoa(id))
	}

	spath := path.Join("wikis/", strconv.Itoa(wikiID), "/attachments")
	resp, err := s.method.Post(ctx, spath, form)
	if err != nil {
		return nil, err
	}

	v := []*model.Attachment{}
	if err := core.DecodeResponse(resp, &v); err != nil {
		return nil, err
	}

	return v, nil
}

func (s *WikiAttachmentService) List(ctx context.Context, wikiID int) ([]*model.Attachment, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments")
	return ListAttachments(ctx, s.method, spath)
}

func (s *WikiAttachmentService) Remove(ctx context.Context, wikiID, attachmentID int) (*model.Attachment, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}
	if err := validate.ValidateAttachmentID(attachmentID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments", strconv.Itoa(attachmentID))
	return RemoveAttachment(ctx, s.method, spath)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func NewPullRequestAttachmentService(method *core.Method) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		method: method,
	}
}

func NewIssueAttachmentService(method *core.Method) *IssueAttachmentService {
	return &IssueAttachmentService{
		method: method,
	}
}

func NewSpaceAttachmentService(method *core.Method) *SpaceAttachmentService {
	return &SpaceAttachmentService{
		method: method,
	}
}

func NewWikiAttachmentService(method *core.Method) *WikiAttachmentService {
	return &WikiAttachmentService{
		method: method,
	}
}
