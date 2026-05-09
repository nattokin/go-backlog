package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type webhookFixtures struct {
	AllEventJSON      string
	AllEvent          *backlog.Webhook
	ActivityTypesJSON string
	ActivityTypes     *backlog.Webhook
	ListJSON          string
	List              []*backlog.Webhook
}

// Webhook provides test fixtures for Webhook-related tests.

var webhookAllEvent = &backlog.Webhook{
	ID:              1,
	Name:            "Example Webhook",
	HookURL:         "https://example.com/webhook",
	AllEvent:        true,
	ActivityTypeIDs: []int{},
	CreatedUser: &backlog.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    backlog.RoleAdministrator,
		Lang:        "ja",
		MailAddress: "eguchi@nulab.example",
	},
	Created: mustTimestamp("2024-01-01T00:00:00Z"),
	UpdatedUser: &backlog.User{
		ID:          2,
		UserID:      "normal_user",
		Name:        "normal_user",
		RoleType:    backlog.RoleNormalUser,
		Lang:        "ja",
		MailAddress: "sato@nulab.example",
	},
	Updated: mustTimestamp("2024-01-02T00:00:00Z"),
}

var webhookActivityTypes = &backlog.Webhook{
	ID:              2,
	Name:            "Issue Events Webhook",
	Description:     "Webhook for issue updates",
	HookURL:         "https://example.com/issues",
	AllEvent:        false,
	ActivityTypeIDs: []int{1, 2, 3, 4},
	CreatedUser: &backlog.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    backlog.RoleAdministrator,
		Lang:        "ja",
		MailAddress: "eguchi@nulab.example",
	},
	Created: mustTimestamp("2024-01-01T00:00:00Z"),
	UpdatedUser: &backlog.User{
		ID:          2,
		UserID:      "normal_user",
		Name:        "normal_user",
		RoleType:    backlog.RoleNormalUser,
		Lang:        "ja",
		MailAddress: "sato@nulab.example",
	},
	Updated: mustTimestamp("2024-01-02T00:00:00Z"),
}

var Webhook = webhookFixtures{
	AllEventJSON: `
{
    "id": 1,
    "name": "Example Webhook",
    "hookUrl": "https://example.com/webhook",
    "allEvent": true,
    "activityTypeIds": [],
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2024-01-01T00:00:00Z",
    "updatedUser": {
        "id": 2,
        "userId": "normal_user",
        "name": "normal_user",
        "roleType": 2,
        "lang": "ja",
        "mailAddress": "sato@nulab.example"
    },
    "updated": "2024-01-02T00:00:00Z"
}
`,
	AllEvent: webhookAllEvent,
	ActivityTypesJSON: `
{
    "id": 2,
    "name": "Issue Events Webhook",
    "description": "Webhook for issue updates",
    "hookUrl": "https://example.com/issues",
    "allEvent": false,
    "activityTypeIds": [1,2,3,4],
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2024-01-01T00:00:00Z",
    "updatedUser": {
        "id": 2,
        "userId": "normal_user",
        "name": "normal_user",
        "roleType": 2,
        "lang": "ja",
        "mailAddress": "sato@nulab.example"
    },
    "updated": "2024-01-02T00:00:00Z"
}
`,
	ActivityTypes: webhookActivityTypes,
	ListJSON: `
[
    {
        "id": 1,
        "name": "Example Webhook",
        "hookUrl": "https://example.com/webhook",
        "allEvent": true,
        "activityTypeIds": [],
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2024-01-01T00:00:00Z",
        "updatedUser": {
            "id": 2,
            "userId": "normal_user",
            "name": "normal_user",
            "roleType": 2,
            "lang": "ja",
            "mailAddress": "sato@nulab.example"
        },
        "updated": "2024-01-02T00:00:00Z"
    },
    {
        "id": 2,
        "name": "Issue Events Webhook",
        "hookUrl": "https://example.com/issues",
        "allEvent": false,
        "activityTypeIds": [1,2,3,4],
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2024-01-01T00:00:00Z",
        "updatedUser": {
            "id": 2,
            "userId": "normal_user",
            "name": "normal_user",
            "roleType": 2,
            "lang": "ja",
            "mailAddress": "sato@nulab.example"
        },
        "updated": "2024-01-02T00:00:00Z"
    }
]
`,
	List: []*backlog.Webhook{
		webhookAllEvent,
		webhookActivityTypes,
	},
}
