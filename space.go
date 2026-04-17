package backlog

import (
	"context"
	"io"
	"time"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
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
