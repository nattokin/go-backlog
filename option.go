package backlog

import (
	"errors"
	"fmt"
	"strconv"
)

type option func(p *requestParams) error

func withActivityTypeIDs(typeIDs []int) option {
	return func(p *requestParams) error {
		for _, id := range typeIDs {
			if id < 1 || 26 < id {
				return errors.New("activityTypeId must be between 1 and 26")
			}
			p.Add("activityTypeId[]", strconv.Itoa(id))
		}
		return nil
	}
}

func withArchived(archived bool) option {
	return func(p *requestParams) error {
		p.Set("archived", strconv.FormatBool(archived))
		return nil
	}
}

func withChartEnabled(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("chartEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

func withContent(content string) option {
	return func(p *requestParams) error {
		if content == "" {
			return errors.New("content must not be empty")
		}
		p.Set("content", content)
		return nil
	}
}

func withCount(count int) option {
	return func(p *requestParams) error {
		if count < 1 || 100 < count {
			return errors.New("count must be between 1 and 100")
		}
		p.Set("count", strconv.Itoa(count))
		return nil
	}
}

func withKey(key string) option {
	return func(p *requestParams) error {
		if key == "" {
			return errors.New("key must not be empty")
		}
		p.Set("key", key)
		return nil
	}
}

func withName(name string) option {
	return func(p *requestParams) error {
		if name == "" {
			return errors.New("name must not be empty")
		}
		p.Set("name", name)
		return nil
	}
}

func withMailAddress(mailAddress string) option {
	// ToDo: validate mailAddress
	return func(p *requestParams) error {
		if mailAddress == "" {
			return errors.New("mailAddress must not be empty")
		}
		p.Set("mailAddress", mailAddress)
		return nil
	}
}

func withMailNotify(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("mailNotify", strconv.FormatBool(enabeld))
		return nil
	}
}

func withMaxID(maxID int) option {
	return func(p *requestParams) error {
		if maxID < 1 {
			return errors.New("maxId must be greater than 1")
		}
		p.Set("maxId", strconv.Itoa(maxID))
		return nil
	}
}

func withMinID(minID int) option {
	return func(p *requestParams) error {
		if minID < 1 {
			return errors.New("minId must be greater than 1")
		}
		p.Set("minId", strconv.Itoa(minID))
		return nil
	}
}

func withOrder(order order) option {
	return func(p *requestParams) error {
		if order != OrderAsc && order != OrderDesc {
			return fmt.Errorf("order must be only '%s' or '%s'", string(OrderAsc), string(OrderDesc))
		}
		p.Set("order", string(order))
		return nil
	}
}

func withPassword(password string) option {
	return func(p *requestParams) error {
		if password == "" {
			return errors.New("password must not be empty")
		}
		p.Set("password", password)
		return nil
	}
}

func withProjectLeaderCanEditProjectLeader(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("projectLeaderCanEditProjectLeader", strconv.FormatBool(enabeld))
		return nil
	}
}

func withRoleType(roleType role) option {
	return func(p *requestParams) error {
		if roleType < 1 || 6 < roleType {
			return errors.New("roleType must be between 1 and 7")
		}
		p.Add("roleType", strconv.Itoa(int(roleType)))
		return nil
	}
}

func withSubtaskingEnabled(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("subtaskingEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

func withTextFormattingRule(format format) option {
	return func(p *requestParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			return fmt.Errorf("format must be only '%s' or '%s'", string(FormatBacklog), string(FormatMarkdown))
		}
		p.Set("textFormattingRule", string(format))
		return nil
	}
}

// ActivityOption is type of functional option for ActivityService.
type ActivityOption option

// WithActivityTypeIDs returns option. the option sets `activityTypeId` for user.
func (*ActivityOptionService) WithActivityTypeIDs(typeIDs []int) ActivityOption {
	return ActivityOption(withActivityTypeIDs(typeIDs))
}

// WithMinID returns option. the option sets `minId` for user.
func (*ActivityOptionService) WithMinID(minID int) ActivityOption {
	return ActivityOption(withMinID(minID))
}

// WithMaxID returns option. the option sets `maxId` for user.
func (*ActivityOptionService) WithMaxID(maxID int) ActivityOption {
	return ActivityOption(withMaxID(maxID))
}

// WithCount returns option. the option sets `count` for user.
func (*ActivityOptionService) WithCount(count int) ActivityOption {
	return ActivityOption(withCount(count))
}

// WithOrder returns option. the option sets `order` for user.
func (*ActivityOptionService) WithOrder(order order) ActivityOption {
	return ActivityOption(withOrder(order))
}

// ProjectOption is type of functional option for ProjectService.
type ProjectOption option

// WithKey returns option. the option sets `key` for project.
func (*ProjectOptionService) WithKey(key string) ProjectOption {
	return ProjectOption(withKey(key))
}

// WithName returns option. the option sets `name` for project.
func (*ProjectOptionService) WithName(name string) ProjectOption {
	return ProjectOption(withName(name))
}

// WithChartEnabled returns option. the option sets `chartEnabled` for project.
func (*ProjectOptionService) WithChartEnabled(enabeld bool) ProjectOption {
	return ProjectOption(withChartEnabled(enabeld))
}

// WithSubtaskingEnabled returns option. the option sets `subtaskingEnabled` for project.
func (*ProjectOptionService) WithSubtaskingEnabled(enabeld bool) ProjectOption {
	return ProjectOption(withSubtaskingEnabled(enabeld))
}

// WithProjectLeaderCanEditProjectLeader returns option. the option sets `projectLeaderCanEditProjectLeader` for project.
func (*ProjectOptionService) WithProjectLeaderCanEditProjectLeader(enabeld bool) ProjectOption {
	return ProjectOption(withProjectLeaderCanEditProjectLeader(enabeld))
}

// WithTextFormattingRule returns option. the option sets `textFormattingRule` for project.
func (*ProjectOptionService) WithTextFormattingRule(format format) ProjectOption {
	return ProjectOption(withTextFormattingRule(format))
}

// WithArchived returns option. the option sets `archived` for project.
func (*ProjectOptionService) WithArchived(archived bool) ProjectOption {
	return ProjectOption(withArchived(archived))
}

// UserOption is type of functional option for UserService.
type UserOption option

// WithPassword returns option. the option sets `password` for user.
func (*UserOptionService) WithPassword(password string) UserOption {
	return UserOption(withPassword(password))
}

// WithName returns option. the option sets `password` for user.
func (*UserOptionService) WithName(name string) UserOption {
	return UserOption(withName(name))
}

// WithMailAddress returns option. the option sets `mailAddress` for user.
func (*UserOptionService) WithMailAddress(mailAddress string) UserOption {
	return UserOption(withMailAddress(mailAddress))
}

// WithRoleType returns option. the option sets `roleType` for user.
func (*UserOptionService) WithRoleType(roleType role) UserOption {
	return UserOption(withRoleType(roleType))
}

// WikiOption is type of functional option for WikiService.
type WikiOption option

// WithName returns option. the option sets `name` for wiki.
func (*WikiOptionService) WithName(name string) WikiOption {
	return WikiOption(withName(name))
}

// WithContent returns option. the option sets `content` for wiki.
func (*WikiOptionService) WithContent(content string) WikiOption {
	return WikiOption(withContent(content))
}

// WithMailNotify returns option. the option sets `mailNotify` for wiki.
func (*WikiOptionService) WithMailNotify(enabeld bool) WikiOption {
	return WikiOption(withMailNotify(enabeld))
}
