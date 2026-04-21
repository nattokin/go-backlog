package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type spaceFixtures struct {
	SpaceJSON        string
	Space            *backlog.Space
	DiskUsageJSON    string
	DiskUsage        *backlog.DiskUsageSpace
	NotificationJSON string
	Notification     *backlog.SpaceNotification
}

// Space provides test fixtures for Space-related tests.
var Space = spaceFixtures{
	SpaceJSON: `
{
    "spaceKey": "nulab",
    "name": "Nulab Inc.",
    "ownerId": 1,
    "lang": "ja",
    "timezone": "Asia/Tokyo",
    "reportSendTime": "08:00:00",
    "textFormattingRule": "markdown",
    "created": "2008-07-06T15:00:00Z",
    "updated": "2013-06-18T07:55:37Z"
}
`,
	Space: &backlog.Space{
		SpaceKey:           "nulab",
		Name:               "Nulab Inc.",
		OwnerID:            1,
		Lang:               "ja",
		Timezone:           "Asia/Tokyo",
		ReportSendTime:     "08:00:00",
		TextFormattingRule: backlog.FormatMarkdown,
		Created:            mustTime("2008-07-06T15:00:00Z"),
		Updated:            mustTime("2013-06-18T07:55:37Z"),
	},
	DiskUsageJSON: `
{
    "capacity": 1073741824,
    "issue": 119511,
    "wiki": 0,
    "file": 0,
    "subversion": 0,
    "git": 0,
    "gitLFS": 0,
    "details": [
        {
            "projectId": 1,
            "issue": 11931,
            "wiki": 0,
            "file": 0,
            "subversion": 0,
            "git": 0,
            "gitLFS": 0
        }
    ]
}
`,
	DiskUsage: &backlog.DiskUsageSpace{
		Capacity:   1073741824,
		Issue:      119511,
		Wiki:       0,
		File:       0,
		Subversion: 0,
		Git:        0,
		GitLFS:     0,
		Details: []*backlog.DiskUsageProject{
			{
				ProjectID:  1,
				Issue:      11931,
				Wiki:       0,
				File:       0,
				Subversion: 0,
				Git:        0,
				GitLFS:     0,
			},
		},
	},
	NotificationJSON: `
{
    "content": "Backlog is a project management tool.",
    "updated": "2013-06-18T07:55:37Z"
}
`,
	Notification: &backlog.SpaceNotification{
		Content: "Backlog is a project management tool.",
		Updated: mustTime("2013-06-18T07:55:37Z"),
	},
}
