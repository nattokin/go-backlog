package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type userFixtures struct {
	SingleJSON string
	Single     *backlog.User
	ListJSON   string
	List       []*backlog.User
}

// User provides test fixtures for User-related tests.
var User = userFixtures{
	SingleJSON: `
{
    "id": 1,
    "userId": "admin",
    "name": "admin",
    "roleType": 1,
    "lang": "ja",
    "mailAddress": "eguchi@nulab.example"
}
`,
	Single: &backlog.User{
		ID:          1,
		UserID:      "admin",
		Name:        "admin",
		RoleType:    backlog.RoleAdministrator,
		Lang:        "ja",
		MailAddress: "eguchi@nulab.example",
	},
	ListJSON: `
[ 
    { 
        "id": 1, 
        "userId": "admin", 
        "name": "admin", 
        "roleType": 1, 
        "lang": "ja", 
        "mailAddress": "eguchi@nulab.example" 
    },
    { 
        "id": 2, 
        "userId": "normal_user", 
        "name": "normal_user", 
        "roleType": 2, 
        "lang": "ja", 
        "mailAddress": "sato@nulab.example" 
    },
    { 
        "id": 3, 
        "userId": "reporter", 
        "name": "reporter", 
        "roleType": 3, 
        "lang": "ja", 
        "mailAddress": "yamada@nulab.example" 
    },
    { 
        "id": 4, 
        "userId": "viewer", 
        "name": "viewer", 
        "roleType": 4, 
        "lang": "ja", 
        "mailAddress": "tanaka@nulab.example" 
    }
] 
`,
	List: []*backlog.User{
		{
			ID:          1,
			UserID:      "admin",
			Name:        "admin",
			RoleType:    backlog.RoleAdministrator,
			Lang:        "ja",
			MailAddress: "eguchi@nulab.example",
		},
		{
			ID:          2,
			UserID:      "normal_user",
			Name:        "normal_user",
			RoleType:    backlog.RoleNormalUser,
			Lang:        "ja",
			MailAddress: "sato@nulab.example",
		},
		{
			ID:          3,
			UserID:      "reporter",
			Name:        "reporter",
			RoleType:    backlog.RoleReporter,
			Lang:        "ja",
			MailAddress: "yamada@nulab.example",
		},
		{
			ID:          4,
			UserID:      "viewer",
			Name:        "viewer",
			RoleType:    backlog.RoleViewer,
			Lang:        "ja",
			MailAddress: "tanaka@nulab.example",
		},
	},
}
