package backlog

import (
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/wiki"
)

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &core.OptionService{}

	// --- Initialize shared option services --------------------------------------
	activityOptionService := &ActivityOptionService{
		base: baseOptionService,
	}

	// --- Initialize IssueService -------------------------------------------------
	c.Issue = &IssueService{
		method: c.core.Method,
		Attachment: &IssueAttachmentService{
			method: c.core.Method,
		},
	}

	// --- Initialize ProjectService ----------------------------------------------
	c.Project = &ProjectService{
		method: c.core.Method,
		Activity: &ProjectActivityService{
			method: c.core.Method,
			Option: activityOptionService,
		},
		User: &ProjectUserService{
			method: c.core.Method,
		},
		Option: &ProjectOptionService{
			base: baseOptionService,
		},
	}

	// --- Initialize PullRequestService ------------------------------------------
	c.PullRequest = &PullRequestService{
		method: c.core.Method,
		Attachment: &PullRequestAttachmentService{
			method: c.core.Method,
		},
	}

	// --- Initialize SpaceService -------------------------------------------------
	c.Space = &SpaceService{
		method: c.core.Method,
		Activity: &SpaceActivityService{
			method: c.core.Method,
			Option: activityOptionService,
		},
		Attachment: &SpaceAttachmentService{
			method: c.core.Method,
		},
	}

	// --- Initialize UserService --------------------------------------------------
	c.User = &UserService{
		method: c.core.Method,
		Activity: &UserActivityService{
			method: c.core.Method,
			Option: activityOptionService,
		},
		Option: &UserOptionService{
			base: baseOptionService,
		},
	}

	// --- Initialize WikiService --------------------------------------------------
	c.Wiki = wiki.NewWikiService(c.core.Method, baseOptionService)
}
