package backlog

// StatusService has methods for Status.
type StatusService struct {
	clientMethod *clientMethod
}

func newStatusService(cm *clientMethod) *StatusService {
	return &StatusService{
		clientMethod: cm,
	}
}
