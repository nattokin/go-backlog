package backlog

import (
	"github.com/nattokin/go-backlog/internal/activity"
	"github.com/nattokin/go-backlog/internal/core"
	"github.com/nattokin/go-backlog/internal/project"
	"github.com/nattokin/go-backlog/internal/user"
	"github.com/nattokin/go-backlog/internal/wiki"
)

type apiParamOptionType = core.APIParamOptionType

type RequestOption = core.RequestOption

type apiParamOption = core.APIParamOption

type ActivityOptionService = activity.ActivityOptionService

type ProjectOptionService = project.ProjectOptionService

type UserOptionService = user.UserOptionService

type WikiOptionService = wiki.WikiOptionService
