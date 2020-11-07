go-backlog
====
[![GoDoc](https://godoc.org/github.com/nattokin/go-backlog?status.svg)](https://godoc.org/github.com/nattokin/go-backlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/nattokin/go-backlog)](https://goreportcard.com/report/github.com/nattokin/go-backlog)
[![Test](https://github.com/nattokin/go-backlog/workflows/Test/badge.svg)](https://github.com/nattokin/go-backlog/actions?query=workflow%3ATest+branch%3Amaster)
[![codecov](https://codecov.io/gh/nattokin/go-backlog/branch/master/graph/badge.svg)](https://codecov.io/gh/nattokin/go-backlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

[Go](https://golang.org) client library for [Nulab Backlog API](https://developer.nulab.com/docs/backlog)

## Feature

- You can request each API endpoint using the Backlog API client created from the API base URL and token.
- Converts API response to a corresponding structure.
- Structures are provided for all endpoints and responses.

## Requirements

- Go >= 1.11

## Installation

```
go get github.com/nattokin/go-backlog
```

## Examples

### Get a wiki

```go
package main

import (
	"fmt"
	"log"

	"github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "BACKLOG_BASE_URL"
	// The tokun for request to Backlog API.
	token := "BACKLOG_TOKEN"

	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}

	r, err := c.Wiki.One(12345)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%#v\n", r)
}
```

### Get wikis list in the project

```go
package main

import (
	"fmt"
	"log"

	"github.com/nattokin/go-backlog"
)

func main() {
	// The base URL of Backlog API.
	baseURL := "BACKLOG_BASE_URL"
	// The tokun for request to Backlog API.
	token := "BACKLOG_TOKEN"
	// Create Backlog API client.
	c, err := backlog.NewClient(baseURL, token)
	if err != nil {
		log.Fatalln(err)
	}
	// ID or Key of the project.
	projectKey := "PROJECTKEY"
	// or
	// projectID := 1234

	r, err := c.Wiki.All(backlog.ProjectKey(projectKey))
	// r, err := c.Wiki.All(backlog.ProjectID(projectID))

	if err != nil {
		log.Fatalln(err)
	}
	for _, w := range r {
		fmt.Printf("%#v\n", w)
	}
}
```

## Supported API endpoints

### (*Client).Space

- [Get Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-recent-updates) - Returns recent updates in your space.

### (*Client).Attachment

- [Post Attachment File](https://developer.nulab-inc.com/docs/backlog/api/2/post-attachment-file/) - Posts an attachment file for issue or wiki. Returns id of the attachment file.

### (*Client).User

- [Get User List](https://developer.nulab.com/docs/backlog/api/2/get-user-list) - Returns list of users in your space.
- [Get User](https://developer.nulab.com/docs/backlog/api/2/get-user) - Returns information about user.
- [Add User](https://developer.nulab.com/docs/backlog/api/2/add-user) - Adds new user to the space. “Project Administrator” cannot add “Admin” user. You can’t use this API at `backlog.com` space.
- [Update User](https://developer.nulab.com/docs/backlog/api/2/update-user) - Updates information about user. You can’t use this API at backlog.com space.
- [Delete User](https://developer.nulab.com/docs/backlog/api/2/delete-user) - Deletes user from the space. You can’t use this API at backlog.com space.
- [Get Own User](https://developer.nulab.com/docs/backlog/api/2/get-own-user) - Returns own information about user.

### (*Client).User.Activity
- [Get User Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates) - Returns user’s recent updates.

### (*Client).Project

- [Get Project List](https://developer.nulab.com/docs/backlog/api/2/get-project-list) - Returns list of projects.
- [Add Project](https://developer.nulab.com/docs/backlog/api/2/add-project) - Adds new project.
- [Get Project](https://developer.nulab.com/docs/backlog/api/2/get-project) - Returns information about project.
- [Update Project](https://developer.nulab.com/docs/backlog/api/2/update-project) - Updates information about project.
- [Delete Project](https://developer.nulab.com/docs/backlog/api/2/delete-project) - Deletes project.

###  (*Client).Project.Activity

- [Get Project Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates) - Returns recent update in the project.

### (*Client).Project.User

- [Add Project User](https://developer.nulab.com/docs/backlog/api/2/add-project-user) - Adds user to list of project members.
- [Get Project User List](https://developer.nulab.com/docs/backlog/api/2/get-project-user-list) - Returns list of project members.
- [Delete Project User](https://developer.nulab.com/docs/backlog/api/2/delete-project-user) - Removes user from list project members.
- [Add Project Administrator](https://developer.nulab.com/docs/backlog/api/2/add-project-administrator) - Adds “Project Administrator” role to user.
- [Get List of Project Administrators](https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators) - Returns list of users who has Project Administrator role.
- [Delete Project Administrator](https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator) - Removes Project Administrator role from user.

### (*Client).Wiki

- [Get Wiki Page List](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page-list/) - Returns list of Wiki pages.
- [Get Wiki Page Tag List](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page-tag-list/) - Returns list of tags that are used in the project.
- [Count Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/count-wiki-page/) - Returns number of Wiki pages.
- [Get Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page/) - Returns information about Wiki page.
- [Add Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/add-wiki-page/) - Adds new Wiki page.
- [Delete Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/delete-wiki-page/) - Deletes Wiki page.

### (*Client).Wiki.Attachment

- [Get List of Wiki attachments](https://developer.nulab-inc.com/docs/backlog/api/2/get-list-of-wiki-attachments/) - Gets list of files attached to Wiki.
- [Attach File to Wiki](https://developer.nulab-inc.com/docs/backlog/api/2/attach-file-to-wiki/) - Attaches file to Wiki
- [Remove Wiki Attachment](https://developer.nulab-inc.com/docs/backlog/api/2/remove-wiki-attachment/) - Removes files attached to Wiki.

## License

The license of this project is [MIT license](https://opensource.org/licenses/MIT).
