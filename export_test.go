package backlog

import (
	"net/http"
	"net/url"
)

type (
	ExportClientMethod  = clientMethod
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
	ExportNewActivityService    = newActivityService
	ExportNewCategoryService    = newCategoryService
	ExportNewCustomFieldService = newCustomFieldService
	ExportNewIssueService       = newIssueService
	ExportNewPriorityService    = newPriorityService
	ExportNewProjectService     = newProjectService
	ExportNewProjectUserService = newProjectUserService
	ExportNewPullRequestService = newPullRequestService
	ExportNewResolutionService  = newResolutionService
	ExportNewSpaceService       = newSpaceService
	ExportNewStatusService      = newStatusService
	ExportNewUserService        = newUserService
	ExportNewVersionService     = newVersionService
	ExportNewWikiService        = newWikiService
)

var (
	ExportNewBaseActivityService = newBaseActivityService
	ExportNewBaseUserService     = newBaseUserService
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
