package backlog

// SpaceService has methods for Space.
type SpaceService struct {
	clientMethod *clientMethod

	Activity *SpaceActivityService
}

func newSpaceService(cm *clientMethod) *SpaceService {
	return &SpaceService{
		clientMethod: cm,
		Activity:     newSpaceActivityService(cm),
	}
}
