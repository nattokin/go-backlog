package backlog

// SpaceService has methods for Space.
type SpaceService struct {
	method *method

	Activity   *SpaceActivityService
	Attachment *SpaceAttachmentService
}
