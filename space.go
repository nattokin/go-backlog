package backlog

// SpaceService has methods for Space.
type SpaceService struct {
	method *method

	Activity *SpaceActivityService
}

// SpaceActivityService has methods for activitys in your space.
type SpaceActivityService struct {
	method *method
}
