package fixture

import backlog "github.com/nattokin/go-backlog"

type issueFixtures struct {
	ListJSON string
	List     []*backlog.Issue
}

// Issue provides test fixtures for Issue-related tests.
var Issue = issueFixtures{
	ListJSON: `
[
    {
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
        "resolutions": null,
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
    {
        "id": 2,
        "projectId": 10,
        "issueKey": "PRJ-2",
        "keyId": 2,
        "issueType": {
            "id": 3,
            "projectId": 10,
            "name": "Task",
            "color": "#7ea800",
            "displayOrder": 1
        },
        "summary": "Second issue",
        "description": "",
        "resolutions": null,
        "priority": {
            "id": 2,
            "name": "High"
        },
        "status": {
            "id": 2,
            "projectId": 10,
            "name": "In Progress",
            "color": "#4488c5",
            "displayOrder": 2000
        },
        "assignee": {
            "id": 2,
            "userId": "user1",
            "name": "User One",
            "roleType": 2,
            "lang": "en",
            "mailAddress": "user1@example.com"
        },
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
            "mailAddress": "admin@example.com"
        },
        "created": "2024-01-11T10:00:00Z",
        "updatedUser": {
            "id": 2,
            "userId": "user1",
            "name": "User One",
            "roleType": 2,
            "lang": "en",
            "mailAddress": "user1@example.com"
        },
        "updated": "2024-01-15T14:30:00Z",
        "customFields": [],
        "attachments": [],
        "sharedFiles": [],
        "stars": []
    }
]
`,
	List: []*backlog.Issue{
		{
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
			Priority:    &backlog.Priority{ID: 3, Name: "Normal"},
			Status: &backlog.Status{
				ID:   1,
				Name: "Open",
			},
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "admin@example.com",
			},
			Created: mustTime("2024-01-10T09:00:00Z"),
			UpdatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "admin@example.com",
			},
			Updated: mustTime("2024-01-10T09:00:00Z"),
		},
		{
			ID:        2,
			ProjectID: 10,
			IssueKey:  "PRJ-2",
			KeyID:     2,
			IssueType: &backlog.IssueType{
				ID:           3,
				ProjectID:    10,
				Name:         "Task",
				Color:        "#7ea800",
				DisplayOrder: 1,
			},
			Summary:  "Second issue",
			Priority: &backlog.Priority{ID: 2, Name: "High"},
			Status: &backlog.Status{
				ID:   2,
				Name: "In Progress",
			},
			Assignee: &backlog.User{
				ID:          2,
				UserID:      "user1",
				Name:        "User One",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "en",
				MailAddress: "user1@example.com",
			},
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "admin@example.com",
			},
			Created: mustTime("2024-01-11T10:00:00Z"),
			UpdatedUser: &backlog.User{
				ID:          2,
				UserID:      "user1",
				Name:        "User One",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "en",
				MailAddress: "user1@example.com",
			},
			Updated: mustTime("2024-01-15T14:30:00Z"),
		},
	},
}
