package backlog

import (
	"net/http"
	"net/url"
)

type (
	ExportRole   = role
	ExportOrder  = order
	ExportFormat = format
)

type (
	ExportMethod        = method
	ExportRequestParams = requestParams
	ExportRequest       = request
	ExportResponse      = response
)

var (
	ExportClientNewReqest = (*Client).newReqest
	ExportClientDo        = (*Client).do
	ExportClientGet       = (*Client).get
	ExportClientPost      = (*Client).post
	ExportClientPatch     = (*Client).patch
	ExportClientDelete    = (*Client).delete
	ExportClientUploade   = (*Client).uploade
)

var (
	ExportNewClientError    = newClientError
	ExportNewRequestParams  = newRequestParams
	ExportNewResponse       = newResponse
	ExportCeckResponseError = checkResponseError
)

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

func (r *response) ExportGetHTTPResponse() *http.Response {
	return r.Response
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
