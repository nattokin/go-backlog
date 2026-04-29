package fixture

import "github.com/nattokin/go-backlog"

type versionFixtures struct {
	SingleJSON string
	Single     *backlog.Version
	ListJSON   string
	List       []*backlog.Version
}

var Version = versionFixtures{
	SingleJSON: `{
		"id": 1,
		"projectId": 100,
		"name": "Version 1.0",
		"description": "Initial release milestone",
		"startDate": null,
		"releaseDueDate": null,
		"archived": false,
		"displayOrder": 1
	}`,
	Single: &backlog.Version{
		ID:           1,
		ProjectID:    100,
		Name:         "Version 1.0",
		Description:  "Initial release milestone",
		Archived:     false,
		DisplayOrder: 1,
	},
	ListJSON: `[
		{
			"id": 1,
			"projectId": 100,
			"name": "Version 1.0",
			"description": "Initial release milestone",
			"startDate": null,
			"releaseDueDate": null,
			"archived": false,
			"displayOrder": 1
		},
		{
			"id": 2,
			"projectId": 100,
			"name": "Version 2.0",
			"description": "Second release milestone",
			"startDate": null,
			"releaseDueDate": null,
			"archived": false,
			"displayOrder": 2
		}
	]`,
	List: []*backlog.Version{
		{
			ID:           1,
			ProjectID:    100,
			Name:         "Version 1.0",
			Description:  "Initial release milestone",
			Archived:     false,
			DisplayOrder: 1,
		},
		{
			ID:           2,
			ProjectID:    100,
			Name:         "Version 2.0",
			Description:  "Second release milestone",
			Archived:     false,
			DisplayOrder: 2,
		},
	},
}
