// Package backlog is a Go client library for the [Backlog API] by Nulab.
//
// # Creating a client
//
// [NewClient] creates a [Client] that authenticates with an API key:
//
//	c, err := backlog.NewClient("https://example.backlog.com", apiKey)
//	if err != nil {
//		return err
//	}
//
// To authenticate with an OAuth 2.0 access token, use [NewOAuthClient] instead.
// The client does not obtain or refresh access tokens; create a new client
// when the token expires.
//
// Requests are sent with [net/http.DefaultClient] unless another HTTP client is
// given with [WithDoer].
//
// # Calling the API
//
// A [Client] has a service for each group of API endpoints, such as c.Issue,
// c.Project and c.Wiki, and some services have sub-services, such as
// c.Issue.Comment. Every method takes a [context.Context]:
//
//	wikis, err := c.Wiki.List(ctx, "MYPROJECT")
//
// # Options
//
// Optional parameters are passed as [RequestOption] values. A service that
// accepts options has an Option field whose methods build them:
//
//	wikis, err := c.Wiki.List(ctx, "MYPROJECT",
//		c.Wiki.Option.WithKeyword("design"),
//	)
//
// The doc comment of each method lists the options it accepts. Passing any
// other option returns an [*InvalidOptionKeyError] without sending a request.
//
// # Pagination
//
// [IssueService.All] and [PullRequestService.All] return an iterator that
// fetches the results page by page:
//
//	issues, err := c.Issue.All(ctx, 100, c.Issue.Option.WithProjectIDs([]int{projectID}))
//	if err != nil {
//		return err
//	}
//	for issue, err := range issues {
//		if err != nil {
//			return err
//		}
//		fmt.Println(issue.IssueKey, issue.Summary)
//	}
//
// # Errors
//
// Errors are returned as the following types. Use [errors.As] to tell them
// apart:
//
//   - [*APIResponseError]: the API returned an error response.
//     [APIResponseError.RateLimitStatus] reports the rate limit status.
//   - [*OAuthError]: the API rejected an OAuth 2.0 access token. It wraps an
//     [*APIResponseError], so check for it first.
//   - [*ValidationError]: an argument or option is invalid. No request is sent.
//   - [*InvalidOptionKeyError]: the method does not accept the option. No
//     request is sent.
//   - [*InvalidOptionError]: the option itself is invalid, such as nil.
//   - [*InternalClientError]: the client is misconfigured, such as with an
//     empty API key or a malformed base URL.
//
// # Dates and times
//
// Date and time fields are [Timestamp] values. Fields that hold a date without
// a time, such as the StartDate and DueDate fields of [Issue], hold midnight
// of the date in the space's time zone, expressed in UTC.
//
// # More examples
//
// The examples directory of the repository contains complete programs:
// https://github.com/nattokin/go-backlog/tree/main/examples
//
// [Backlog API]: https://developer.nulab.com/docs/backlog/
package backlog
