// Package validate provides input validation helpers shared across service packages.
package validate

import (
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
	return ValidateIDOrKey("issueIDOrKey", issueIDOrKey)
}

func ValidateProjectID(projectID int) *core.ValidationError {
	return ValidatePositiveInt("projectID", projectID)
}

func ValidateProjectIDOrKey(projectIDOrKey string) *core.ValidationError {
	return ValidateIDOrKey("projectIDOrKey", projectIDOrKey)
}

func ValidatePRNumber(prNumber int) *core.ValidationError {
	return ValidatePositiveInt("prNumber", prNumber)
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) *core.ValidationError {
	return ValidateIDOrKey("repositoryIDOrName", repositoryIDOrName)
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
