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

const (
	_ optionType = iota
	optionActivityTypeIDs
	optionAll
	optionArchived
	optionChartEnabled
	optionContent
	optionCount
	optionKey
	optionKeyword
	optionName
	optionMailAddress
	optionMailNotify
	optionMaxID
	optionMinID
	optionOrder
	optionPassword
	optionProjectLeaderCanEditProjectLeader
	optionRoleType
	optionSubtaskingEnabled
	optionTextFormattingRule
)
