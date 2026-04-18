package core

import (
	"fmt"
	"net/url"
	"strconv"

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
	return &APIParamOption{
		Type: ParamCount,
		CheckFunc: func() error {
			if count < 1 || 100 < count {
				return NewValidationError("count must be between 1 and 100")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamCount.Value(), strconv.Itoa(count))
			return nil
		},
	}
}

// WithKeyword returns an option to set the `keyword` parameter.
func (s *OptionService) WithKeyword(keyword string) RequestOption {
	return &APIParamOption{
		Type: ParamKeyword,
		SetFunc: func(v url.Values) error {
			v.Set(ParamKeyword.Value(), keyword)
			return nil
		},
	}
}

// WithName returns a option that sets the `name` field.
func (s *OptionService) WithName(name string) RequestOption {
	return &APIParamOption{
		Type: ParamName,
		CheckFunc: func() error {
			if name == "" {
				return NewValidationError("name must not be empty")
			}
			return nil
		},
		SetFunc: func(v url.Values) error {
			v.Set(ParamName.Value(), name)
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
