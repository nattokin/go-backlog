package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/pullrequest"
)

// PullRequest represents pull request of Backlog git.
type PullRequest struct {
	ID           int
	ProjectID    int
	RepositoryID int
	Number       int
	Summary      string
	Description  string
	Base         string
	Branch       string
	Status       *Status
	Assignee     *User
	Issue        *Issue
	BaseCommit   string
	BranchCommit string
	CloseAt      time.Time
	MergeAt      time.Time
	CreatedUser  *User
	Created      time.Time
	UpdatedUser  *User
	Updated      time.Time
	Attachments  []*Attachment
	Stars        []*Star
}

// ──────────────────────────────────────────────────────────────
//  PullRequestService
// ──────────────────────────────────────────────────────────────

// PullRequestService handles communication with the pull request-related methods of the Backlog API.
type PullRequestService struct {
	base *pullrequest.Service

	Attachment *PullRequestAttachmentService
}

// PullRequestAttachmentService handles communication with the pull request attachment-related methods of the Backlog API.
type PullRequestAttachmentService struct {
	base *attachment.PullRequestService
}

// List returns a list of all attachments in the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment
func (s *PullRequestAttachmentService) List(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int) ([]*Attachment, error) {
	v, err := s.base.List(ctx, projectIDOrKey, repositoryIDOrName, prNumber)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes a file attached to the pull request.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments
func (s *PullRequestAttachmentService) Remove(ctx context.Context, projectIDOrKey string, repositoryIDOrName string, prNumber int, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, projectIDOrKey, repositoryIDOrName, prNumber, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newPullRequestService(method *core.Method) *PullRequestService {
	return &PullRequestService{
		base:       pullrequest.NewService(method),
		Attachment: newPullRequestAttachmentService(method),
	}
}

func newPullRequestAttachmentService(method *core.Method) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		base: attachment.NewPullRequestService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func pullRequestFromModel(m *model.PullRequest) *PullRequest {
	if m == nil {
		return nil
	}
	attachments := make([]*Attachment, len(m.Attachments))
	for i, v := range m.Attachments {
		attachments[i] = attachmentFromModel(v)
	}
	stars := make([]*Star, len(m.Stars))
	for i, v := range m.Stars {
		stars[i] = starFromModel(v)
	}
	return &PullRequest{
		ID:           m.ID,
		ProjectID:    m.ProjectID,
		RepositoryID: m.RepositoryID,
		Number:       m.Number,
		Summary:      m.Summary,
		Description:  m.Description,
		Base:         m.Base,
		Branch:       m.Branch,
		Status:       statusFromModel(m.Status),
		Assignee:     userFromModel(m.Assignee),
		Issue:        issueFromModel(m.Issue),
		BaseCommit:   m.BaseCommit,
		BranchCommit: m.BranchCommit,
		CloseAt:      m.CloseAt,
		MergeAt:      m.MergeAt,
		CreatedUser:  userFromModel(m.CreatedUser),
		Created:      m.Created,
		UpdatedUser:  userFromModel(m.UpdatedUser),
		Updated:      m.Updated,
		Attachments:  attachments,
		Stars:        stars,
	}
}
