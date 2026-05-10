package core

import (
	"fmt"
	"net/url"
	"strconv"
)

func (s *OptionService) WithAssigneeID(id int) RequestOption {
	return positiveIntOption(ParamAssigneeID, id)
}

func (s *OptionService) WithCommentID(id int) RequestOption {
	return positiveIntOption(ParamCommentID, id)
}

// WithCount sets `count`. Valid range: 1–100.
func (s *OptionService) WithCount(count int) RequestOption {
	return intRangeOption(ParamCount, count, 1, 100)
}

func (s *OptionService) WithFieldType(fieldType int) RequestOption {
	return positiveIntOption(ParamTypeID, int(fieldType))
}

// WithInitialShift sets `initialShift` for Date type custom fields.
// Used when initialValueType is 2 (today + N days).
func (s *OptionService) WithInitialShift(days int) RequestOption {
	return &APIParamOption{
		Type:    ParamInitialShift,
		SetFunc: setIntFunc(ParamInitialShift, days),
	}
}

// WithInitialValueType sets `initialValueType` for Date type custom fields.
// 1: Today, 2: Today + initialShift days, 3: Specified date.
func (s *OptionService) WithInitialValueType(initialValueType int) RequestOption {
	return intRangeOption(ParamInitialValueType, initialValueType, 1, 3)
}

func (s *OptionService) WithIssueID(id int) RequestOption {
	return positiveIntOption(ParamIssueID, id)
}

func (s *OptionService) WithIssueTypeID(id int) RequestOption {
	return positiveIntOption(ParamIssueTypeID, id)
}

// WithMaxActivityTypeID sets `maxId` for activity type filtering. Valid range: 1–26.
func (s *OptionService) WithMaxActivityTypeID(id int) RequestOption {
	return intRangeOption(ParamMaxID, id, 1, MaxActivityTypeID)
}

// WithMinActivityTypeID sets `minId` for activity type filtering. Valid range: 1–26.
func (s *OptionService) WithMinActivityTypeID(id int) RequestOption {
	return intRangeOption(ParamMinID, id, 1, MaxActivityTypeID)
}

func (s *OptionService) WithMaxID(id int) RequestOption {
	return positiveIntOption(ParamMaxID, id)
}

func (s *OptionService) WithMinID(id int) RequestOption {
	return positiveIntOption(ParamMinID, id)
}

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

// WithParentChild sets `parentChild`.
// 0: All, 1: Exclude Child Issue, 2: Child Issue, 3: Neither Parent nor Child, 4: Parent Issue.
func (s *OptionService) WithParentChild(parentChild int) RequestOption {
	return intRangeOption(ParamParentChild, parentChild, 0, 4)
}

func (s *OptionService) WithParentIssueID(id int) RequestOption {
	return positiveIntOption(ParamParentIssueID, id)
}

func (s *OptionService) WithPriorityID(id int) RequestOption {
	return positiveIntOption(ParamPriorityID, id)
}

func (s *OptionService) WithPullRequestCommentID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestCommentID, id)
}

func (s *OptionService) WithPullRequestID(id int) RequestOption {
	return positiveIntOption(ParamPullRequestID, id)
}

func (s *OptionService) WithResolutionID(id int) RequestOption {
	return positiveIntOption(ParamResolutionID, id)
}

// WithRoleType sets `roleType`. Valid range: 1–6.
func (s *OptionService) WithRoleType(roleType int) RequestOption {
	return intRangeOption(ParamRoleType, int(roleType), 1, 6)
}

func (s *OptionService) WithStatusID(id int) RequestOption {
	return positiveIntOption(ParamStatusID, id)
}

func (s *OptionService) WithUserID(id int) RequestOption {
	return positiveIntOption(ParamUserID, id)
}

func (s *OptionService) WithWikiID(id int) RequestOption {
	return positiveIntOption(ParamWikiID, id)
}

// intRangeOption builds a RequestOption that validates an int is within [min, max] and sets it.
func intRangeOption(paramType APIParamOptionType, value, min, max int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value < min || value > max {
				return NewValidationError(fmt.Sprintf("%s must be between %d and %d", paramType.Value(), min, max))
			}
			return nil
		},
		SetFunc: setIntFunc(paramType, value),
	}
}

// positiveIntOption builds a RequestOption that validates an int is >= 1 and sets it.
func positiveIntOption(paramType APIParamOptionType, value int) RequestOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() error {
			if value < 1 {
				return NewValidationError(fmt.Sprintf("invalid %s: must not be less than 1", paramType.Value()))
			}
			return nil
		},
		SetFunc: setIntFunc(paramType, value),
	}
}

func setIntFunc(key APIParamOptionType, value int) func(url.Values) error {
	return func(v url.Values) error {
		v.Set(key.Value(), strconv.Itoa(value))
		return nil
	}
}
