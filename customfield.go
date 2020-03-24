package backlog

// CustomFieldService has methods for CustomField.
type CustomFieldService struct {
	clientMethod *clientMethod
}

func newCustomFieldService(cm *clientMethod) *CustomFieldService {
	return &CustomFieldService{
		clientMethod: cm,
	}
}
