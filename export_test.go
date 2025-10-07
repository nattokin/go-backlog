package backlog

import (
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
	ExportCreateFormFile         = createFormFile
	ExportCopy                   = copy
)

var ExportQueryParamsWithOptions = (*QueryParams).withOptions

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

func (c *Client) ExportSetWrapper(wrapper *wrapper) {
	c.wrapper = wrapper
}

func (p *FormParams) ExportURLValues() *url.Values {
	return p.Values
}

func (s *SpaceAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *CategoryService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *CustomFieldService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *IssueService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *IssueAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *PriorityService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *ProjectService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *ProjectActivityService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *ProjectUserService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *PullRequestService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *PullRequestAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *ResolutionService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *SpaceService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *SpaceActivityService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *StatusService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *UserService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *UserActivityService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *VersionService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *WikiService) ExportSetMethod(m *method) {
	s.method = m
}

func (s *WikiAttachmentService) ExportSetMethod(m *method) {
	s.method = m
}
