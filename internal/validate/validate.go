// Package validate provides input validation helpers shared across service packages.
package validate

import "github.com/nattokin/go-backlog/internal/core"

func ValidateActivityID(activityID int) *core.ValidationError {
	if activityID < 1 {
		return core.NewValidationError("activityID", "activityID must not be less than 1")
	}
	return nil
}

func ValidateAttachmentID(attachmentID int) *core.ValidationError {
	if attachmentID < 1 {
		return core.NewValidationError("attachmentID", "attachmentID must not be less than 1")
	}
	return nil
}

func ValidateCommentID(commentID int) *core.ValidationError {
	if commentID < 1 {
		return core.NewValidationError("commentID", "commentID must not be less than 1")
	}
	return nil
}

func ValidateCustomFieldID(customFieldID int) *core.ValidationError {
	if customFieldID < 1 {
		return core.NewValidationError("customFieldID", "customFieldID must not be less than 1")
	}
	return nil
}

func ValidateIssueIDOrKey(issueIDOrKey string) *core.ValidationError {
	if issueIDOrKey == "" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be '0'")
	}
	return nil
}

func ValidateProjectID(projectID int) *core.ValidationError {
	if projectID < 1 {
		return core.NewValidationError("projectID", "projectID must not be less than 1")
	}
	return nil
}

func ValidateProjectIDOrKey(projectIDOrKey string) *core.ValidationError {
	if projectIDOrKey == "" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be empty")
	}
	if projectIDOrKey == "0" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be '0'")
	}
	return nil
}

func ValidatePRNumber(prNumber int) *core.ValidationError {
	if prNumber < 1 {
		return core.NewValidationError("prNumber", "prNumber must not be less than 1")
	}
	return nil
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) *core.ValidationError {
	if repositoryIDOrName == "" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be '0'")
	}
	return nil
}

func ValidateSharedFileID(fileID int) *core.ValidationError {
	if fileID < 1 {
		return core.NewValidationError("fileID", "fileID must not be less than 1")
	}
	return nil
}

func ValidateStarID(starID int) *core.ValidationError {
	if starID < 1 {
		return core.NewValidationError("starID", "starID must not be less than 1")
	}
	return nil
}

func ValidateUserID(userID int) *core.ValidationError {
	if userID < 1 {
		return core.NewValidationError("userID", "userID must not be less than 1")
	}
	return nil
}

func ValidateVersionID(versionID int) *core.ValidationError {
	if versionID < 1 {
		return core.NewValidationError("versionID", "versionID must not be less than 1")
	}
	return nil
}

func ValidateWebhookID(webhookID int) *core.ValidationError {
	if webhookID < 1 {
		return core.NewValidationError("webhookID", "webhookID must not be less than 1")
	}
	return nil
}

func ValidateWikiID(wikiID int) *core.ValidationError {
	if wikiID < 1 {
		return core.NewValidationError("wikiID", "wikiID must not be less than 1")
	}
	return nil
}
