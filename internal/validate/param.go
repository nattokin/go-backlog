package validate

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/nattokin/go-backlog/internal/core"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var validIssueSorts = []string{
	"issueType", "category", "version", "milestone", "summary", "status",
	"priority", "attachment", "sharedFile", "created", "createdUser",
	"updated", "updatedUser", "assignee", "startDate", "dueDate",
	"estimatedHours", "actualHours", "childIssue",
}

var validTextFormattingRules = []string{"backlog", "markdown"}

// ValidatePositiveInt validates that value is not less than 1.
func ValidatePositiveInt(field string, value int) *core.ValidationError {
	if value < 1 {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must not be less than 1", field))
	}
	return nil
}

// ValidatePositiveInts validates that every element of values is not less than 1.
func ValidatePositiveInts(field string, values []int) *core.ValidationError {
	for _, v := range values {
		if v < 1 {
			return core.NewValidationError(field, fmt.Sprintf("invalid %s: %d must not be less than 1", field, v))
		}
	}
	return nil
}

// ValidateIntRange validates that value is within [min, max].
func ValidateIntRange(field string, value, min, max int) *core.ValidationError {
	if value < min || value > max {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be between %d and %d", field, min, max))
	}
	return nil
}

// ValidatePositiveFloat64 validates that value is greater than 0.
func ValidatePositiveFloat64(field string, value float64) *core.ValidationError {
	if value <= 0 {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be greater than 0", field))
	}
	return nil
}

// ValidateDateFormat validates that date is formatted as yyyy-MM-dd.
func ValidateDateFormat(field, date string) *core.ValidationError {
	if !datePattern.MatchString(date) {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be formatted as yyyy-MM-dd, got %q", field, date))
	}
	return nil
}

// ValidateNonEmptyString validates that value is not empty or whitespace-only.
func ValidateNonEmptyString(field, value string) *core.ValidationError {
	if strings.TrimSpace(value) == "" {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must not be empty", field))
	}
	return nil
}

// ValidateIDOrKey validates that value is not empty/whitespace-only and not the literal "0".
func ValidateIDOrKey(field, value string) *core.ValidationError {
	if strings.TrimSpace(value) == "" {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must not be empty", field))
	}
	if value == "0" {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must not be '0'", field))
	}
	return nil
}

// ValidateEmail validates that value is a single, well-formed email address.
func ValidateEmail(field, value string) *core.ValidationError {
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: not a valid email address", field))
	}
	return nil
}

// ValidateOrder validates that order is "asc" or "desc".
func ValidateOrder(field, order string) *core.ValidationError {
	if order != "asc" && order != "desc" {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be only 'asc' or 'desc'", field))
	}
	return nil
}

// ValidatePassword validates that password is at least 8 characters long.
func ValidatePassword(field, password string) *core.ValidationError {
	if len(password) < 8 {
		return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be at least 8 characters long", field))
	}
	return nil
}

// ValidateIssueSort validates that sort is one of the Backlog issue-list sort keys.
func ValidateIssueSort(field, sort string) *core.ValidationError {
	for _, v := range validIssueSorts {
		if sort == v {
			return nil
		}
	}
	return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be a valid sort value", field))
}

// ValidateTextFormattingRule validates that format is "backlog" or "markdown".
func ValidateTextFormattingRule(field, format string) *core.ValidationError {
	for _, v := range validTextFormattingRules {
		if format == v {
			return nil
		}
	}
	return core.NewValidationError(field, fmt.Sprintf("invalid %s: must be only 'backlog' or 'markdown'", field))
}
