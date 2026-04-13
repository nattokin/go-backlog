package space

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
)

// SpaceService handles communication with the space-related methods of the Backlog API.
type SpaceService struct {
	method *core.Method

	Activity   *activity.SpaceActivityService
	Attachment *attachment.SpaceAttachmentService
}

// ──────────────────────────────────────────────────────────────
//  Constructors
// ──────────────────────────────────────────────────────────────

// NewWikiService returns a new WikiService.
func NewSpaceService(method *core.Method, option *core.OptionService) *SpaceService {
	return &SpaceService{
		method:     method,
		Activity:   activity.NewSpaceActivityService(method, option),
		Attachment: attachment.NewSpaceAttachmentService(method),
	}
}
