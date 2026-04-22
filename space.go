package backlog

import (
	"context"
	"io"
	"time"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/space"
)

// ──────────────────────────────────────────────────────────────
//  Space models
// ──────────────────────────────────────────────────────────────

// DiskUsageProject represents project's disk usage.
type DiskUsageProject struct {
	ProjectID  int
	Issue      int
	Wiki       int
	File       int
	Subversion int
	Git        int
	GitLFS     int
}

// DiskUsageSpace represents space's disk usage.
type DiskUsageSpace struct {
	Capacity   int
	Issue      int
	Wiki       int
	File       int
	Subversion int
	Git        int
	GitLFS     int
	Details    []*DiskUsageProject
}

// Space represents space of Backlog.
type Space struct {
	SpaceKey           string
	Name               string
	OwnerID            int
	Lang               string
	Timezone           string
	ReportSendTime     string
	TextFormattingRule Format
	Created            time.Time
	Updated            time.Time
}

// SpaceNotification represents a notification of Space.
type SpaceNotification struct {
	Content string
	Updated time.Time
}

// ──────────────────────────────────────────────────────────────
//  SpaceService
// ──────────────────────────────────────────────────────────────

// SpaceService handles communication with the space-related methods of the Backlog API.
type SpaceService struct {
	base *space.Service

	Activity   *SpaceActivityService
	Attachment *SpaceAttachmentService
}

// One returns information about your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space
func (s *SpaceService) One(ctx context.Context) (*Space, error) {
	v, err := s.base.One(ctx)
	return spaceFromModel(v), convertError(err)
}

// DiskUsage returns information about the disk usage of your space.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space-disk-usage
func (s *SpaceService) DiskUsage(ctx context.Context) (*DiskUsageSpace, error) {
	v, err := s.base.DiskUsage(ctx)
	return diskUsageSpaceFromModel(v), convertError(err)
}

// Notification returns the space notification.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-space-notification
func (s *SpaceService) Notification(ctx context.Context) (*SpaceNotification, error) {
	v, err := s.base.Notification(ctx)
	return spaceNotificationFromModel(v), convertError(err)
}

// UpdateNotification updates the space notification.
//
// content must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/update-space-notification
func (s *SpaceService) UpdateNotification(ctx context.Context, content string) (*SpaceNotification, error) {
	v, err := s.base.UpdateNotification(ctx, content)
	return spaceNotificationFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  SpaceActivityService
// ──────────────────────────────────────────────────────────────

// SpaceActivityService handles communication with the space activities-related methods of the Backlog API.
type SpaceActivityService struct {
	base *activity.SpaceService

	Option *ActivityOptionService
}

// List returns a list of activities in your space.
//
// This method supports options returned by methods in "*Client.Space.Activity.Option",
// such as:
//   - WithActivityTypeIDs
//   - WithCount
//   - WithMaxID
//   - WithMinID
//   - WithOrder
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/get-recent-updates
func (s *SpaceActivityService) List(ctx context.Context, opts ...RequestOption) ([]*Activity, error) {
	v, err := s.base.List(ctx, toCoreOptions(opts)...)
	return activitiesFromModel(v), convertError(err)
}

// SpaceAttachmentService handles communication with the space attachment-related methods of the Backlog API.
type SpaceAttachmentService struct {
	base *attachment.SpaceService
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *SpaceAttachmentService) Upload(ctx context.Context, fileName string, r io.Reader) (*Attachment, error) {
	v, err := s.base.Upload(ctx, fileName, r)
	return attachmentFromModel(v), convertError(err)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newSpaceService(method *core.Method, option *core.OptionService) *SpaceService {
	return &SpaceService{
		base:       space.NewService(method),
		Activity:   newSpaceActivityService(method, option),
		Attachment: newSpaceAttachmentService(method),
	}
}

func newSpaceActivityService(method *core.Method, option *core.OptionService) *SpaceActivityService {
	return &SpaceActivityService{
		base:   activity.NewSpaceService(method),
		Option: newActivityOptionService(option),
	}
}

func newSpaceAttachmentService(method *core.Method) *SpaceAttachmentService {
	return &SpaceAttachmentService{
		base: attachment.NewSpaceService(method),
	}
}

// ──────────────────────────────────────────────────────────────
//  Model converters
// ──────────────────────────────────────────────────────────────

func spaceFromModel(m *model.Space) *Space {
	if m == nil {
		return nil
	}
	return &Space{
		SpaceKey:           m.SpaceKey,
		Name:               m.Name,
		OwnerID:            m.OwnerID,
		Lang:               m.Lang,
		Timezone:           m.Timezone,
		ReportSendTime:     m.ReportSendTime,
		TextFormattingRule: Format(m.TextFormattingRule),
		Created:            m.Created,
		Updated:            m.Updated,
	}
}

func diskUsageProjectFromModel(m *model.DiskUsageProject) *DiskUsageProject {
	if m == nil {
		return nil
	}
	return &DiskUsageProject{
		ProjectID:  m.ProjectID,
		Issue:      m.Issue,
		Wiki:       m.Wiki,
		File:       m.File,
		Subversion: m.Subversion,
		Git:        m.Git,
		GitLFS:     m.GitLFS,
	}
}

func diskUsageSpaceFromModel(m *model.DiskUsageSpace) *DiskUsageSpace {
	if m == nil {
		return nil
	}
	details := make([]*DiskUsageProject, len(m.Details))
	for i, v := range m.Details {
		details[i] = diskUsageProjectFromModel(v)
	}
	return &DiskUsageSpace{
		Capacity:   m.Capacity,
		Issue:      m.Issue,
		Wiki:       m.Wiki,
		File:       m.File,
		Subversion: m.Subversion,
		Git:        m.Git,
		GitLFS:     m.GitLFS,
		Details:    details,
	}
}

func spaceNotificationFromModel(m *model.SpaceNotification) *SpaceNotification {
	if m == nil {
		return nil
	}
	return &SpaceNotification{
		Content: m.Content,
		Updated: m.Updated,
	}
}
