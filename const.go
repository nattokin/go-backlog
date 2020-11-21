package backlog

// Order type
const (
	OrderAsc  order = "asc"
	OrderDesc order = "desc"
)

// Fomat of Backlog wiki
const (
	FormatMarkdown format = "markdown"
	FormatBacklog  format = "backlog"
)

// Role type
const (
	_ role = iota
	RoleAdministrator
	RoleNormalUser
	RoleReporter
	RoleViewer
	RoleGuestReporter
	RoleGuestViewer
)

// TODO: activityId

const (
	_ formType = iota
	formArchived
	formChartEnabled
	formContent
	formKey
	formName
	formMailAddress
	formMailNotify
	formPassword
	formProjectLeaderCanEditProjectLeader
	formRoleType
	formSubtaskingEnabled
	formTextFormattingRule
)

const (
	_ queryType = iota
	queryActivityTypeIDs
	queryAll
	queryArchived
	queryCount
	queryKey
	queryKeyword
	queryMaxID
	queryMinID
	queryOrder
)
