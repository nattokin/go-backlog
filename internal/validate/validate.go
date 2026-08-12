// Package validate provides input validation helpers shared across service packages.
package validate

import (
	"strings"

	"github.com/nattokin/go-backlog/internal/core"
)

func ValidateActivityID(activityID int) *core.ValidationError {
	return ValidatePositiveInt("activityID", activityID)
}

func ValidateAttachmentID(attachmentID int) *core.ValidationError {
	return ValidatePositiveInt("attachmentID", attachmentID)
}

func ValidateCommentID(commentID int) *core.ValidationError {
	return ValidatePositiveInt("commentID", commentID)
}

func ValidateCustomFieldID(customFieldID int) *core.ValidationError {
	return ValidatePositiveInt("customFieldID", customFieldID)
}

func ValidateIssueIDOrKey(issueIDOrKey string) *core.ValidationError {
	if strings.TrimSpace(issueIDOrKey) == "" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be '0'")
	}
	return nil
}

func ValidateProjectID(projectID int) *core.ValidationError {
	return ValidatePositiveInt("projectID", projectID)
}

func ValidateProjectIDOrKey(projectIDOrKey string) *core.ValidationError {
	if strings.TrimSpace(projectIDOrKey) == "" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be empty")
	}
	if projectIDOrKey == "0" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be '0'")
	}
	return nil
}

func ValidatePRNumber(prNumber int) *core.ValidationError {
	return ValidatePositiveInt("prNumber", prNumber)
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) *core.ValidationError {
	if strings.TrimSpace(repositoryIDOrName) == "" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be '0'")
	}
	return nil
}

func ValidateSharedFileID(fileID int) *core.ValidationError {
	return ValidatePositiveInt("fileID", fileID)
}

func ValidateStarID(starID int) *core.ValidationError {
	return ValidatePositiveInt("starID", starID)
}

func ValidateUserID(userID int) *core.ValidationError {
	return ValidatePositiveInt("userID", userID)
}

func ValidateVersionID(versionID int) *core.ValidationError {
	return ValidatePositiveInt("versionID", versionID)
}

func ValidateWebhookID(webhookID int) *core.ValidationError {
	return ValidatePositiveInt("webhookID", webhookID)
}

func ValidateWikiID(wikiID int) *core.ValidationError {
	return ValidatePositiveInt("wikiID", wikiID)
}
