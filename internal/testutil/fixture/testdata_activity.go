package fixture

import (
	"github.com/nattokin/go-backlog"
)

type activityFixtures struct {
	ListJSON   string
	List       []*backlog.Activity
	SingleJSON string
	Single     *backlog.Activity
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
                "content": "test comment",
                "changeLog": [
                    {
                        "field": "status",
                        "newValue": "4",
                        "originalValue": "1"
                    },
                    {
                        "field": "milestone",
                        "newValue": "R2014-07-23",
                        "originalValue": ""
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
                "created": "2014-07-21T06:48:40Z",
                "updated": "2014-07-21T06:48:40Z",
                "stars": [
                    {
                        "id": 75,
                        "comment": "ok",
                        "url": "https://xx.backlogtool.com/view/BLG-1",
                        "title": "[BLG-1] first issue | Show issue - Backlog",
                        "presenter": {
                            "id": 1,
                            "userId": "admin",
                            "name": "admin",
                            "roleType": 1,
                            "lang": "ja",
                            "mailAddress": "eguchi@nulab.example"
                        },
                        "created": "2014-01-23T10:55:19Z"
                    }
                ],
                "notifications": [
                    {
                        "id": 25,
                        "alreadyRead": false,
                        "reason": 2,
                        "resourceAlreadyRead": false
                    }
                ]
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
				ID:          4809,
				KeyID:       121,
				Summary:     "Comment",
				Description: "",
				Comment: &backlog.Comment{
					ID:      7237,
					Content: "test comment",
					ChangeLogs: []*backlog.ChangeLog{
						{Field: "status", NewValue: "4", OriginalValue: "1"},
						{Field: "milestone", NewValue: "R2014-07-23", OriginalValue: ""},
					},
					CreatedUser: &backlog.User{
						ID:          1,
						UserID:      "admin",
						Name:        "admin",
						RoleType:    backlog.RoleAdministrator,
						Lang:        "ja",
						MailAddress: "eguchi@nulab.example",
					},
					Created: mustTime("2014-07-21T06:48:40Z"),
					Updated: mustTime("2014-07-21T06:48:40Z"),
					Stars: []*backlog.Star{
						{
							ID:      75,
							Comment: "ok",
							URL:     "https://xx.backlogtool.com/view/BLG-1",
							Title:   "[BLG-1] first issue | Show issue - Backlog",
							Presenter: &backlog.User{
								ID:          1,
								UserID:      "admin",
								Name:        "admin",
								RoleType:    backlog.RoleAdministrator,
								Lang:        "ja",
								MailAddress: "eguchi@nulab.example",
							},
							Created: mustTime("2014-01-23T10:55:19Z"),
						},
					},
					Notifications: []*backlog.Notification{
						{
							ID:          25,
							AlreadyRead: false,
							Reason:      2,
						},
					},
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
		},
	},
	SingleJSON: `{
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
        "description": ""
    },
    "notifications": [],
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2014-07-21T06:48:40Z"
}`,
	Single: &backlog.Activity{
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
			ID:          4809,
			KeyID:       121,
			Summary:     "Comment",
			Description: "",
		},
		Notifications: []*backlog.Notification{},
		CreatedUser: &backlog.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    backlog.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
	},
}
