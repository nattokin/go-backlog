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

var (
	ExportNewActivityService              = newActivityService
	ExportNewAttachmentService            = newAttachmentService
	ExportNewCategoryService              = newCategoryService
	ExportNewCustomFieldService           = newCustomFieldService
	ExportNewIssueService                 = newIssueService
	ExportNewIssueAttachmentService       = newIssueAttachmentService
	ExportNewPriorityService              = newPriorityService
	ExportNewProjectService               = newProjectService
	ExportNewProjectActivityService       = newProjectActivityService
	ExportNewProjectUserService           = newProjectUserService
	ExportNewPullRequestService           = newPullRequestService
	ExportNewPullRequestAttachmentService = newPullRequestAttachmentService
	ExportNewResolutionService            = newResolutionService
	ExportNewSpaceService                 = newSpaceService
	ExportNewSpaceActivityService         = newSpaceActivityService
	ExportNewStatusService                = newStatusService
	ExportNewUserService                  = newUserService
	ExportNewUserActivityService          = newUserActivityService
	ExportNewVersionService               = newVersionService
	ExportNewWikiService                  = newWikiService
	ExportNewWikiAttachmentService        = newWikiAttachmentService
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
