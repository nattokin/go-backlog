package backlog

// SpaceService handles communication with the space-related methods of the Backlog API.
type SpaceService struct {
	method *method

	Activity   *SpaceActivityService
	Attachment *SpaceAttachmentService
}
