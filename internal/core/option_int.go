package core

import "github.com/nattokin/go-backlog/internal/model"

// WithActualHours returns an option to set the `actualHours` parameter.
func (s *OptionService) WithActualHours(hours int) RequestOption {
	return positiveIntOption(ParamActualHours, hours)
}

// WithAssigneeID returns an option to set the `assigneeId` parameter.
func (s *OptionService) WithAssigneeID(id int) RequestOption {
	return positiveIntOption(ParamAssigneeID, id)
}

// WithCommentID returns an option to set the `commentId` parameter.
func (s *OptionService) WithCommentID(id int) RequestOption {
	return positiveIntOption(ParamCommentID, id)
}

// WithCount returns an option to set the `count` parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return intRangeOption(ParamCount, count, 1, 100)
}

// WithEstimatedHours returns an option to set the `estimatedHours` parameter.
func (s *OptionService) WithEstimatedHours(hours int) RequestOption {
	return positiveIntOption(ParamEstimatedHours, hours)
}

// WithIssueID returns an option to set the `issueId` parameter.
func (s *OptionService) WithIssueID(id int) RequestOption {
	return positiveIntOption(ParamIssueID, id)
}

// WithIssueTypeID returns an option to set the `issueTypeId` parameter.
func (s *OptionService) WithIssueTypeID(id int) RequestOption {
	return positiveIntOption(ParamIssueTypeID, id)
}

// WithMaxID returns an option to set the `maxId` parameter.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return intRangeOption(ParamMaxID, id, 1, MaxActivityTypeID)
}

// WithMinID returns an option to set the `minId` parameter.
func (s *OptionService) WithMinID(id int) RequestOption {
	return intRangeOption(ParamMinID, id, 1, MaxActivityTypeID)
}

// WithOffset returns an option to set the `offset` parameter.
func (s *OptionService) WithOffset(offset int) RequestOption {
	return &APIParamOption{
		Type: ParamOffset,
		CheckFunc: func() error {
			if offset < 0 {
				return NewValidationError("offset must not be negative")
			}
			return nil
		},
		SetFunc: setIntFunc(ParamOffset, offset),
	}
}

// WithParentChild returns an option to set the `parentChild` parameter.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *OptionService) WithParentChild(parentChild int) RequestOption {
	return intRangeOption(ParamParentChild, parentChild, 0, 4)
}

// WithParentIssueID returns an option to set the `parentIssueId` parameter.
func (s *OptionService) WithParentIssueID(id int) RequestOption {
	return positiveIntOption(ParamParentIssueID, id)
}

// WithPriorityID returns an option to set the `priorityId` parameter.
func (s *OptionService) WithPriorityID(id int) RequestOption {
	return positiveIntOption(ParamPriorityID, id)
}

// WithPullRequestCommentID returns an option to set the `pullRequestCommentId` parameter.
func (s *OptionService) WithPullRequestCommentID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestCommentID, id)
}

// WithPullRequestID returns an option to set the `pullRequestId` parameter.
func (s *OptionService) WithPullRequestID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestID, id)
}

// WithResolutionID returns an option to set the `resolutionId` parameter.
func (s *OptionService) WithResolutionID(id int) RequestOption {
	return positiveIntOption(ParamResolutionID, id)
}

// WithRoleType returns a option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType model.Role) RequestOption {
	return intRangeOption(ParamRoleType, int(roleType), 1, 6)
}

// WithStatusID returns an option to set the `statusId` parameter.
func (s *OptionService) WithStatusID(id int) RequestOption {
	return positiveIntOption(ParamStatusID, id)
}

// WithUserID returns a option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return positiveIntOption(ParamUserID, id)
}

// WithWikiID returns an option to set the `wikiPageId` parameter.
func (s *OptionService) WithWikiID(id int) RequestOption {
	return positiveIntOption(ParamWikiID, id)
}
