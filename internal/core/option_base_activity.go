package core

import (
	"fmt"
	"net/url"
	"strconv"

	"github.com/nattokin/go-backlog/internal/model"
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

// WithOrder returns an option to set the `order` parameter.
func (s *OptionService) WithOrder(order model.Order) RequestOption {
	return &APIParamOption{
		Type: ParamOrder,
		CheckFunc: func() error {
			if order != model.OrderAsc && order != model.OrderDesc {
				msg := fmt.Sprintf("order must be only '%s' or '%s'", string(model.OrderAsc), string(model.OrderDesc))
				return NewValidationError(msg)
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamOrder.Value(), string(order))
			return nil
		},
	}
}
