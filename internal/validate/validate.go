// Package validate provides input validation helpers shared across service packages.
package validate

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"

	"github.com/nattokin/go-backlog/internal/validation"
)

var datePattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)

var validIssueSorts = []string{
	"issueType", "category", "version", "milestone", "summary", "status",
	"priority", "attachment", "sharedFile", "created", "createdUser",
	"updated", "updatedUser", "assignee", "startDate", "dueDate",
	"estimatedHours", "actualHours", "childIssue",
}

var validTextFormattingRules = []string{"backlog", "markdown"}

// ValidateDateFormat validates that date is formatted as yyyy-MM-dd.
func ValidateDateFormat(field, date string) *validation.Error {
	if !datePattern.MatchString(date) {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must be formatted as yyyy-MM-dd, got %q", field, date))
	}
	return nil
}

// ValidateEmail validates that value is a single, well-formed email address.
func ValidateEmail(field, value string) *validation.Error {
	addr, err := mail.ParseAddress(value)
	if err != nil || addr.Address != value {
		return validation.NewError(field, fmt.Sprintf("invalid %s: not a valid email address", field))
	}
	return nil
}

// ValidateIDOrKey validates that value is not empty/whitespace-only and not the literal "0".
func ValidateIDOrKey(field, value string) *validation.Error {
	if strings.TrimSpace(value) == "" {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must not be empty", field))
	}
	if value == "0" {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must not be '0'", field))
	}
	return nil
}

// ValidateIntRange validates that value is within [min, max].
func ValidateIntRange(field string, value, min, max int) *validation.Error {
	if value < min || value > max {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must be between %d and %d", field, min, max))
	}
	return nil
}

// ValidateIssueSort validates that sort is one of the Backlog issue-list sort keys.
func ValidateIssueSort(field, sort string) *validation.Error {
	for _, v := range validIssueSorts {
		if sort == v {
			return nil
		}
	}
	return validation.NewError(field, fmt.Sprintf("invalid %s: must be a valid sort value", field))
}

// ValidateNonEmptyString validates that value is not empty or whitespace-only.
func ValidateNonEmptyString(field, value string) *validation.Error {
	if strings.TrimSpace(value) == "" {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must not be empty", field))
	}
	return nil
}

// ValidateOrder validates that order is "asc" or "desc".
func ValidateOrder(field, order string) *validation.Error {
	if order != "asc" && order != "desc" {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must be only 'asc' or 'desc'", field))
	}
	return nil
}

// ValidatePassword validates that password is at least 8 characters long.
func ValidatePassword(field, password string) *validation.Error {
	if len(password) < 8 {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must be at least 8 characters long", field))
	}
	return nil
}

// ValidatePositiveFloat64 validates that value is greater than 0.
func ValidatePositiveFloat64(field string, value float64) *validation.Error {
	if value <= 0 {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must be greater than 0", field))
	}
	return nil
}

// ValidatePositiveInt validates that value is not less than 1.
func ValidatePositiveInt(field string, value int) *validation.Error {
	if value < 1 {
		return validation.NewError(field, fmt.Sprintf("invalid %s: must not be less than 1", field))
	}
	return nil
}

// ValidatePositiveInts validates that every element of values is not less than 1.
func ValidatePositiveInts(field string, values []int) *validation.Error {
	for _, v := range values {
		if v < 1 {
			return validation.NewError(field, fmt.Sprintf("invalid %s: %d must not be less than 1", field, v))
		}
	}
	return nil
}

// ValidateTextFormattingRule validates that format is "backlog" or "markdown".
func ValidateTextFormattingRule(field, format string) *validation.Error {
	for _, v := range validTextFormattingRules {
		if format == v {
			return nil
		}
	}
	return validation.NewError(field, fmt.Sprintf("invalid %s: must be only 'backlog' or 'markdown'", field))
}

func ValidateActivityID(activityID int) *validation.Error {
	return ValidatePositiveInt("activityID", activityID)
}

func ValidateAttachmentID(attachmentID int) *validation.Error {
	return ValidatePositiveInt("attachmentID", attachmentID)
}

func ValidateCommentID(commentID int) *validation.Error {
	return ValidatePositiveInt("commentID", commentID)
}

func ValidateCustomFieldID(customFieldID int) *validation.Error {
	return ValidatePositiveInt("customFieldID", customFieldID)
}

func ValidateIssueIDOrKey(issueIDOrKey string) *validation.Error {
	return ValidateIDOrKey("issueIDOrKey", issueIDOrKey)
}

func ValidatePRNumber(prNumber int) *validation.Error {
	return ValidatePositiveInt("prNumber", prNumber)
}

func ValidateProjectID(projectID int) *validation.Error {
	return ValidatePositiveInt("projectID", projectID)
}

func ValidateProjectIDOrKey(projectIDOrKey string) *validation.Error {
	return ValidateIDOrKey("projectIDOrKey", projectIDOrKey)
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) *validation.Error {
	return ValidateIDOrKey("repositoryIDOrName", repositoryIDOrName)
}

func ValidateSharedFileID(fileID int) *validation.Error {
	return ValidatePositiveInt("fileID", fileID)
}

func ValidateStarID(starID int) *validation.Error {
	return ValidatePositiveInt("starID", starID)
}

func ValidateUserID(userID int) *validation.Error {
	return ValidatePositiveInt("userID", userID)
}

func ValidateVersionID(versionID int) *validation.Error {
	return ValidatePositiveInt("versionID", versionID)
}

func ValidateWebhookID(webhookID int) *validation.Error {
	return ValidatePositiveInt("webhookID", webhookID)
}

func ValidateWikiID(wikiID int) *validation.Error {
	return ValidatePositiveInt("wikiID", wikiID)
}
