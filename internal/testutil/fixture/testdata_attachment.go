package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type attachmentFixtures struct {
	SingleJSON     string
	Single         *backlog.Attachment
	ListJSON       string
	List           []*backlog.Attachment
	SingleListJSON string
	SingleList     []*backlog.Attachment
	UploadJSON     string
	Upload         *backlog.Attachment
}

// Attachment provides test fixtures for Attachment-related tests.
var Attachment = attachmentFixtures{
	SingleJSON: `
{
    "id": 8,
    "name": "IMG0088.png",
    "size": 5563,
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2014-10-28T09:24:43Z"
}
`,
	Single: &backlog.Attachment{
		ID:   8,
		Name: "IMG0088.png",
		Size: 5563,
		CreatedUser: &backlog.User{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    backlog.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		Created: mustTime("2014-10-28T09:24:43Z"),
	},
	ListJSON: `
[
    {
        "id": 2,
        "name": "A.png",
        "size": 196186,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    },
    {
        "id": 5,
        "name": "B.png",
        "size": 201257,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    }
]
`,
	List: []*backlog.Attachment{
		{
			ID:   2,
			Name: "A.png",
			Size: 196186,
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2014-07-11T06:26:05Z"),
		},
		{
			ID:   5,
			Name: "B.png",
			Size: 201257,
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2014-07-11T06:26:05Z"),
		},
	},
	SingleListJSON: `
[
    {
        "id": 2,
        "name": "A.png",
        "size": 196186,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    }
]
`,
	SingleList: []*backlog.Attachment{
		{
			ID:   2,
			Name: "A.png",
			Size: 196186,
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTime("2014-07-11T06:26:05Z"),
		},
	},
	UploadJSON: `
{
    "id": 1,
    "name": "test.txt",
    "size": 8857
}
`,
	Upload: &backlog.Attachment{
		ID:   1,
		Name: "test.txt",
		Size: 8857,
	},
}
