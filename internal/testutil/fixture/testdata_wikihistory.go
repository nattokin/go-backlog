package fixture

import (
	"github.com/nattokin/go-backlog"
)

type wikiHistoryFixtures struct {
	ListJSON string
	List     []*backlog.WikiHistory
}

// WikiHistory provides test fixtures for WikiHistory-related tests.
var WikiHistory = wikiHistoryFixtures{
	ListJSON: `
[
    {
        "pageId": 34,
        "version": 2,
        "name": "Home",
        "content": "## Updated content",
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-06-10T09:00:00Z"
    },
    {
        "pageId": 34,
        "version": 1,
        "name": "Home",
        "content": "## Initial content",
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-06-01T09:00:00Z"
    }
]
`,
	List: []*backlog.WikiHistory{
		{
			PageID:  34,
			Version: 2,
			Name:    "Home",
			Content: "## Updated content",
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    1,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTimestamp("2014-06-10T09:00:00Z"),
		},
		{
			PageID:  34,
			Version: 1,
			Name:    "Home",
			Content: "## Initial content",
			CreatedUser: &backlog.User{
				ID:          1,
				UserID:      "admin",
				Name:        "admin",
				RoleType:    1,
				Lang:        "ja",
				MailAddress: "eguchi@nulab.example",
			},
			Created: mustTimestamp("2014-06-01T09:00:00Z"),
		},
	},
}
