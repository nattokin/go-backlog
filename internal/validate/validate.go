// Package validate provides input validation helpers shared across service packages.
package validate

import "github.com/nattokin/go-backlog/internal/core"

func ValidateActivityID(activityID int) error {
	if activityID < 1 {
		return core.NewValidationError("activityID", "activityID must not be less than 1")
	}
	return nil
}

func ValidateAttachmentID(attachmentID int) error {
	if attachmentID < 1 {
		return core.NewValidationError("attachmentID", "attachmentID must not be less than 1")
	}
	return nil
}

func ValidateCommentID(commentID int) error {
	if commentID < 1 {
		return core.NewValidationError("commentID", "commentID must not be less than 1")
	}
	return nil
}

func ValidateCustomFieldID(customFieldID int) error {
	if customFieldID < 1 {
		return core.NewValidationError("customFieldID", "customFieldID must not be less than 1")
	}
	return nil
}

func ValidateIssueIDOrKey(issueIDOrKey string) error {
	if issueIDOrKey == "" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return core.NewValidationError("issueIDOrKey", "issueIDOrKey must not be '0'")
	}
	return nil
}

func ValidateProjectID(projectID int) error {
	if projectID < 1 {
		return core.NewValidationError("projectID", "projectID must not be less than 1")
	}
	return nil
}

func ValidateProjectIDOrKey(projectIDOrKey string) error {
	if projectIDOrKey == "" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be empty")
	}
	if projectIDOrKey == "0" {
		return core.NewValidationError("projectIDOrKey", "projectIDOrKey must not be '0'")
	}
	return nil
}

func ValidatePRNumber(prNumber int) error {
	if prNumber < 1 {
		return core.NewValidationError("prNumber", "prNumber must not be less than 1")
	}
	return nil
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) error {
	if repositoryIDOrName == "" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return core.NewValidationError("repositoryIDOrName", "repositoryIDOrName must not be '0'")
	}
	return nil
}

func ValidateSharedFileID(fileID int) error {
	if fileID < 1 {
		return core.NewValidationError("fileID", "fileID must not be less than 1")
	}
	return nil
}

func ValidateStarID(starID int) error {
	if starID < 1 {
		return core.NewValidationError("starID", "starID must not be less than 1")
	}
	return nil
}

func ValidateUserID(userID int) error {
	if userID < 1 {
		return core.NewValidationError("userID", "userID must not be less than 1")
	}
	return nil
}

func ValidateVersionID(versionID int) error {
	if versionID < 1 {
		return core.NewValidationError("versionID", "versionID must not be less than 1")
	}
	return nil
}

func ValidateWebhookID(webhookID int) error {
	if webhookID < 1 {
		return core.NewValidationError("webhookID", "webhookID must not be less than 1")
	}
	return nil
}

func ValidateWikiID(wikiID int) error {
	if wikiID < 1 {
		return core.NewValidationError("wikiID", "wikiID must not be less than 1")
	}
	return nil
}
