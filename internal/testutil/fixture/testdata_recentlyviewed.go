package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type recentlyViewedFixtures struct {
	IssueListJSON   string
	IssueList       []*backlog.Issue
	IssueSingleJSON string
	IssueSingle     *backlog.Issue

	ProjectListJSON   string
	ProjectList       []*backlog.Project
	ProjectSingleJSON string
	ProjectSingle     *backlog.Project

	WikiListJSON   string
	WikiList       []*backlog.Wiki
	WikiSingleJSON string
	WikiSingle     *backlog.Wiki
}

// RecentlyViewed provides test fixtures for RecentlyViewed-related tests.
var RecentlyViewed = recentlyViewedFixtures{
	IssueListJSON: `
[
    {
        "id": 1,
        "projectId": 1,
        "issueKey": "TEST-1",
        "keyId": 1,
        "summary": "first issue",
        "status": {"id": 1, "name": "Open"}
    },
    {
        "id": 2,
        "projectId": 1,
        "issueKey": "TEST-2",
        "keyId": 2,
        "summary": "second issue",
        "status": {"id": 1, "name": "Open"}
    }
]
`,
	IssueSingleJSON: `
{
    "id": 1,
    "projectId": 1,
    "issueKey": "TEST-1",
    "keyId": 1,
    "summary": "first issue",
    "status": {"id": 1, "name": "Open"}
}
`,
	ProjectListJSON: `
[
    {
        "id": 1,
        "projectKey": "TEST",
        "name": "Test Project",
        "chartEnabled": false,
        "subtaskingEnabled": false,
        "projectLeaderCanEditProjectLeader": false,
        "textFormattingRule": "markdown"
    },
    {
        "id": 2,
        "projectKey": "DEMO",
        "name": "Demo Project",
        "chartEnabled": true,
        "subtaskingEnabled": true,
        "projectLeaderCanEditProjectLeader": false,
        "textFormattingRule": "backlog"
    }
]
`,
	ProjectSingleJSON: `
{
    "id": 1,
    "projectKey": "TEST",
    "name": "Test Project",
    "chartEnabled": false,
    "subtaskingEnabled": false,
    "projectLeaderCanEditProjectLeader": false,
    "textFormattingRule": "markdown"
}
`,
	WikiListJSON: `
[
    {
        "id": 10,
        "projectId": 1,
        "name": "Home",
        "content": "Welcome!",
        "tags": [],
        "attachments": [],
        "sharedFiles": [],
        "stars": []
    },
    {
        "id": 11,
        "projectId": 1,
        "name": "Guide",
        "content": "How to use.",
        "tags": [],
        "attachments": [],
        "sharedFiles": [],
        "stars": []
    }
]
`,
	WikiSingleJSON: `
{
    "id": 10,
    "projectId": 1,
    "name": "Home",
    "content": "Welcome!",
    "tags": [],
    "attachments": [],
    "sharedFiles": [],
    "stars": []
}
`,
}
