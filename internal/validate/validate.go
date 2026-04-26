package validate

import "github.com/nattokin/go-backlog/internal/core"

func ValidateActivityID(activityID int) error {
	if activityID < 1 {
		return core.NewValidationError("activityID must not be less than 1")
	}
	return nil
}

func ValidateAttachmentID(attachmentID int) error {
	if attachmentID < 1 {
		return core.NewValidationError("attachmentID must not be less than 1")
	}
	return nil
}

func ValidateCommentID(commentID int) error {
	if commentID < 1 {
		return core.NewValidationError("commentID must not be less than 1")
	}
	return nil
}

func ValidateIssueIDOrKey(issueIDOrKey string) error {
	if issueIDOrKey == "" {
		return core.NewValidationError("issueIDOrKey must not be empty")
	}
	if issueIDOrKey == "0" {
		return core.NewValidationError("issueIDOrKey must not be '0'")
	}
	return nil
}

func ValidateProjectID(projectID int) error {
	if projectID < 1 {
		return core.NewValidationError("projectID must not be less than 1")
	}
	return nil
}

func ValidateProjectIDOrKey(projectIDOrKey string) error {
	if projectIDOrKey == "" {
		return core.NewValidationError("projectIDOrKey must not be empty")
	}
	if projectIDOrKey == "0" {
		return core.NewValidationError("projectIDOrKey must not be '0'")
	}
	return nil
}

func ValidatePRNumber(prNumber int) error {
	if prNumber < 1 {
		return core.NewValidationError("prNumber must not be less than 1")
	}
	return nil
}

func ValidateRepositoryIDOrName(repositoryIDOrName string) error {
	if repositoryIDOrName == "" {
		return core.NewValidationError("repositoryIDOrName must not be empty")
	}
	if repositoryIDOrName == "0" {
		return core.NewValidationError("repositoryIDOrName must not be '0'")
	}
	return nil
}

func ValidateSharedFileID(fileID int) error {
	if fileID < 1 {
		return core.NewValidationError("fileID must not be less than 1")
	}
	return nil
}

func ValidateStarID(starID int) error {
	if starID < 1 {
		return core.NewValidationError("starID must not be less than 1")
	}
	return nil
}

func ValidateUserID(userID int) error {
	if userID < 1 {
		return core.NewValidationError("userID must not be less than 1")
	}
	return nil
}

func ValidateWikiID(wikiID int) error {
	if wikiID < 1 {
		return core.NewValidationError("wikiID must not be less than 1")
	}
	return nil
}
