package backlog

import (
	"errors"
	"fmt"
	"strconv"
)

type optionType int

func (o optionType) String() string {
	switch o {
	case optionActivityTypeIDs:
		return "ActivityTypeIDs"
	case optionAll:
		return "All"
	case optionArchived:
		return "Archived"
	case optionChartEnabled:
		return "ChartEnabled"
	case optionContent:
		return "Content"
	case optionCount:
		return "Count"
	case optionKey:
		return "Key"
	case optionKeyword:
		return "Keyword"
	case optionName:
		return "Name"
	case optionMailAddress:
		return "MailAddress"
	case optionMailNotify:
		return "MailNotify"
	case optionMaxID:
		return "MaxID"
	case optionMinID:
		return "MinID"
	case optionOrder:
		return "Order"
	case optionPassword:
		return "Password"
	case optionProjectLeaderCanEditProjectLeader:
		return "ProjectLeaderCanEditProjectLeader"
	case optionRoleType:
		return "RoleType"
	case optionSubtaskingEnabled:
		return "SubtaskingEnabled"
	case optionTextFormattingRule:
		return "TextFormattingRule"
	default:
		return "unknown"
	}
}

type optionFunc func(params *requestParams) error

type option struct {
	t optionType
	f optionFunc
}

func (o *option) validate(validTypes []optionType) error {
	for _, valid := range validTypes {
		if o.t == valid {
			return nil
		}
	}
	return newInvalidOptionError(o.t, validTypes)
}

func (o *option) set(params *requestParams) error {
	return o.f(params)
}

func withActivityTypeIDs(typeIDs []int) *option {
	return &option{optionActivityTypeIDs, func(params *requestParams) error {
		for _, id := range typeIDs {
			if id < 1 || 26 < id {
				return errors.New("activityTypeId must be between 1 and 26")
			}
			params.Add("activityTypeId[]", strconv.Itoa(id))
		}
		return nil
	}}
}

func withAll(enabeld bool) *option {
	return &option{optionAll, func(params *requestParams) error {
		params.Set("all", strconv.FormatBool(enabeld))
		return nil
	}}
}

func withArchived(archived bool) *option {
	return &option{optionArchived, func(params *requestParams) error {
		params.Set("archived", strconv.FormatBool(archived))
		return nil
	}}
}

func withChartEnabled(enabeld bool) *option {
	return &option{optionChartEnabled, func(params *requestParams) error {
		params.Set("chartEnabled", strconv.FormatBool(enabeld))
		return nil
	}}
}

func withContent(content string) *option {
	return &option{optionContent, func(params *requestParams) error {
		if content == "" {
			return errors.New("content must not be empty")
		}
		params.Set("content", content)
		return nil
	}}
}

func withCount(count int) *option {
	return &option{optionCount, func(params *requestParams) error {
		if count < 1 || 100 < count {
			return errors.New("count must be between 1 and 100")
		}
		params.Set("count", strconv.Itoa(count))
		return nil
	}}
}

func withKey(key string) *option {
	return &option{optionKey, func(params *requestParams) error {
		if key == "" {
			return errors.New("key must not be empty")
		}
		params.Set("key", key)
		return nil
	}}
}

func withKeyword(keyword string) *option {
	return &option{optionKeyword, func(params *requestParams) error {
		params.Set("keyword", keyword)
		return nil
	}}
}

func withName(name string) *option {
	return &option{optionName, func(params *requestParams) error {
		if name == "" {
			return errors.New("name must not be empty")
		}
		params.Set("name", name)
		return nil
	}}
}

func withMailAddress(mailAddress string) *option {
	// ToDo: validate mailAddress
	return &option{optionMailAddress, func(params *requestParams) error {
		if mailAddress == "" {
			return errors.New("mailAddress must not be empty")
		}
		params.Set("mailAddress", mailAddress)
		return nil
	}}
}

func withMailNotify(enabeld bool) *option {
	return &option{optionMailNotify, func(params *requestParams) error {
		params.Set("mailNotify", strconv.FormatBool(enabeld))
		return nil
	}}
}

func withMaxID(maxID int) *option {
	return &option{optionMaxID, func(params *requestParams) error {
		if maxID < 1 {
			return errors.New("maxId must be greater than 1")
		}
		params.Set("maxId", strconv.Itoa(maxID))
		return nil
	}}
}

func withMinID(minID int) *option {
	return &option{optionMinID, func(params *requestParams) error {
		if minID < 1 {
			return errors.New("minId must be greater than 1")
		}
		params.Set("minId", strconv.Itoa(minID))
		return nil
	}}
}

func withOrder(order order) *option {
	return &option{optionOrder, func(params *requestParams) error {
		if order != OrderAsc && order != OrderDesc {
			return fmt.Errorf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
		}
		params.Set("order", string(order))
		return nil
	}}
}

func withPassword(password string) *option {
	return &option{optionPassword, func(params *requestParams) error {
		if password == "" {
			return errors.New("password must not be empty")
		}
		params.Set("password", password)
		return nil
	}}
}

func withProjectLeaderCanEditProjectLeader(enabeld bool) *option {
	return &option{optionProjectLeaderCanEditProjectLeader, func(params *requestParams) error {
		params.Set("projectLeaderCanEditProjectLeader", strconv.FormatBool(enabeld))
		return nil
	}}
}

func withRoleType(roleType role) *option {
	return &option{optionRoleType, func(params *requestParams) error {
		if roleType < 1 || 6 < roleType {
			return errors.New("roleType must be between 1 and 7")
		}
		params.Add("roleType", strconv.Itoa(int(roleType)))
		return nil
	}}
}

