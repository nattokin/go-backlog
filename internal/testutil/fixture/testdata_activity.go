package fixture

import (
	"time"

	backlog "github.com/nattokin/go-backlog"
)

type activityFixtures struct {
	ListJSON string
	List     []*backlog.Activity
}

// Activity provides test fixtures for Activity-related tests.
var Activity = activityFixtures{
	ListJSON: `
[
    {
        "id": 3153,
        "project": {
            "id": 92,
            "projectKey": "SUB",
            "name": "Subtasking",
            "chartEnabled": true,
            "subtaskingEnabled": true,
            "projectLeaderCanEditProjectLeader": false,
            "textFormattingRule": null,
            "archived": false,
            "displayOrder": 0
        },
        "type": 2,
        "content": {
            "id": 4809,
            "key_id": 121,
            "summary": "Comment",
            "description": "",
            "comment": {
                "id": 7237,
                "content": ""
            },
            "changes": [
                {
                    "field": "milestone",
                    "new_value": "R2014-07-23",
                    "old_value": "",
                    "type": "standard"
                },
                {
                    "field": "status",
                    "new_value": "4",
                    "old_value": "1",
                    "type": "standard"
                }
            ]
        },
        "notifications": [
            {
                "id": 25,
                "alreadyRead": false,
                "reason": 2,
                "user": {
                    "id": 5686,
                    "userId": "takada",
                    "name": "takada",
                    "roleType": 2,
                    "lang": "ja",
                    "mailAddress": "takada@nulab.example"
                },
                "resourceAlreadyRead": false
            }
        ],
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-21T06:48:40Z"
    }
]
`,
	List: []*backlog.Activity{
		{
			ID: 3153,
			Project: &backlog.Project{
				ID:                92,
				ProjectKey:        "SUB",
				Name:              "Subtasking",
				ChartEnabled:      true,
				SubtaskingEnabled: true,
			},
			Type: 2,
			Content: &backlog.ActivityContent{
				ID:      4809,
				KeyID:   121,
				Summary: "Comment",
				Comment: &backlog.Comment{
					ID: 7237,
				},
			},
			Notifications: []*backlog.Notification{
				{
					ID:          25,
					AlreadyRead: false,
					Reason:      2,
				},
			},
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    backlog.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			// Created は time.Time のゼロ値比較が難しいため省略
		},
	},
}

// activityCreated is the parsed time for the activity list fixture.
var activityCreated, _ = time.Parse(time.RFC3339, "2014-07-21T06:48:40Z")
