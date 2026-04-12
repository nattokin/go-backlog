package backlog

const (
	testdataActivityListJSON = `
[
    {
        "id": 3153,
        "project": {
            "id": 92,
            "projectKey": "SUB",
            "name": "Subtasking",
            "chartEnabled": true,
            "subtaskingEnabled": true,
            "projectLeaderCanEditProjectLeader": false,
            "textFormattingRule": null,
            "archived": false,
            "displayOrder": 0
        },
        "type": 2,
        "content": {
            "id": 4809,
            "key_id": 121,
            "summary": "Comment",
            "description": "",
            "comment": {
                "id": 7237,
                "content": ""
            },
            "changes": [
                {
                    "field": "milestone",
                    "new_value": "R2014-07-23",
                    "old_value": "",
                    "type": "standard"
                },
                {
                    "field": "status",
                    "new_value": "4",
                    "old_value": "1",
                    "type": "standard"
                }
            ]
        },
        "notifications": [
            {
                "id": 25,
                "alreadyRead": false,
                "reason": 2,
                "user": {
                    "id": 5686,
                    "userId": "takada",
                    "name": "takada",
                    "roleType": 2,
                    "lang": "ja",
                    "mailAddress": "takada@nulab.example"
                },
                "resourceAlreadyRead": false
            }
        ],
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-21T06:48:40Z"
    }
]
`

	testdataAttachmentSingleListJSON = `
[
    {
        "id": 2,
        "name": "A.png",
        "size": 196186,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    }
]
`

	testdataAttachmentListJSON = `
[
    {
        "id": 2,
        "name": "A.png",
        "size": 196186,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    },
    {
        "id": 5,
        "name": "B.png",
        "size": 201257,
        "createdUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": null,
            "mailAddress": "eguchi@nulab.example"
        },
        "created": "2014-07-11T06:26:05Z"
    }
]
`

	testdataAttachmentUploadJSON = `
{
    "id": 1,
    "name": "test.txt",
    "size": 8857
}
`

	testdataAttachmentJSON = `
{
    "id": 8,
    "name": "IMG0088.png",
    "size": 5563,
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2014-10-28T09:24:43Z"
}
`

	testdataInvalidJSON = `
{invalid}
`

	testdataProjectListJSON = `
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
`

	testdataProjectJSON = `
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
`

	testdataUserListJSON = `
[ 
    { 
        "id": 1, 
        "userId": "admin", 
        "name": "admin", 
        "roleType": 1, 
        "lang": "ja", 
        "mailAddress": "eguchi@nulab.example" 
    },
    { 
        "id": 2, 
        "userId": "normal_user", 
        "name": "normal_user", 
        "roleType": 2, 
        "lang": "ja", 
        "mailAddress": "sato@nulab.example" 
    },
    { 
        "id": 3, 
        "userId": "reporter", 
        "name": "reporter", 
        "roleType": 3, 
        "lang": "ja", 
        "mailAddress": "yamada@nulab.example" 
    },
    { 
        "id": 4, 
        "userId": "viewer", 
        "name": "viewer", 
        "roleType": 4, 
        "lang": "ja", 
        "mailAddress": "tanaka@nulab.example" 
    }
] 
`

	testdataUserJSON = `
{
    "id": 1,
    "userId": "admin",
    "name": "admin",
    "roleType": 1,
    "lang": "ja",
    "mailAddress": "eguchi@nulab.example"
}
`
)
