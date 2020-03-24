package backlog

// ResolutionService has methods for Resolution.
type ResolutionService struct {
	clientMethod *clientMethod
}

func newResolutionService(cm *clientMethod) *ResolutionService {
	return &ResolutionService{
		clientMethod: cm,
	}
}
