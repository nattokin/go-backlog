package fixture

import backlog "github.com/nattokin/go-backlog"

type priorityFixtures struct {
	ListJSON string
	List     []*backlog.Priority
}

// Priority provides test fixtures for Priority-related tests.
var Priority = priorityFixtures{
	ListJSON: `
[
    {
        "id": 2,
        "name": "High"
    },
    {
        "id": 3,
        "name": "Normal"
    },
    {
        "id": 4,
        "name": "Low"
    }
]
`,
	List: []*backlog.Priority{
		{ID: 2, Name: "High"},
		{ID: 3, Name: "Normal"},
		{ID: 4, Name: "Low"},
	},
}
