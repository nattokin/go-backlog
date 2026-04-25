package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type sharedFileFixtures struct {
	SingleJSON     string
	Single         *backlog.SharedFile
	SingleListJSON string
	SingleList     []*backlog.SharedFile
	ListJSON       string
	List           []*backlog.SharedFile
}

// SharedFile provides test fixtures for SharedFile-related tests.
var SharedFile = sharedFileFixtures{
	SingleJSON: `
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
`,
	Single: &backlog.SharedFile{
		ID:   454403,
		Type: "file",
		Dir:  "/icon/",
		Name: "01_buz.png",
		Size: 2735,
		CreatedUser: &backlog.User{
			ID:          5686,
			UserID:      "takada",
			Name:        "takada",
			RoleType:    backlog.RoleNormalUser,
			Lang:        "ja",
			MailAddress: "takada@nulab.example",
		},
		Created: mustTime("2009-02-27T03:26:15Z"),
		UpdatedUser: &backlog.User{
			ID:          5686,
			UserID:      "takada",
			Name:        "takada",
			RoleType:    backlog.RoleNormalUser,
			Lang:        "ja",
			MailAddress: "takada@nulab.example",
		},
		Updated: mustTime("2009-03-03T16:57:47Z"),
	},
	SingleListJSON: `
[
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
]
`,
	SingleList: []*backlog.SharedFile{
		{
			ID:   454403,
			Type: "file",
			Dir:  "/icon/",
			Name: "01_buz.png",
			Size: 2735,
			CreatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Created: mustTime("2009-02-27T03:26:15Z"),
			UpdatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Updated: mustTime("2009-03-03T16:57:47Z"),
		},
	},
	ListJSON: `
[
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
    },
    {
        "id": 454404,
        "type": "file",
        "dir": "/docs/",
        "name": "readme.md",
        "size": 512,
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
]
`,
	List: []*backlog.SharedFile{
		{
			ID:   454403,
			Type: "file",
			Dir:  "/icon/",
			Name: "01_buz.png",
			Size: 2735,
			CreatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Created: mustTime("2009-02-27T03:26:15Z"),
			UpdatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Updated: mustTime("2009-03-03T16:57:47Z"),
		},
		{
			ID:   454404,
			Type: "file",
			Dir:  "/docs/",
			Name: "readme.md",
			Size: 512,
			CreatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Created: mustTime("2009-02-27T03:26:15Z"),
			UpdatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Updated: mustTime("2009-03-03T16:57:47Z"),
		},
	},
}
