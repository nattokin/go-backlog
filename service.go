package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/pullrequest"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

// ──────────────────────────────────────────────────────────────
//  Implemented services (aliases to internal packages)
// ──────────────────────────────────────────────────────────────

type IssueAttachmentService = attachment.IssueAttachmentService

type IssueService = issue.IssueService

type PullRequestAttachmentService = attachment.PullRequestAttachmentService

type PullRequestService = pullrequest.PullRequestService

type SpaceActivityService = activity.SpaceActivityService

type SpaceAttachmentService = attachment.SpaceAttachmentService

type SpaceService = space.SpaceService

type UserActivityService = activity.UserActivityService

type UserService = user.UserService

type WikiAttachmentService = attachment.WikiAttachmentService

type WikiService = wiki.WikiService

// ──────────────────────────────────────────────────────────────
//  Unimplemented service stubs (to be extracted to internal/*)
// ──────────────────────────────────────────────────────────────

// CategoryService handles communication with the category-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type CategoryService struct {
	method *core.Method
}

// CustomFieldService handles communication with the custom field-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type CustomFieldService struct {
	method *core.Method
}

// PriorityService handles communication with the priority-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type PriorityService struct {
	method *core.Method
}

// RepositoryService handles communication with the repository-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type RepositoryService struct {
	method *core.Method
}

// ResolutionService handles communication with the resolution-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type ResolutionService struct {
	method *core.Method
}

// StatusService handles communication with the status-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type StatusService struct {
	method *core.Method
}

// VersionService handles communication with the version-related methods of the Backlog API.
//
//nolint:unused // API not implemented yet
type VersionService struct {
	method *core.Method
}
