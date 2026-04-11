package backlog

import "github.com/nattokin/go-backlog/internal/core"

// SpaceService handles communication with the space-related methods of the Backlog API.
type SpaceService struct {
	method *core.Method

	Activity   *SpaceActivityService
	Attachment *SpaceAttachmentService
}
