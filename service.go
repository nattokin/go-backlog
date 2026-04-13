package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/attachment"
	"github.com/nattokin/go-backlog/internal/issue"
	"github.com/nattokin/go-backlog/internal/pullrequest"
	"github.com/nattokin/go-backlog/internal/space"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

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
