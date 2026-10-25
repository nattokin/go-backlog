package option

import (
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/validate"
	"github.com/nattokin/go-backlog/internal/validation"
)

func (s *OptionService) WithAssigneeID(id int) *APIParamOption {
	return positiveIntOption(ParamAssigneeID, id)
}

func (s *OptionService) WithCommentID(id int) *APIParamOption {
	return positiveIntOption(ParamCommentID, id)
}

// WithCount sets `count`. Valid range: 1–100.
func (s *OptionService) WithCount(count int) *APIParamOption {
	return intRangeOption(ParamCount, count, 1, 100)
}

// WithSharedFileCount sets `count` for shared file list endpoints. Valid range: 1–1000.
func (s *OptionService) WithSharedFileCount(count int) *APIParamOption {
	return intRangeOption(ParamCount, count, 1, 1000)
}

func (s *OptionService) WithFieldType(fieldType int) *APIParamOption {
	return positiveIntOption(ParamTypeID, int(fieldType))
}

// WithInitialShift sets `initialShift` for Date type custom fields.
func (s *OptionService) WithInitialShift(days int) *APIParamOption {
	return &APIParamOption{
		Type:    ParamInitialShift,
		SetFunc: setIntFunc(ParamInitialShift, days),
	}
}

// WithInitialValueType sets `initialValueType` for Date type custom fields.
// 1: Today, 2: Today + initialShift days, 3: Specified date.
func (s *OptionService) WithInitialValueType(initialValueType int) *APIParamOption {
	return intRangeOption(ParamInitialValueType, initialValueType, 1, 3)
}

func (s *OptionService) WithIssueID(id int) *APIParamOption {
	return positiveIntOption(ParamIssueID, id)
}

func (s *OptionService) WithIssueTypeID(id int) *APIParamOption {
	return positiveIntOption(ParamIssueTypeID, id)
}

func (s *OptionService) WithMaxID(id int) *APIParamOption {
	return positiveIntOption(ParamMaxID, id)
}

func (s *OptionService) WithMinID(id int) *APIParamOption {
	return positiveIntOption(ParamMinID, id)
}

func (s *OptionService) WithOffset(offset int) *APIParamOption {
	return &APIParamOption{
		Type: ParamOffset,
		CheckFunc: func() *validation.Error {
			if offset < 0 {
				return validation.NewError(ParamOffset.Value(), "offset must not be negative")
			}
			return nil
		},
		SetFunc: setIntFunc(ParamOffset, offset),
	}
}

// WithParentChild sets `parentChild`. Valid range: 0–10.
// 0: All, 1: Not child, 2: Child or grandchild, 3: Standalone, 4: Has children,
// 5: Grandchild only, 6: Child only, 7: Top level only, 8: Exclude grandchild,
// 9: Exclude top level, 10: Leaf only.
func (s *OptionService) WithParentChild(parentChild int) *APIParamOption {
	return intRangeOption(ParamParentChild, parentChild, 0, 10)
}

func (s *OptionService) WithParentIssueID(id int) *APIParamOption {
	return positiveIntOption(ParamParentIssueID, id)
}

func (s *OptionService) WithPriorityID(id int) *APIParamOption {
	return positiveIntOption(ParamPriorityID, id)
}

func (s *OptionService) WithPullRequestCommentID(id int) *APIParamOption {
	return positiveIntOption(ParamPullRequestCommentID, id)
}

func (s *OptionService) WithPullRequestID(id int) *APIParamOption {
	return positiveIntOption(ParamPullRequestID, id)
}

func (s *OptionService) WithResolutionID(id int) *APIParamOption {
	return positiveIntOption(ParamResolutionID, id)
}

// WithRoleType sets `roleType`. Valid range: 1–6.
func (s *OptionService) WithRoleType(roleType int) *APIParamOption {
	return intRangeOption(ParamRoleType, int(roleType), 1, 6)
}

func (s *OptionService) WithStatusID(id int) *APIParamOption {
	return positiveIntOption(ParamStatusID, id)
}

func (s *OptionService) WithUserID(id int) *APIParamOption {
	return positiveIntOption(ParamUserID, id)
}

func (s *OptionService) WithWikiID(id int) *APIParamOption {
	return positiveIntOption(ParamWikiID, id)
}

// intRangeOption builds a *APIParamOption that validates an int is within [min, max].
func intRangeOption(paramType APIParamOptionType, value, min, max int) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			return validate.ValidateIntRange(paramType.Value(), value, min, max)
		},
		SetFunc: setIntFunc(paramType, value),
	}
}

// positiveIntOption builds a *APIParamOption that validates an int is >= 1.
func positiveIntOption(paramType APIParamOptionType, value int) *APIParamOption {
	return &APIParamOption{
		Type: paramType,
		CheckFunc: func() *validation.Error {
			return validate.ValidatePositiveInt(paramType.Value(), value)
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
