package core

import "github.com/nattokin/go-backlog/internal/model"

// WithCount returns an option to set the `count` parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return intRangeOption(ParamCount, count, 1, 100)
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

// WithRoleType returns a option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType model.Role) RequestOption {
	return intRangeOption(ParamRoleType, int(roleType), 1, 6)
}

// WithUserID returns a option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return positiveIntOption(ParamUserID, id)
}
