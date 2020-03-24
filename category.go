package backlog

// CategoryService has methods for Category.
type CategoryService struct {
	clientMethod *clientMethod
}

func newCategoryService(cm *clientMethod) *CategoryService {
	return &CategoryService{
		clientMethod: cm,
	}
}
