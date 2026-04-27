package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type projectFixtures struct {
	SingleJSON string
	Single     *backlog.Project
	ListJSON   string
	List       []*backlog.Project
}

type categoryFixtures struct {
	SingleJSON string
	Single     *backlog.Category
	ListJSON   string
	List       []*backlog.Category
}

// Project provides test fixtures for Project-related tests.
var Project = projectFixtures{
	SingleJSON: `
{
    "id": 6,
    "projectKey": "TEST",
    "name": "test",
    "chartEnabled": false,
    "subtaskingEnabled": false,
    "projectLeaderCanEditProjectLeader": false,
    "textFormattingRule": "markdown",
    "archived": false
}
`,
	Single: &backlog.Project{
		ID:                 6,
		ProjectKey:         "TEST",
		Name:               "test",
		TextFormattingRule: backlog.FormatMarkdown,
	},
	ListJSON: `
[
    {
        "id": 1,
        "projectKey": "TEST",
        "name": "test",
        "chartEnabled": false,
        "subtaskingEnabled": false,
        "projectLeaderCanEditProjectLeader": false,
        "textFormattingRule": "markdown",
        "archived": false
    },
    {
        "id": 2,
        "projectKey": "TEST2",
        "name": "test2",
        "chartEnabled": true,
        "subtaskingEnabled": false,
        "projectLeaderCanEditProjectLeader": true,
        "textFormattingRule": "markdown",
        "archived": false
    },
    {
        "id": 3,
        "projectKey": "TEST3",
        "name": "test3",
        "chartEnabled": false,
        "subtaskingEnabled": false,
        "projectLeaderCanEditProjectLeader": false,
        "textFormattingRule": "markdown",
        "archived": false
    }
]
`,
	List: []*backlog.Project{
		{
			ID:                 1,
			ProjectKey:         "TEST",
			Name:               "test",
			TextFormattingRule: backlog.FormatMarkdown,
		},
		{
			ID:                                2,
			ProjectKey:                        "TEST2",
			Name:                              "test2",
			ChartEnabled:                      true,
			ProjectLeaderCanEditProjectLeader: true,
			TextFormattingRule:                backlog.FormatMarkdown,
		},
		{
			ID:                 3,
			ProjectKey:         "TEST3",
			Name:               "test3",
			TextFormattingRule: backlog.FormatMarkdown,
		},
	},
}

// Category provides test fixtures for Category-related tests.
var Category = categoryFixtures{
	SingleJSON: `
{
    "id": 12,
    "name": "Bug",
    "displayOrder": 0
}
`,
	Single: &backlog.Category{
		ID:   12,
		Name: "Bug",
	},
	ListJSON: `
[
    {
        "id": 12,
        "name": "Bug",
        "displayOrder": 0
    },
    {
        "id": 13,
        "name": "Feature",
        "displayOrder": 1
    }
]
`,
	List: []*backlog.Category{
		{
			ID:   12,
			Name: "Bug",
		},
		{
			ID:           13,
			Name:         "Feature",
			DisplayOrder: 1,
		},
	},
}
