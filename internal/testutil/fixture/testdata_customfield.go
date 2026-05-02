package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type customFieldFixtures struct {
	SingleJSON string
	Single     *backlog.CustomField
	ListJSON   string
	List       []*backlog.CustomField
}

// CustomField provides test fixtures for CustomField-related tests.
var CustomField = customFieldFixtures{
	SingleJSON: `
{
    "id": 1,
    "typeId": 1,
    "name": "Sprint",
    "description": "",
    "required": false,
    "applicableIssueTypes": [],
    "allowAddItem": false,
    "items": []
}
`,
	Single: &backlog.CustomField{
		ID:     1,
		TypeID: 1,
		Name:   "Sprint",
	},
	ListJSON: `
[
    {
        "id": 1,
        "typeId": 1,
        "name": "Sprint",
        "description": "",
        "required": false,
        "applicableIssueTypes": [],
        "allowAddItem": false,
        "items": []
    },
    {
        "id": 2,
        "typeId": 5,
        "name": "Priority Label",
        "description": "",
        "required": true,
        "applicableIssueTypes": [],
        "allowAddItem": true,
        "items": [
            {"id": 10, "name": "High", "displayOrder": 0},
            {"id": 11, "name": "Low", "displayOrder": 1}
        ]
    }
]
`,
	List: []*backlog.CustomField{
		{
			ID:     1,
			TypeID: 1,
			Name:   "Sprint",
		},
		{
			ID:           2,
			TypeID:       5,
			Name:         "Priority Label",
			Required:     true,
			AllowAddItem: true,
			Items: []*backlog.CustomFieldItem{
				{ID: 10, Name: "High"},
				{ID: 11, Name: "Low", DisplayOrder: 1},
			},
		},
	},
}
