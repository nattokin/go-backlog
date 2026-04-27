package fixture

import backlog "github.com/nattokin/go-backlog"

type repositoryFixtures struct {
	SingleJSON string
	Single     backlog.Repository
	ListJSON   string
	List       []*backlog.Repository
}

// Repository provides test fixtures for Repository-related tests.
var Repository = repositoryFixtures{
	SingleJSON: `
{
    "id": 5,
    "projectId": 3,
    "name": "foo",
    "description": "test repo",
    "hookUrl": null,
    "httpUrl": "https://example.backlog.com/git/TEST/foo.git",
    "sshUrl": "git@example.backlog.com:TEST/foo.git",
    "displayOrder": 0,
    "pushedAt": null,
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
    "updated": "2024-01-10T09:00:00Z"
}
`,
	Single: backlog.Repository{
		ID:           5,
		ProjectID:    3,
		Name:         "foo",
		Description:  "test repo",
		HTTPURL:      "https://example.backlog.com/git/TEST/foo.git",
		SSHURL:       "git@example.backlog.com:TEST/foo.git",
		DisplayOrder: 0,
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
        "id": 5,
        "projectId": 3,
        "name": "foo",
        "description": "test repo",
        "hookUrl": null,
        "httpUrl": "https://example.backlog.com/git/TEST/foo.git",
        "sshUrl": "git@example.backlog.com:TEST/foo.git",
        "displayOrder": 0,
        "pushedAt": null,
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
        "updated": "2024-01-10T09:00:00Z"
    },
    {
        "id": 6,
        "projectId": 3,
        "name": "bar",
        "description": "second repo",
        "hookUrl": null,
        "httpUrl": "https://example.backlog.com/git/TEST/bar.git",
        "sshUrl": "git@example.backlog.com:TEST/bar.git",
        "displayOrder": 1,
        "pushedAt": null,
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
        "updated": "2024-01-11T10:00:00Z"
    }
]
`,
	List: []*backlog.Repository{
		{
			ID:           5,
			ProjectID:    3,
			Name:         "foo",
			Description:  "test repo",
			HTTPURL:      "https://example.backlog.com/git/TEST/foo.git",
			SSHURL:       "git@example.backlog.com:TEST/foo.git",
			DisplayOrder: 0,
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
			ID:           6,
			ProjectID:    3,
			Name:         "bar",
			Description:  "second repo",
			HTTPURL:      "https://example.backlog.com/git/TEST/bar.git",
			SSHURL:       "git@example.backlog.com:TEST/bar.git",
			DisplayOrder: 1,
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
