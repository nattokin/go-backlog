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
- Converts API responses to idiomatic Go structs.
- Structs are provided for all API endpoints and responses.

## Requirements

- Go >= 1.14

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
    // The base URL for the Backlog API.
    baseURL := "BACKLOG_BASE_URL"
    // The token for requests to the Backlog API.
    token := "BACKLOG_TOKEN"

    // Create Backlog API client.
    c, err := backlog.NewClient(baseURL, token)
    if err != nil {
        log.Fatalln(err)
    }

    // The wiki ID.
    wikiID := 12345
    r, err := c.Wiki.One(wikiID)
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
    // The base URL for the Backlog API.
    baseURL := "BACKLOG_BASE_URL"
    // The token for requests to the Backlog API.
    token := "BACKLOG_TOKEN"

    // Create Backlog API client.
    c, err := backlog.NewClient(baseURL, token)
    if err != nil {
        log.Fatalln(err)
    }

    // The project ID or Key.
    projectIDOrKey := "PROJECTKEY"
    r, err := c.Wiki.All(projectIDOrKey)

    if err != nil {
        log.Fatalln(err)
    }
    for _, w := range r {
        fmt.Printf("%#v\n", w)
    }
}
```

## Supported API endpoints

### Client.Space.[Activity](https://godoc.org/github.com/nattokin/go-backlog#SpaceActivityService)

- [Get Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-recent-updates) - Returns recent updates in the space.

### Client.Space.[Attachment](https://godoc.org/github.com/nattokin/go-backlog#SpaceAttachmentService)

- [Post Attachment File](https://developer.nulab-inc.com/docs/backlog/api/2/post-attachment-file/) - Posts an attachment file for issue or wiki, and returns its ID.

### Client.[User](https://godoc.org/github.com/nattokin/go-backlog#UserService)

- [Get User List](https://developer.nulab.com/docs/backlog/api/2/get-user-list) - Returns a list of users in your space.
- [Get User](https://developer.nulab.com/docs/backlog/api/2/get-user) - Returns information about a specific user.
- [Add User](https://developer.nulab.com/docs/backlog/api/2/add-user) - Adds new user to the space. “Project Administrator” cannot add “Admin” user. You can’t use this API at `backlog.com` space.
- [Update User](https://developer.nulab.com/docs/backlog/api/2/update-user) - Updates information about a user (Note: Not available at backlog.com).
- [Delete User](https://developer.nulab.com/docs/backlog/api/2/delete-user) - Deletes a user from the space (Note: Not available at backlog.com).
- [Get Own User](https://developer.nulab.com/docs/backlog/api/2/get-own-user) - Returns information about the currently authenticated user.

### Client.User.[Activity](https://godoc.org/github.com/nattokin/go-backlog#UserActivityService)
- [Get User Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates) - Returns a user’s recent updates.

### Client.[Project](https://godoc.org/github.com/nattokin/go-backlog#ProjectService)

- [Get Project List](https://developer.nulab.com/docs/backlog/api/2/get-project-list) - Returns a list of projects.
- [Add Project](https://developer.nulab.com/docs/backlog/api/2/add-project) - Adds a new project.
- [Get Project](https://developer.nulab.com/docs/backlog/api/2/get-project) - Returns information about a project.
- [Update Project](https://developer.nulab.com/docs/backlog/api/2/update-project) - Updates information about project.
- [Delete Project](https://developer.nulab.com/docs/backlog/api/2/delete-project) - Deletes a project.

### Client.Project.[Activity](https://godoc.org/github.com/nattokin/go-backlog#ProjectActivityService)

- [Get Project Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates) - Returns recent updates in the project.

### Client.Project.[User](https://godoc.org/github.com/nattokin/go-backlog#ProjectUserService)

- [Add Project User](https://developer.nulab.com/docs/backlog/api/2/add-project-user) - Adds a user to the list of project members.
- [Get Project User List](https://developer.nulab.com/docs/backlog/api/2/get-project-user-list) - Returns a list of project members.
- [Delete Project User](https://developer.nulab.com/docs/backlog/api/2/delete-project-user) - Removes a user from the list of project members.
- [Add Project Administrator](https://developer.nulab.com/docs/backlog/api/2/add-project-administrator) - Adds the Project Administrator role to a user.
- [Get List of Project Administrators](https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators) - Returns a list of users with the Project Administrator role.
- [Delete Project Administrator](https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator) - Removes the Project Administrator role from a user.

### Client.[Wiki](https://godoc.org/github.com/nattokin/go-backlog#WikiService)

- [Get Wiki Page List](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page-list/) - Returns a list of Wiki pages.
- [Get Wiki Page Tag List](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page-tag-list/) - Returns a list of tags used in the project.
- [Count Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/count-wiki-page/) - Returns the number of Wiki pages.
- [Get Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/get-wiki-page/) - Returns information about a Wiki page.
- [Add Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/add-wiki-page/) - Adds a new Wiki page.
- [Delete Wiki Page](https://developer.nulab-inc.com/docs/backlog/api/2/delete-wiki-page/) - Deletes a Wiki page.

### Client.Wiki.[Attachment](https://godoc.org/github.com/nattokin/go-backlog#WikiAttachmentService)

- [Get List of Wiki attachments](https://developer.nulab-inc.com/docs/backlog/api/2/get-list-of-wiki-attachments/) - Gets a list of files attached to a Wiki.
- [Attach File to Wiki](https://developer.nulab-inc.com/docs/backlog/api/2/attach-file-to-wiki/) - Attaches file to Wiki
- [Remove Wiki Attachment](https://developer.nulab-inc.com/docs/backlog/api/2/remove-wiki-attachment/) - Removes files attached to a Wiki.

## License

The license of this project is [MIT license](https://opensource.org/licenses/MIT).
