package backlog_test

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
	testdataWikiListJSON = `
[
    {
        "id": 112,
        "projectId": 56,
        "name": "test1",
        "tags": [
            {
                "id": 12,
                "name": "proceedings"
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
        "created": "2013-05-30T09:11:36Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "updated": "2013-05-30T09:11:36Z"
    },
    {
        "id": 115,
        "projectId": 56,
        "name": "test2",
        "tags": [
            {
                "id": 12,
                "name": "proceedings"
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
        "created": "2013-05-30T09:11:36Z",
        "updatedUser": {
            "id": 1,
            "userId": "admin",
            "name": "admin",
            "roleType": 1,
            "lang": "ja",
            "mailAddress": "eguchi@nulab.example"
        },
        "updated": "2013-05-30T09:11:36Z"
    }
]
`
	testdataWikiMaximumJSON = `
{
    "id": 34,
    "projectId": 56,
    "name": "Maximum Wiki Page",
    "content": "This is a muximal wiki page.",
    "tags": [
        {
            "id": 12,
            "name": "proceedings"
        }
    ],
    "attachments": [
        {
            "id": 23,
            "name": "test.json",
            "size": 8857,
            "createdUser": {
                "id": 1,
                "userId": "admin",
                "name": "admin",
                "roleType": 1,
                "lang": "ja",
                "mailAddress": "eguchi@nulab.example"
            },
            "created": "2014-01-06T11:10:45Z"
        }
    ],
    "sharedFiles": [
        {
            "id": 454403,
            "type": "file",
            "dir": "/icon/",
            "name": "01_buz.png",
            "size": 2735,
            "createdUser": {
                "id": 5686,
                "userId": "takada",
                "name": "takada",
                "roleType": 2,
                "lang": "ja",
                "mailAddress": "takada@nulab.example"
            },
            "created": "2009-02-27T03:26:15Z",
            "updatedUser": {
                "id": 5686,
                "userId": "takada",
                "name": "takada",
                "roleType": 2,
                "lang": "ja",
                "mailAddress": "takada@nulab.example"
            },
            "updated": "2009-03-03T16:57:47Z"
        }
    ],
    "stars": [
        { 
            "id":75, 
            "comment":null, 
            "url": "https://xx.backlogtool.com/view/BLG-1", 
            "title": "[BLG-1] first issue | Show issue - Backlog", 
            "presenter":{ 
                "id":1, 
                "userId": "admin", 
                "name":"admin", 
                "roleType":1, 
                "lang":"ja", 
                "mailAddress":"eguchi@nulab.example" 
            }, 
            "created":"2014-01-23T10:55:19Z" 
        },
        { 
            "id":76, 
            "comment":"ok", 
            "url": "https://xx.backlogtool.com/view/BLG-1", 
            "title": "[BLG-1] first issue | Show issue - Backlog", 
            "presenter":{ 
                "id":1, 
                "userId": "admin", 
                "name":"admin", 
                "roleType":1, 
                "lang":"ja", 
                "mailAddress":"eguchi@nulab.example" 
            }, 
            "created":"2014-01-23T10:55:19Z" 
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
    "created": "2012-07-23T06:09:48Z",
    "updatedUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "updated": "2012-07-23T06:09:48Z"
}
`
	testdataWikiMinimumJSON = `
{
    "id": 34,
    "projectId": 56,
    "name": "Minimum Wiki Page",
    "content": "This is a minimal wiki page.",
    "tags": [
        {
            "id": 12,
            "name": "proceedings"
        }
    ],
    "attachments": [],
    "sharedFiles": [],
    "stars": [],
    "createdUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "created": "2012-07-23T06:09:48Z",
    "updatedUser": {
        "id": 1,
        "userId": "admin",
        "name": "admin",
        "roleType": 1,
        "lang": "ja",
        "mailAddress": "eguchi@nulab.example"
    },
    "updated": "2012-07-23T06:09:48Z"
}
`
)
