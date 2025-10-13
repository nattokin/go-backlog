package backlog

import (
	"net/url"
)

type (
	ExportFormType  = formType
	ExportQueryType = queryType
)

const (
	ExportQueryActivityTypeIDs = queryActivityTypeIDs
	ExportQueryAll             = queryAll
	ExportQueryArchived        = queryArchived
	ExportQueryCount           = queryCount
	ExportQueryKey             = queryKey
	ExportQueryKeyword         = queryKeyword
	ExportQueryOrder           = queryOrder
)

var ExportNewProjectOptionService = newProjectOptionService
var ExportNewActivityOptionService = newActivityOptionService
var ExportNewUserOptionService = newUserOptionService

const (
	ExportFormArchived                          = formArchived
	ExportFormChartEnabled                      = formChartEnabled
	ExportFormContent                           = formContent
	ExportFormKey                               = formKey
	ExportFormName                              = formName
	ExportFormMailAddress                       = formMailAddress
	ExportFormMailNotify                        = formMailNotify
	ExportFormPassword                          = formPassword
	ExportFormProjectLeaderCanEditProjectLeader = formProjectLeaderCanEditProjectLeader
	ExportFormRoleType                          = formRoleType
	ExportFormSubtaskingEnabled                 = formSubtaskingEnabled
	ExportFormTextFormattingRule                = formTextFormattingRule
)

type (
	ExportMethod  = method
	ExportWrapper = wrapper
)

var (
	ExportClientNewRequest = (*Client).newRequest
	ExportClientDo         = (*Client).do
)

var (
	ExportQueryOptionSet      = (*QueryOption).set
	ExportQueryOptionValidate = (*QueryOption).validate
	ExportFormOptionSet       = (*FormOption).set
	ExportFormOptionValidate  = (*FormOption).validate
)

var (
	ExportNewInternalClientError = newInternalClientError
	ExportCheckResponse          = checkResponse
)

func ExportNewQueryOption(queryType queryType, checkFunc optionCheckFunc, queryFunc queryOptionFunc) *QueryOption {
	return &QueryOption{t: queryType, checkFunc: checkFunc, setFunc: queryFunc}
}

func ExportNewFormOption(formType formType, checkFunc optionCheckFunc, fFunc formOptionFunc) *FormOption {
	return &FormOption{t: formType, checkFunc: checkFunc, setFunc: fFunc}
}

func (p *FormParams) ExportURLValues() *url.Values {
	return p.Values
}

// --- Test Helper Methods for Setting *method (Alphabetical Order) ---

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *CategoryService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *CustomFieldService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *IssueService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *IssueAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *PriorityService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *ProjectService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *ProjectActivityService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *ProjectUserService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *PullRequestService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *PullRequestAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *ResolutionService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *SpaceService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *SpaceActivityService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *SpaceAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *StatusService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *UserService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *UserActivityService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *VersionService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *WikiAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

// --- Sub-Service Helpers (Attachment, Activity, User) ---

// ExportNewIssueAttachmentService returns a test instance of IssueAttachmentService.
func ExportNewIssueAttachmentService() *IssueAttachmentService {
	return &IssueAttachmentService{
		method: newClientMethodMock(),
	}
}

// ExportNewPullRequestAttachmentService returns a test instance of PullRequestAttachmentService.
func ExportNewPullRequestAttachmentService() *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		method: newClientMethodMock(),
	}
}

// ExportNewSpaceAttachmentService returns a test instance of SpaceAttachmentService.
func ExportNewSpaceAttachmentService() *SpaceAttachmentService {
	return &SpaceAttachmentService{
		method: newClientMethodMock(),
	}
}

// ExportNewWikiAttachmentService returns a test instance of WikiAttachmentService.
func ExportNewWikiAttachmentService() *WikiAttachmentService {
	return &WikiAttachmentService{
		method: newClientMethodMock(),
	}
}

// ExportNewProjectUserService returns a test instance of ProjectUserService.
func ExportNewProjectUserService() *ProjectUserService {
	return &ProjectUserService{
		method: newClientMethodMock(),
	}
}

// ExportNewProjectActivityService returns a test instance of ProjectActivityService.
func ExportNewProjectActivityService() *ProjectActivityService {
	return &ProjectActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// ExportNewSpaceActivityService returns a test instance of SpaceActivityService.
func ExportNewSpaceActivityService() *SpaceActivityService {
	return &SpaceActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// ExportNewUserActivityService returns a test instance of UserActivityService.
func ExportNewUserActivityService() *UserActivityService {
	return &UserActivityService{
		method: newClientMethodMock(),
		Option: newActivityOptionService(),
	}
}

// --- Main Service Helpers ---

// ExportNewIssueService returns a test instance of IssueService.
func ExportNewIssueService() *IssueService {
	return &IssueService{
		method:     newClientMethodMock(),
		Attachment: ExportNewIssueAttachmentService(),
	}
}

// ExportNewPullRequestService returns a test instance of PullRequestService.
func ExportNewPullRequestService() *PullRequestService {
	return &PullRequestService{
		method:     newClientMethodMock(),
		Attachment: ExportNewPullRequestAttachmentService(),
	}
}

// ExportNewSpaceService returns a test instance of SpaceService.
func ExportNewSpaceService() *SpaceService {
	return &SpaceService{
		method:     newClientMethodMock(),
		Activity:   ExportNewSpaceActivityService(),
		Attachment: ExportNewSpaceAttachmentService(),
	}
}

// ExportNewUserService returns a test instance of UserService.
func ExportNewUserService() *UserService {
	return &UserService{
		method:   newClientMethodMock(),
		Activity: ExportNewUserActivityService(),
		Option:   ExportNewUserOptionService(),
	}
}
