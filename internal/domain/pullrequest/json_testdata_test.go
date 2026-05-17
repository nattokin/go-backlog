package pullrequest_test

// pullRequestLastPageJSON is a single-element array used to simulate
// the last page of a paginated response in All tests.
const pullRequestLastPageJSON = `
[
    {
        "id": 4,
        "projectId": 3,
        "repositoryId": 5,
        "number": 3,
        "summary": "last PR",
        "description": "",
        "base": "main",
        "branch": "feature/baz",
        "status": {
            "id": 1,
            "name": "Open"
        },
        "assignee": null,
        "issue": null,
        "baseCommit": null,
        "branchCommit": null,
        "mergeCommit": null,
        "closeAt": null,
        "mergeAt": null,
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
        "attachments": [],
        "stars": []
    }
]
`
