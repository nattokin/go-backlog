package backlog

// SpaceService has methods for Space.
type SpaceService struct {
	method *method

	Activity *SpaceActivityService
}

func newSpaceService(m *method) *SpaceService {
	return &SpaceService{
		method:   m,
		Activity: newSpaceActivityService(m),
	}
}

// SpaceActivityService has methods for activitys in your space.
type SpaceActivityService struct {
	method *method
}

func newSpaceActivityService(m *method) *SpaceActivityService {
	return &SpaceActivityService{
		method: m,
	}
}

// AttachmentService hs methods for attachment.
type AttachmentService struct {
	method *method
}

func newAttachmentService(m *method) *AttachmentService {
	return &AttachmentService{
		method: m,
	}
}

// ActivityService has methods for Activitys.
type ActivityService struct {
	method *method
	Option *ActivityOptionService
}

func newActivityService(m *method) *ActivityService {
	return &ActivityService{
		method: m,
		Option: &ActivityOptionService{},
	}
}

// ActivityOptionService has methods to make functional option for ActivityService.
type ActivityOptionService struct {
}

// CategoryService has methods for Category.
type CategoryService struct {
	method *method
}

func newCategoryService(m *method) *CategoryService {
	return &CategoryService{
		method: m,
	}
}

// CustomFieldService has methods for CustomField.
type CustomFieldService struct {
	method *method
}

func newCustomFieldService(m *method) *CustomFieldService {
	return &CustomFieldService{
		method: m,
	}
}

// VersionService has methods for Version.
type VersionService struct {
	method *method
}

func newVersionService(m *method) *VersionService {
	return &VersionService{
		method: m,
	}
}

// PriorityService has methods for Priority.
type PriorityService struct {
	method *method
}

func newPriorityService(m *method) *PriorityService {
	return &PriorityService{
		method: m,
	}
}

// ResolutionService has methods for Resolution.
type ResolutionService struct {
	method *method
}

func newResolutionService(m *method) *ResolutionService {
	return &ResolutionService{
		method: m,
	}
}

// IssueService has methods for Issue.
type IssueService struct {
	method *method

	Attachment *IssueAttachmentService
}

func newIssueService(m *method) *IssueService {
	return &IssueService{
		method:     m,
		Attachment: newIssueAttachmentService(m),
	}
}

// IssueAttachmentService hs methods for attachment file of issue.
type IssueAttachmentService struct {
	*AttachmentService
}

func newIssueAttachmentService(m *method) *IssueAttachmentService {
	return &IssueAttachmentService{
		AttachmentService: newAttachmentService(m),
	}
}

// StatusService has methods for Status.
type StatusService struct {
	method *method
}

func newStatusService(m *method) *StatusService {
	return &StatusService{
		method: m,
	}
}

// UserService has methods for user
type UserService struct {
	method *method

	Activity *UserActivityService
	Option   *UserOptionService
}

func newUserService(m *method) *UserService {
	return &UserService{
		method:   m,
		Activity: newUserActivityService(m),
		Option:   &UserOptionService{},
	}
}

// UserActivityService has methods for user activitys.
type UserActivityService struct {
	method *method
}

func newUserActivityService(m *method) *UserActivityService {
	return &UserActivityService{
		method: m,
	}
}

// UserOptionService has methods to make functional option for UserService.
type UserOptionService struct {
}

// ProjectService has methods for Project.
type ProjectService struct {
	method *method

	Activity *ProjectActivityService
	User     *ProjectUserService
	Option   *ProjectOptionService
}

func newProjectService(m *method) *ProjectService {
	return &ProjectService{
		method:   m,
		Activity: newProjectActivityService(m),
		User:     newProjectUserService(m),
		Option:   &ProjectOptionService{},
	}
}

// ProjectActivityService has methods for activitys of the project.
type ProjectActivityService struct {
	method *method
}

func newProjectActivityService(m *method) *ProjectActivityService {
	return &ProjectActivityService{
		method: m,
	}
}

// ProjectUserService has methods for user of project.
type ProjectUserService struct {
	method *method
}

func newProjectUserService(m *method) *ProjectUserService {
	return &ProjectUserService{
		method: m,
	}
}

// ProjectOptionService has methods to make functional option for ProjectService.
type ProjectOptionService struct {
}

// WikiService has methods for Wiki.
type WikiService struct {
	method *method

	Attachment *WikiAttachmentService
	Option     *WikiOptionService
}

func newWikiService(m *method) *WikiService {
	return &WikiService{
		method:     m,
		Attachment: newWikiAttachmentService(m),
		Option:     &WikiOptionService{},
	}
}

// WikiAttachmentService hs methods for attachment file of wiki.
type WikiAttachmentService struct {
	*AttachmentService
}

func newWikiAttachmentService(m *method) *WikiAttachmentService {
	return &WikiAttachmentService{
		AttachmentService: newAttachmentService(m),
	}

}

// WikiOptionService has methods to make functional option for WikiService.
type WikiOptionService struct {
}

// PullRequestService has methods for Issue.
type PullRequestService struct {
	method *method

	Attachment *PullRequestAttachmentService
}

func newPullRequestService(m *method) *PullRequestService {
	return &PullRequestService{
		method:     m,
		Attachment: newPullRequestAttachmentService(m),
	}
}

// PullRequestAttachmentService hs methods for attachment file of pull request.
type PullRequestAttachmentService struct {
	*AttachmentService
}

func newPullRequestAttachmentService(m *method) *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		AttachmentService: newAttachmentService(m),
	}
}
