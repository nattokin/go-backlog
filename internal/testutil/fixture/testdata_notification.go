package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type notificationFixtures struct {
	ListJSON  string
	List      []*backlog.Notification
	CountJSON string

	// ResetUnreadCountJSON is the response of Reset Unread Notification Count.
	ResetUnreadCountJSON string
}

// Notification provides test fixtures for Notification-related tests.
var Notification = notificationFixtures{
	ListJSON: `
[
    {
        "id": 299,
        "alreadyRead": true,
        "reason": 2,
        "resourceAlreadyRead": true,
        "project": {
            "id": 2,
            "projectKey": "TEST2",
            "name": "test2",
            "chartEnabled": true,
            "useResolvedForChart": true,
            "subtaskingEnabled": true,
            "projectLeaderCanEditProjectLeader": false,
            "useWiki": true,
            "useDocument": true,
            "useFileSharing": true,
            "useWikiTreeView": true,
            "useSubversion": true,
            "useGit": true,
            "useOriginalImageSizeAtWiki": false,
            "textFormattingRule": "backlog",
            "archived": false,
            "displayOrder": 3,
            "useDevAttributes": true
        },
        "issue": {
            "id": 4531,
            "projectId": 2,
            "issueKey": "TEST2-17",
            "keyId": 17,
            "issueType": {
                "id": 7,
                "projectId": 2,
                "name": "Bug",
                "color": "#990000",
                "displayOrder": 0
            },
            "summary": "aaa",
            "description": "",
            "resolution": null,
            "priority": {
                "id": 3,
                "name": "Normal"
            },
            "status": {
                "id": 1,
                "projectId": 2,
                "name": "Open",
                "color": "#ed8077",
                "displayOrder": 1000
            },
            "assignee": {
                "id": 2,
                "userId": "eguchi",
                "name": "eguchi",
                "roleType": 2,
                "lang": null,
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "category": [],
            "versions": [],
            "milestone": [],
            "startDate": "2013-08-29T15:00:00Z",
            "dueDate": "2013-09-03T15:00:00Z",
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
            "created": "2013-04-23T07:38:59Z",
            "updatedUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "updated": "2013-09-06T09:25:41Z",
            "customFields": [],
            "attachments": [],
            "sharedFiles": [],
            "stars": []
        },
        "comment": {
            "id": 7007,
            "projectId": 5,
            "issueId": 50,
            "content": "hoge",
            "changeLog": null,
            "createdUser": {
                "id": 2,
                "userId": "eguchi",
                "name": "eguchi",
                "roleType": 2,
                "lang": null,
                "mailAddress": "eguchi@nulab.example",
                "lastLoginTime": "2022-09-01T06:35:39Z"
            },
            "created": "2013-10-31T06:58:58Z",
            "updated": "2013-10-31T06:58:58Z",
            "stars": [],
            "notifications": []
        },
        "repository": null,
        "pullRequest": null,
        "pullRequestComment": null,
        "sender": {
            "id": 2,
            "userId": "eguchi",
            "name": "eguchi",
            "roleType": 2,
            "lang": null,
            "mailAddress": "eguchi@nulab.example",
            "lastLoginTime": "2022-09-01T06:35:39Z"
        },
        "created": "2013-10-31T06:58:59Z"
    }
]
`,
	List: []*backlog.Notification{
		{
			ID:                  299,
			AlreadyRead:         true,
			Reason:              2,
			ResourceAlreadyRead: true,
			Project: &backlog.Project{
				ID:                  2,
				ProjectKey:          "TEST2",
				Name:                "test2",
				ChartEnabled:        true,
				UseResolvedForChart: true,
				SubtaskingEnabled:   true,
				UseWiki:             true,
				UseFileSharing:      true,
				UseWikiTreeView:     true,
				UseSubversion:       true,
				UseGit:              true,
				TextFormattingRule:  backlog.FormatBacklog,
				DisplayOrder:        3,
				UseDevAttributes:    true,
			},
			Issue: &backlog.Issue{
				ID:        4531,
				ProjectID: 2,
				IssueKey:  "TEST2-17",
				KeyID:     17,
				IssueType: &backlog.IssueType{
					ID:           7,
					ProjectID:    2,
					Name:         "Bug",
					Color:        "#990000",
					DisplayOrder: 0,
				},
				Summary:  "aaa",
				Priority: &backlog.Priority{ID: 3, Name: "Normal"},
				Status: &backlog.Status{
					ID:           1,
					ProjectID:    2,
					Name:         "Open",
					Color:        "#ed8077",
					DisplayOrder: 1000,
				},
				Assignee: &backlog.User{
					ID:            2,
					UserID:        "eguchi",
					Name:          "eguchi",
					RoleType:      backlog.RoleNormalUser,
					MailAddress:   "eguchi@nulab.example",
					LastLoginTime: mustTimestamp("2022-09-01T06:35:39Z"),
				},
				Category:  []*backlog.Category{},
				Versions:  []*backlog.Version{},
				Milestone: []*backlog.Version{},
				StartDate: mustTimestamp("2013-08-29T15:00:00Z"),
				DueDate:   mustTimestamp("2013-09-03T15:00:00Z"),
				CreatedUser: &backlog.User{
					ID:            1,
					UserID:        "admin",
					Name:          "admin",
					RoleType:      backlog.RoleAdministrator,
					Lang:          "ja",
					MailAddress:   "eguchi@nulab.example",
					LastLoginTime: mustTimestamp("2022-09-01T06:35:39Z"),
				},
				Created: mustTimestamp("2013-04-23T07:38:59Z"),
				UpdatedUser: &backlog.User{
					ID:            1,
					UserID:        "admin",
					Name:          "admin",
					RoleType:      backlog.RoleAdministrator,
					Lang:          "ja",
					MailAddress:   "eguchi@nulab.example",
					LastLoginTime: mustTimestamp("2022-09-01T06:35:39Z"),
				},
				Updated:      mustTimestamp("2013-09-06T09:25:41Z"),
				CustomFields: []*backlog.CustomField{},
				Attachments:  []*backlog.Attachment{},
				SharedFiles:  []*backlog.SharedFile{},
				Stars:        []*backlog.Star{},
			},
			Comment: &backlog.Comment{
				ID:         7007,
				Content:    "hoge",
				ChangeLogs: []*backlog.ChangeLog{},
				CreatedUser: &backlog.User{
					ID:            2,
					UserID:        "eguchi",
					Name:          "eguchi",
					RoleType:      backlog.RoleNormalUser,
					MailAddress:   "eguchi@nulab.example",
					LastLoginTime: mustTimestamp("2022-09-01T06:35:39Z"),
				},
				Created:       mustTimestamp("2013-10-31T06:58:58Z"),
				Updated:       mustTimestamp("2013-10-31T06:58:58Z"),
				Stars:         []*backlog.Star{},
				Notifications: []*backlog.Notification{},
				ProjectID:     5,
				IssueID:       50,
			},
			Sender: &backlog.User{
				ID:            2,
				UserID:        "eguchi",
				Name:          "eguchi",
				RoleType:      backlog.RoleNormalUser,
				MailAddress:   "eguchi@nulab.example",
				LastLoginTime: mustTimestamp("2022-09-01T06:35:39Z"),
			},
			Created: mustTimestamp("2013-10-31T06:58:59Z"),
		},
	},
	CountJSON:            `{"count": 138}`,
	ResetUnreadCountJSON: `{"count": 4}`,
}
