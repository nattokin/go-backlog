package fixture

import backlog "github.com/nattokin/go-backlog"

type watchingFixtures struct {
	SingleJSON string
	Single     backlog.Watching
	ListJSON   string
	List       []*backlog.Watching
	CountJSON  string
}

// Watching provides test fixtures for Watching-related tests.
// The watched issue has the same content as Issue.SingleJSON.
var Watching = watchingFixtures{
	SingleJSON: `
{
    "id": 1,
    "alreadyRead": true,
    "note": "Check before release",
    "type": "issue",
    "issue": {
        "id": 1,
        "projectId": 10,
        "issueKey": "PRJ-1",
        "keyId": 1,
        "issueType": {
            "id": 2,
            "projectId": 10,
            "name": "Bug",
            "color": "#990000",
            "displayOrder": 0
        },
        "summary": "First issue",
        "description": "Description of first issue",
        "resolution": {
            "id": 0,
            "name": "Fixed"
        },
        "priority": {
            "id": 3,
            "name": "Normal"
        },
        "status": {
            "id": 1,
            "projectId": 10,
            "name": "Open",
            "color": "#ed8077",
            "displayOrder": 1000
        },
        "assignee": null,
        "category": [],
        "versions": [],
        "milestone": [],
        "startDate": "2024-01-14T15:00:00Z",
        "dueDate": "2024-01-30T15:00:00Z",
        "estimatedHours": null,
        "actualHours": null,
        "parentIssueId": null,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "admin@example.com"
        },
        "created": "2024-01-10T09:00:00Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "admin@example.com"
        },
        "updated": "2024-01-10T09:00:00Z",
        "customFields": [],
        "attachments": [],
        "sharedFiles": [],
        "stars": []
    },
    "lastContentUpdated": "2024-01-10T09:00:00Z",
    "created": "2024-01-11T09:00:00Z",
    "updated": "2024-01-12T09:00:00Z"
}
`,
	Single: backlog.Watching{
		ID:                 1,
		AlreadyRead:        true,
		Note:               "Check before release",
		Type:               "issue",
		Issue:              watchedIssue(),
		LastContentUpdated: mustTimestamp("2024-01-10T09:00:00Z"),
		Created:            mustTimestamp("2024-01-11T09:00:00Z"),
		Updated:            mustTimestamp("2024-01-12T09:00:00Z"),
	},
	ListJSON: `
[
    {
        "id": 1,
        "alreadyRead": true,
        "note": "Check before release",
        "type": "issue",
        "issue": {
            "id": 1,
            "projectId": 10,
            "issueKey": "PRJ-1",
            "keyId": 1,
            "issueType": {
                "id": 2,
                "projectId": 10,
                "name": "Bug",
                "color": "#990000",
                "displayOrder": 0
            },
            "summary": "First issue",
            "description": "Description of first issue",
            "resolution": {
                "id": 0,
                "name": "Fixed"
            },
            "priority": {
                "id": 3,
                "name": "Normal"
            },
            "status": {
                "id": 1,
                "projectId": 10,
                "name": "Open",
                "color": "#ed8077",
                "displayOrder": 1000
            },
            "assignee": null,
            "category": [],
            "versions": [],
            "milestone": [],
            "startDate": "2024-01-14T15:00:00Z",
            "dueDate": "2024-01-30T15:00:00Z",
            "estimatedHours": null,
            "actualHours": null,
            "parentIssueId": null,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "admin@example.com"
            },
            "created": "2024-01-10T09:00:00Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "admin@example.com"
            },
            "updated": "2024-01-10T09:00:00Z",
            "customFields": [],
            "attachments": [],
            "sharedFiles": [],
            "stars": []
        },
        "lastContentUpdated": "2024-01-10T09:00:00Z",
        "created": "2024-01-11T09:00:00Z",
        "updated": "2024-01-12T09:00:00Z"
    },
    {
        "id": 2,
        "alreadyRead": false,
        "note": "",
        "type": "issue",
        "issue": {
            "id": 1,
            "projectId": 10,
            "issueKey": "PRJ-1",
            "keyId": 1,
            "issueType": {
                "id": 2,
                "projectId": 10,
                "name": "Bug",
                "color": "#990000",
                "displayOrder": 0
            },
            "summary": "First issue",
            "description": "Description of first issue",
            "resolution": {
                "id": 0,
                "name": "Fixed"
            },
            "priority": {
                "id": 3,
                "name": "Normal"
            },
            "status": {
                "id": 1,
                "projectId": 10,
                "name": "Open",
                "color": "#ed8077",
                "displayOrder": 1000
            },
            "assignee": null,
            "category": [],
            "versions": [],
            "milestone": [],
            "startDate": "2024-01-14T15:00:00Z",
            "dueDate": "2024-01-30T15:00:00Z",
            "estimatedHours": null,
            "actualHours": null,
            "parentIssueId": null,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "admin@example.com"
            },
            "created": "2024-01-10T09:00:00Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "admin@example.com"
            },
            "updated": "2024-01-10T09:00:00Z",
            "customFields": [],
            "attachments": [],
            "sharedFiles": [],
            "stars": []
        },
        "lastContentUpdated": "2024-01-10T09:00:00Z",
        "created": "2024-01-13T09:00:00Z",
        "updated": "2024-01-13T09:00:00Z"
    }
]
`,
	List: []*backlog.Watching{
		{
			ID:                 1,
			AlreadyRead:        true,
			Note:               "Check before release",
			Type:               "issue",
			Issue:              watchedIssue(),
			LastContentUpdated: mustTimestamp("2024-01-10T09:00:00Z"),
			Created:            mustTimestamp("2024-01-11T09:00:00Z"),
			Updated:            mustTimestamp("2024-01-12T09:00:00Z"),
		},
		{
			ID:                 2,
			AlreadyRead:        false,
			Note:               "",
			Type:               "issue",
			Issue:              watchedIssue(),
			LastContentUpdated: mustTimestamp("2024-01-10T09:00:00Z"),
			Created:            mustTimestamp("2024-01-13T09:00:00Z"),
			Updated:            mustTimestamp("2024-01-13T09:00:00Z"),
		},
	},
	CountJSON: `{"count": 138}`,
}

