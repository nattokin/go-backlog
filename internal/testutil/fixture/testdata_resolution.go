package fixture

import backlog "github.com/nattokin/go-backlog"

type resolutionFixtures struct {
	ListJSON string
	List     []*backlog.Resolution
}

// Resolution provides test fixtures for Resolution-related tests.
var Resolution = resolutionFixtures{
	ListJSON: `
[
    {
        "id": 0,
        "name": "Fixed"
    },
    {
        "id": 1,
        "name": "Won't Fix"
    },
    {
        "id": 2,
        "name": "Invalid"
    },
    {
        "id": 3,
        "name": "Duplication"
    },
    {
        "id": 4,
        "name": "Cannot Reproduce"
    }
]
`,
	List: []*backlog.Resolution{
		{ID: 0, Name: "Fixed"},
		{ID: 1, Name: "Won't Fix"},
		{ID: 2, Name: "Invalid"},
		{ID: 3, Name: "Duplication"},
		{ID: 4, Name: "Cannot Reproduce"},
	},
}
