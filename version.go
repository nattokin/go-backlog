package backlog

// VersionService has methods for Version.
type VersionService struct {
	clientMethod *clientMethod
}

func newVersionService(cm *clientMethod) *VersionService {
	return &VersionService{
		clientMethod: cm,
	}
}