func withSubtaskingEnabled(enabeld bool) *option {
	return &option{optionSubtaskingEnabled, func(params *requestParams) error {
		params.Set("subtaskingEnabled", strconv.FormatBool(enabeld))
		return nil
	}}
}

func withTextFormattingRule(format format) *option {
	return &option{optionTextFormattingRule, func(params *requestParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			return fmt.Errorf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
		}
		params.Set("textFormattingRule", string(format))
		return nil
	}}
}

// ActivityOption is type of option for ActivityService.
type ActivityOption struct {
	*option
}

// ActivityOptionService has methods to make option for ActivityService.
type ActivityOptionService struct {
}

// WithActivityTypeIDs returns option. the option sets `activityTypeId` for user.
func (*ActivityOptionService) WithActivityTypeIDs(typeIDs []int) *ActivityOption {
	return &ActivityOption{withActivityTypeIDs(typeIDs)}
}

// WithMinID returns option. the option sets `minId` for user.
func (*ActivityOptionService) WithMinID(minID int) *ActivityOption {
	return &ActivityOption{withMinID(minID)}
}

// WithMaxID returns option. the option sets `maxId` for user.
func (*ActivityOptionService) WithMaxID(maxID int) *ActivityOption {
	return &ActivityOption{withMaxID(maxID)}
}

// WithCount returns option. the option sets `count` for user.
func (*ActivityOptionService) WithCount(count int) *ActivityOption {
	return &ActivityOption{withCount(count)}
}

// WithOrder returns option. the option sets `order` for user.
func (*ActivityOptionService) WithOrder(order order) *ActivityOption {
	return &ActivityOption{withOrder(order)}
}

// ProjectOption is type of option for ProjectService.
type ProjectOption struct {
	*option
}

// ProjectOptionService has methods to make option for ProjectService.
type ProjectOptionService struct {
}

// WithAll returns option. the option sets `all` for project.
func (*ProjectOptionService) WithAll(enabeld bool) *ProjectOption {
	return &ProjectOption{withAll(enabeld)}
}

// WithKey returns option. the option sets `key` for project.
func (*ProjectOptionService) WithKey(key string) *ProjectOption {
	return &ProjectOption{withKey(key)}
}

// WithName returns option. the option sets `name` for project.
func (*ProjectOptionService) WithName(name string) *ProjectOption {
	return &ProjectOption{withName(name)}
}

// WithChartEnabled returns option. the option sets `chartEnabled` for project.
func (*ProjectOptionService) WithChartEnabled(enabeld bool) *ProjectOption {
	return &ProjectOption{withChartEnabled(enabeld)}
}

// WithSubtaskingEnabled returns option. the option sets `subtaskingEnabled` for project.
func (*ProjectOptionService) WithSubtaskingEnabled(enabeld bool) *ProjectOption {
	return &ProjectOption{withSubtaskingEnabled(enabeld)}
}

// WithProjectLeaderCanEditProjectLeader returns option. the option sets `projectLeaderCanEditProjectLeader` for project.
func (*ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabeld bool) *ProjectOption {
	return &ProjectOption{withProjectLeaderCanEditProjectLeader(enabeld)}
}

// WithTextFormattingRule returns option. the option sets `textFormattingRule` for project.
func (*ProjectOptionService) WithTextFormattingRule(format format) *ProjectOption {
	return &ProjectOption{withTextFormattingRule(format)}
}

// WithArchived returns option. the option sets `archived` for project.
func (*ProjectOptionService) WithArchived(archived bool) *ProjectOption {
	return &ProjectOption{withArchived(archived)}
}

// UserOption is type of option for UserService.
type UserOption struct {
	*option
}

// UserOptionService has methods to make option for UserService.
type UserOptionService struct {
}

// WithPassword returns option. the option sets `password` for user.
func (*UserOptionService) WithPassword(password string) *UserOption {
	return &UserOption{withPassword(password)}
}

// WithName returns option. the option sets `password` for user.
func (*UserOptionService) WithName(name string) *UserOption {
	return &UserOption{withName(name)}
}

// WithMailAddress returns option. the option sets `mailAddress` for user.
func (*UserOptionService) WithMailAddress(mailAddress string) *UserOption {
	return &UserOption{withMailAddress(mailAddress)}
}

// WithRoleType returns option. the option sets `roleType` for user.
func (*UserOptionService) WithRoleType(roleType role) *UserOption {
	return &UserOption{withRoleType(roleType)}
}

// WikiOption is type of option for WikiService.
type WikiOption struct {
	*option
}

// WikiOptionService has methods to make option for WikiService.
type WikiOptionService struct {
}

// WithKeyword returns option. the option sets `keyword` for wiki.
func (*WikiOptionService) WithKeyword(keyword string) *WikiOption {
	return &WikiOption{withKeyword(keyword)}
}

// WithName returns option. the option sets `name` for wiki.
func (*WikiOptionService) WithName(name string) *WikiOption {
	return &WikiOption{withName(name)}
}

// WithContent returns option. the option sets `content` for wiki.
func (*WikiOptionService) WithContent(content string) *WikiOption {
	return &WikiOption{withContent(content)}
}

// WithMailNotify returns option. the option sets `mailNotify` for wiki.
func (*WikiOptionService) WithMailNotify(enabeld bool) *WikiOption {
	return &WikiOption{withMailNotify(enabeld)}
}
