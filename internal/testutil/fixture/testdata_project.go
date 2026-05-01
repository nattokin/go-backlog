package fixture

import backlog "github.com/nattokin/go-backlog"

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

type statusFixtures struct {
	SingleJSON string
	Single     *backlog.Status
	ListJSON   string
	List       []*backlog.Status
}

// Project provides test fixtures for Project-related tests.
var Project = projectFixtures{
	SingleJSON: `
{
    "id": 1,
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
		ID:                                1,
		ProjectKey:                        "TEST",
		Name:                              "test",
		ChartEnabled:                      false,
		SubtaskingEnabled:                 false,
		ProjectLeaderCanEditProjectLeader: false,
		TextFormattingRule:                backlog.FormatMarkdown,
		Archived:                          false,
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
        "subtaskingEnabled": true,
        "projectLeaderCanEditProjectLeader": true,
        "textFormattingRule": "backlog",
        "archived": true
    }
]
`,
	List: []*backlog.Project{
		{
			ID:                                1,
			ProjectKey:                        "TEST",
			Name:                              "test",
			ChartEnabled:                      false,
			SubtaskingEnabled:                 false,
			ProjectLeaderCanEditProjectLeader: false,
			TextFormattingRule:                backlog.FormatMarkdown,
			Archived:                          false,
		},
		{
			ID:                                2,
			ProjectKey:                        "TEST2",
			Name:                              "test2",
			ChartEnabled:                      true,
			SubtaskingEnabled:                 true,
			ProjectLeaderCanEditProjectLeader: true,
			TextFormattingRule:                backlog.FormatBacklog,
			Archived:                          true,
		},
	},
}

// Category provides test fixtures for Category-related tests.
var Category = categoryFixtures{
	SingleJSON: `
{
    "id": 1,
    "name": "Bug",
    "displayOrder": 0
}
`,
	Single: &backlog.Category{
		ID:           1,
		Name:         "Bug",
		DisplayOrder: 0,
	},
	ListJSON: `
[
    {
        "id": 1,
        "name": "Bug",
        "displayOrder": 0
    },
    {
        "id": 2,
        "name": "Document",
        "displayOrder": 1
    }
]
`,
	List: []*backlog.Category{
		{
			ID:           1,
			Name:         "Bug",
			DisplayOrder: 0,
		},
		{
			ID:           2,
			Name:         "Document",
			DisplayOrder: 1,
		},
	},
}

// Status provides test fixtures for Status-related tests.
var Status = statusFixtures{
	SingleJSON: `
{
    "id": 1,
    "projectId": 6,
    "name": "Open",
    "color": "#ed8077",
    "displayOrder": 1000
}
`,
	Single: &backlog.Status{
		ID:           1,
		ProjectID:    6,
		Name:         "Open",
		Color:        "#ed8077",
		DisplayOrder: 1000,
	},
	ListJSON: `
[
    {
        "id": 1,
        "projectId": 6,
        "name": "Open",
        "color": "#ed8077",
        "displayOrder": 1000
    },
    {
        "id": 2,
        "projectId": 6,
        "name": "In Progress",
        "color": "#f5ab35",
        "displayOrder": 2000
    }
]
`,
	List: []*backlog.Status{
		{
			ID:           1,
			ProjectID:    6,
			Name:         "Open",
			Color:        "#ed8077",
			DisplayOrder: 1000,
		},
		{
			ID:           2,
			ProjectID:    6,
			Name:         "In Progress",
			Color:        "#f5ab35",
			DisplayOrder: 2000,
		},
	},
}
