package backlog

// PriorityService has methods for Priority.
type PriorityService struct {
	clientMethod *clientMethod
}

func newPriorityService(cm *clientMethod) *PriorityService {
	return &PriorityService{
		clientMethod: cm,
	}
}
