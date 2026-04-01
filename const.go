package backlog

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
	formArchived                          formType = "archived"
	formChartEnabled                      formType = "chartEnabled"
	formContent                           formType = "content"
	formKey                               formType = "key"
	formName                              formType = "name"
	formMailAddress                       formType = "mailAddress"
	formMailNotify                        formType = "mailNotify"
	formPassword                          formType = "password"
	formProjectLeaderCanEditProjectLeader formType = "projectLeaderCanEditProjectLeader"
	formRoleType                          formType = "roleType"
	formSubtaskingEnabled                 formType = "subtaskingEnabled"
	formTextFormattingRule                formType = "textFormattingRule"
	formUserID                            formType = "userId"
	formSendMail                          formType = "sendMail"
)

const (
	queryActivityTypeIDs queryType = "activityTypeId[]"
	queryAll             queryType = "all"
	queryArchived        queryType = "archived"
	queryCount           queryType = "count"
	queryKey             queryType = "key"
	queryKeyword         queryType = "keyword"
	queryMaxID           queryType = "maxId"
	queryMinID           queryType = "minId"
	queryOrder           queryType = "order"
)
