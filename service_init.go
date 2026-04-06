package backlog

// ──────────────────────────────────────────────────────────────
//  Service initialization
// ──────────────────────────────────────────────────────────────

func initServices(c *Client) {
	baseOptionService := &OptionService{}

	// --- Initialize shared option services --------------------------------------
	activityOptionService := &ActivityOptionService{
		base: baseOptionService,
	}

	// --- Initialize IssueService -------------------------------------------------
	c.Issue = &IssueService{
		method: c.method,
		Attachment: &IssueAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize ProjectService ----------------------------------------------
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
			base: baseOptionService,
		},
	}

	// --- Initialize PullRequestService ------------------------------------------
	c.PullRequest = &PullRequestService{
		method: c.method,
		Attachment: &PullRequestAttachmentService{
			method: c.method,
		},
	}

	// --- Initialize SpaceService -------------------------------------------------
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
	c.User = &UserService{
		method: c.method,
		Activity: &UserActivityService{
			method: c.method,
			Option: activityOptionService,
		},
		Option: &UserOptionService{
			base: baseOptionService,
		},
	}

	// --- Initialize WikiService --------------------------------------------------
	c.Wiki = &WikiService{
		method: c.method,
		Attachment: &WikiAttachmentService{
			method: c.method,
		},
		Option: &WikiOptionService{
			base: baseOptionService,
		},
	}
}
