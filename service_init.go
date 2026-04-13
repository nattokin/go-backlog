package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/pullrequest"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &core.OptionService{}

	// --- Initialize shared option services --------------------------------------
	activityOptionService := activity.NewActivityOptionService(baseOptionService)

	// --- Initialize IssueService -------------------------------------------------
	c.Issue = issue.NewIssueService(c.core.Method, baseOptionService)

	// --- Initialize ProjectService ----------------------------------------------
	c.Project = &ProjectService{
		method: c.core.Method,
		Activity: &ProjectActivityService{
			method: c.core.Method,
			Option: activityOptionService,
		},
		User: user.NewProjectUserService(c.core.Method, baseOptionService),
		Option: &ProjectOptionService{
			base: baseOptionService,
		},
	}

	// --- Initialize PullRequestService ------------------------------------------
	c.PullRequest = pullrequest.NewPullRequestService(c.core.Method)

	// --- Initialize SpaceService -------------------------------------------------
	c.Space = space.NewSpaceService(c.core.Method, baseOptionService)

	// --- Initialize UserService --------------------------------------------------
	c.User = user.NewUserService(c.core.Method, baseOptionService)

	// --- Initialize WikiService --------------------------------------------------
	c.Wiki = wiki.NewWikiService(c.core.Method, baseOptionService)
}
