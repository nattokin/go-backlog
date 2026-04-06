package backlog

const (
	apiVersion = "v2"
)

// Order defines the sort order (ascending or descending).
const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Format defines the text formatting rule for the Backlog wiki.
const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
)

// Role defines the type of user role within a project.
const (
	_ Role = iota
	RoleAdministrator
	RoleNormalUser
	RoleReporter
	RoleViewer
	RoleGuestReporter
	RoleGuestViewer
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
