package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type recentlyViewedFixtures struct {
	IssueListJSON   string
	IssueList       []*backlog.RecentlyViewedIssue
	IssueSingleJSON string
	IssueSingle     *backlog.Issue

	ProjectListJSON   string
	ProjectList       []*backlog.RecentlyViewedProject
	ProjectSingleJSON string
	ProjectSingle     *backlog.Project

	WikiListJSON   string
	WikiList       []*backlog.RecentlyViewedWiki
	WikiSingleJSON string
	WikiSingle     *backlog.Wiki
}

// RecentlyViewed provides test fixtures for RecentlyViewed-related tests.
var RecentlyViewed = recentlyViewedFixtures{
	IssueListJSON: `
[
    {
        "issue": {
            "id": 1,
            "projectId": 1,
            "issueKey": "TEST-1",
            "keyId": 1,
            "issueType": {
                "id": 2,
                "projectId": 1,
                "name": "Task",
                "color": "#7ea800",
                "displayOrder": 0
            },
            "summary": "first issue",
            "description": "",
            "resolution": null,
            "priority": {
                "id": 3,
                "name": "Normal"
            },
            "status": {
                "id": 1,
                "projectId": 1,
                "name": "Open",
                "color": "#ed8077",
                "displayOrder": 1000
            },
            "assignee": null,
            "category": [],
            "versions": [],
            "milestone": [],
            "startDate": null,
            "dueDate": null,
            "estimatedHours": null,
            "actualHours": null,
            "parentIssueId": null,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "created": "2012-07-23T06:10:15Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "updated": "2013-02-07T08:09:49Z",
            "customFields": [],
            "attachments": [],
            "sharedFiles": [],
            "stars": []
        },
        "updated": "2014-07-11T06:21:07Z"
    },
    {
        "issue": {
            "id": 2,
            "projectId": 1,
            "issueKey": "TEST-2",
            "keyId": 2,
            "issueType": {
                "id": 2,
                "projectId": 1,
                "name": "Task",
                "color": "#7ea800",
                "displayOrder": 0
            },
            "summary": "second issue",
            "description": "",
            "resolution": null,
            "priority": {
                "id": 3,
                "name": "Normal"
            },
            "status": {
                "id": 1,
                "projectId": 1,
                "name": "Open",
                "color": "#ed8077",
                "displayOrder": 1000
            },
            "assignee": null,
            "category": [],
            "versions": [],
            "milestone": [],
            "startDate": null,
            "dueDate": null,
            "estimatedHours": null,
            "actualHours": null,
            "parentIssueId": null,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "created": "2012-07-23T06:10:15Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "updated": "2013-02-07T08:09:49Z",
            "customFields": [],
            "attachments": [],
            "sharedFiles": [],
            "stars": []
        },
        "updated": "2014-07-10T03:15:42Z"
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
        "project": {
            "id": 1,
            "projectKey": "TEST",
            "name": "Test Project",
            "chartEnabled": false,
            "useResolvedForChart": false,
            "subtaskingEnabled": false,
            "projectLeaderCanEditProjectLeader": false,
            "useWiki": true,
            "useFileSharing": true,
            "useWikiTreeView": true,
            "useOriginalImageSizeAtWiki": false,
            "textFormattingRule": "markdown",
            "archived": false,
            "displayOrder": 2147483646,
            "useDevAttributes": true
        },
        "updated": "2014-07-11T06:21:07Z"
    },
    {
        "project": {
            "id": 2,
            "projectKey": "DEMO",
            "name": "Demo Project",
            "chartEnabled": true,
            "useResolvedForChart": false,
            "subtaskingEnabled": true,
            "projectLeaderCanEditProjectLeader": false,
            "useWiki": true,
            "useFileSharing": true,
            "useWikiTreeView": true,
            "useOriginalImageSizeAtWiki": false,
            "textFormattingRule": "backlog",
            "archived": false,
            "displayOrder": 2147483646,
            "useDevAttributes": true
        },
        "updated": "2014-07-10T03:15:42Z"
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
        "page": {
            "id": 10,
            "projectId": 1,
            "name": "Home",
            "tags": [
                {
                    "id": 12,
                    "name": "proceedings"
                }
            ],
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "created": "2013-05-30T09:11:36Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "updated": "2013-05-30T09:11:36Z"
        },
        "updated": "2014-07-11T06:21:07Z"
    },
    {
        "page": {
            "id": 11,
            "projectId": 1,
            "name": "Guide",
            "tags": [
                {
                    "id": 12,
                    "name": "proceedings"
                }
            ],
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "created": "2013-05-30T09:11:36Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "updated": "2013-05-30T09:11:36Z"
        },
        "updated": "2014-07-10T03:15:42Z"
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
