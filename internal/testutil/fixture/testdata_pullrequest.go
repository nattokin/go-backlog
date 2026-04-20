package fixture

import backlog "github.com/nattokin/go-backlog"

type pullRequestFixtures struct {
	SingleJSON string
	Single     backlog.PullRequest
	ListJSON   string
	List       []*backlog.PullRequest
}

// PullRequest provides test fixtures for PullRequest-related tests.
var PullRequest = pullRequestFixtures{
	SingleJSON: `
{
    "id": 2,
    "projectId": 3,
    "repositoryId": 5,
    "number": 1,
    "summary": "test PR",
    "description": "test description",
    "base": "main",
    "branch": "feature/foo",
    "status": {
        "id": 1,
        "name": "Open"
    },
    "assignee": null,
    "issue": null,
    "baseCommit": null,
    "branchCommit": null,
    "closeAt": null,
    "mergeAt": null,
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
    "attachments": [],
    "stars": []
}
`,
	Single: backlog.PullRequest{
		ID:           2,
		ProjectID:    3,
		RepositoryID: 5,
		Number:       1,
		Summary:      "test PR",
		Description:  "test description",
		Base:         "main",
		Branch:       "feature/foo",
		Status:       &backlog.Status{ID: 1, Name: "Open"},
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
	ListJSON: `
[
    {
        "id": 2,
        "projectId": 3,
        "repositoryId": 5,
        "number": 1,
        "summary": "test PR",
        "description": "test description",
        "base": "main",
        "branch": "feature/foo",
        "status": {
            "id": 1,
            "name": "Open"
        },
        "assignee": null,
        "issue": null,
        "baseCommit": null,
        "branchCommit": null,
        "closeAt": null,
        "mergeAt": null,
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
        "attachments": [],
        "stars": []
    },
    {
        "id": 3,
        "projectId": 3,
        "repositoryId": 5,
        "number": 2,
        "summary": "second PR",
        "description": "",
        "base": "main",
        "branch": "feature/bar",
        "status": {
            "id": 2,
            "name": "Closed"
        },
        "assignee": null,
        "issue": null,
        "baseCommit": null,
        "branchCommit": null,
        "closeAt": null,
        "mergeAt": null,
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
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "admin@example.com"
        },
        "updated": "2024-01-11T10:00:00Z",
        "attachments": [],
        "stars": []
    }
]
`,
	List: []*backlog.PullRequest{
		{
			ID:           2,
			ProjectID:    3,
			RepositoryID: 5,
			Number:       1,
			Summary:      "test PR",
			Description:  "test description",
			Base:         "main",
			Branch:       "feature/foo",
			Status:       &backlog.Status{ID: 1, Name: "Open"},
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
			ID:           3,
			ProjectID:    3,
			RepositoryID: 5,
			Number:       2,
			Summary:      "second PR",
			Base:         "main",
			Branch:       "feature/bar",
			Status:       &backlog.Status{ID: 2, Name: "Closed"},
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
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "admin@example.com",
			},
			Updated: mustTime("2024-01-11T10:00:00Z"),
		},
	},
}
