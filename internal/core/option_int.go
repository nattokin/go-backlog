package core

import "strconv"

// WithActivityTypeIDs sets the activity type IDs filter.
func (s *OptionService) WithActivityTypeIDs(ids []int) RequestOption {
	return intSliceOption(ParamActivityTypeIDs, "activityTypeId", ids)
}

// WithAssigneeID sets the assignee user ID.
func (s *OptionService) WithAssigneeID(id int) RequestOption {
	return positiveIntOption(ParamAssigneeID, id)
}

// WithCommentID sets the comment ID for Add Star.
func (s *OptionService) WithCommentID(id int) RequestOption {
	return positiveIntOption(ParamCommentID, id)
}

// WithCount sets the number of results to return.
func (s *OptionService) WithCount(count int) RequestOption {
	return intRangeOption(ParamCount, count, 1, 100)
}

// WithIssueID sets the issue ID for Add Star.
func (s *OptionService) WithIssueID(id int) RequestOption {
	return positiveIntOption(ParamIssueID, id)
}

// WithIssueTypeID sets the issue type ID.
func (s *OptionService) WithIssueTypeID(id int) RequestOption {
	return positiveIntOption(ParamIssueTypeID, id)
}

// WithMaxID sets the maximum activity ID.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return positiveIntOption(ParamMaxID, id)
}

// WithMinID sets the minimum activity ID.
func (s *OptionService) WithMinID(id int) RequestOption {
	return positiveIntOption(ParamMinID, id)
}

// WithOffset sets the offset for pagination.
func (s *OptionService) WithOffset(offset int) RequestOption {
	return &APIParamOption{
		Type: ParamOffset,
		CheckFunc: func() error {
			if offset < 0 {
				return NewValidationError("offset must not be negative")
			}
			return nil
		},
		SetFunc: func(v Values) error {
			v.Set(ParamOffset.Value(), strconv.Itoa(offset))
			return nil
		},
	}
}

// WithParentIssueID sets the parent issue ID.
func (s *OptionService) WithParentIssueID(id int) RequestOption {
	return positiveIntOption(ParamParentIssueID, id)
}

// WithPriorityID sets the priority ID.
func (s *OptionService) WithPriorityID(id int) RequestOption {
	return positiveIntOption(ParamPriorityID, id)
}

// WithPullRequestCommentID sets the pull request comment ID for Add Star.
func (s *OptionService) WithPullRequestCommentID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestCommentID, id)
}

// WithPullRequestID sets the pull request ID for Add Star.
func (s *OptionService) WithPullRequestID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestID, id)
}

// WithResolutionID sets the resolution ID.
func (s *OptionService) WithResolutionID(id int) RequestOption {
	return positiveIntOption(ParamResolutionID, id)
}

// WithRoleType sets the role type.
func (s *OptionService) WithRoleType(role Role) RequestOption {
	return positiveIntOption(ParamRoleType, int(role))
}

// WithStarID sets the star ID for Remove Star.
func (s *OptionService) WithStarID(id int) RequestOption {
	return positiveIntOption(ParamStarID, id)
}

// WithUserID sets the user ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return positiveIntOption(ParamUserID, id)
}

// WithVersionIDs sets the version IDs filter.
func (s *OptionService) WithVersionIDs(ids []int) RequestOption {
	return intSliceOption(ParamVersionIDs, "versionId", ids)
}

// WithWikiPageID sets the wiki page ID for Add Star.
func (s *OptionService) WithWikiPageID(id int) RequestOption {
	return positiveIntOption(ParamWikiPageID, id)
}
