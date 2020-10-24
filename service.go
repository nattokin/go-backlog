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

// SpaceActivityService has methods for activitys in your space.
type SpaceActivityService struct {
	clientMethod *clientMethod
}

func newSpaceActivityService(cm *clientMethod) *SpaceActivityService {
	return &SpaceActivityService{
		clientMethod: cm,
	}
}

// AttachmentService hs methods for attachment.
type AttachmentService struct {
	clientMethod *clientMethod
}

func newAttachmentService(cm *clientMethod) *AttachmentService {
	return &AttachmentService{
		clientMethod: cm,
	}
}

// ActivityService has methods for Activitys.
type ActivityService struct {
	clientMethod *clientMethod
}

func newActivityService(cm *clientMethod) *ActivityService {
	return &ActivityService{
		clientMethod: cm,
	}
}

// ActivityOptionService has methods to make functional option for ActivityService.
type ActivityOptionService struct {
}

// CategoryService has methods for Category.
type CategoryService struct {
	clientMethod *clientMethod
}

func newCategoryService(cm *clientMethod) *CategoryService {
	return &CategoryService{
		clientMethod: cm,
	}
}

// CustomFieldService has methods for CustomField.
type CustomFieldService struct {
	clientMethod *clientMethod
}

func newCustomFieldService(cm *clientMethod) *CustomFieldService {
	return &CustomFieldService{
		clientMethod: cm,
	}
}

// VersionService has methods for Version.
type VersionService struct {
	clientMethod *clientMethod
}

func newVersionService(cm *clientMethod) *VersionService {
	return &VersionService{
		clientMethod: cm,
	}
}

// PriorityService has methods for Priority.
type PriorityService struct {
	clientMethod *clientMethod
}

func newPriorityService(cm *clientMethod) *PriorityService {
	return &PriorityService{
		clientMethod: cm,
	}
}

// ResolutionService has methods for Resolution.
type ResolutionService struct {
	clientMethod *clientMethod
}

func newResolutionService(cm *clientMethod) *ResolutionService {
	return &ResolutionService{
		clientMethod: cm,
	}
}

// IssueService has methods for Issue.
type IssueService struct {
	clientMethod *clientMethod

	Attachment *IssueAttachmentService
}

func newIssueService(cm *clientMethod) *IssueService {
	return &IssueService{
		clientMethod: cm,
		Attachment:   newIssueAttachmentService(cm),
	}
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	*AttachmentService
}

func newIssueAttachmentService(cm *clientMethod) *IssueAttachmentService {
	return &IssueAttachmentService{
		AttachmentService: newAttachmentService(cm),
	}
}

// StatusService has methods for Status.
type StatusService struct {
	clientMethod *clientMethod
}

func newStatusService(cm *clientMethod) *StatusService {
	return &StatusService{
		clientMethod: cm,
	}
}

// UserService has methods for user
type UserService struct {
	clientMethod *clientMethod

	Activity *UserActivityService
	Option   *UserOptionService
}

func newUserService(cm *clientMethod) *UserService {
	return &UserService{
		clientMethod: cm,
		Activity:     newUserActivityService(cm),
		Option:       &UserOptionService{},
	}
}

// UserActivityService has methods for user activitys.
type UserActivityService struct {
	clientMethod *clientMethod
}

func newUserActivityService(cm *clientMethod) *UserActivityService {
	return &UserActivityService{
		clientMethod: cm,
	}
}

// UserOptionService has methods to make functional option for UserService.
type UserOptionService struct {
}

// ProjectService has methods for Project.
type ProjectService struct {
	clientMethod *clientMethod

	Activity *ProjectActivityService
	User     *ProjectUserService
	Option   *ProjectOptionService
}

func newProjectService(cm *clientMethod) *ProjectService {
	return &ProjectService{
		clientMethod: cm,
		Activity:     newProjectActivityService(cm),
		User:         newProjectUserService(cm),
		Option:       &ProjectOptionService{},
	}
}

// ProjectActivityService has methods for activitys of the project.
type ProjectActivityService struct {
	clientMethod *clientMethod
}

func newProjectActivityService(cm *clientMethod) *ProjectActivityService {
	return &ProjectActivityService{
		clientMethod: cm,
	}
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	clientMethod *clientMethod
}

func newProjectUserService(cm *clientMethod) *ProjectUserService {
	return &ProjectUserService{
		clientMethod: cm,
	}
}

// ProjectOptionService has methods to make functional option for ProjectService.
type ProjectOptionService struct {
}

// WikiService has methods for Wiki.
type WikiService struct {
	clientMethod *clientMethod

	Attachment *WikiAttachmentService
	option     *wikiOptionService
}

func newWikiService(cm *clientMethod) *WikiService {
	return &WikiService{
		clientMethod: cm,
		Attachment:   newWikiAttachmentService(cm),
		option:       &wikiOptionService{},
	}
}

// WikiAttachmentService hs methods for attachment file of wiki.
type WikiAttachmentService struct {
	*AttachmentService
}

func newWikiAttachmentService(cm *clientMethod) *WikiAttachmentService {
	return &WikiAttachmentService{
		AttachmentService: newAttachmentService(cm),
	}

}

// wikiOptionService has methods to make functional option for WikiService.
type wikiOptionService struct {
}

// PullRequestService has methods for Issue.
type PullRequestService struct {
	clientMethod *clientMethod

	Attachment *PullRequestAttachmentService
}

func newPullRequestService(cm *clientMethod) *PullRequestService {
	return &PullRequestService{
		clientMethod: cm,
		Attachment:   newPullRequestAttachmentService(cm),
	}
}

// PullRequestAttachmentService hs methods for attachment file of pull request.
type PullRequestAttachmentService struct {
	*AttachmentService
}

func newPullRequestAttachmentService(cm *clientMethod) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		AttachmentService: newAttachmentService(cm),
	}
}
