package fixture

import (
	"github.com/nattokin/go-backlog/internal/model"
)

type wikiFixtures struct {
	MaximumJSON string
	Maximum     model.Wiki
	MinimumJSON string
	Minimum     model.Wiki
	ListJSON    string
	List        []*model.Wiki
}

// Wiki provides test fixtures for Wiki-related tests.
var Wiki = wikiFixtures{
	MaximumJSON: `
{
    "id": 34,
    "projectId": 56,
    "name": "Maximum Wiki Page",
    "content": "This is a muximal wiki page.",
    "tags": [
        {
            "id": 12,
            "name": "proceedings"
        }
    ],
    "attachments": [
        {
            "id": 23,
            "name": "test.json",
            "size": 8857,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example"
            },
            "created": "2014-01-06T11:10:45Z"
        }
    ],
    "sharedFiles": [
        {
            "id": 454403,
            "type": "file",
            "dir": "/icon/",
            "name": "01_buz.png",
            "size": 2735,
            "createdUser": {
                "id": 5686,
                "userId": "takada",
                "name": "takada",
                "roleType": 2,
                "lang": "ja",
                "mailAddress": "takada@nulab.example"
            },
            "created": "2009-02-27T03:26:15Z",
            "updatedUser": {
                "id": 5686,
                "userId": "takada",
                "name": "takada",
                "roleType": 2,
                "lang": "ja",
                "mailAddress": "takada@nulab.example"
            },
            "updated": "2009-03-03T16:57:47Z"
        }
    ],
    "stars": [
        { 
            "id":75, 
            "comment":null, 
            "url": "https://xx.backlogtool.com/view/BLG-1", 
            "title": "[BLG-1] first issue | Show issue - Backlog", 
            "presenter":{ 
                "id":1, 
                "userId": "admin", 
                "name":"admin", 
                "roleType":1, 
                "lang":"ja", 
                "mailAddress":"eguchi@nulab.example" 
            }, 
            "created":"2014-01-23T10:55:19Z" 
        },
        { 
            "id":76, 
            "comment":"ok", 
            "url": "https://xx.backlogtool.com/view/BLG-1", 
            "title": "[BLG-1] first issue | Show issue - Backlog", 
            "presenter":{ 
                "id":1, 
                "userId": "admin", 
                "name":"admin", 
                "roleType":1, 
                "lang":"ja", 
                "mailAddress":"eguchi@nulab.example" 
            }, 
            "created":"2014-01-23T10:55:19Z" 
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
    "created": "2012-07-23T06:09:48Z",
    "updatedUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "updated": "2012-07-23T06:09:48Z"
}
`,
	Maximum: model.Wiki{
		ID:        34,
		ProjectID: 56,
		Name:      "Maximum Wiki Page",
		Content:   "This is a muximal wiki page.",
		Tags: []*model.Tag{
			{ID: 12, Name: "proceedings"},
		},
		Attachments: []*model.Attachment{
			{
				ID:   23,
				Name: "test.json",
				Size: 8857,
				CreatedUser: &model.User{
					ID:          1,
					UserID:      "admin",
					Name:        "admin",
					RoleType:    model.RoleAdministrator,
					Lang:        "ja",
					MailAddress: "eguchi@nulab.example",
				},
				Created: mustTime("2014-01-06T11:10:45Z"),
			},
		},
		SharedFiles: []*model.SharedFile{
			{
				ID:   454403,
				Type: "file",
				Dir:  "/icon/",
				Name: "01_buz.png",
				Size: 2735,
				CreatedUser: &model.User{
					ID:          5686,
					UserID:      "takada",
					Name:        "takada",
					RoleType:    model.RoleNormalUser,
					Lang:        "ja",
					MailAddress: "takada@nulab.example",
				},
				Created: mustTime("2009-02-27T03:26:15Z"),
				UpdatedUser: &model.User{
					ID:          5686,
					UserID:      "takada",
					Name:        "takada",
					RoleType:    model.RoleNormalUser,
					Lang:        "ja",
					MailAddress: "takada@nulab.example",
				},
				Updated: mustTime("2009-03-03T16:57:47Z"),
			},
		},
		Stars: []*model.Star{
			{
				ID:    75,
				URL:   "https://xx.backlogtool.com/view/BLG-1",
				Title: "[BLG-1] first issue | Show issue - Backlog",
				Presenter: &model.User{
					ID:          1,
					UserID:      "admin",
					Name:        "admin",
					RoleType:    model.RoleAdministrator,
					Lang:        "ja",
					MailAddress: "eguchi@nulab.example",
				},
				Created: mustTime("2014-01-23T10:55:19Z"),
			},
			{
				ID:      76,
				Comment: "ok",
				URL:     "https://xx.backlogtool.com/view/BLG-1",
				Title:   "[BLG-1] first issue | Show issue - Backlog",
				Presenter: &model.User{
					ID:          1,
					UserID:      "admin",
					Name:        "admin",
					RoleType:    model.RoleAdministrator,
					Lang:        "ja",
					MailAddress: "eguchi@nulab.example",
				},
				Created: mustTime("2014-01-23T10:55:19Z"),
			},
		},
		CreatedUser: &model.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    model.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		Created: mustTime("2012-07-23T06:09:48Z"),
		UpdatedUser: &model.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    model.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		Updated: mustTime("2012-07-23T06:09:48Z"),
	},
	MinimumJSON: `
{
    "id": 34,
    "projectId": 56,
    "name": "Minimum Wiki Page",
    "content": "This is a minimal wiki page.",
    "tags": [
        {
            "id": 12,
            "name": "proceedings"
        }
    ],
    "attachments": [],
    "sharedFiles": [],
    "stars": [],
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2012-07-23T06:09:48Z",
    "updatedUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "updated": "2012-07-23T06:09:48Z"
}
`,
	Minimum: model.Wiki{
		ID:        34,
		ProjectID: 56,
		Name:      "Minimum Wiki Page",
		Content:   "This is a minimal wiki page.",
		Tags: []*model.Tag{
			{ID: 12, Name: "proceedings"},
		},
		Attachments: []*model.Attachment{},
		SharedFiles: []*model.SharedFile{},
		Stars:       []*model.Star{},
		CreatedUser: &model.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    model.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		Created: mustTime("2012-07-23T06:09:48Z"),
		UpdatedUser: &model.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    model.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		Updated: mustTime("2012-07-23T06:09:48Z"),
	},
	ListJSON: `
[
    {
        "id": 112,
        "projectId": 56,
        "name": "test1",
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
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2013-05-30T09:11:36Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "updated": "2013-05-30T09:11:36Z"
    },
    {
        "id": 115,
        "projectId": 56,
        "name": "test2",
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
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2013-05-30T09:11:36Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "updated": "2013-05-30T09:11:36Z"
    }
]
`,
	List: []*model.Wiki{
		{
			ID:        112,
			ProjectID: 56,
			Name:      "test1",
			Tags:      []*model.Tag{{ID: 12, Name: "proceedings"}},
			CreatedUser: &model.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    model.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2013-05-30T09:11:36Z"),
			UpdatedUser: &model.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    model.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Updated: mustTime("2013-05-30T09:11:36Z"),
		},
		{
			ID:        115,
			ProjectID: 56,
			Name:      "test2",
			Tags:      []*model.Tag{{ID: 12, Name: "proceedings"}},
			CreatedUser: &model.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    model.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2013-05-30T09:11:36Z"),
			UpdatedUser: &model.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    model.RoleAdministrator,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Updated: mustTime("2013-05-30T09:11:36Z"),
		},
	},
}
