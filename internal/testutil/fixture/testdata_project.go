package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type projectFixtures struct {
	SingleJSON    string
	Single        *backlog.Project
	ListJSON      string
	List          []*backlog.Project
	DiskUsageJSON string
	DiskUsage     *backlog.DiskUsageProject
}

type categoryFixtures struct {
	SingleJSON string
	Single     *backlog.Category
	ListJSON   string
	List       []*backlog.Category
}

type issueTypeFixtures struct {
	SingleJSON string
	Single     *backlog.IssueType
	ListJSON   string
	List       []*backlog.IssueType
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
    "id": 6,
    "projectKey": "TEST",
    "name": "test",
    "chartEnabled": false,
    "subtaskingEnabled": false,
    "projectLeaderCanEditProjectLeader": false,
    "textFormattingRule": "markdown",
    "archived": false,
    "useResolvedForChart": true,
    "useWiki": true,
    "useFileSharing": true,
    "useWikiTreeView": true,
    "useSubversion": false,
    "useGit": true,
    "useOriginalImageSizeAtWiki": false,
    "displayOrder": 2147483646,
    "useDevAttributes": true
}
`,
	Single: &backlog.Project{
		ID:                  6,
		ProjectKey:          "TEST",
		Name:                "test",
		TextFormattingRule:  backlog.FormatMarkdown,
		UseResolvedForChart: true,
		UseWiki:             true,
		UseFileSharing:      true,
		UseWikiTreeView:     true,
		UseGit:              true,
		DisplayOrder:        2147483646,
		UseDevAttributes:    true,
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
	DiskUsageJSON: `
{
	"projectId": 1,
	"issue": 11931,
	"wiki": 0,
	"file": 0,
	"subversion": 0,
	"git": 0,
	"gitLFS": 0
}
`,
	DiskUsage: &backlog.DiskUsageProject{
		ProjectID:  1,
		Issue:      11931,
		Wiki:       0,
		File:       0,
		Subversion: 0,
		Git:        0,
		GitLFS:     0,
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

// IssueType provides test fixtures for IssueType-related tests.
var IssueType = issueTypeFixtures{
	SingleJSON: `
{
    "id": 1,
    "projectId": 6,
    "name": "Bug",
    "color": "#e30000",
    "displayOrder": 0
}
`,
	Single: &backlog.IssueType{
		ID:        1,
		ProjectID: 6,
		Name:      "Bug",
		Color:     "#e30000",
	},
	ListJSON: `
[
    {
        "id": 1,
        "projectId": 6,
        "name": "Bug",
        "color": "#e30000",
        "displayOrder": 0
    },
    {
        "id": 2,
        "projectId": 6,
        "name": "Task",
        "color": "#7ea800",
        "displayOrder": 1
    }
]
`,
	List: []*backlog.IssueType{
		{
			ID:        1,
			ProjectID: 6,
			Name:      "Bug",
			Color:     "#e30000",
		},
		{
			ID:           2,
			ProjectID:    6,
			Name:         "Task",
			Color:        "#7ea800",
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
