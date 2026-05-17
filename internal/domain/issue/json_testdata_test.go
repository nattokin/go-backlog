package issue_test

// issueLastPageJSON is a single-element array used to simulate
// the last page of a paginated response in All tests.
const issueLastPageJSON = `
[
    {
        "id": 3,
        "projectId": 10,
        "issueKey": "PRJ-3",
        "keyId": 3,
        "issueType": {
            "id": 2,
            "projectId": 10,
            "name": "Bug",
            "color": "#990000",
            "displayOrder": 0
        },
        "summary": "last issue",
        "description": "",
        "resolutions": null,
        "priority": {
            "id": 3,
            "name": "Normal"
        },
        "status": {
            "id": 1,
            "projectId": 10,
            "name": "Open",
            "color": "#ed8077",
            "displayOrder": 1000
        },
        "assignee": null,
        "category": [],
        "versions": [],
        "milestone": [],
        "startDate": null,
        "dueDate": null,
        "estimatedHours": null,
        "actualHours": null,
        "parentIssueId": null,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "admin@example.com"
        },
        "created": "2024-01-12T10:00:00Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "admin@example.com"
        },
        "updated": "2024-01-12T10:00:00Z",
        "customFields": [],
        "attachments": [],
        "sharedFiles": [],
        "stars": []
    }
]
`
