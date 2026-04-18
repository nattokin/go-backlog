package core

import (
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/model"
)

// WithMailAddress returns a option that sets the `mailAddress` field.
func (s *OptionService) WithMailAddress(mailAddress string) RequestOption {
	// ToDo: validate mailAddress (Note: The validation remains as simple not-empty check)
	return nonEmptyStringOption(ParamMailAddress, mailAddress)
}

// WithPassword returns a option that sets the `password` field.
func (s *OptionService) WithPassword(password string) RequestOption {
	return &APIParamOption{
		Type: ParamPassword,
		CheckFunc: func() error {
			if len(password) < 8 {
				return NewValidationError("password must be at least 8 characters long")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamPassword.Value(), password)
			return nil
		},
	}
}

// WithRoleType returns a option that sets the `roleType` field.
func (s *OptionService) WithRoleType(roleType model.Role) RequestOption {
	return &APIParamOption{
		Type: ParamRoleType,
		CheckFunc: func() error {
			if roleType < 1 || 6 < roleType {
				return NewValidationError("roleType must be between 1 and 6")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamRoleType.Value(), strconv.Itoa(int(roleType)))
			return nil
		},
	}
}

// WithSendMail returns a option to specify whether to send an invitation email.
func (s *OptionService) WithSendMail(enabled bool) RequestOption {
	return boolOption(ParamSendMail, enabled)
}

// WithUserID returns a option to set the user's ID.
func (s *OptionService) WithUserID(id int) RequestOption {
	return positiveIntOption(ParamUserID, id)
}
