package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type commentFixtures struct {
	SingleJSON string
	Single     *backlog.Comment
	ListJSON   string
	List       []*backlog.Comment
}

// Comment provides test fixtures for Comment-related tests.
var Comment = commentFixtures{
	SingleJSON: `
{
    "id": 1,
    "content": "This is a comment.",
    "changeLog": [],
    "createdUser": {
        "id": 5686,
        "userId": "takada",
        "name": "takada",
        "roleType": 2,
        "lang": "ja",
        "mailAddress": "takada@nulab.example"
    },
    "created": "2013-08-05T06:15:06Z",
    "updated": "2013-08-05T06:15:06Z",
    "stars": [],
    "notifications": []
}
`,
	Single: &backlog.Comment{
		ID:      1,
		Content: "This is a comment.",
		CreatedUser: &backlog.User{
			ID:          5686,
			UserID:      "takada",
			Name:        "takada",
			RoleType:    backlog.RoleNormalUser,
			Lang:        "ja",
			MailAddress: "takada@nulab.example",
		},
		Created: mustTime("2013-08-05T06:15:06Z"),
		Updated: mustTime("2013-08-05T06:15:06Z"),
	},
	ListJSON: `
[
    {
        "id": 1,
        "content": "This is a comment.",
        "changeLog": [],
        "createdUser": {
            "id": 5686,
            "userId": "takada",
            "name": "takada",
            "roleType": 2,
            "lang": "ja",
            "mailAddress": "takada@nulab.example"
        },
        "created": "2013-08-05T06:15:06Z",
        "updated": "2013-08-05T06:15:06Z",
        "stars": [],
        "notifications": []
    },
    {
        "id": 2,
        "content": "Second comment.",
        "changeLog": [],
        "createdUser": {
            "id": 5686,
            "userId": "takada",
            "name": "takada",
            "roleType": 2,
            "lang": "ja",
            "mailAddress": "takada@nulab.example"
        },
        "created": "2013-08-06T06:15:06Z",
        "updated": "2013-08-06T06:15:06Z",
        "stars": [],
        "notifications": []
    }
]
`,
	List: []*backlog.Comment{
		{
			ID:      1,
			Content: "This is a comment.",
			CreatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Created: mustTime("2013-08-05T06:15:06Z"),
			Updated: mustTime("2013-08-05T06:15:06Z"),
		},
		{
			ID:      2,
			Content: "Second comment.",
			CreatedUser: &backlog.User{
				ID:          5686,
				UserID:      "takada",
				Name:        "takada",
				RoleType:    backlog.RoleNormalUser,
				Lang:        "ja",
				MailAddress: "takada@nulab.example",
			},
			Created: mustTime("2013-08-06T06:15:06Z"),
			Updated: mustTime("2013-08-06T06:15:06Z"),
		},
	},
}
