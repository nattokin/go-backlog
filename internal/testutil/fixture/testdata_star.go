package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type starFixtures struct {
	ListJSON  string
	List      []*backlog.Star
	CountJSON string
}

// Star provides test fixtures for Star-related tests.
var Star = starFixtures{
	ListJSON: `
[
    {
        "id": 10,
        "comment": null,
        "url": "https://example.backlog.com/view/TEST-1",
        "title": "[TEST-1] first issue",
        "presenter": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2024-01-15T10:00:00Z"
    },
    {
        "id": 20,
        "comment": "nice!",
        "url": "https://example.backlog.com/view/TEST-2",
        "title": "[TEST-2] second issue",
        "presenter": {
            "id": 2,
            "userId": "normal_user",
            "name": "normal_user",
            "roleType": 2,
            "lang": "ja",
            "mailAddress": "sato@nulab.example"
        },
        "created": "2024-02-20T12:00:00Z"
    }
]
`,
	List: []*backlog.Star{
		{
			ID:      10,
			Comment: "",
			URL:     "https://example.backlog.com/view/TEST-1",
			Title:   "[TEST-1] first issue",
			Presenter: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2024-01-15T10:00:00Z"),
		},
		{
			ID:      20,
			Comment: "nice!",
			URL:     "https://example.backlog.com/view/TEST-2",
			Title:   "[TEST-2] second issue",
			Presenter: &backlog.User{
				ID:          2,
				UserID:      "normal_user",
				Name:        "normal_user",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "sato@nulab.example",
			},
			Created: mustTime("2024-02-20T12:00:00Z"),
		},
	},
	CountJSON: `{"count": 42}`,
}
