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
//  IssueService
// ──────────────────────────────────────────────────────────────

type IssueService struct {
	method *core.Method
}

func (s *IssueService) List(ctx context.Context, issueIDOrKey string) ([]*model.Attachment, error) {
	if err := validate.ValidateIssueIDOrKey(issueIDOrKey); err != nil {
		return nil, err
	}

	spath := path.Join("issues", issueIDOrKey, "attachments")
	return ListAttachments(ctx, s.method, spath)
}

func (s *IssueService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*model.Attachment, error) {
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
//  PullRequestService
// ──────────────────────────────────────────────────────────────

type PullRequestService struct {
	method *core.Method
}

func (s *PullRequestService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*model.Attachment, error) {
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

func (s *PullRequestService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*model.Attachment, error) {
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
//  SpaceService
// ──────────────────────────────────────────────────────────────

type SpaceService struct {
	method *core.Method
}

func (s *SpaceService) Upload(ctx context.Context, fileName string, r io.Reader) (*model.Attachment, error) {
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
//  WikiService
// ──────────────────────────────────────────────────────────────

type WikiService struct {
	method *core.Method
}

func (s *WikiService) Attach(ctx context.Context, wikiID int, attachmentIDs []int) ([]*model.Attachment, error) {
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

func (s *WikiService) List(ctx context.Context, wikiID int) ([]*model.Attachment, error) {
	if err := validate.ValidateWikiID(wikiID); err != nil {
		return nil, err
	}

	spath := path.Join("wikis", strconv.Itoa(wikiID), "attachments")
	return ListAttachments(ctx, s.method, spath)
}

func (s *WikiService) Remove(ctx context.Context, wikiID, attachmentID int) (*model.Attachment, error) {
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

func NewPullRequestService(method *core.Method) *PullRequestService {
	return &PullRequestService{
		method: method,
	}
}

func NewIssueService(method *core.Method) *IssueService {
	return &IssueService{
		method: method,
	}
}

func NewSpaceService(method *core.Method) *SpaceService {
	return &SpaceService{
		method: method,
	}
}

func NewWikiService(method *core.Method) *WikiService {
	return &WikiService{
		method: method,
	}
}
