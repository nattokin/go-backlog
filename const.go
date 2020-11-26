package backlog

// Order type
const (
	OrderAsc  Order = "asc"
	OrderDesc Order = "desc"
)

// Fomat of Backlog wiki
const (
	FormatMarkdown Format = "markdown"
	FormatBacklog  Format = "backlog"
)

// Role type
const (
	_ Role = iota
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
