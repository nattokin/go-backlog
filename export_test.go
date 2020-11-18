package backlog

import (
	"net/http"
	"net/url"
)

type (
	ExportRole       = role
	ExportOrder      = order
	ExportFormat     = format
	ExportOptionType = optionType
)

const (
	ExportOptionActivityTypeIDs                   = optionActivityTypeIDs
	ExportOptionAll                               = optionAll
	ExportOptionArchived                          = optionArchived
	ExportOptionChartEnabled                      = optionChartEnabled
	ExportOptionContent                           = optionContent
	ExportOptionCount                             = optionCount
	ExportOptionKey                               = optionKey
	ExportOptionKeyword                           = optionKeyword
	ExportOptionName                              = optionName
	ExportOptionMailAddress                       = optionMailAddress
	ExportOptionMailNotify                        = optionMailNotify
	ExportOptionMaxID                             = optionMaxID
	ExportOptionMinID                             = optionMinID
	ExportOptionOrder                             = optionOrder
	ExportOptionPassword                          = optionPassword
	ExportOptionProjectLeaderCanEditProjectLeader = optionProjectLeaderCanEditProjectLeader
	ExportOptionRoleType                          = optionRoleType
	ExportOptionSubtaskingEnabled                 = optionSubtaskingEnabled
	ExportOptionTextFormattingRule                = optionTextFormattingRule
)

type (
	ExportMethod        = method
	ExportRequestParams = requestParams
	ExportWrapper       = wrapper
)

var (
	ExportClientNewReqest = (*Client).newReqest
	ExportClientDo        = (*Client).do
	ExportClientGet       = (*Client).get
	ExportClientPost      = (*Client).post
	ExportClientPatch     = (*Client).patch
	ExportClientDelete    = (*Client).delete
	ExportClientUpload    = (*Client).upload
)

var (
	ExportActivityOptionSet = (*ActivityOption).set
	ExportProjectOptionSet  = (*ProjectOption).set
	ExportUserOptionSet     = (*UserOption).set
	ExportWikiOptionSet     = (*WikiOption).set
)

var (
	ExportNewClientError   = newClientError
	ExportNewRequestParams = newRequestParams
	ExportCeckResponse     = checkResponse
	ExportCreateFormFile   = createFormFile
	ExportCopy             = copy
)

func ExportNewWikiOption(optionType optionType, optionFunc optionFunc) *WikiOption {
	return &WikiOption{&option{optionType, optionFunc}}
}

func ExportNewActivityOption(optionType optionType, optionFunc optionFunc) *ActivityOption {
	return &ActivityOption{&option{optionType, optionFunc}}
}

func ExportNewProjectOption(optionType optionType, optionFunc optionFunc) *ProjectOption {
	return &ProjectOption{&option{optionType, optionFunc}}
}

func ExportNewUserOption(optionType optionType, optionFunc optionFunc) *UserOption {
	return &UserOption{&option{optionType, optionFunc}}
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

func (p *requestParams) ExportURLValues() *url.Values {
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
