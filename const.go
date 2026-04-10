package backlog

import "github.com/nattokin/go-backlog/internal/model"

// Order defines the sort order (ascending or descending).
const (
	OrderAsc  = model.OrderAsc
	OrderDesc = model.OrderDesc
)

// Format defines the text formatting rule for the Backlog wiki.
const (
	FormatMarkdown = model.FormatMarkdown
	FormatBacklog  = model.FormatBacklog
)

// Role defines the type of user role within a project.
const (
	RoleAdministrator = model.RoleAdministrator
	RoleNormalUser    = model.RoleNormalUser
	RoleReporter      = model.RoleReporter
	RoleViewer        = model.RoleViewer
	RoleGuestReporter = model.RoleGuestReporter
	RoleGuestViewer   = model.RoleGuestViewer
)

const (
	paramActivityTypeIDs                   apiParamOptionType = "activityTypeId[]"
	paramAll                               apiParamOptionType = "all"
	paramArchived                          apiParamOptionType = "archived"
	paramChartEnabled                      apiParamOptionType = "chartEnabled"
	paramContent                           apiParamOptionType = "content"
	paramCount                             apiParamOptionType = "count"
	paramKey                               apiParamOptionType = "key"
	paramKeyword                           apiParamOptionType = "keyword"
	paramMailAddress                       apiParamOptionType = "mailAddress"
	paramMailNotify                        apiParamOptionType = "mailNotify"
	paramMaxID                             apiParamOptionType = "maxId"
	paramMinID                             apiParamOptionType = "minId"
	paramName                              apiParamOptionType = "name"
	paramOrder                             apiParamOptionType = "order"
	paramPassword                          apiParamOptionType = "password"
	paramProjectLeaderCanEditProjectLeader apiParamOptionType = "projectLeaderCanEditProjectLeader"
	paramRoleType                          apiParamOptionType = "roleType"
	paramSendMail                          apiParamOptionType = "sendMail"
	paramSubtaskingEnabled                 apiParamOptionType = "subtaskingEnabled"
	paramTextFormattingRule                apiParamOptionType = "textFormattingRule"
	paramUserID                            apiParamOptionType = "userId"
)

const maxActivityTypeID = 26
