go-backlog
====
[![Go Reference](https://pkg.go.dev/badge/github.com/nattokin/go-backlog.svg)](https://pkg.go.dev/github.com/nattokin/go-backlog)
[![Go Report Card](https://goreportcard.com/badge/github.com/nattokin/go-backlog)](https://goreportcard.com/report/github.com/nattokin/go-backlog)
[![Test](https://github.com/nattokin/go-backlog/workflows/Test/badge.svg)](https://github.com/nattokin/go-backlog/actions?query=workflow%3ATest+branch%3Amain)
[![codecov](https://codecov.io/gh/nattokin/go-backlog/branch/main/graph/badge.svg)](https://codecov.io/gh/nattokin/go-backlog)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

[Go](https://golang.org) client library for [Nulab Backlog API](https://developer.nulab.com/docs/backlog)

## Features

- **Type-safe option builders** — Filter and configure requests using strongly-typed option methods (e.g. `WithKeyword`, `WithCount`, `WithOrder`), avoiding raw string/map parameters.
- **Idiomatic Go structs** — API responses are mapped to Go structs with proper types (`time.Time`, typed constants, etc.) instead of raw JSON.
- **Context support** — Every API method accepts `context.Context` for cancellation and timeout control.
- **Structured error types** — Errors are returned as typed values (e.g. `*APIResponseError` for API errors, `*ValidationError` for invalid arguments), enabling precise handling with `errors.As`.

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
wikis, err := c.Wiki.List(context.Background(), "MYPROJECT")

// Filter by keyword using an option.
wikis, err = c.Wiki.List(context.Background(), "MYPROJECT",
    c.Wiki.Option.WithKeyword("design"),
)
```

More examples can be found in the [examples/](examples/) directory and on [pkg.go.dev](https://pkg.go.dev/github.com/nattokin/go-backlog).

## Supported API endpoints

### Client.[Issue](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueService)

- [Get Issue List](https://developer.nulab.com/docs/backlog/api/2/get-issue-list) - Returns a list of issues.
- [Count Issue](https://developer.nulab.com/docs/backlog/api/2/count-issue) - Returns the number of issues.
- [Get Issue](https://developer.nulab.com/docs/backlog/api/2/get-issue) - Returns information about a specific issue.
- [Add Issue](https://developer.nulab.com/docs/backlog/api/2/add-issue) - Adds a new issue.
- [Update Issue](https://developer.nulab.com/docs/backlog/api/2/update-issue) - Updates information about an issue.
- [Delete Issue](https://developer.nulab.com/docs/backlog/api/2/delete-issue) - Deletes an issue.
- [Get Issue Participants](https://developer.nulab.com/docs/backlog/api/2/get-issue-participant-list) - Returns a list of users participating in an issue.

### Client.Issue.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueAttachmentService)

- [Get List of Issue Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-issue-attachments) - Returns a list of files attached to an issue.
- [Get Issue Attachment](https://developer.nulab.com/docs/backlog/api/2/get-issue-attachment/) - Downloads a file attached to an issue.
- [Delete Issue Attachment](https://developer.nulab.com/docs/backlog/api/2/delete-issue-attachment) - Removes a file attached to an issue.

### Client.Issue.[Comment](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueCommentService)

- [Get Comment List](https://developer.nulab.com/docs/backlog/api/2/get-comment-list) - Returns a list of comments on an issue.
- [Add Comment](https://developer.nulab.com/docs/backlog/api/2/add-comment) - Adds a comment to an issue.
- [Count Comment](https://developer.nulab.com/docs/backlog/api/2/count-comment) - Returns the number of comments on an issue.
- [Get Comment](https://developer.nulab.com/docs/backlog/api/2/get-comment) - Returns information about a specific comment.
- [Update Comment](https://developer.nulab.com/docs/backlog/api/2/update-comment) - Updates a comment on an issue.
- [Delete Comment](https://developer.nulab.com/docs/backlog/api/2/delete-comment) - Deletes a comment from an issue.
- [Get List of Comment Notifications](https://developer.nulab.com/docs/backlog/api/2/get-list-of-comment-notifications) - Returns notifications for a comment.
- [Add Comment Notification](https://developer.nulab.com/docs/backlog/api/2/add-comment-notification) - Sends notifications for a comment.

### Client.Issue.[SharedFile](https://pkg.go.dev/github.com/nattokin/go-backlog#IssueSharedFileService)

- [Get List of Linked Shared Files](https://developer.nulab.com/docs/backlog/api/2/get-list-of-linked-shared-files) - Returns shared files linked to an issue.
- [Link Shared Files to Issue](https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-issue) - Links shared files to an issue.
- [Remove Link to Shared File from Issue](https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-issue) - Removes a shared file link from an issue.

### Client.[Project](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectService)

- [Get Project List](https://developer.nulab.com/docs/backlog/api/2/get-project-list) - Returns a list of projects.
- [Add Project](https://developer.nulab.com/docs/backlog/api/2/add-project) - Adds a new project.
- [Get Project](https://developer.nulab.com/docs/backlog/api/2/get-project) - Returns information about a project.
- [Update Project](https://developer.nulab.com/docs/backlog/api/2/update-project) - Updates information about project.
- [Delete Project](https://developer.nulab.com/docs/backlog/api/2/delete-project) - Deletes a project.
- [Get Project Disk Usage](https://developer.nulab.com/docs/backlog/api/2/get-project-disk-usage/) - Returns disk usage of a project.
- [Get Project Icon](https://developer.nulab.com/docs/backlog/api/2/get-project-icon/) - Returns the icon of a project.

### Client.Project.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectActivityService)

- [Get Project Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-project-recent-updates) - Returns recent updates in the project.

### Client.Project.[Category](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectCategoryService)

- [Get Category List](https://developer.nulab.com/docs/backlog/api/2/get-category-list/) - Returns a list of categories in a project.
- [Add Category](https://developer.nulab.com/docs/backlog/api/2/add-category/) - Adds a new category to a project.
- [Update Category](https://developer.nulab.com/docs/backlog/api/2/update-category/) - Updates a category in a project.
- [Delete Category](https://developer.nulab.com/docs/backlog/api/2/delete-category/) - Deletes a category from a project.

### Client.Project.[CustomField](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectCustomFieldService)

- [Get Custom Field List](https://developer.nulab.com/docs/backlog/api/2/get-custom-field-list/) - Returns a list of custom fields in a project.
- [Add Custom Field](https://developer.nulab.com/docs/backlog/api/2/add-custom-field/) - Adds a new custom field to a project.
- [Update Custom Field](https://developer.nulab.com/docs/backlog/api/2/update-custom-field/) - Updates a custom field in a project.
- [Delete Custom Field](https://developer.nulab.com/docs/backlog/api/2/delete-custom-field/) - Deletes a custom field from a project.
- [Add List Item for List Type Custom Field](https://developer.nulab.com/docs/backlog/api/2/add-list-item-for-list-type-custom-field/) - Adds a list item to a list type custom field.
- [Update List Item for List Type Custom Field](https://developer.nulab.com/docs/backlog/api/2/update-list-item-for-list-type-custom-field/) - Updates a list item in a list type custom field.
- [Delete List Item for List Type Custom Field](https://developer.nulab.com/docs/backlog/api/2/delete-list-item-for-list-type-custom-field/) - Deletes a list item from a list type custom field.

### Client.Project.[IssueType](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectIssueTypeService)

- [Get Issue Type List](https://developer.nulab.com/docs/backlog/api/2/get-issue-type-list/) - Returns a list of issue types in a project.
- [Add Issue Type](https://developer.nulab.com/docs/backlog/api/2/add-issue-type/) - Adds a new issue type to a project.
- [Update Issue Type](https://developer.nulab.com/docs/backlog/api/2/update-issue-type/) - Updates an issue type in a project.
- [Delete Issue Type](https://developer.nulab.com/docs/backlog/api/2/delete-issue-type/) - Deletes an issue type from a project.

### Client.Project.[Status](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectStatusService)

- [Get Status List of Project](https://developer.nulab.com/docs/backlog/api/2/get-status-list-of-project/) - Returns a list of statuses in a project.
- [Add Status](https://developer.nulab.com/docs/backlog/api/2/add-status/) - Adds a new status to a project.
- [Update Status](https://developer.nulab.com/docs/backlog/api/2/update-status/) - Updates a status in a project.
- [Delete Status](https://developer.nulab.com/docs/backlog/api/2/delete-status/) - Deletes a status from a project.
- [Update Order of Status](https://developer.nulab.com/docs/backlog/api/2/update-order-of-status/) - Updates the display order of statuses in a project.

### Client.Project.[SharedFile](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectSharedFileService)

- [Get List of Shared Files](https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files/) - Returns a list of shared files in a project.
- [Get File](https://developer.nulab.com/docs/backlog/api/2/get-file/) - Downloads a shared file.

### Client.Project.[User](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectUserService)

- [Add Project User](https://developer.nulab.com/docs/backlog/api/2/add-project-user) - Adds a user to the list of project members.
- [Get Project User List](https://developer.nulab.com/docs/backlog/api/2/get-project-user-list) - Returns a list of project members.
- [Delete Project User](https://developer.nulab.com/docs/backlog/api/2/delete-project-user) - Removes a user from the list of project members.
- [Add Project Administrator](https://developer.nulab.com/docs/backlog/api/2/add-project-administrator) - Adds the Project Administrator role to a user.
- [Get List of Project Administrators](https://developer.nulab.com/docs/backlog/api/2/get-list-of-project-administrators) - Returns a list of users with the Project Administrator role.
- [Delete Project Administrator](https://developer.nulab.com/docs/backlog/api/2/delete-project-administrator) - Removes the Project Administrator role from a user.

### Client.Project.[Webhook](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectWebhookService)

- [Get List of Webhooks](https://developer.nulab.com/docs/backlog/api/2/get-list-of-webhooks/) - Returns a list of webhooks in a project.
- [Get Webhook](https://developer.nulab.com/docs/backlog/api/2/get-webhook/) - Returns information about a specific webhook.
- [Add Webhook](https://developer.nulab.com/docs/backlog/api/2/add-webhook/) - Adds a new webhook to a project.
- [Update Webhook](https://developer.nulab.com/docs/backlog/api/2/update-webhook/) - Updates a webhook in a project.
- [Delete Webhook](https://developer.nulab.com/docs/backlog/api/2/delete-webhook/) - Deletes a webhook from a project.

### Client.Project.[Version](https://pkg.go.dev/github.com/nattokin/go-backlog#ProjectVersionService)

- [Get Version/Milestone List](https://developer.nulab.com/docs/backlog/api/2/get-version-milestone-list/) - Returns a list of versions/milestones in a project.
- [Add Version/Milestone](https://developer.nulab.com/docs/backlog/api/2/add-version-milestone/) - Adds a new version/milestone to a project.
- [Update Version/Milestone](https://developer.nulab.com/docs/backlog/api/2/update-version-milestone/) - Updates a version/milestone in a project.
- [Delete Version](https://developer.nulab.com/docs/backlog/api/2/delete-version/) - Deletes a version/milestone from a project.

### Client.[PullRequest](https://pkg.go.dev/github.com/nattokin/go-backlog#PullRequestService)

- [Get Pull Request List](https://developer.nulab.com/docs/backlog/api/2/get-pull-request-list) - Returns a list of pull requests.
- [Get Number of Pull Requests](https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-requests) - Returns the number of pull requests.
- [Get Pull Request](https://developer.nulab.com/docs/backlog/api/2/get-pull-request) - Returns information about a specific pull request.
- [Add Pull Request](https://developer.nulab.com/docs/backlog/api/2/add-pull-request) - Creates a new pull request.
- [Update Pull Request](https://developer.nulab.com/docs/backlog/api/2/update-pull-request) - Updates information about a pull request.

### Client.PullRequest.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#PullRequestAttachmentService)

- [Get List of Pull Request Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-pull-request-attachment) - Returns a list of files attached to a pull request.
- [Download Pull Request Attachment](https://developer.nulab.com/docs/backlog/api/2/download-pull-request-attachment/) - Downloads a file attached to a pull request.
- [Delete Pull Request Attachments](https://developer.nulab.com/docs/backlog/api/2/delete-pull-request-attachments) - Removes a file attached to a pull request.

### Client.PullRequest.[Comment](https://pkg.go.dev/github.com/nattokin/go-backlog#PullRequestCommentService)

- [Get Pull Request Comment](https://developer.nulab.com/docs/backlog/api/2/get-pull-request-comment/) - Returns a list of comments on a pull request.
- [Add Pull Request Comment](https://developer.nulab.com/docs/backlog/api/2/add-pull-request-comment/) - Adds a comment to a pull request.
- [Get Number of Pull Request Comments](https://developer.nulab.com/docs/backlog/api/2/get-number-of-pull-request-comments/) - Returns the number of comments on a pull request.
- [Update Pull Request Comment Information](https://developer.nulab.com/docs/backlog/api/2/update-pull-request-comment-information/) - Updates a comment on a pull request.

### Client.[RecentlyViewed](https://pkg.go.dev/github.com/nattokin/go-backlog#RecentlyViewedService)

- [Get Recently Viewed Issues](https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-issues) - Returns a list of recently viewed issues.
- [Get Recently Viewed Projects](https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-projects) - Returns a list of recently viewed projects.
- [Get Recently Viewed Wikis](https://developer.nulab.com/docs/backlog/api/2/get-list-of-recently-viewed-wikis/) - Returns a list of recently viewed wiki pages.
- [Add Recently Viewed Issue](https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-issue/) - Adds an issue to recently viewed list.
- [Add Recently Viewed Wiki](https://developer.nulab.com/docs/backlog/api/2/add-recently-viewed-wiki/) - Adds a wiki to recently viewed list.

### Client.[Repository](https://pkg.go.dev/github.com/nattokin/go-backlog#RepositoryService)

- [Get List of Git Repositories](https://developer.nulab.com/docs/backlog/api/2/get-list-of-git-repositories/) - Returns a list of Git repositories in a project.
- [Get Git Repository](https://developer.nulab.com/docs/backlog/api/2/get-git-repository/) - Returns information about a specific Git repository.

### Client.[Space](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceService)

- [Get Space](https://developer.nulab.com/docs/backlog/api/2/get-space) - Returns information about your space.
- [Get Space Disk Usage](https://developer.nulab.com/docs/backlog/api/2/get-space-disk-usage) - Returns disk usage of your space.
- [Get Space Notification](https://developer.nulab.com/docs/backlog/api/2/get-space-notification) - Returns the space notification.
- [Update Space Notification](https://developer.nulab.com/docs/backlog/api/2/update-space-notification) - Updates the space notification.

### Client.Space.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceActivityService)

- [Get Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-recent-updates) - Returns recent updates in the space.
- [Get Activity](https://developer.nulab.com/docs/backlog/api/2/get-activity/) - Returns a specific activity in the space.

### Client.Space.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#SpaceAttachmentService)

- [Post Attachment File](https://developer.nulab.com/docs/backlog/api/2/post-attachment-file/) - Posts an attachment file for issue or wiki, and returns its ID.

### Client.[Star](https://pkg.go.dev/github.com/nattokin/go-backlog#StarService)

- [Add Star](https://developer.nulab.com/docs/backlog/api/2/add-star) - Adds a star to a resource.
- [Remove Star](https://developer.nulab.com/docs/backlog/api/2/remove-star) - Removes a star by its ID.

### Client.[User](https://pkg.go.dev/github.com/nattokin/go-backlog#UserService)

- [Get User List](https://developer.nulab.com/docs/backlog/api/2/get-user-list) - Returns a list of users in your space.
- [Get User](https://developer.nulab.com/docs/backlog/api/2/get-user) - Returns information about a specific user.
- [Add User](https://developer.nulab.com/docs/backlog/api/2/add-user) - Adds new user to the space. "Project Administrator" cannot add "Admin" user. You can't use this API at `backlog.com` space.
- [Update User](https://developer.nulab.com/docs/backlog/api/2/update-user) - Updates information about a user (Note: Not available at backlog.com).
- [Delete User](https://developer.nulab.com/docs/backlog/api/2/delete-user) - Deletes a user from the space (Note: Not available at backlog.com).
- [Get Own User](https://developer.nulab.com/docs/backlog/api/2/get-own-user) - Returns information about the currently authenticated user.
- [Get User Icon](https://developer.nulab.com/docs/backlog/api/2/get-user-icon/) - Returns the icon of a user.

### Client.User.[Activity](https://pkg.go.dev/github.com/nattokin/go-backlog#UserActivityService)

- [Get User Recent Updates](https://developer.nulab.com/docs/backlog/api/2/get-user-recent-updates) - Returns a user's recent updates.

### Client.User.[Star](https://pkg.go.dev/github.com/nattokin/go-backlog#UserStarService)

- [Get Received Star List](https://developer.nulab.com/docs/backlog/api/2/get-received-star-list) - Returns a list of stars received by a user.
- [Count User Received Stars](https://developer.nulab.com/docs/backlog/api/2/count-user-received-stars) - Returns the number of stars received by a user.

### Client.[Wiki](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiService)

- [Get Wiki Page List](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-list/) - Returns a list of Wiki pages.
- [Get Wiki Page Tag List](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-tag-list/) - Returns a list of tags used in the project.
- [Count Wiki Page](https://developer.nulab.com/docs/backlog/api/2/count-wiki-page/) - Returns the number of Wiki pages.
- [Get Wiki Page](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page/) - Returns information about a Wiki page.
- [Add Wiki Page](https://developer.nulab.com/docs/backlog/api/2/add-wiki-page/) - Adds a new Wiki page.
- [Update Wiki Page](https://developer.nulab.com/docs/backlog/api/2/update-wiki-page/) - Updates a Wiki page.
- [Delete Wiki Page](https://developer.nulab.com/docs/backlog/api/2/delete-wiki-page/) - Deletes a Wiki page.

### Client.Wiki.[Attachment](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiAttachmentService)

- [Get List of Wiki Attachments](https://developer.nulab.com/docs/backlog/api/2/get-list-of-wiki-attachments/) - Gets a list of files attached to a Wiki.
- [Get Wiki Page Attachment](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-attachment/) - Downloads a file attached to a Wiki page.
- [Attach File to Wiki](https://developer.nulab.com/docs/backlog/api/2/attach-file-to-wiki/) - Attaches file to Wiki.
- [Remove Wiki Attachment](https://developer.nulab.com/docs/backlog/api/2/remove-wiki-attachment/) - Removes files attached to a Wiki.

### Client.Wiki.[History](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiHistoryService)

- [Get Wiki Page History](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-history/) - Returns version history of a wiki page.

### Client.Wiki.[SharedFile](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiSharedFileService)

- [Get List of Shared Files on Wiki](https://developer.nulab.com/docs/backlog/api/2/get-list-of-shared-files-on-wiki) - Returns shared files linked to a wiki page.
- [Link Shared Files to Wiki](https://developer.nulab.com/docs/backlog/api/2/link-shared-files-to-wiki) - Links shared files to a wiki page.
- [Remove Link to Shared File from Wiki](https://developer.nulab.com/docs/backlog/api/2/remove-link-to-shared-file-from-wiki) - Removes a shared file link from a wiki page.

### Client.Wiki.[Star](https://pkg.go.dev/github.com/nattokin/go-backlog#WikiStarService)

- [Get Wiki Page Star](https://developer.nulab.com/docs/backlog/api/2/get-wiki-page-star) - Returns stars on a wiki page.

## License

The license of this project is [MIT license](https://opensource.org/licenses/MIT).
