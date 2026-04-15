package backlog

import (
	"context"
	"io"

	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/model"
	"github.com/nattokin/go-backlog/internal/space"
)

// ──────────────────────────────────────────────────────────────
//  SpaceService
// ──────────────────────────────────────────────────────────────

// SpaceService handles communication with the space-related methods of the Backlog API.
type SpaceService struct {
	base *space.SpaceService

	Activity   *SpaceActivityService
	Attachment *SpaceAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  SpaceActivityService
// ──────────────────────────────────────────────────────────────

// SpaceActivityService handles communication with the space activities-related methods of the Backlog API.
type SpaceActivityService struct {
	base *activity.SpaceActivityService

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
func (s *SpaceActivityService) List(ctx context.Context, opts ...core.RequestOption) ([]*model.Activity, error) {
	return s.base.List(ctx, opts...)
}

// ──────────────────────────────────────────────────────────────
//  SpaceAttachmentService
// ──────────────────────────────────────────────────────────────

// SpaceAttachmentService handles communication with the space attachment-related methods of the Backlog API.
type SpaceAttachmentService struct {
	base *attachment.SpaceAttachmentService
}

// Upload uploads any file to the space.
//
// The file name must not be empty.
//
// Backlog API docs: https://developer.nulab.com/docs/backlog/api/2/post-attachment-file
func (s *SpaceAttachmentService) Upload(ctx context.Context, fileName string, r io.Reader) (*model.Attachment, error) {
	return s.base.Upload(ctx, fileName, r)
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

func newSpaceService(method *core.Method, option *core.OptionService) *SpaceService {
	return &SpaceService{
		base:       space.NewSpaceService(method, option),
		Activity:   newSpaceActivityService(method, option),
		Attachment: newSpaceAttachmentService(method),
	}
}

func newSpaceActivityService(method *core.Method, option *core.OptionService) *SpaceActivityService {
	return &SpaceActivityService{
		base:   activity.NewSpaceActivityService(method, option),
		Option: &ActivityOptionService{},
	}
}

func newSpaceAttachmentService(method *core.Method) *SpaceAttachmentService {
	return &SpaceAttachmentService{
		base: attachment.NewSpaceAttachmentService(method),
	}
}
