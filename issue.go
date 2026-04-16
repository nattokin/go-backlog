package backlog

import (
	"context"
	"time"

	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/model"
)

// ──────────────────────────────────────────────────────────────
//  Issue models
// ──────────────────────────────────────────────────────────────

// Category represents an issue category.
type Category struct {
	ID           int
	Name         string
	DisplayOrder int
}

// Issue represents a issue of Backlog.
type Issue struct {
	ID             int
	ProjectID      int
	IssueKey       string
	KeyID          int
	IssueType      *IssueType
	Summary        string
	Description    string
	Resolutions    []*Resolution
	Priority       *Priority
	Status         *Status
	Assignee       *User
	Category       []*Category
	Versions       []*Version
	Milestone      []*Version
	StartDate      time.Time
	DueDate        time.Time
	EstimatedHours int
	ActualHours    int
	ParentIssueID  int
	CreatedUser    *User
	Created        time.Time
	UpdatedUser    *User
	Updated        time.Time
	CustomFields   []*CustomField
	Attachments    []*Attachment
	SharedFiles    []*SharedFile
	Stars          []*Star
}

// IssueType represents type of Issue.
type IssueType struct {
	ID           int
	ProjectID    int
	Name         string
	Color        string
	DisplayOrder int
}

// Priority represents a priority.
type Priority struct {
	ID   int
	Name string
}

// Resolution represents a resolution.
type Resolution struct {
	ID   int
	Name string
}

// ──────────────────────────────────────────────────────────────
//  IssueService
// ──────────────────────────────────────────────────────────────

// IssueService handles communication with the issue-related methods of the Backlog API.
type IssueService struct {
	base *issue.Service

	Attachment *IssueAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  IssueAttachmentService
// ──────────────────────────────────────────────────────────────

// IssueAttachmentService handles communication with the issue attachment-related methods of the Backlog API.
type IssueAttachmentService struct {
	base *attachment.IssueService
}

// List returns a list of all attachments in the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments
func (s *IssueAttachmentService) List(ctx context.Context, issueIDOrKey string) ([]*Attachment, error) {
	v, err := s.base.List(ctx, issueIDOrKey)
	return attachmentsFromModel(v), convertError(err)
}

// Remove removes a file attached to the issue.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment
func (s *IssueAttachmentService) Remove(ctx context.Context, issueIDOrKey string, attachmentID int) (*Attachment, error) {
	v, err := s.base.Remove(ctx, issueIDOrKey, attachmentID)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newIssueService(method *core.Method) *IssueService {
	return &IssueService{
		base:       issue.NewService(method),
		Attachment: newIssueAttachmentService(method),
	}
}

func newIssueAttachmentService(method *core.Method) *IssueAttachmentService {
	return &IssueAttachmentService{
		base: attachment.NewIssueService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Helpers
// ──────────────────────────────────────────────────────────────

func issueFromModel(m *model.Issue) *Issue {
	if m == nil {
		return nil
	}
	resolutions := make([]*Resolution, len(m.Resolutions))
	for i, v := range m.Resolutions {
		if v == nil {
			continue
		}
		resolutions[i] = &Resolution{ID: v.ID, Name: v.Name}
	}
	categories := make([]*Category, len(m.Category))
	for i, v := range m.Category {
		if v == nil {
			continue
		}
		categories[i] = &Category{ID: v.ID, Name: v.Name, DisplayOrder: v.DisplayOrder}
	}
	versions := make([]*Version, len(m.Versions))
	for i, v := range m.Versions {
		versions[i] = versionFromModel(v)
	}
	milestone := make([]*Version, len(m.Milestone))
	for i, v := range m.Milestone {
		milestone[i] = versionFromModel(v)
	}
	customFields := make([]*CustomField, len(m.CustomFields))
	for i, v := range m.CustomFields {
		customFields[i] = customFieldFromModel(v)
	}
	sharedFiles := make([]*SharedFile, len(m.SharedFiles))
	for i, v := range m.SharedFiles {
		sharedFiles[i] = sharedFileFromModel(v)
	}
	stars := make([]*Star, len(m.Stars))
	for i, v := range m.Stars {
		stars[i] = starFromModel(v)
	}
	var issueType *IssueType
	if m.IssueType != nil {
		issueType = &IssueType{
			ID:           m.IssueType.ID,
			ProjectID:    m.IssueType.ProjectID,
			Name:         m.IssueType.Name,
			Color:        m.IssueType.Color,
			DisplayOrder: m.IssueType.DisplayOrder,
		}
	}
	var priority *Priority
	if m.Priority != nil {
		priority = &Priority{ID: m.Priority.ID, Name: m.Priority.Name}
	}
	return &Issue{
		ID:             m.ID,
		ProjectID:      m.ProjectID,
		IssueKey:       m.IssueKey,
		KeyID:          m.KeyID,
		IssueType:      issueType,
		Summary:        m.Summary,
		Description:    m.Description,
		Resolutions:    resolutions,
		Priority:       priority,
		Status:         statusFromModel(m.Status),
		Assignee:       userFromModel(m.Assignee),
		Category:       categories,
		Versions:       versions,
		Milestone:      milestone,
		StartDate:      m.StartDate,
		DueDate:        m.DueDate,
		EstimatedHours: m.EstimatedHours,
		ActualHours:    m.ActualHours,
		ParentIssueID:  m.ParentIssueID,
		CreatedUser:    userFromModel(m.CreatedUser),
		Created:        m.Created,
		UpdatedUser:    userFromModel(m.UpdatedUser),
		Updated:        m.Updated,
		CustomFields:   customFields,
		Attachments:    attachmentsFromModel(m.Attachments),
		SharedFiles:    sharedFiles,
		Stars:          stars,
	}
}
