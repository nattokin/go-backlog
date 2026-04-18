package core

import (
	"fmt"

	"github.com/nattokin/go-backlog/internal/model"
)

//
// ──────────────────────────────────────────────────────────────
//  OptionService — unified builder
// ──────────────────────────────────────────────────────────────
//

// OptionService provides builders for request options.
// Each XxxOptionService selectively exposes only the valid methods.
type OptionService struct{}

// WithCount returns an option to set the `count` parameter.
func (s *OptionService) WithCount(count int) RequestOption {
	return intRangeOption(ParamCount, count, 1, 100)
}

// WithKeyword returns an option to set the `keyword` parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &APIParamOption{
		Type:    ParamKeyword,
		SetFunc: setStringFunc(ParamKeyword, keyword),
	}
}

// WithName returns a option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return nonEmptyStringOption(ParamName, name)
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
		SetFunc: setStringFunc(ParamOrder, string(order)),
	}
}
