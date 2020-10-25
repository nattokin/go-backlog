package backlog

import (
	"errors"
	"fmt"
	"strconv"
)

type option func(p *requestParams) error

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

func withContent(content string) option {
	return func(p *requestParams) error {
		if content == "" {
			return errors.New("content must not be empty")
		}
		p.Set("content", content)
		return nil
	}
}

func withMailNotify(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("mailNotify", strconv.FormatBool(enabeld))
		return nil
	}
}

func withChartEnabled(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("chartEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

func withSubtaskingEnabled(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("subtaskingEnabled", strconv.FormatBool(enabeld))
		return nil
	}
}

func withProjectLeaderCanEditProjectLeader(enabeld bool) option {
	return func(p *requestParams) error {
		p.Set("projectLeaderCanEditProjectLeader", strconv.FormatBool(enabeld))
		return nil
	}
}

func withTextFormattingRule(format string) option {
	return func(p *requestParams) error {
		if format != FormatBacklog && format != FormatMarkdown {
			return fmt.Errorf("format must be only '%s' or '%s'", FormatBacklog, FormatMarkdown)
		}
		p.Set("textFormattingRule", format)
		return nil
	}
}

func withArchived(archived bool) option {
	return func(p *requestParams) error {
		p.Set("archived", strconv.FormatBool(archived))
		return nil
	}
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
func (*ProjectOptionService) WithTextFormattingRule(format string) ProjectOption {
	return ProjectOption(withTextFormattingRule(format))
}

// WithArchived returns option. the option sets `archived` for project.
func (*ProjectOptionService) WithArchived(archived bool) ProjectOption {
	return ProjectOption(withArchived(archived))
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
