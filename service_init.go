package backlog

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	// --- Initialize shared option services --------------------------------------
	// Option services provide reusable form and query parameter builders
	// used across multiple Backlog API services.
	optionRegistry := &optionRegistry{
		query: &QueryOptionService{},
		form:  &FormOptionService{},
	}

	// ActivityOptionService wraps shared optionRegistry to be reused
	// by activity-related services such as ProjectActivityService or SpaceActivityService.
	activityOptionService := &ActivityOptionService{
		registry: optionRegistry,
	}

	// --- Initialize IssueService -------------------------------------------------
	// Provides methods for issue management and file attachment operations.
	c.Issue = &IssueService{
		method: c.method,
		Attachment: &IssueAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize ProjectService ----------------------------------------------
	// Includes sub-services for project activities, users, and project options.
	c.Project = &ProjectService{
		method: c.method,
		Activity: &ProjectActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		User: &ProjectUserService{
			method: c.method,
		},
		Option: &ProjectOptionService{
			registry: optionRegistry,
		},
	}

	// --- Initialize PullRequestService ------------------------------------------
	// Handles pull request operations and related file attachments.
	c.PullRequest = &PullRequestService{
		method: c.method,
		Attachment: &PullRequestAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize SpaceService -------------------------------------------------
	// Provides access to space-level APIs including activities and attachments.
	c.Space = &SpaceService{
		method: c.method,
		Activity: &SpaceActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		Attachment: &SpaceAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize UserService --------------------------------------------------
	// Provides APIs related to user activities and user option settings.
	c.User = &UserService{
		method: c.method,
		Activity: &UserActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		Option: &UserOptionService{
			registry: optionRegistry,
		},
	}

	// --- Initialize WikiService --------------------------------------------------
	// Provides wiki page APIs, including file attachments and option configurations.
	c.Wiki = &WikiService{
		method: c.method,
		Attachment: &WikiAttachmentService{
			method: c.method,
		},
		Option: &WikiOptionService{
			registry: optionRegistry,
		},
	}
}
