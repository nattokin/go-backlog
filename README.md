go-backlog
====
[![Go Reference](https://pkg.go.dev/badge/github.com/nattokin/go-backlog.svg)](https://pkg.go.dev/github.com/nattokin/go-backlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/nattokin/go-backlog)](https://goreportcard.com/report/github.com/nattokin/go-backlog)
[![Test](https://github.com/nattokin/go-backlog/workflows/Test/badge.svg)](https://github.com/nattokin/go-backlog/actions?query=workflow%3ATest+branch%3Amain)
[![codecov](https://codecov.io/gh/nattokin/go-backlog/branch/main/graph/badge.svg)](https://codecov.io/gh/nattokin/go-backlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

[Go](https://golang.org) client library for [Nulab Backlog API](https://developer.nulab.com/docs/backlog)

## Feature

- You can request each API endpoint using the Backlog API client created from the API base URL and token.
- Converts API responses to idiomatic Go structs.
- Structs are provided for all API endpoints and responses.

## Requirements

- Go >= 1.23

## Installation

```
go get github.com/nattokin/go-backlog
```

## Examples

```go
c, err := backlog.NewClient(
    os.Getenv("BACKLOG_BASE_URL"),
    os.Getenv("BACKLOG_TOKEN"),
)

// Get all Wiki pages in a project.
wikis, err := c.Wiki.All(context.Background(), "MYPROJECT")

// Filter by keyword using an option.
wikis, err = c.Wiki.All(context.Background(), "MYPROJECT",
    c.Wiki.Option.WithKeyword("design"),
)
```

More examples can be found in the [examples/](examples/) directory and on [pkg.go.dev](https://pkg.go.dev/github.com/nattokin/go-backlog).

## Supported API endpoints

### Client.[Space](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceService)

- [Get Space](https://developer.nulab.com/docs/backlog/api/2/get-space) - Returns information about your space.
- [Get Space Disk Usage](https://developer.nulab.com/docs/backlog/api/2/get-space-disk-usage) - Returns disk usage of your space.
- [Get Space Notification](https://developer.nulab.com/docs/backlog/api/2/get-space-notification) - Returns the space notification.
- [Update Space Notification](https://developer.nulab.com/docs/backlog/api/2/update-space-notification) - Updates the space notification.

### Client.Space.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceActivityService)

- [Get Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-recent-updates) - Returns recent updates in the space.

### Client.Space.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceAttachmentService)

- [Post Attachment File](https://developer.nulab.com/docs/backlog/api/2/post-attachment-file/) - Posts an attachment file for issue or wiki, and returns its ID.

### Client.[User](https://pkg.go.dev/github.com/nattokin/go-backlog#UserService)

- [Get User List](https://developer.nulab.com/docs/backlog/api/2/get-user-list) - Returns a list of users in your space.
- [Get User](https://developer.nulab.com/docs/backlog/api/2/get-user) - Returns information about a specific user.
- [Add User](https://developer.nulab.com/docs/backlog/api/2/add-user) - Adds new user to the space. "Project Administrator" cannot add "Admin" user. You can't use this API at `backlog.com` space.
- [Update User](https://developer.nulab.com/docs/backlog/api/2/update-user) - Updates information about a user (Note: Not available at backlog.com).
- [Delete User](https://developer.nulab.com/docs/backlog/api/2/delete-user) - Deletes a user from the space (Note: Not available at backlog.com).
- [Get Own User](https://developer.nulab.com/docs/backlog/api/2/get-own-user) - Returns information about the currently authenticated user.

### Client.User.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#UserActivityService)
- [Get User Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates) - Returns a user's recent updates.

### Client.[Project](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectService)

- [Get Project List](https://developer.nulab.com/docs/backlog/api/2/get-project-list) - Returns a list of projects.
- [Add Project](https://developer.nulab.com/docs/backlog/api/2/add-project) - Adds a new project.
- [Get Project](https://developer.nulab.com/docs/backlog/api/2/get-project) - Returns information about a project.
- [Update Project](https://developer.nulab.com/docs/backlog/api/2/update-project) - Updates information about project.
- [Delete Project](https://developer.nulab.com/docs/backlog/api/2/delete-project) - Deletes a project.

### Client.Project.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectActivityService)

- [Get Project Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates) - Returns recent updates in the project.

### Client.Project.[User](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectUserService)

- [Add Project User](https://developer.nulab.com/docs/backlog/api/2/add-project-user) - Adds a user to the list of project members.
- [Get Project User List](https://developer.nulab.com/docs/backlog/api/2/get-project-user-list) - Returns a list of project members.
- [Delete Project User](https://developer.nulab.com/docs/backlog/api/2/delete-project-user) - Removes a user from the list of project members.
- [Add Project Administrator](https://developer.nulab.com/docs/backlog/api/2/add-project-administrator) - Adds the Project Administrator role to a user.
- [Get List of Project Administrators](https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators) - Returns a list of users with the Project Administrator role.
- [Delete Project Administrator](https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator) - Removes the Project Administrator role from a user.

### Client.[Issue](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueService)

- [Get Issue List](https://developer.nulab.com/docs/backlog/api/2/get-issue-list) - Returns a list of issues.
- [Count Issue](https://developer.nulab.com/docs/backlog/api/2/count-issue) - Returns the number of issues.
- [Get Issue](https://developer.nulab.com/docs/backlog/api/2/get-issue) - Returns information about a specific issue.
- [Add Issue](https://developer.nulab.com/docs/backlog/api/2/add-issue) - Adds a new issue.
- [Update Issue](https://developer.nulab.com/docs/backlog/api/2/update-issue) - Updates information about an issue.
- [Delete Issue](https://developer.nulab.com/docs/backlog/api/2/delete-issue) - Deletes an issue.

### Client.Issue.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueAttachmentService)

- [Get List of Issue Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments) - Returns a list of files attached to an issue.
- [Delete Issue Attachment](https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment) - Removes a file attached to an issue.

### Client.[PullRequest](https://pkg.go.dev/github.com/nattokin/go-backlog#PullRequestService)

- [Get Pull Request List](https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list) - Returns a list of pull requests.
- [Get Number of Pull Requests](https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests) - Returns the number of pull requests.
- [Get Pull Request](https://developer.nulab.com/docs/backlog/api/2/get-pull-request) - Returns information about a specific pull request.
- [Add Pull Request](https://developer.nulab.com/docs/backlog/api/2/add-pull-request) - Creates a new pull request.
- [Update Pull Request](https://developer.nulab.com/docs/backlog/api/2/update-pull-request) - Updates information about a pull request.

### Client.PullRequest.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#PullRequestAttachmentService)

- [Get List of Pull Request Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment) - Returns a list of files attached to a pull request.
- [Delete Pull Request Attachments](https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments) - Removes a file attached to a pull request.

### Client.[Wiki](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiService)

- [Get Wiki Page List](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list/) - Returns a list of Wiki pages.
- [Get Wiki Page Tag List](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-tag-list/) - Returns a list of tags used in the project.
- [Count Wiki Page](https://developer.nulab.com/docs/backlog/api/2/count-wiki-page/) - Returns the number of Wiki pages.
- [Get Wiki Page](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page/) - Returns information about a Wiki page.
- [Add Wiki Page](https://developer.nulab.com/docs/backlog/api/2/add-wiki-page/) - Adds a new Wiki page.
- [Delete Wiki Page](https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page/) - Deletes a Wiki page.

### Client.Wiki.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiAttachmentService)

- [Get List of Wiki Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments/) - Gets a list of files attached to a Wiki.
- [Attach File to Wiki](https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki/) - Attaches file to Wiki.
- [Remove Wiki Attachment](https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment/) - Removes files attached to a Wiki.

## License

The license of this project is [MIT license](https://opensource.org/licenses/MIT).
