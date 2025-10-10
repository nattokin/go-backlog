package backlog

import (
	"errors"
	"net/http"
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
var ExportNewWikiOptionService = newWikiOptionService
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
	ExportMethod        = method
	ExportRequestParams = FormParams
	ExportWrapper       = wrapper
)

var (
	ExportClientNewRequest = (*Client).newRequest
	ExportClientDo         = (*Client).do
	ExportClientGet        = (*Client).get
	ExportClientPost       = (*Client).post
	ExportClientPatch      = (*Client).patch
	ExportClientDelete     = (*Client).delete
	ExportClientUpload     = (*Client).upload
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

func ExportNewQueryOption(qType queryType, qFunc queryOptionFunc) *QueryOption {
	return &QueryOption{qType, qFunc}
}

func ExportNewFormOption(formType formType, fFunc formOptionFunc) *FormOption {
	return &FormOption{formType, fFunc}
}

func (c *Client) ExportURL() *url.URL {
	return c.url
}

func (c *Client) ExportSetURL(url *url.URL) {
	c.url = url
}

func (c *Client) ExportHTTPClient() *http.Client {
	return c.httpClient
}

func (c *Client) ExportSetHTTPClient(httpClient *http.Client) {
	c.httpClient = httpClient
}

func (c *Client) ExportToken() string {
	return c.token
}

func (c *Client) ExportSetToken(token string) {
	c.token = token
}

func (c *Client) ExportSetWrapper(wrapper wrapper) {
	c.wrapper = wrapper
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
func (s *WikiService) ExportSetMethod(m *method) {
	s.method = m
}

// ExportSetMethod sets the internal 'method' field for testing purposes.
func (s *WikiAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

// newTestClientMethod creates a mock 'method' that returns an error by default for API calls.
func newTestClientMethod() *method {
	return &method{
		Get: func(spath string, query *QueryParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Post: func(spath string, form *ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Patch: func(spath string, form *ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
		Delete: func(spath string, form *ExportRequestParams) (*http.Response, error) {
			return nil, errors.New("default mock not implemented")
		},
	}
}

// --- Sub-Service Helpers (Attachment, Activity, User) ---

// ExportNewIssueAttachmentService returns a test instance of IssueAttachmentService.
func ExportNewIssueAttachmentService() *IssueAttachmentService {
	return &IssueAttachmentService{
		method: newTestClientMethod(),
	}
}

// ExportNewPullRequestAttachmentService returns a test instance of PullRequestAttachmentService.
func ExportNewPullRequestAttachmentService() *PullRequestAttachmentService {
	return &PullRequestAttachmentService{
		method: newTestClientMethod(),
	}
}

// ExportNewSpaceAttachmentService returns a test instance of SpaceAttachmentService.
func ExportNewSpaceAttachmentService() *SpaceAttachmentService {
	return &SpaceAttachmentService{
		method: newTestClientMethod(),
	}
}

// ExportNewWikiAttachmentService returns a test instance of WikiAttachmentService.
func ExportNewWikiAttachmentService() *WikiAttachmentService {
	return &WikiAttachmentService{
		method: newTestClientMethod(),
	}
}

// ExportNewProjectUserService returns a test instance of ProjectUserService.
func ExportNewProjectUserService() *ProjectUserService {
	return &ProjectUserService{
		method: newTestClientMethod(),
	}
}

// ExportNewProjectActivityService returns a test instance of ProjectActivityService.
func ExportNewProjectActivityService() *ProjectActivityService {
	return &ProjectActivityService{
		method: newTestClientMethod(),
		Option: newActivityOptionService(),
	}
}

// ExportNewSpaceActivityService returns a test instance of SpaceActivityService.
func ExportNewSpaceActivityService() *SpaceActivityService {
	return &SpaceActivityService{
		method: newTestClientMethod(),
		Option: newActivityOptionService(),
	}
}

// ExportNewUserActivityService returns a test instance of UserActivityService.
func ExportNewUserActivityService() *UserActivityService {
	return &UserActivityService{
		method: newTestClientMethod(),
		Option: newActivityOptionService(),
	}
}

// --- Main Service Helpers ---

// ExportNewIssueService returns a test instance of IssueService.
func ExportNewIssueService() *IssueService {
	return &IssueService{
		method:     newTestClientMethod(),
		Attachment: ExportNewIssueAttachmentService(),
	}
}

// ExportNewProjectService returns a test instance of ProjectService.
func ExportNewProjectService() *ProjectService {
	return &ProjectService{
		method:   newTestClientMethod(),
		Activity: ExportNewProjectActivityService(),
		User:     ExportNewProjectUserService(),
		Option:   newProjectOptionService(),
	}
}

// ExportNewPullRequestService returns a test instance of PullRequestService.
func ExportNewPullRequestService() *PullRequestService {
	return &PullRequestService{
		method:     newTestClientMethod(),
		Attachment: ExportNewPullRequestAttachmentService(),
	}
}

// ExportNewSpaceService returns a test instance of SpaceService.
func ExportNewSpaceService() *SpaceService {
	return &SpaceService{
		method:     newTestClientMethod(),
		Activity:   ExportNewSpaceActivityService(),
		Attachment: ExportNewSpaceAttachmentService(),
	}
}

// ExportNewUserService returns a test instance of UserService.
func ExportNewUserService() *UserService {
	return &UserService{
		method:   newTestClientMethod(),
		Activity: ExportNewUserActivityService(),
		Option:   ExportNewUserOptionService(),
	}
}

// ExportNewWikiService returns a test instance of WikiService.
func ExportNewWikiService() *WikiService {
	return &WikiService{
		method:     newTestClientMethod(),
		Attachment: ExportNewWikiAttachmentService(),
		Option:     ExportNewWikiOptionService(),
	}
}
