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

//

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
