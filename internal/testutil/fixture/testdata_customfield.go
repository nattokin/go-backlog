package fixture

import (
	backlog "github.com/nattokin/go-backlog"
)

type customFieldFixtures struct {
	SingleJSON string
	Single     *backlog.CustomField
	ListJSON   string
	List       []*backlog.CustomField

	NumberJSON   string
	Number       *backlog.CustomField
	DateJSON     string
	Date         *backlog.CustomField
	CheckBoxJSON string
	CheckBox     *backlog.CustomField
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
	NumberJSON: `
{
    "id": 3,
    "projectId": 1,
    "typeId": 3,
    "name": "Estimate",
    "description": "",
    "required": false,
    "useIssueType": true,
    "applicableIssueTypes": [],
    "displayOrder": 2,
    "version": 3,
    "min": 0.0,
    "max": 100.0,
    "initialValue": 8.0,
    "unit": "h"
}
`,
	Number: &backlog.CustomField{
		ID:                     3,
		TypeID:                 3,
		Name:                   "Estimate",
		ApplicableIssueTypeIDs: []int{},
		Items:                  []*backlog.CustomFieldItem{},

		Min:          ptrFloat64(0),
		Max:          ptrFloat64(100),
		InitialValue: ptrFloat64(8),
		Unit:         "h",

		ProjectID:    1,
		UseIssueType: true,
		DisplayOrder: 2,
		Version:      3,
	},
	DateJSON: `
{
    "id": 4,
    "projectId": 1,
    "typeId": 4,
    "name": "Deadline",
    "description": "",
    "required": false,
    "applicableIssueTypes": [],
    "min": "2024-01-14T15:00:00Z",
    "max": "2024-12-30T15:00:00Z",
    "initialDate": {
        "id": 2,
        "shift": 7
    }
}
`,
	Date: &backlog.CustomField{
		ID:                     4,
		TypeID:                 4,
		Name:                   "Deadline",
		ApplicableIssueTypeIDs: []int{},
		Items:                  []*backlog.CustomFieldItem{},

		MinDate:     mustTimestamp("2024-01-14T15:00:00Z"),
		MaxDate:     mustTimestamp("2024-12-30T15:00:00Z"),
		InitialDate: &backlog.CustomFieldInitialDate{InitialValueType: 2, Shift: 7},

		ProjectID: 1,
	},
	CheckBoxJSON: `
{
    "id": 5,
    "projectId": 1,
    "typeId": 7,
    "name": "Platform",
    "description": "",
    "required": false,
    "applicableIssueTypes": [],
    "allowAddItem": true,
    "allowInput": true,
    "items": [
        {"id": 20, "name": "Windows", "displayOrder": 0}
    ]
}
`,
	CheckBox: &backlog.CustomField{
		ID:                     5,
		TypeID:                 7,
		Name:                   "Platform",
		ApplicableIssueTypeIDs: []int{},
		AllowAddItem:           true,
		Items: []*backlog.CustomFieldItem{
			{ID: 20, Name: "Windows"},
		},

		AllowInput: true,

		ProjectID: 1,
	},
}
