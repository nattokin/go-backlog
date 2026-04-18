package core

import (
	"fmt"
	"net/url"
	"strconv"
)

// WithActivityTypeIDs returns an option to set multiple `activityTypeId[]` parameters.
func (s *OptionService) WithActivityTypeIDs(typeIDs []int) RequestOption {
	return &APIParamOption{
		Type: ParamActivityTypeIDs,
		CheckFunc: func() error {
			for _, id := range typeIDs {
				if err := validateActivityID(id, "activityTypeIds"); err != nil {
					return err
				}
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			for _, id := range typeIDs {
				v.Add(ParamActivityTypeIDs.Value(), strconv.Itoa(id))
			}
			return nil
		},
	}
}

// WithMaxID returns an option to set the `maxId` parameter.
func (s *OptionService) WithMaxID(id int) RequestOption {
	return &APIParamOption{
		Type: ParamMaxID,
		CheckFunc: func() error {
			return validateActivityID(id, "maxID")
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamMaxID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// WithMinID returns an option to set the `minId` parameter.
func (s *OptionService) WithMinID(id int) RequestOption {
	return &APIParamOption{
		Type: ParamMinID,
		CheckFunc: func() error {
			return validateActivityID(id, "minID")
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamMinID.Value(), strconv.Itoa(id))
			return nil
		},
	}
}

// validateActivityID ensures that the given activity ID is within the valid range [1, 26].
func validateActivityID(id int, key string) error {
	if id < 1 || id > MaxActivityTypeID {
		return NewValidationError(fmt.Sprintf("invalid %s: must be between 1 and %d", key, MaxActivityTypeID))
	}
	return nil
}