// watchedIssue returns the issue in Issue.SingleJSON as decoded by the client.
func watchedIssue() *backlog.Issue {
	admin := func() *backlog.User {
		return &backlog.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    backlog.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "admin@example.com",
		}
	}
	return &backlog.Issue{
		ID:        1,
		ProjectID: 10,
		IssueKey:  "PRJ-1",
		KeyID:     1,
		IssueType: &backlog.IssueType{
			ID:           2,
			ProjectID:    10,
			Name:         "Bug",
			Color:        "#990000",
			DisplayOrder: 0,
		},
		Summary:     "First issue",
		Description: "Description of first issue",
		Resolution:  &backlog.Resolution{ID: 0, Name: "Fixed"},
		Priority:    &backlog.Priority{ID: 3, Name: "Normal"},
		Status: &backlog.Status{
			ID:           1,
			ProjectID:    10,
			Name:         "Open",
			Color:        "#ed8077",
			DisplayOrder: 1000,
		},
		Category:     []*backlog.Category{},
		Versions:     []*backlog.Version{},
		Milestone:    []*backlog.Version{},
		StartDate:    mustTimestamp("2024-01-14T15:00:00Z"),
		DueDate:      mustTimestamp("2024-01-30T15:00:00Z"),
		CreatedUser:  admin(),
		Created:      mustTimestamp("2024-01-10T09:00:00Z"),
		UpdatedUser:  admin(),
		Updated:      mustTimestamp("2024-01-10T09:00:00Z"),
		CustomFields: []*backlog.CustomField{},
		Attachments:  []*backlog.Attachment{},
		SharedFiles:  []*backlog.SharedFile{},
		Stars:        []*backlog.Star{},
	}
}
